package Utils

import (
	"net/url"
	"testing"
)

func TestHTTPGet(t *testing.T) {
	res, err := HTTPGet(url.URL{Scheme: "https", Host: "www.baidu.com"})
	if err != nil {
		t.Error(err.Error())
		t.Fail()
	}
	t.Log(string(res))
}
