package single_flight

import (
	"context"
	"sync"
	"time"
)

// MultiSome 传入的多个task, 在规定时间内 第一个完成且没有错误 就返回该任务的结果,其他task 直接结束
type MultiSome interface {
	Exec(timeout time.Duration, fun ...TaskFun) (interface{}, error)
}

// TaskFun 处理逻辑中需要通过 ctx 来判断是否超时，如果超时需要在taskFun中返回
type TaskFun func(ctx context.Context) (interface{}, error)

type DefaultMultiSome struct {
	sync.Mutex
	doneChan   chan struct{}
	Timeout    time.Duration
	downOnce   sync.Once
	assignOnce sync.Once
}

func NewMultiSome() MultiSome {
	return &DefaultMultiSome{
		doneChan: make(chan struct{}),
	}
}

// Exec 在规定时间内，第一个完成且没有错误 就返回该任务的结果,其他task 直接结束
func (m *DefaultMultiSome) Exec(timeout time.Duration, fun ...TaskFun) (r interface{}, err error) {
	m.Timeout = timeout
	if len(fun) == 0 {
		return
	}
	ctx := context.Background()
	for _, f := range fun {
		go func(f TaskFun) {
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()
			go func() {
				internalR, internalErr := f(ctx)
				if internalErr == nil { // 在规定时间内，有可能多个任务同时成功，同时走到这里，只有第一次赋值允许
					m.assignOnce.Do(func() {
						m.done()
						m.Lock()
						r = internalR
						m.Unlock()
					})
				}
			}()
			select {
			case <-time.After(timeout): // 任务全部超时
				m.done()
				return
			case <-m.doneChan:
				return
			}
		}(f)
	}
	<-m.doneChan
	return
}

func (m *DefaultMultiSome) done() {
	m.downOnce.Do(func() {
		close(m.doneChan)
	})
}
