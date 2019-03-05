package Config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
	"tipu.com/go-framework/ArrValueTranslate"
)

var (
	dir string
	cfg = GetInstance()
)

type DBConfig struct {
	Master  string        `yaml:"master"`
	Slave   string        `yaml:"slave"`
	Timeout time.Duration `yml:"timeout"`
}

type Conf struct {
	Domain  string        `yaml:"domain"`
	IsDev   bool          `yaml:"is_dev"`
	Timeout time.Duration `yml:"timeout"`
	DB      DBConfig      `yaml:"db"`
}

func TestYAMLLoader_Load(t *testing.T) {
	baseConfig := `
# config.yml
domain: example.com
db:
  master:  rw@/example
  slave:   ro@/example
  timeout: 0.5s
`
	localConfig := `
# config_local.yml
domain: dev.example.com
is_dev: true
`
	conf := &Conf{}
	baseConfigYaml, _ := genConfigFile("config.yml", baseConfig)         // /path/to/config.yml
	localConfigYaml, _ := genConfigFile("config_local.yml", localConfig) // /path/to/config_local.yml

	cfg.SetDecoder(YAML)
	err := cfg.Load(conf, baseConfigYaml, localConfigYaml)
	if err != nil {
		t.Error(err.Error())
		t.Fail()
	}
	str, _ := json.MarshalIndent(conf, "", "  ")
	t.Log(string(str))
}

func TestJSONLoader_Load(t *testing.T) {
	a, err := genConfigFile("a.json", `{
  "foo": "bar",
	"list": ["123", "2319", "23131"]
}`)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	b, err := genConfigFile("b.json", `{
  "bar": "baz"
}`)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	c := make(map[string]interface{})
	cfg.SetDecoder(JSON)
	err = cfg.Load(&c, a, b)
	if err != nil {
		t.Error(err)
	}
	if c["foo"] != "bar" {
		t.Error("foo 应该为 bar")
		t.Fail()
	}

	list := make([]string, 3)
	err = ArrValueTranslate.ArrValueTranslate(c["list"], &list)
	if err != nil{
		t.Error(err.Error())
	}

	if list[2] != "23131" {
		t.Error("json里面的list取得失败")
		t.Fail()
	}
	if c["bar"] != "baz" {
		t.Error("bar 应该为 baz！")
		t.Fail()
	}
}

func genConfigFile(name string, config string) (string, error) {
	path := filepath.Join(dir, name)
	io, err := os.Create(path)
	if err != nil {
		return "", err
	}
	if _, err := io.WriteString(config); err != nil {
		return "", err
	}
	return path, nil
}
