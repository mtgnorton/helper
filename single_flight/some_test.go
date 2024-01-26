package single_flight

import (
	"context"
	"testing"
	"time"
)

func TestDefaultMultiSome_Exec(t *testing.T) {

	type want struct {
		r   interface{}
		err error
	}

	testData := []struct {
		name    string
		timeout time.Duration
		tasks   []TaskFun
		wants   []want
	}{
		{
			name:    "MultiSome-time",
			timeout: 1 * time.Second,
			tasks: []TaskFun{
				func(ctx context.Context) (interface{}, error) {
					time.Sleep(time.Millisecond)
					return 1, nil
				},
				func(ctx context.Context) (interface{}, error) {
					time.Sleep(time.Millisecond * 2)
					return 2, nil
				},
				func(ctx context.Context) (interface{}, error) {
					time.Sleep(time.Millisecond * 3)
					return 3, nil
				},
			},
			wants: []want{
				{
					r:   1,
					err: nil,
				},
			},
		},
		{
			name:    "MultiSome-timeout",
			timeout: 1 * time.Second,
			tasks: []TaskFun{
				func(ctx context.Context) (interface{}, error) {
					time.Sleep(time.Second * 2)
					return 1, nil
				},
				func(ctx context.Context) (interface{}, error) {
					time.Sleep(time.Second * 2)
					return 2, nil
				},
				func(ctx context.Context) (interface{}, error) {
					time.Sleep(time.Second * 3)
					return 3, nil
				},
			},
			wants: []want{
				{
					r:   nil,
					err: nil,
				},
			},
		},
		{
			name:    "MultiSome-concurrency",
			timeout: 1 * time.Second,
			tasks:   getTasks(),
			wants: []want{
				{
					r:   1,
					err: nil,
				},
			},
		},
	}
	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMultiSome()
			r, _ := m.Exec(tt.timeout, tt.tasks...)
			if r != tt.wants[0].r {
				t.Errorf("Exec got = %v, want %v", r, tt.wants[0].r)
			}
		})
	}

}

func getTasks() (tasks []TaskFun) {
	tasks = make([]TaskFun, 1000)
	for i := 0; i < 1000; i++ {
		tasks[i] = func(ctx context.Context) (interface{}, error) {
			time.Sleep(time.Millisecond * 10)
			return 1, nil
		}
	}
	return
}
