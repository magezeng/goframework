package Logger

import (
	"testing"
)

func TestLogger_Log(t *testing.T) {
	l := NewLogger()
	cw := NewLogConsoleWriter()
	// 注册命令行写入器到日志器中
	l.AddWriter(cw)
	// 创建文件写入器
	fw := NewLogFileWriter()
	// 设置文件名
	if err := fw.SetFile("app.log"); err != nil {
		t.Error(err.Error())
		t.Fail()
	}
	// 注册文件写入器到日志器中
	l.AddWriter(fw)
	// 设置了error，只有error级别的可以输出
	l.SetLevel(WARN)
	l.Warn("This is a test")
}
