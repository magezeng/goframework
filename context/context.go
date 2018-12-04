package context

import (
	DO "go-framework/model/do"

	"github.com/go-xorm/xorm"
)

// Context 上下文类型，这个上下文可以携带数据
type Context map[string]interface{}

var instance Context

func init() {
	context := Context{}
	director := director{Builder: defaultContextBuilder{context}}
	instance = director.defaultBuild()
}

// GetInstance 获得初始化后的上下文实例
func GetInstance() *Context {
	return &instance
}

// Config 获得初始化后的上下文中的config数据
func (c *Context) Config() (config DO.Configuration) {
	configInterface := instance["config"]
	if configInterface != nil {
		config = configInterface.(DO.Configuration)
	}
	return
}

// DB 获得初始化后的上下文中的db实例
func (c *Context) DB() (engine *xorm.Engine) {
	dbInterface := instance["db"]
	if dbInterface != nil {
		engine = dbInterface.(*xorm.Engine)
	}
	return
}
