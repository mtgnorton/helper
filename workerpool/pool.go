package workerpool

import (
	"fmt"
	"github.com/mtgnorton/helper/mlogger"
	"github.com/mtgnorton/helper/workerpool/statistics"
	"github.com/mtgnorton/helper/workerpool/task"
	"github.com/mtgnorton/helper/workerpool/util"
	"github.com/mtgnorton/helper/workerpool/worker"
	"sync"
	"time"
)

type Status int

const (
	NoStarted Status = iota
	Running
	Stopping
)

// Pool 工作池
// 工作池是一个协程安全的工作池,可以提交任务,并发执行任务
type Pool interface {
	Init(...Option) error
	Options() Options
	SubmitTask(task ...task.Task)
	Run()
	Stop()
	GetStatus() Status
	GetStatistics() statistics.Statistics
	Wait() <-chan struct{} // 等待工作池停止
}

type defaultPool struct {
	opts       Options
	workerChan worker.TransmitWorkerChan // worker通道
	taskChan   worker.TransmitTaskChan   // 任务通道
	stopChan   chan struct{}
	resultChan worker.TransmitResultChan // 结果通道
	workers    []worker.Worker
	status     Status
	sync.Mutex
}

// NewPool 创建一个工作池
func NewPool(opts ...Option) Pool {
	options := Options{
		WorkerNumber:      10,
		logger:            mlogger.DefaultLogger,
		CreateWorkerFunc:  worker.DefaultCreateWorkerFunc,
		displayProcessGap: 5 * time.Second,
		statistics:        statistics.NewStatistics(),
	}
	for _, o := range opts {
		o(&options)
	}
	if options.name == "" {
		// 随机名称
		options.name = "pool-" + util.GenerateRandomName(5)
	}

	c := &defaultPool{
		opts:       options,
		workerChan: make(worker.TransmitWorkerChan),
		taskChan:   make(worker.TransmitTaskChan),
		resultChan: make(worker.TransmitResultChan),
		stopChan:   make(chan struct{}),
		status:     NoStarted,
	}
	return c
}

// Init 初始化工作池
func (c *defaultPool) Init(opts ...Option) error {
	for _, o := range opts {
		o(&c.opts)
	}
	return nil
}

// Options 返回工作池的配置
func (c *defaultPool) Options() Options {
	return c.opts
}

// SubmitTask 提交任务
func (c *defaultPool) SubmitTask(task ...task.Task) {

	for _, t := range task {
		select {
		case c.taskChan <- t:
		case <-c.stopChan:
			return
		}
	}
}

// Stop 停止工作池
func (c *defaultPool) Stop() {
	c.Lock()
	if c.status == Stopping {
		c.Unlock()
		return
	}
	c.status = Stopping
	c.Unlock()
	for _, w := range c.workers {
		w.Stop()
	}
	for {
		select {
		case <-c.stopChan:
			return
		default:
			close(c.stopChan)
		}
	}
}

// Run 运行工作池
func (c *defaultPool) Run() {
	c.Lock()
	if c.status == Running {
		c.Unlock()
		return
	}
	c.status = Running
	c.Unlock()
	go c.run()
	c.CreateWorkers()
}

// GetStatus 获取工作池状态 Running, Stopping, NoStarted
func (c *defaultPool) GetStatus() Status {
	c.Lock()
	defer c.Unlock()
	return c.status
}

// GetStatistics 获取工作池统计信息
func (c *defaultPool) GetStatistics() statistics.Statistics {
	return c.opts.statistics
}

// Wait 等待工作池停止
func (c *defaultPool) Wait() <-chan struct{} {
	return c.stopChan
}

