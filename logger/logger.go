package logger

import (
	"fmt"
	Config "go-framework/config"
	Utils "go-framework/utils"
	"log"
	"os"
	"time"
)

var (
	logger     *log.Logger
	levelFlags = []string{"DEBUG", "INFO", "WARN", "ERROR"}
)

// 日志级别
const (
	DEBUG = iota
	INFO
	WARNING
	ERROR
)

// Init 日志调用一次即可，多次import只初始化第一次
func Init() {
	// 此处返回的为相对路径
	filePath := getLogFileFullPath()
	var file *os.File
	if ok, _ := Utils.IsFileExsist(filePath); ok {
		file = Utils.OpenLocalFile(filePath)
	} else {
		file = Utils.CreateLocalFile(filePath)
	}
	logger = log.New(file, "", log.LstdFlags)
}

// getLogFilePath 获得日志文件位置
func getLogFilePath() string {
	logSavePath := Config.GetInstance().Logger.SavePath
	return fmt.Sprintf("%s", logSavePath)
}

// getLogFileFullPath 返回日志文件完整路径
func getLogFileFullPath() string {
	logSaveName := Config.GetInstance().Logger.FileName
	logFileExt := Config.GetInstance().Logger.FileExt
	timeFormat := Config.GetInstance().Logger.TimeFormat
	beforePath := getLogFilePath()
	afterPath := fmt.Sprintf("%s%s.%s", logSaveName, time.Now().Format(timeFormat), logFileExt)
	return fmt.Sprintf("%s%s", beforePath, afterPath)
}

// Debug 输出Debug级别日志
func Debug(v ...interface{}) {
	setPrefix(DEBUG)
	logger.Println(v)
}

// Info 输出Info级别日志
func Info(v ...interface{}) {
	setPrefix(INFO)
	logger.Println(v)
}

// Warn 输出Warn级别日志
func Warn(v ...interface{}) {
	setPrefix(WARNING)
	logger.Println(v)
}

// Error 输出Error级别日志
func Error(v ...interface{}) {
	setPrefix(ERROR)
	logger.Println(v)
}

// 设置日志每一行前方的日志标签
func setPrefix(level int) {
	logPrefix := fmt.Sprintf("[%s]", levelFlags[level])
	logger.SetPrefix(logPrefix)
}
