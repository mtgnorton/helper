package timeout

import (
	"github.com/pkg/errors"
	"time"
)

type Called func() Result

type Result struct {
	R   interface{}
	Err error
}

var Err = errors.New("called timeout")

// Do 在指定时间内执行函数，如果超时则返回错误
func Do(called Called, timeDuration time.Duration) Result {
	ch := make(chan Result, 1)
	go func() {
		ch <- called()
	}()
	select {
	case <-ch:
		return Result{}
	case <-time.After(timeDuration):
		return Result{Err: Err}
	}
}
