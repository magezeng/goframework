package utils

import (
	"bytes"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
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

// HTTPGet Get请求方式
func HTTPGet(url url.URL) (body []byte, err error) {
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response, err := httpClient.Do(req)
	defer func() {
		err = response.Body.Close()
		return
	}()
	if err != nil && response == nil {
		return
	}

	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	return
}

// HTTPPost Post请求方式，返回字节流
func HTTPPost(url url.URL, reqBody []byte) (body []byte, err error) {
	req, err := http.NewRequest("POST", url.String(), bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(req)
	defer func() {
		err = response.Body.Close()
		return
	}()
	if err != nil {
		return
	}

	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	return
}
