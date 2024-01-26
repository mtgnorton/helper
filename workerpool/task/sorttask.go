package task

// SortTask use QuickSort  sort task by priority
func SortTask(tasks []Task, low, high int) {
	if high <= low {
		return
	}
	pivot := partition(tasks, low, high)
	SortTask(tasks, low, pivot-1)
	SortTask(tasks, pivot+1, high)
}

func partition(tasks []Task, low, high int) int {
	pivot := tasks[low]
	for low < high {
		for low < high && tasks[high].Priority() >= pivot.Priority() {
			high--
		}
		tasks[low] = tasks[high]
		for low < high && tasks[low].Priority() <= pivot.Priority() {
			low++
		}
		tasks[high] = tasks[low]
	}
	tasks[low] = pivot
	return low
}
