package Config

import (
	"github.com/pkg/errors"
	"io/ioutil"
	"sync"
)
// 反序列化方法
type unmarshalFunc func([]byte, interface{}) error

const (
	JSON = iota
	YAML
	TOML
	INI
)

var (
	once     sync.Once
	instance Config
)

type Config struct {
	decoder int
	loader  LoaderInterface
}

func GetInstance() Config {
	once.Do(func() {
		instance = Config{loader: new(YAMLLoader)}
	})
	return instance
}

func (c Config) SetDecoder(decoder int) {
	switch decoder {
	case JSON:
		c.loader = new(JSONLoader)
	case YAML:
		c.loader = new(YAMLLoader)
	default:
		panic(errors.New("未实现解析器!"))
	}
}

func (c Config) Load(configCls interface{}, configPaths ...string) error {
	return c.loader.Load(configCls, configPaths...)
}

func loadWithFunc(conf interface{}, configPaths []string, unmarshal unmarshalFunc) error {
	for _, configPath := range configPaths {
		err := loadConfig(configPath, conf, unmarshal)
		if err != nil {
			return err
		}
	}
	return nil
}

func loadConfig(configPath string, conf interface{}, unmarshal unmarshalFunc) error {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return errors.Wrapf(err, "%s 文件读取失败", configPath)
	}
	if err := unmarshal(data, conf); err != nil {
		return errors.Wrapf(err, "%s 配置转换失败", configPath)
	}
	return nil
}
