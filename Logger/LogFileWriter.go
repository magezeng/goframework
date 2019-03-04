package Logger

import (
	"errors"
	"fmt"
	"os"
)

// 文件写入器
type logFileWriter struct {
	file *os.File
}

// SetFile 设置文件写入器写入的文件名
func (f *logFileWriter) SetFile(filename string) (err error) {
	// 如果文件已经打开, 关闭前一个文件
	if f.file != nil {
		_ = f.file.Close()
	}
	f.file, err = os.Create(filename)
	return err
}

// Write 实现接口的Write方法
func (f *logFileWriter) Write(prefix string, data interface{}) error {
	// 日志文件可能没有创建成功
	if f.file == nil {
		return errors.New("日志文件未创建成功")
	}
	str := fmt.Sprintf("%v\n", data)
	_, err := f.file.Write([]byte(prefix + str))
	return err
}

// 创建文件写入器实例
func NewLogFileWriter() *logFileWriter {
	return &logFileWriter{}
}
