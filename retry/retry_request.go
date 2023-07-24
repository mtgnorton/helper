package retry

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// RetryRequestWithTimeout 通过 context 上下文对象重试请求, 指定超时时间
func RetryRequestWithTimeout(ctx context.Context, timeout time.Duration, fn func(ctx context.Context) error) error {
	return RetryRequest(ctx, timeout, 3, 0, fn)
}

// RetryRequestWithTimeoutAndTimes 通过 context 上下文对象重试请求,指定超时时间,重试次数
func RetryRequestWithTimeoutAndTimes(ctx context.Context, timeout time.Duration, retryTimes int64, fn func(ctx context.Context) error) error {
	return RetryRequest(ctx, timeout, retryTimes, 0, fn)
}

// RetryRequest 通过 context 上下文对象重试请求, 指定超时时间,重试次数,重试间隔
func RetryRequest(ctx context.Context, timeout time.Duration, retryTimes int64, interval time.Duration, fn func(ctx context.Context) error) error {

	var i int64
	for i = 0; i < retryTimes; i++ {
		ctx, _ := context.WithTimeout(ctx, timeout)
		err := fn(ctx)
		// 当err为context deadline exceeded 时，重试
		fmt.Println(i, err)
		if err != nil && strings.Contains(err.Error(), "context deadline exceeded") {
			if interval > 0 {
				time.Sleep(interval)
			}
			continue
		}
		return err
	}
	return fmt.Errorf("retry fail")
}
