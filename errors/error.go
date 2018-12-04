package tipuError

type IErr interface {
	Code() int
	Message() string
}

// ParseErrCode 类型转换错误结构
type ParseError struct {
	InternalErr error
}

func (e *ParseError) Code() int {
	return ParseErrCode
}

func (e *ParseError) Message() string {
	return "类型转换发生了错误：" + e.InternalErr.Error()
}

// SQLError 数据库错误结构
type SQLError struct {
	InternalErr error
}

func (e *SQLError) Code() int {
	return SQLErrCode
}

func (e *SQLError) Message() string {
	return "数据库操作发生了错误：" + e.InternalErr.Error()
}

// NetError http错误结构
type NetError struct {
	InternalErr error
}

func (e *NetError) Code() int {
	return NetErrCode
}

func (e *NetError) Message() string {
	return "网络发生了错误：" + e.InternalErr.Error()
}

// IOError 输入输出错误结构
type IOError struct {
	InternalErr error
}

func (e *IOError) Code() int {
	return IOErrCode
}

func (e *IOError) Message() string {
	return "文件操作发生了错误：" + e.InternalErr.Error()
}

// RuntimeError 运行时错误结构
type RuntimeError struct {
	InternalErr error
}

func (e *RuntimeError) Code() int {
	return RuntimeErrCode
}

func (e *RuntimeError) Message() string {
	return "运行中发生了错误：" + e.InternalErr.Error()
}

// UnknownError 未知错误结构
type UnknownError struct {
	InternalErr error
}

func (e *UnknownError) Code() int {
	return UnknownErrCode
}

func (e *UnknownError) Message() string {
	return "发生了未知错误：" + e.InternalErr.Error()
}
