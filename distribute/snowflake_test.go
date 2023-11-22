package distribute

import (
	"sync"
	"testing"
	"time"
)

func TestNode_Generate(t *testing.T) {
	node, err := NewNode(1)
	if err != nil {
		t.Fatal(err)
	}
	// 多个协程
	wg := sync.WaitGroup{}
	wg.Add(100)
	ch := make(chan int64, 100)
	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				ch <- node.Generate()
			}

		}()
	}
	go func() {
		m := make(map[int64]bool)
		for v := range ch {
			t.Log(v)
			if _, ok := m[v]; ok {
				t.Fatal("重复")
			}
		}
	}()
	wg.Wait()

	time.Sleep(time.Second * 2)

}
