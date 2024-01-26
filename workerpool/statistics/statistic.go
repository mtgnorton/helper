package statistics

import "github.com/mtgnorton/helper/workerpool/task"

// Statistics 统计任务的完成信息
type Statistics interface {
	AddReceiveNumber(int)
	AddTaskResult(result task.Result)
	AllResult() []task.Result
	ReceiveTaskNumber() int
	FailTaskNumber() int
	SuccessTaskNumber() int
	FinishTaskNumber() int
}

type DefaultStatistics struct {
	opts              Options
	receiveTaskNumber int // 接受到的任务数
	failTaskNumber    int // 错误的任务数
	successTaskNumber int // 成功的任务数
	result            []task.Result
}

// NewStatistics 创建一个统计信息器
func NewStatistics(options ...Option) Statistics {
	opts := Options{}
	for _, option := range options {
		option(&opts)
	}
	return &DefaultStatistics{
		opts:   opts,
		result: make([]task.Result, 0),
	}
}

// AddReceiveNumber 添加接受到的任务数
func (s *DefaultStatistics) AddReceiveNumber(i int) {
	s.receiveTaskNumber++
}

// AddTaskResult 添加任务结果,任务结果默认不保存，避免内存占用过大,如果需要保存结果,请设置 WithSaveResult(true)
func (s *DefaultStatistics) AddTaskResult(result task.Result) {
	if result.Err() == nil {
		s.successTaskNumber++
	} else {
		s.failTaskNumber++
	}
	if s.opts.saveResult {
		s.result = append(s.result, result)
	}
}

// AllResult 返回所有的结果
func (s *DefaultStatistics) AllResult() []task.Result {
	return s.result
}

// ReceiveTaskNumber 返回接受到的任务数
func (s *DefaultStatistics) ReceiveTaskNumber() int {

	return s.receiveTaskNumber
}

// FailTaskNumber 返回失败的任务数
func (s *DefaultStatistics) FailTaskNumber() int {
	return s.failTaskNumber
}

// SuccessTaskNumber 返回成功的任务数
func (s *DefaultStatistics) SuccessTaskNumber() int {

	return s.successTaskNumber
}

// FinishTaskNumber 返回完成的任务数
func (s *DefaultStatistics) FinishTaskNumber() int {
	return s.successTaskNumber + s.failTaskNumber
}
