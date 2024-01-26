package timeout

import (
	"fmt"
	"testing"
	"time"
)

func Test_Do(t *testing.T) {
	r := Do(func() Result {
		time.Sleep(time.Second * 2)
		return Result{}
	}, time.Second)
	fmt.Println(r)
}
