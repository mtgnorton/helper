package task

type Result interface {
	Task() Task              // 当前任务
	Result() interface{}     // 任务执行结果
	Err() error              // 任务执行错误
	AdditionalTasks() []Task // 附加任务
}

type DefaultResult struct {
	task            Task
	rs              interface{}
	err             error
	additionalTasks []Task
}

func NewDefaultResult(task Task, rs interface{}, additionalTasks []Task, err error) Result {
	return &DefaultResult{
		task:            task,
		rs:              rs,
		additionalTasks: additionalTasks,
		err:             err,
	}
}
func (r *DefaultResult) Task() Task {
	return r.task
}
func (r *DefaultResult) Result() interface{} {
	return r.rs
}
func (r *DefaultResult) Err() error {
	return r.err
}
func (r *DefaultResult) AdditionalTasks() []Task {
	return r.additionalTasks
}
