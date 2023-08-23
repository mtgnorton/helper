package single_flight

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestFlightGroup_Do(t *testing.T) {

	writeCount := 0

	g := NewRWSingleFlight()

	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r, err := g.Do("test", func() (interface{}, error) {
				time.Sleep(time.Millisecond)
				writeCount++
				return writeCount, nil
			}, func(w interface{}) (interface{}, error) {
				return w, nil
			})
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(r)
			}
		}()
	}
	wg.Wait()

}
