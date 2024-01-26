package task

// Task 任务接口
type Task interface {
	Priority() int     // 获取优先级，越大越优先
	Data() interface{} // 获取数据
}

type defaultTask struct {
	data     interface{}
	priority int
}

// NewDefaultTask 创建默认任务
func NewDefaultTask(data interface{}, priority ...int) Task {

	if len(priority) == 0 {
		return &defaultTask{
			priority: 0,
			data:     data,
		}
	}

	return &defaultTask{
		priority: priority[0],
		data:     data,
	}
}

// Priority 获取优先级
func (t *defaultTask) Priority() int {
	return t.priority
}

// Data 获取数据
func (t *defaultTask) Data() interface{} {
	return t.data
}