func (c *defaultPool) run() {
	logger := c.opts.logger
	var (
		tasks   []task.Task
		workers []worker.Worker
	)
	progressTimer := time.NewTicker(c.opts.displayProcessGap)
	defer progressTimer.Stop()
	for {
		activeWorker, activeTask := getActiveWorkerAndTask(tasks, workers)

		select {
		case result := <-c.resultChan:
			c.opts.statistics.AddTaskResult(result)
			if c.reachFinishTaskMax() {
				c.Stop() // 如果完成任务数达到最大任务数,则停止工作池
				continue
			}
			if len(result.AdditionalTasks()) > 0 {
				go c.SubmitTask(result.AdditionalTasks()...)
			}
		case t := <-c.taskChan:
			if c.reachReceiveTaskMax() {
				continue
			}
			c.opts.statistics.AddReceiveNumber(1)
			tasks = append(tasks, t)
		case w := <-c.workerChan:
			workers = append(workers, w)

		case <-c.stopChan:
			c.displayProcess()
			logger.Logf(mlogger.DebugLevel, "[%s] pool stop", c.opts.name)
			return

		case <-progressTimer.C:
			c.displayProcess()
		default:
			if activeWorker != nil && activeTask != nil {
				tasks, workers = removeActiveWorkerAndTask(tasks, workers)
				activeWorker.TaskChan() <- activeTask
			}
		}
	}
}

// ResultChan 返回结果通道
func (c *defaultPool) ResultChan() worker.TransmitResultChan {
	return c.resultChan
}

// WorkerChan 返回 worker 通道
func (c *defaultPool) WorkerChan() worker.TransmitWorkerChan {
	return c.workerChan
}

func (c *defaultPool) CreateWorkers() {
	for i := 0; i < c.opts.WorkerNumber; i++ {
		w := c.opts.CreateWorkerFunc(c)
		c.workers = append(c.workers, w)
		w.Run()
	}
}
func (c *defaultPool) reachFinishTaskMax() bool {
	if c.opts.statistics.FinishTaskNumber() >= c.opts.MaxTaskNumber {
		c.opts.logger.Logf(mlogger.DebugLevel, "[%s] finish task number is max,max:%v", c.opts.name, c.opts.MaxTaskNumber)
		return true
	}
	return false
}

func (c *defaultPool) reachReceiveTaskMax() (reach bool) {

	if c.opts.MaxTaskNumber > 0 && c.opts.statistics.ReceiveTaskNumber() >= c.opts.MaxTaskNumber {
		//c.opts.logger.Logf(mlogger.DebugLevel, "[%s] receive task number is max,max:%v", c.opts.name, c.opts.MaxTaskNumber)
		return true
	}
	return
}

func (c *defaultPool) displayProcess() {
	logger := c.opts.logger

	var (
		receiveInfo = fmt.Sprintf("receive-task-number:%d ", c.opts.statistics.ReceiveTaskNumber())
		maxInfo     = fmt.Sprintf("max-task-number:%d ", c.opts.MaxTaskNumber)
		finishInfo  = fmt.Sprintf("finish-task-number:%d ", c.opts.statistics.FinishTaskNumber())
		successInfo = fmt.Sprintf("success-task-number:%d ", c.opts.statistics.SuccessTaskNumber())
		failInfo    = fmt.Sprintf("fail-task-number:%d ", c.opts.statistics.FailTaskNumber())
	)

	var successPercentInfo, finishPercentInfo string
	if c.opts.statistics.FinishTaskNumber() > 0 {
		successPercentInfo = fmt.Sprintf("success-task-percent: %.2f %% ", float64(c.opts.statistics.SuccessTaskNumber())/float64(c.opts.statistics.FinishTaskNumber())*100)
	}
	if c.opts.MaxTaskNumber > 0 {
		finishPercentInfo = fmt.Sprintf("finish-task-percent: %.2f %% ", float64(c.opts.statistics.FinishTaskNumber())/float64(c.opts.MaxTaskNumber)*100)
	}

	logger.Log(mlogger.InfoLevel, fmt.Sprintf("[%s]", c.opts.name), receiveInfo, maxInfo, finishInfo, successInfo, failInfo, successPercentInfo, finishPercentInfo)
}

func getActiveWorkerAndTask(tasks []task.Task, workers []worker.Worker) (worker.Worker, task.Task) {
	if len(tasks) > 0 && len(workers) > 0 {
		task.SortTask(tasks, 0, len(tasks)-1)
		return workers[0], tasks[len(tasks)-1]
	}
	return nil, nil
}

func removeActiveWorkerAndTask(tasks []task.Task, workers []worker.Worker) ([]task.Task, []worker.Worker) {
	tasks = tasks[0 : len(tasks)-1]
	workers = workers[1:]
	return tasks, workers
}
