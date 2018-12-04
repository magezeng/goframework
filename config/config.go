package config

import (
	"encoding/json"
	Errors "go-framework/errors"
	DO "go-framework/model/do"
	Utils "go-framework/utils"
	"os"
)

var instance *DO.Configuration

// GetInstance 返回配置实例
func GetInstance() DO.Configuration {
	if instance == nil {
		configFileName := getConfigFileNameByMode()
		file := Utils.OpenLocalFile(configFileName)
		defer Utils.CloseFile(file)
		decoder := json.NewDecoder(file)
		// 全局变量的冲突问题
		instance = &DO.Configuration{}
		err := decoder.Decode(instance)
		if err != nil {
			Errors.Raise(err, Errors.IOErrCode)
		}
	}
	return *instance
}

// 从环境变量中决定应该使用哪一个配置文件
func getConfigFileNameByMode() string {
	switch mode := os.Getenv("mode"); mode {
	case "prod":
		return "config.prod.json"
	case "test":
		return "config.test.json"
	default:
		return "config.dev.json"
	}
}
