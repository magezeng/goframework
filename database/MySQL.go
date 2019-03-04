package common

import (
	"fmt"
	Config "go-framework/config"
	Errors "go-framework/errors"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

var instance *xorm.Engine

// 获得对mysql的连接字符串，此处采用简单默认值
func getConnectStr() string {
	config := Config.GetInstance()
	host := config.Db.Host
	port := config.Db.Port
	user := config.Db.User
	password := config.Db.Password
	database := config.Db.Database
	connectStr := fmt.Sprintf(`%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local`, user, password, host, port, database)
	return connectStr
}

// GetInstance 获得对数据库操作的实例
func GetInstance() *xorm.Engine {
	if instance == nil {
		connectStr := getConnectStr()
		db, err := xorm.NewEngine("mysql", connectStr)
		if err != nil {
			Errors.Raise(err, Errors.SQLErrCode)
		}
		db.ShowSQL(true)
		db.ShowExecTime(true)
		instance = db
	}
	return instance
}
