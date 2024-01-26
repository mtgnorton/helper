package worker

import (
	"fmt"
	"github.com/mtgnorton/helper/workerpool/task"
	"github.com/pkg/errors"
)

var DefaultProcess Process = func(task task.Task) (interface{}, []task.Task, error) {
	if task == nil {
		panic("task is nil")
	}
	fmt.Println(task)
	return nil, nil, nil
}

var ErrTimeout = errors.New("process timeout")

type (
	TransmitTaskChan   chan task.Task
	TransmitWorkerChan chan Worker

	TransmitResultChan chan task.Result
)

type Process func(task task.Task) (interface{}, []task.Task, error)

type PrepareWorkFinishTask interface {
	WorkerChan() TransmitWorkerChan // worker准备完成后通过通道传输到pool
	ResultChan() TransmitResultChan // 完成的task通过通道传输到pool

}
