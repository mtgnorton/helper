package request

import "testing"

func Test_Get(t *testing.T) {
	//code, r, err := Get("https://www.baidu.com")
	//if err != nil {
	//	t.Fatal(err)
	//}
	//t.Log(code, string(r))
	code, r1, err := Get("https://search.mtapi.io/Search", 20)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(code, string(r1))
}
