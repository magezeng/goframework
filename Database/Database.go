package Database

// 这个模块应该是被动接受了config位置后再去读文件
// 自己不能默认config位置

import (
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"go-framework/Config"
	"go-framework/Models"
)

var (
	instance = make(map[string]*gorm.DB)
)

func GetEngine(engineName string, configPath ...string) (engine *gorm.DB, err error) {
	existedEngine, isExist := instance[engineName]
	if !isExist {
		engine, err = getEngine(engineName, configPath...)
		if err != nil {
			return
		}
		instance[engineName] = engine
	} else {
		engine = existedEngine
	}
	return
}

func ReleaseEngine(engineName string) (err error) {
	existedEngine, isExist := instance[engineName]
	if isExist {
		err = existedEngine.Close()
	} else {
		err = errors.New("请求关闭的数据库未创建！")
	}
	return
}

func getConnectStr(engineName string, configPath ...string) (connectStr string, err error) {
	conf := Models.Config{}
	err = Config.GetInstance().MustGetData(&conf, configPath...)
	if err != nil {
		return
	}
	for _, c := range conf.DB {
		if c.Name == engineName {
			connectStr = c.ConnectStr
		}
	}
	return
}

func getEngine(engineName string, configPath ...string) (engine *gorm.DB, err error) {
	connectStr, err := getConnectStr(engineName, configPath...)
	if err != nil {
		return
	}
	engine, err = gorm.Open(engineName, connectStr)
	if err != nil {
		err = errors.New("数据库连接失败: " + err.Error())
		engine = nil
	} else {
		engine.LogMode(false)
	}
	return
}
