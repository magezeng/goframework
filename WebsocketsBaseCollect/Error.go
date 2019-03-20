package WebsocketsBaseCollect

import (
	"fmt"
	"runtime"
	"strings"
)

const MaxStackDepth = 50

// stack 新的错误结构
// 适用于调用层次比较深的error来正确的展示调用过程
type stack struct {
	msg        string
	stackFrame []uintptr
}

// callers 获取调用栈
func callers() []uintptr {
	var pcs [MaxStackDepth]uintptr
	// 跳过前三层
	n := runtime.Callers(3, pcs[:])
	st := pcs[0:n]
	return st
}

// Error 实现Error方法
func (e *stack) Error() string {
	var b strings.Builder
	if e.msg != "" {
		b.WriteString(e.msg)
		b.WriteString("\n")
	}
	b.WriteString("调用过程:")
	for _, pc := range e.stackFrame {
		fn := runtime.FuncForPC(pc)
		b.WriteString("\n")
		f, n := fn.FileLine(pc)
		b.WriteString(fmt.Sprintf("%s:%d", f, n))
	}
	b.WriteString("\n")
	return b.String()
}

func NewTraceWithMsg(message string) error {
	return &stack{msg: message, stackFrame: callers()}
}

func NewTrace() error {
	return &stack{stackFrame: callers()}
}
