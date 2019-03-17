package Logger

import (
	"testing"
)

func TestLogger_Log(t *testing.T) {
	NewLogger().SetFileWriter("shit.log").SetLevel(INFO).Info("This is a test01", "这是中文日志01")
	NewLogger().SetFileWriter("fuck.log").SetLevel(WARN).Info("This is a test02", "这是中文日志02")
}
