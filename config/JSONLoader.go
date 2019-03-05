package Config

import (
	"encoding/json"
)

type JSONLoader struct {
}

// Load 把多个yaml配置文件load到conf的class中
func (l *JSONLoader) Load(conf interface{}, configPaths ...string) error {
	return loadWithFunc(conf, configPaths, json.Unmarshal)
}