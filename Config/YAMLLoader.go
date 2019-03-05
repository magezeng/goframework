package Config

import (
	"gopkg.in/yaml.v2"
)

type YAMLLoader struct {
}

// Load 把多个yaml配置文件load到conf的class中
func (l *YAMLLoader) Load(conf interface{}, configPaths ...string) error {
	return loadWithFunc(conf, configPaths, yaml.Unmarshal)
}
