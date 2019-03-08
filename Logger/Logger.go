package Logger

import (
	"encoding/json"
	"sync"
	"time"
)

var (
	instance *Logger
	once     sync.Once
)

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

func (l *Logger) SetFileWriter(path string) *Logger {
	fw := NewLogFileWriter()
	err := fw.setFile(path)
	if err != nil {
		panic(err)
	}
	return l.addWriter(fw)
}

func (l *Logger) SetLevel(level int) *Logger {
	l.level = level
	return l
}

func (l *Logger) SetTimeFormatter(timeFormatter string) *Logger {
	l.timeFormatter = timeFormatter
	return l
}

func (l *Logger) Debug(data ...interface{}) {
	if l.level <= DEBUG {
		l.log("[DEBUG] ", data...)
	}
}

func (l *Logger) Info(data ...interface{}) {
	if l.level <= INFO {
		l.log("[INFO] ", data...)
	}
}

func (l *Logger) Warn(data ...interface{}) {
	if l.level <= WARN {
		l.log("[WARN] ", data...)
	}
}

func (l *Logger) Error(data ...interface{}) {
	if l.level <= ERROR {
		l.log("[ERROR] ", data...)
	}
}

// Log 将一个data类型的数据写入日志
func (l *Logger) log(prefix string, data ...interface{}) {
	now := time.Now()
	timeStr := now.Format(l.timeFormatter)
	var logStr string
	for _, item := range data {
		if str, ok := item.(string); ok {
			logStr += str
		} else {
			res, _ := json.Marshal(item)
			logStr += string(res)
		}
		logStr += " "
	}
	for _, writer := range l.writerList {
		// 将日志输出到每一个写入器中
		_ = writer.Write(timeStr+" "+prefix, logStr)
	}
}

// RegisterWriter 注册一个日志写入器
func (l *Logger) addWriter(writer LogWriterInterface) *Logger {
	l.writerList = append(l.writerList, writer)
	return l
}

func (l *Logger) withConsoleWriter() *Logger {
	l.addWriter(NewLogConsoleWriter())
	return l
}

func NewLogger() *Logger {
	instance = &Logger{level: DEBUG, timeFormatter: "2006-01-01 15:04:05"}
	return instance.withConsoleWriter()
}
