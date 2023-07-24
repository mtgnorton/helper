package retry

import (
	"context"
	"fmt"
	"time"
)

// RetryFixedInterval 在指定的 retryPeriod 时间内, 按照 retryIntervals 重试
// retryIntervals 重试间隔,如[]int{1,2,3} 表示第一次重试间隔1秒,第二次重试间隔2秒,第三次重试间隔3秒
// retryPeriod 重试周期,如 retryPeriod=10,表示10秒内重试,超过10秒后,重试次数清零,如果为 0,重试完成后直接退出
// callback 重试回调函数
func RetryFixedInterval(ctx context.Context, retryIntervals []int64, retryPeriod int64, callback func(ctx context.Context) error) (r error) {

	if len(retryIntervals) == 0 {
		return
	}

	lastConnectTime := time.Now().Unix()
	retryTimes := 0
	for {
		now := time.Now().Unix()
		if retryTimes >= len(retryIntervals) {
			if retryPeriod <= 0 {
				return
			}
			if now-lastConnectTime < retryPeriod {
				time.Sleep(time.Second * 10)
				continue
			} else {
				retryTimes = 0
				lastConnectTime = now
			}
		}
		if err := callback(ctx); err != nil {
			r = fmt.Errorf("%w", fmt.Errorf("%d times retry fail,error:%s", retryTimes, err.Error()))
		}
		// 休眠
		time.Sleep(time.Second * time.Duration(retryIntervals[retryTimes]))
		retryTimes++
	}
}
