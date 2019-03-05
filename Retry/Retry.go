package Retry

import (
	"fmt"
	"time"
)

func Retry(attempts int, sleep time.Duration, callback func(...interface{}) error, args ...interface{}) (err error) {
	var internalErr error
	for i := 0; ; i++ {
		internalErr = callback(args...)
		if err == nil {
			return
		}
		if i >= (attempts - 1) {
			break
		}
		time.Sleep(sleep)
	}
	err = fmt.Errorf("经过了 %d 次重试, 仍然出错: %s\n", attempts, internalErr.Error())
	return
}
