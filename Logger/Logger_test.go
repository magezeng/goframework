package Logger

import (
	"testing"
)

func TestLogger_Log(t *testing.T) {
	l := NewLogger()
	cw := NewLogConsoleWriter()
	// 创建文件写入器
	fw := NewLogFileWriter()
	// 设置文件名
	if err := fw.SetFile("app.log"); err != nil {
		t.Error(err.Error())
		t.Fail()
	}
	// 设置了error，只有error级别的可以输出
	// 注册文件写入器到日志器中
	l.AddWriter(fw).AddWriter(cw).SetLevel(WARN).Warn("This is a test")
}
