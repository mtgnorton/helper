package workerpool

import (
	"context"
	"github.com/mtgnorton/helper/mlogger"
	"github.com/mtgnorton/helper/workerpool/statistics"
	"github.com/mtgnorton/helper/workerpool/worker"
	"time"
)

type Option func(o *Options)

// Options pool options
type Options struct {
	name              string
	logger            mlogger.Logger
	Context           context.Context
	CreateWorkerFunc  worker.CreateWorkerFunc // 创建 worker 的方法
	WorkerNumber      int                     // worker 数量
	MaxTaskNumber     int                     // 最大任务数,如果设置了最大任务数,则当任务数达到最大任务数时,pool 会退出
	displayProcessGap time.Duration           // 显示进度的间隔
	statistics        statistics.Statistics   // 统计信息
}

func WithName(name string) Option {
	return func(o *Options) {
		o.name = name
	}
}

func WithLogger(logger mlogger.Logger) Option {
	return func(o *Options) {
		o.logger = logger
	}
}

func WithContext(ctx context.Context) Option {
	return func(o *Options) {
		o.Context = ctx
	}
}

func WithCreateWorkerFunc(createWorkerFunc worker.CreateWorkerFunc) Option {
	return func(o *Options) {
		o.CreateWorkerFunc = createWorkerFunc
	}
}

func WithWorkerNumber(workerNumber int) Option {
	return func(o *Options) {
		o.WorkerNumber = workerNumber
	}
}

func WithMaxTaskNumber(maxTaskNumber int) Option {
	return func(o *Options) {
		o.MaxTaskNumber = maxTaskNumber
	}
}

func WithDisplayProcessGap(gap time.Duration) Option {
	return func(o *Options) {
		o.displayProcessGap = gap
	}
}

func WithStatistics(statistics statistics.Statistics) Option {
	return func(o *Options) {
		o.statistics = statistics
	}
}
