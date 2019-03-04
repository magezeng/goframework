package errors

import (
	"bytes"
	"fmt"
	"reflect"
	"runtime"
)

// 最大显示的栈深度
const MaxStackDepth = 50

// Error里面包含了栈错误的相关信息及错误本身
type Error struct {
	// 错误主体
	Err error
	// 每一行的栈帧指针
	stack []uintptr
	// 栈帧
	frames []StackFrame
}

// New 创建一个错误
func New(e interface{}) *Error {
	if e == nil {
		return nil
	}

	var err error
	// 要创建的错误来源可能有多种类型
	switch e := e.(type) {
	case *Error:
		return e
	case error:
		err = e
	default:
		err = fmt.Errorf("%v", e)
	}

	stack := make([]uintptr, MaxStackDepth)
	length := runtime.Callers(2, stack[:])
	return &Error{
		Err:   err,
		stack: stack[:length],
	}
}

// 判断两个错误是不是一种
func Is(e error, original error) bool {
	if e == original {
		return true
	}
	if e, ok := e.(*Error); ok {
		return Is(e.Err, original)
	}
	if original, ok := original.(*Error); ok {
		return Is(e, original.Err)
	}
	return false
}

// Error 返回错误的字符串
func (err *Error) Error() string {
	return err.Err.Error()
}

// Stack 返回类似runtime/debug.printStack()的结果
func (err *Error) Stack() []byte {
	buf := bytes.Buffer{}
	for _, frame := range err.StackFrames() {
		buf.WriteString(frame.String())
	}
	return buf.Bytes()
}

// ErrorStack 返回打印message和栈
func (err *Error) ErrorStack() string {
	return err.TypeName() + " " + err.Error() + "\n" + string(err.Stack())
}

// StackFrames 返回错误栈的数组
func (err *Error) StackFrames() []StackFrame {
	if err.frames == nil {
		err.frames = make([]StackFrame, len(err.stack))
		for i, pc := range err.stack {
			err.frames[i] = NewStackFrame(pc)
		}
	}
	return err.frames
}

// 返回错误的类型，例如*errors.stringError.
func (err *Error) TypeName() string {
	return reflect.TypeOf(err.Err).String()
}
