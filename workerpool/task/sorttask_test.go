package task

import (
	"testing"
)

func Test_sortTask(t *testing.T) {

	testData := []struct {
		name  string
		tasks []Task
	}{
		{
			name: "test1",
			tasks: []Task{
				NewDefaultTask(3, 3),
				NewDefaultTask(1, 1),
				NewDefaultTask(7, 7),
				NewDefaultTask(2, 2),
				NewDefaultTask(4, 4),
			},
		},
		{
			name: "test2",
			tasks: []Task{
				NewDefaultTask(3, 3),
			},
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			SortTask(tt.tasks, 0, len(tt.tasks)-1)
			if len(tt.tasks) <= 1 {
				return
			}
			for i := 0; i < len(tt.tasks)-1; i++ {
				if tt.tasks[i].Priority() > tt.tasks[i+1].Priority() {
					t.Errorf("SortTask() error = %v, want %v", tt.tasks[i], tt.tasks[i+1])
				}
			}
		})
	}
}
