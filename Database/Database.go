package Database

import (
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"sync"
	"tipu.com/go-framework/Config"
	"tipu.com/go-framework/Models"
)

var (
	instance *Database
	once     sync.Once
)

type Database struct {
	engineMap map[string]*gorm.DB
}

func GetInstance() *Database {
	once.Do(func() {
		instance = &Database{engineMap: make(map[string]*gorm.DB)}
	})
	return instance
}

func (d *Database) GetEngine(engineName string) (engine *gorm.DB, err error) {
	existedEngine, isExist := d.engineMap[engineName]
	if !isExist {
		engine, err = d.getEngine(engineName)
		if err != nil {
			return
		}
		d.engineMap[engineName] = engine
	} else {
		engine = existedEngine
	}
	return
}

func (d *Database) getConnectStr(engineName string) (connectStr string, err error) {
	conf := Models.Config{}
	err = Config.GetInstance().Load(&conf, "./config.yml")
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

func (d *Database) getEngine(engineName string) (engine *gorm.DB, err error) {
	connectStr, err := d.getConnectStr(engineName)
	if err != nil {
		return
	}
	engine, err = gorm.Open(engineName, connectStr)
	if err != nil {
		err = errors.New("数据库连接失败: " + err.Error())
		engine = nil
	}
	return
}
