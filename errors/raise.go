package tipuError

import (
	"errors"
	"strconv"
)

// Raise 按照错误类别抛出错误
func Raise(errParam error, code int) {
	var err IErr
	switch code {
	case SQLErrCode:
		err = &SQLError{InternalErr: errors.New(errParam.Error())}
	case IOErrCode:
		err = &IOError{InternalErr: errors.New(errParam.Error())}
	case NetErrCode:
		err = &NetError{InternalErr: errors.New(errParam.Error())}
	case RuntimeErrCode:
		err = &RuntimeError{InternalErr: errors.New(errParam.Error())}
	default:
		err = &UnknownError{InternalErr: errors.New(errParam.Error())}
	}
	errInfo := err.Message() + "\n错误代码:" + strconv.Itoa(err.Code())
	panic(errInfo)
}
