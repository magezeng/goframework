package Logger

import (
	"fmt"
	"os"
)

// 命令行写入器
type logConsoleWriter struct {
}

// 实现LogWriter的Write()方法
func (f *logConsoleWriter) Write(prefix string, data interface{}) error {
	// 将数据序列化为字符串
	str := fmt.Sprintf("%v\n", data)
	// 将数据以字节数组写入命令行中
	_, err := os.Stdout.Write([]byte(prefix + str))
	return err
}

// 创建命令行写入器实例
func NewLogConsoleWriter() *logConsoleWriter {
	return &logConsoleWriter{}
}
