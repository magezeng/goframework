package context

import (
	"os"
)

// 监工（一定要保证装配正确顺序）
type director struct {
	Builder iContextBuilder
}

func (dr *director) defaultBuild() (context Context) {
	env := os.Getenv("env")
	if env == "" {
		env = "test"
	}
	// 先设置环境变量
	dr.Builder.SetEnv(env)
	dr.Builder.SetConfig()
	dr.Builder.SetDB()
	dr.Builder.SetLogger()
	context = dr.Builder.Build()
	return
}
