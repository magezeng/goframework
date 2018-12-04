package utils

import (
	"bytes"
	Errors "go-framework/errors"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

var httpClient *http.Client

func init() {
	httpClient = createHTTPClient()
}

// Todo:数据可以从配置文件里面读取
const (
	maxIdleConns        int           = 100
	maxIdleConnsPerHost int           = 100
	idleConnTimeout     time.Duration = 90
)

func createHTTPClient() *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:        maxIdleConns,
			MaxIdleConnsPerHost: maxIdleConnsPerHost,
			IdleConnTimeout:     idleConnTimeout * time.Second,
		},
	}
	return client
}

// HTTPGet Get请求方式，返回字节流，现阶段有错直接抛
func HTTPGet(url string) []byte {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		Errors.Raise(err, Errors.NetErrCode)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response, err := httpClient.Do(req)
	if err != nil && response == nil {
		Errors.Raise(err, Errors.NetErrCode)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		Errors.Raise(err, Errors.ParseErrCode)
	}
	return body
}

// HTTPPost Post请求方式，返回字节流，现阶段有错直接抛
func HTTPPost(url string, reqBody []byte) []byte {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		Errors.Raise(err, Errors.NetErrCode)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Errors.Raise(err, Errors.ParseErrCode)
	}
	return body
}
