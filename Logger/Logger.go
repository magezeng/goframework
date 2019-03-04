package Logger

import "time"

type Logger struct {
	// 日志写入器的列表
	writerList []LogWriterInterface
	// 级别
	level int
	// 时间格式
	timeFormatter string
}

const (
	DEBUG = iota
	INFO
	WARN
	ERROR
)

// RegisterWriter 注册一个日志写入器
func (l *Logger) AddWriter(writer LogWriterInterface) {
	l.writerList = append(l.writerList, writer)
}

func (l *Logger) SetLevel(level int) {
	l.level = level
}

func (l *Logger) SetTimeFormatter(timeFormatter string) {
	l.timeFormatter = timeFormatter
}

func (l *Logger) Debug(data interface{}) {
	if l.level <= DEBUG {
		l.log("[DEBUG] ", data)
	}
}

func (l *Logger) Info(data interface{}) {
	if l.level <= INFO {
		l.log("[INFO] ", data)
	}
}

func (l *Logger) Warn(data interface{}) {
	if l.level <= WARN {
		l.log("[WARN] ", data)
	}
}

func (l *Logger) Error(data interface{}) {
	if l.level <= ERROR {
		l.log("[ERROR] ", data)
	}
}

// Log 将一个data类型的数据写入日志
func (l *Logger) log(prefix string, data interface{}) {
	now := time.Now()
	timeStr := now.Format(l.timeFormatter)
	for _, writer := range l.writerList {
		// 将日志输出到每一个写入器中
		_ = writer.Write(timeStr+" "+prefix, data)
	}
}

// NewLogger 创建日志的实例
func NewLogger() *Logger {
	return &Logger{level: DEBUG, timeFormatter: "2006-01-01 15:04:05"}
}
