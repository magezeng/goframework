package Logger

import (
	"testing"
)

func TestLogger_Log(t *testing.T) {
	NewLogger().SetFileWriter("app.log").SetFileWriter("shit.log").SetLevel(INFO).Info("This is a test", "这是中文日志")
}
