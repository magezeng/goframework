package context

// 上下文建造器的策略接口
type iContextBuilder interface {
	SetEnv(env string)
	SetLogger()
	SetConfig()
	SetDB()
	Build() Context
}
