package Database

import (
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
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
	engineMap map[string]*xorm.Engine
}

func GetInstance() *Database {
	once.Do(func() {
		instance = &Database{engineMap: make(map[string]*xorm.Engine)}
	})
	return instance
}

func (d *Database) GetEngine(engineName string) (engine *xorm.Engine, err error) {
	existedEngine, isExist := d.engineMap[engineName]
	if !isExist {
		engine, err = d.getEngine(engineName)
		if err != nil {
			return
		}
		d.engineMap[engineName] = engine
	}else{
		engine = existedEngine
	}
	return
}

// isDBSurvival 数据库是否正常在线
func (d *Database) checkDBSurvival(engine *xorm.Engine) error {
	return engine.Ping()
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

func (d *Database) getEngine(engineName string) (engine *xorm.Engine, err error) {
	connectStr, err := d.getConnectStr(engineName)
	if err != nil {
		return
	}
	engine, err = xorm.NewEngine(engineName, connectStr)
	if err != nil {
		return
	}
	err = d.checkDBSurvival(engine)
	if err != nil {
		err = errors.New("数据库Ping失败: " + err.Error())
		engine = nil
	} else {
		engine.ShowSQL(true)
		engine.ShowExecTime(true)
	}
	return
}
