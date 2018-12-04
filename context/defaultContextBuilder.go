package context

import (
	Config "go-framework/config"
	Database "go-framework/database"
	Logger "go-framework/logger"
)

type defaultContextBuilder struct {
	InternalContext Context
}

func (de defaultContextBuilder) SetDB() {
	de.InternalContext["db"] = Database.GetInstance()
}

func (de defaultContextBuilder) SetConfig() {
	de.InternalContext["config"] = Config.GetInstance()
}

func (de defaultContextBuilder) SetLogger() {
	Logger.Init()
}

func (de defaultContextBuilder) SetEnv(env string) {
	de.InternalContext["env"] = env
}

func (de defaultContextBuilder) Build() Context {
	return de.InternalContext
}
