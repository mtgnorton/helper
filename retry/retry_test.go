package retry

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestRetry(t *testing.T) {
	ctx := context.Background()

	err := RetryWithTimeout(ctx, 1*time.Millisecond, func(ctx context.Context) error {

		//构建一个http请求,传入ctx
		req, _ := http.NewRequest("GET", "https://api.github.com/users/helei112g", nil)
		req = req.WithContext(ctx)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer func() {
			err = resp.Body.Close()
			return
		}()

		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(body))
		return nil
	})

	if err != nil && strings.Contains(err.Error(), "retry fail") {
		t.Log("success")
	} else {
		t.Fatal("fail")
	}
}
