package search

// Check
// return 1  if  target > v
// return 0  if  target == v
// return -1 if target < v
type Check func(v interface{}) int

// 1,3,5,6,7
func binarySearch(list []interface{}, check Check) int {
	low, high := 0, len(list)-1
	for low <= high {
		mid := low + (high-low)/2 // low
		if check(list[mid]) == 0 {
			return mid
		} else if check(list[mid]) > 0 {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	return -1
}
