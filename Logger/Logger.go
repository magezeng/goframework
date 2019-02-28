package Logger

import (
	"fmt"
	"log"
	"os"
	"sync"
)

// 日志级别
const (
	DEBUG   = "DEBUG"
	INFO    = "INFO"
	WARNING = "WARN"
	ERROR   = "ERROR"
)

var (
	once        sync.Once
	logInstance logT
	defaultDir  = "./app.log"
)

type logT struct {
	dir    string
	logger *log.Logger
}

func getInstance() *logT {
	once.Do(func() {
		file, err := os.OpenFile(defaultDir, os.O_APPEND|os.O_CREATE, 0755)
		if err != nil {
			log.Fatalln("创建日志文件失败")
		}
		logInstance = logT{logger: log.New(file, "", log.LstdFlags)}
	})
	return &logInstance
}

// Debug 输出Debug级别日志
func Debug(v ...interface{}) {
	setPrefix(DEBUG)
	getInstance().logger.Println(v)
}

// Info 输出Info级别日志
func Info(v ...interface{}) {
	setPrefix(INFO)
	getInstance().logger.Println(v)
}

// Warn 输出Warn级别日志
func Warn(v ...interface{}) {
	setPrefix(WARNING)
	getInstance().logger.Println(v)
}

// Error 输出Error级别日志
func Error(v ...interface{}) {
	setPrefix(ERROR)
	getInstance().logger.Println(v)
}

// 设置日志每一行前方的日志标签
func setPrefix(level string) {
	logPrefix := fmt.Sprintf("[%s]", level)
	getInstance().logger.SetPrefix(logPrefix)
}
