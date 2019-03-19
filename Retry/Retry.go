package Retry

import (
	"fmt"
	"runtime"
	"time"
)

// Retry重试，指定时间和次数，运行回调方法
func Retry(attempts int, sleep time.Duration, callback func(args ...interface{}) (result interface{}, err error), args ...interface{}) (result interface{}, err error) {
	pc, _, _, _ := runtime.Caller(1)
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
	funcName := runtime.FuncForPC(pc).Name()
	err = fmt.Errorf("方法: %s , 经过了 %d 次重试, 仍然出错: %s\n", funcName, attempts, err.Error())
	return
}
