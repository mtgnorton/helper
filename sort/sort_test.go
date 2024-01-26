package sort

import (
	"testing"
)

type user struct {
	name string
	age  int
}

func TestQuickSort(t *testing.T) {
	list := []interface{}{2, 44, 4, 8, 33, 1, 22, -11, 6, 34, 55, 54, 9}
	QuickSort(list, func(a, b interface{}) int {
		aInt := a.(int)
		bInt := b.(int)
		if aInt < bInt {
			return -1
		} else if aInt == bInt {
			return 0
		}
		return 1
	})
	t.Log(list)

	list = []interface{}{
		user{
			name: "a",
			age:  16,
		},
		user{
			name: "b",
			age:  14,
		},
		user{
			name: "c",
			age:  22,
		},
	}
	QuickSort(list, func(a, b interface{}) int {
		aUser := a.(user)
		bUser := b.(user)
		if aUser.age < bUser.age {
			return -1
		} else if aUser.age == bUser.age {
			return 0
		}
		return 1
	})
	t.Log(list)

}
