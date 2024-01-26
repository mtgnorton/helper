package search

import "testing"

func Test_binarySearch(t *testing.T) {
	type args struct {
		list  []interface{}
		check Check
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "test1",
			args: args{
				list: []interface{}{1, 3, 5, 6, 7},
				check: func(a interface{}) int {
					if 5 > a.(int) {
						return 1
					}
					if a.(int) == 5 {
						return 0
					}
					return -1
				},
			},
			want: 2,
		},
		{
			name: "test2",
			args: args{
				list: []interface{}{1, 3, 5, 6, 7},
				check: func(a interface{}) int {
					if 8 > a.(int) {
						return 1
					}
					if a.(int) == 8 {
						return 0
					}
					return -1
				},
			},
			want: -1,
		},
		{
			name: "test3",
			args: args{
				list: []interface{}{},
				check: func(a interface{}) int {
					if 8 > a.(int) {
						return 1
					}
					if a.(int) == 8 {
						return 0
					}
					return -1
				},
			},
			want: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := binarySearch(tt.args.list, tt.args.check); got != tt.want {
				t.Errorf("binarySearch() = %v, want %v", got, tt.want)
			}
		})
	}
}
