package errors

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"runtime"
	"strings"
)

// 一个栈帧包含了这一行调用过程的重要信息
type StackFrame struct {
	// 文件名
	File string
	// 行号
	LineNumber int
	// 方法名
	Name string
	// 包名
	Package string
	// 当前调用的指针
	ProgramCounter uintptr
}

// NewStackFrame 根据当前调用的指针创建一个帧
func NewStackFrame(pc uintptr) (frame StackFrame) {
	frame = StackFrame{ProgramCounter: pc}
	if frame.Func() == nil {
		return
	}
	frame.Package, frame.Name = packageAndName(frame.Func())
	frame.File, frame.LineNumber = frame.Func().FileLine(pc - 1)
	return
}

// Func 返回哪个方法含有当前这个栈帧
func (frame *StackFrame) Func() *runtime.Func {
	if frame.ProgramCounter == 0 {
		return nil
	}
	return runtime.FuncForPC(frame.ProgramCounter)
}

// 返回该帧的字符串
func (frame *StackFrame) String() string {
	str := fmt.Sprintf("%s:%d (0x%x)\n", frame.File, frame.LineNumber, frame.ProgramCounter)
	source, err := frame.SourceLine()
	if err != nil {
		return str
	}
	return str + fmt.Sprintf("\t%s: %s\n", frame.Name, source)
}

// SourceLine 获取当前源代码的行位置
func (frame *StackFrame) SourceLine() (string, error) {
	data, err := ioutil.ReadFile(frame.File)
	if err != nil {
		return "", New(err)
	}
	lines := bytes.Split(data, []byte{'\n'})
	if frame.LineNumber <= 0 || frame.LineNumber >= len(lines) {
		return "???", nil
	}
	return string(bytes.Trim(lines[frame.LineNumber-1], " \t")), nil
}

// packageAndName 获取当前的包名和方法名
func packageAndName(fn *runtime.Func) (string, string) {
	name := fn.Name()
	pkg := ""
	if lastSlash := strings.LastIndex(name, "/"); lastSlash >= 0 {
		pkg += name[:lastSlash] + "/"
		name = name[lastSlash+1:]
	}
	if period := strings.Index(name, "."); period >= 0 {
		pkg += name[:period]
		name = name[period+1:]
	}
	name = strings.Replace(name, "·", ".", -1)
	return pkg, name
}
