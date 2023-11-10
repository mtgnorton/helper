package request

import (
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
)

// Get 发起GET请求
// url:请求地址
// maxScanTokenSizeMul:  bufio.MaxScanTokenSize 为 64kb,如果返回的数据超过64kb,需要调整此参数,比如设置为2,则最大支持128kb,以此类推
// params:请求参数,可选
func Get(url string, maxScanTokenSizeMul int, params ...map[string]string) (statusCode int, content []byte, err error) {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, nil, errors.Wrapf(err, "http.NewRequest err,url:%v", url)
	}
	q := req.URL.Query()

	if len(params) == 1 {
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
	//body, err := io.ReadAll(resp.Body)
	respContent := make([]byte, 0)
	scanner := bufio.NewScanner(resp.Body)
	buf := make([]byte, 0, bufio.MaxScanTokenSize*maxScanTokenSizeMul)
	scanner.Buffer(buf, cap(buf))
	for scanner.Scan() {
		line := scanner.Bytes()
		respContent = append(respContent, line...)
	}
	fmt.Println(len(respContent))
	if scanner.Err() != nil {
		return 0, nil, errors.Wrapf(err, "scanner.Err err,url:%v", url)
	}
	err = resp.Body.Close()

	if err != nil {
		return 0, nil, errors.Wrapf(err, "body close err")
	}
	return resp.StatusCode, respContent, nil
}
