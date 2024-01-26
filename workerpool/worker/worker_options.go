package worker

import (
	"github.com/mtgnorton/helper/mlogger"
	"time"
)

type Options struct {
	name                  string
	logger                mlogger.Logger
	process               Process
	prepareWorkFinishTask PrepareWorkFinishTask
	timeout               time.Duration // 超时时间,为 0 不限制
	retryNumber           int           // 重试次数,为 0 不重试

}

type Option func(o *Options)

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

func WithProcess(process Process) Option {
	return func(o *Options) {
		o.process = process
	}
}

func WithPrepareWorkFinishTask(prepareWorkFinishTask PrepareWorkFinishTask) Option {
	return func(o *Options) {
		o.prepareWorkFinishTask = prepareWorkFinishTask
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(o *Options) {
		o.timeout = timeout
	}
}

func WithRetry(retryNumber int) Option {
	return func(o *Options) {
		o.retryNumber = retryNumber
	}
}
