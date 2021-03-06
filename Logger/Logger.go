package Logger

import (
	"time"
	"fmt"
)

type Logger struct {
	// 日志写入器的列表
	writerList []LogWriterInterface
	// 级别
	level int
	// 时间格式
	timeFormatter string
	// 前缀
	prefix string
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
	l.addWriter(fw)
	return l
}

func (l *Logger) SetLevel(level int) *Logger {
	l.level = level
	return l
}

func (l *Logger) SetTimeFormatter(timeFormatter string) *Logger {
	l.timeFormatter = timeFormatter
	return l
}

func (l *Logger) SetPrefix(prefix string) *Logger{
	l.prefix = prefix
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
			res := fmt.Sprint(item)
			// res, _ := json.Marshal()
			logStr += string(res)
		}
		logStr += " "
	}
	for _, writer := range l.writerList {
		// 将日志输出到每一个写入器中
		_ = writer.Write(timeStr+" "+prefix+l.prefix, logStr)
	}
}

// RegisterWriter 注册一个日志写入器
func (l *Logger) addWriter(writer LogWriterInterface) {
	l.writerList = append(l.writerList, writer)
}

func (l *Logger) withConsoleWriter() *Logger {
	l.addWriter(NewLogConsoleWriter())
	return l
}

// NewLogger 在线程安全的前提下取得logger
func NewLogger() *Logger {
	instance := &Logger{level: DEBUG, timeFormatter: "2006-01-01 15:04:05"}
	return instance.withConsoleWriter()
}
