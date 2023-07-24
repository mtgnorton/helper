package retry

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func Test_RetryFixedInterval(t *testing.T) {
	retryIntervals := []int64{
		1,
		4,
		10,
		30,
		60,
	}
	var retryPeriod int64 = 60
	ctx := context.Background()
	RetryFixedInterval(ctx, retryIntervals, retryPeriod, func(ctx context.Context) {
		fmt.Printf("timestr: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	})
}
