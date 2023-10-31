package request

import (
	"github.com/pkg/errors"
	"io"
	"net/http"
)

// Get 发起GET请求
func Get(url string, params ...map[string]string) (statusCode int, content []byte, err error) {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, nil, errors.Wrapf(err, "http.NewRequest err,url:%v", url)
	}
	q := req.URL.Query()

	if len(params) > 0 {
		for k, v := range params[0] {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}
	var resp *http.Response
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return 0, nil, errors.Wrapf(err, "http.DefaultClient.Do err,url:%v", url)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return 0, nil, errors.Wrapf(err, "io.ReadAll err,url:%v", url)
	}

	err = resp.Body.Close()
	if err != nil {
		return 0, nil, errors.Wrapf(err, "body close err")
	}
	return resp.StatusCode, body, nil
}
