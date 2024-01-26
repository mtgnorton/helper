package worker

import (
	"github.com/mtgnorton/helper/mlogger"
	"github.com/mtgnorton/helper/workerpool/task"
	"github.com/mtgnorton/helper/workerpool/util"
	"go.uber.org/multierr"
	"time"
)

// DefaultCreateWorkerFunc 默认创建工作者的方法
var DefaultCreateWorkerFunc CreateWorkerFunc = func(prepareFinish PrepareWorkFinishTask) Worker {
	return NewWorker(prepareFinish)
}

// Worker 工作者
type Worker interface {
	Name() string
	TaskChan() TransmitTaskChan // 返回接受任务的通道
	Run()
	Stop()
}

// CreateWorkerFunc 创建工作者的方法
type CreateWorkerFunc func(prepareFinish PrepareWorkFinishTask) Worker

type defaultWorker struct {
	opts     Options
	taskChan TransmitTaskChan
	stopChan chan struct{}
}

// NewWorker 创建一个工作者
func NewWorker(prepareFinish PrepareWorkFinishTask, opts ...Option) Worker {
	logger := mlogger.NewLogger(mlogger.WithLevel(mlogger.DebugLevel))
	options := Options{
		logger:                logger,
		process:               DefaultProcess,
		prepareWorkFinishTask: prepareFinish,
	}
	for _, o := range opts {
		o(&options)
	}
	if options.name == "" {
		options.name = "worker-" + util.GenerateRandomName(5)
	}
	return &defaultWorker{
		opts:     options,
		taskChan: make(TransmitTaskChan),
		stopChan: make(chan struct{}),
	}
}

// Name 返回工作者的名称
func (w *defaultWorker) Name() string {
	return w.opts.name
}

// TaskChan 返回接受任务的通道
func (w *defaultWorker) TaskChan() TransmitTaskChan {
	return w.taskChan
}

// Run 运行工作者
func (w *defaultWorker) Run() {
	go w.run()
}

// Stop 停止工作者
func (w *defaultWorker) Stop() {
	for {
		select {
		case <-w.stopChan:
			return
		default:
			close(w.stopChan)
		}
	}
}

func (w *defaultWorker) run() {
	logger := w.opts.logger
	for {
		select {
		case <-w.stopChan:
			logger.Logf(mlogger.DebugLevel, "[%s] worker stop", w.opts.name)
			return
		default:
			select {
			case w.opts.prepareWorkFinishTask.WorkerChan() <- w:
			case <-time.After(50 * time.Millisecond):
				continue
			}
			var t task.Task
			if t = w.receiveTask(); t == nil {
				continue
			}
			r, additionalTasks, err := w.callProcess(t)
			w.opts.prepareWorkFinishTask.ResultChan() <- task.NewDefaultResult(t, r, additionalTasks, err)
		}
	}
}

func (w *defaultWorker) callProcess(t task.Task) (interface{}, []task.Task, error) {
	var (
		r               interface{}
		additionalTasks []task.Task
		err             error
	)
	w.opts.logger.Logf(mlogger.DebugLevel, "[%s] worker process task:%v", w.opts.name, t)

	retryNumber := 1
	if w.opts.retryNumber > 0 {
		retryNumber = w.opts.retryNumber
	}
	for i := 0; i < retryNumber; i++ {
		var subErr error
		if w.opts.timeout == 0 {
			r, additionalTasks, subErr = w.opts.process(t)
		} else {
			r, additionalTasks, subErr = w.processTimeout(t)
		}
		if subErr != nil {
			err = multierr.Append(err, subErr)
		} else {
			break
		}
	}
	return r, additionalTasks, err
}

func (w *defaultWorker) processTimeout(t task.Task) (r interface{}, additionalTasks []task.Task, err error) {
	finishChan := make(chan struct{}, 1)

	go func() {
		r, additionalTasks, err = w.opts.process(t)
		finishChan <- struct{}{}
	}()
	select {
	case <-time.After(w.opts.timeout):
		return nil, nil, ErrTimeout
	case <-finishChan:
		return
	}
}

func (w *defaultWorker) receiveTask() task.Task {
	select {
	case t := <-w.taskChan:
		return t
	case <-time.After(50 * time.Millisecond):
		return nil
	}
}
