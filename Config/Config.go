package Config

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"reflect"
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
	data    map[string]interface{}
}

func GetInstance() Config {
	once.Do(func() {
		instance = Config{loader: new(YAMLLoader), data: make(map[string]interface{})}
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

func (c Config) load(configCls interface{}, configPaths ...string) (err error) {
	err = c.loader.Load(configCls, configPaths...)
	if err != nil {
		return
	}
	c.data[reflect.TypeOf(configCls).String()] = configCls
	return
}

func (c Config) MustGetData(configCls interface{}, configPaths ...string) (err error) {
	if target, ok := c.data[reflect.TypeOf(configCls).String()]; ok {
		// fmt.Println("[配置] 缓存已经被命中!")
		targetByte, _ := json.Marshal(target)
		err = json.Unmarshal(targetByte, configCls)
	} else {
		// fmt.Println("[配置] 正在准备读文件...")
		err = c.load(configCls, configPaths...)
	}
	return err
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
