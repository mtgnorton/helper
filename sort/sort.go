package sort

// Compare the size of two values
// return 1 if a > b
// return 0 if a == b
// return -1 if a < b
type Compare func(a, b interface{}) int

func partition(list []interface{}, low, high int, compare Compare) int {
	pivot := list[low]
	for low < high {
		for low < high && compare(pivot, list[high]) <= 0 {
			high--
		}
		list[low] = list[high]
		for low < high && compare(pivot, list[low]) >= 0 {
			low++
		}
		list[high] = list[low]
	}
	list[low] = pivot
	return low
}
func quickSort(list []interface{}, low, high int, compare Compare) {
	if high > low {
		pivot := partition(list, low, high, compare)
		quickSort(list, low, pivot-1, compare)
		quickSort(list, pivot+1, high, compare)
	}
}

func QuickSort(list []interface{}, compare Compare) {
	quickSort(list, 0, len(list)-1, compare)
}
