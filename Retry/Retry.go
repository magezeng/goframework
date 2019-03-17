package Retry

import (
	"fmt"
	"time"
)

// Retry重试，指定时间和次数，运行回调方法
func Retry(attempts int, sleep time.Duration, callback func(args ...interface{}) (result interface{}, err error), args ...interface{}) (result interface{}, err error) {
	for i := 0; ; i++ {
		result, err = callback(args...)
		if err == nil {
			return
		}
		if i >= (attempts - 1) {
			break
		}
		time.Sleep(sleep)
	}
	err = fmt.Errorf("经过了 %d 次重试, 仍然出错: %s\n", attempts, err.Error())
	return
}
