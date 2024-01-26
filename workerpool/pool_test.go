package workerpool

import (
	"fmt"
	statistics2 "github.com/mtgnorton/helper/workerpool/statistics"
	"github.com/mtgnorton/helper/workerpool/task"
	"github.com/mtgnorton/helper/workerpool/worker"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestPool_default(t *testing.T) {

	pool := NewPool(
		WithMaxTaskNumber(100),
	)
	pool.Run()
	for i := 0; i < 100; i++ {
		pool.SubmitTask(task.NewDefaultTask("task-" + fmt.Sprint(i)))
	}
	<-pool.Wait()

	statistic := pool.GetStatistics()
	if statistic.ReceiveTaskNumber() != 100 {
		t.Errorf("receive task number = %d, want 100", statistic.ReceiveTaskNumber())
	}
	if statistic.FailTaskNumber() != 0 {
		t.Errorf("fail task number = %d, want 0", statistic.FailTaskNumber())
	}
	if statistic.SuccessTaskNumber() != 100 {
		t.Errorf("success task number = %d, want 100", statistic.SuccessTaskNumber())
	}
}

func TestPool_stop(t *testing.T) {
	var wf worker.CreateWorkerFunc = func(prepareFinish worker.PrepareWorkFinishTask) worker.Worker {
		return worker.NewWorker(prepareFinish,
			worker.WithProcess(func(task task.Task) (interface{}, []task.Task, error) {
				// 随机休眠 800ms - 2000ms
				randomNumber := rand.Intn(1200) + 800
				time.Sleep(time.Duration(randomNumber) * time.Millisecond)

				return nil, nil, nil
			}),
		)
	}
	pool := NewPool(
		WithName("test-stop-pool"),
		WithCreateWorkerFunc(wf),
	)
	pool.Run()
	for i := 0; i < 100; i++ {
		pool.SubmitTask(task.NewDefaultTask("task-" + fmt.Sprint(i)))
	}
	go func() {
		time.Sleep(time.Second)
		pool.Stop()
	}()
	<-pool.Wait()
}

func TestPool_retry(t *testing.T) {
	var retryStatistics = make(map[int]int)
	var mu sync.Mutex

	var wf worker.CreateWorkerFunc = func(prepareFinish worker.PrepareWorkFinishTask) worker.Worker {
		return worker.NewWorker(prepareFinish,
			worker.WithProcess(func(task task.Task) (interface{}, []task.Task, error) {
				d := task.Data().(int)
				defer func() {
					mu.Lock()
					retryStatistics[d]++
					mu.Unlock()
				}()
				if d&1 == 1 {
					return nil, nil, nil
				}

				return nil, nil, fmt.Errorf("test error")
			}),
			worker.WithRetry(3),
		)
	}
	pool := NewPool(
		WithName("test-retry-pool"),
		WithCreateWorkerFunc(wf),
		WithMaxTaskNumber(100),
	)
	pool.Run()
	for i := 0; i < 100; i++ {
		pool.SubmitTask(task.NewDefaultTask(i))
	}
	<-pool.Wait()
	for k, v := range retryStatistics {
		if k&1 == 1 {
			if v != 1 {
				t.Errorf("retry %d times, want 3 times", v)
			}
		} else {
			if v != 3 {
				t.Errorf("retry %d times, want 1 times", v)
			}
		}
	}
}

func TestPool_timeout(t *testing.T) {
	var wf worker.CreateWorkerFunc = func(prepareFinish worker.PrepareWorkFinishTask) worker.Worker {
		return worker.NewWorker(prepareFinish,
			worker.WithProcess(func(task task.Task) (interface{}, []task.Task, error) {
				t := task.Data().(int)
				if t&1 == 1 { // 奇数休眠100ms
					time.Sleep(100 * time.Millisecond)
					return nil, nil, nil
				}
				return nil, nil, nil
			}),
			worker.WithTimeout(time.Millisecond*50),
		)
	}

	pool := NewPool(
		WithName("test-timeout-pool"),
		WithMaxTaskNumber(100),
		WithCreateWorkerFunc(wf),
		WithStatistics(statistics2.NewStatistics(statistics2.WithSaveResult(true))),
	)
	pool.Run()
	for i := 0; i < 100; i++ {
		pool.SubmitTask(task.NewDefaultTask(i))
	}
	<-pool.Wait()
	statistics := pool.GetStatistics()
	errNumber := 0
	for _, v := range statistics.AllResult() {
		if v.Err() != nil {
			errNumber++
			if v.Err() != worker.ErrTimeout {
				t.Errorf("err = %v, want %v", v.Err(), worker.ErrTimeout)
			}
		}
	}
	if errNumber != 50 {
		t.Errorf("err number = %d, want 50", errNumber)
	}
}

func TestPool_AdditionalTask(t *testing.T) {
	var wf worker.CreateWorkerFunc = func(prepareFinish worker.PrepareWorkFinishTask) worker.Worker {
		return worker.NewWorker(prepareFinish,
			worker.WithProcess(func(t task.Task) (interface{}, []task.Task, error) {
				d := t.Data().(int)
				if d&1 == 1 { // 奇数返回一个新任务
					tasks := []task.Task{task.NewDefaultTask(1000 + d)}
					return nil, tasks, nil
				}
				return nil, nil, nil
			}),
			// worker.WithTimeout(time.Millisecond*50),
		)
	}
	pool := NewPool(
		WithName("test-AdditionalTask-pool"),
		WithCreateWorkerFunc(wf),
		WithMaxTaskNumber(150),
		WithStatistics(statistics2.NewStatistics(statistics2.WithSaveResult(true))),
	)
	pool.Run()
	for i := 0; i < 100; i++ {
		pool.SubmitTask(task.NewDefaultTask(i))
	}

	<-pool.Wait()

	statistics := pool.GetStatistics()
	if statistics.FinishTaskNumber() != 150 {
		t.Errorf("finish task number = %d, want 150", statistics.FinishTaskNumber())
	}

}

func TestPool_deadlock(t *testing.T) {

	for i := 0; i < 100; i++ {
		var wf worker.CreateWorkerFunc = func(prepareFinish worker.PrepareWorkFinishTask) worker.Worker {
			return worker.NewWorker(prepareFinish,
				worker.WithProcess(func(t task.Task) (interface{}, []task.Task, error) {
					d := t.Data().(int)
					if d&1 == 1 { // 奇数返回一个新任务
						tasks := []task.Task{task.NewDefaultTask(1000 + d)}
						return nil, tasks, nil
					}
					return nil, nil, nil
				}),
			)
		}
		pool := NewPool(
			WithCreateWorkerFunc(wf),
			WithMaxTaskNumber(50),
			WithStatistics(statistics2.NewStatistics(statistics2.WithSaveResult(true))),
		)
		pool.Run()
		for i := 0; i < 100; i++ {
			pool.SubmitTask(task.NewDefaultTask(i))
		}

		<-pool.Wait()

		statistic := pool.GetStatistics()
		if statistic.FinishTaskNumber() != 50 {
			t.Errorf("finish task number = %d, want 50", statistic.FinishTaskNumber())
		}
	}

}
