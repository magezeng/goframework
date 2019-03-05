package Database

import (
	"os"
	"path/filepath"
	"testing"
	"tipu.com/go-framework/Config"
	"tipu.com/go-framework/Models"
)

const configContent = `
# config.yml
http:
db:
  - connectStr:  root:123456@/test?charset=utf8
    name: mysql
  - connectStr:  ./test.db3
    name: sqlite3
`

var database = GetInstance()

func TestGetConnectStr(t *testing.T) {
	path, err := genConfigFile("config.yml", configContent)
	if err != nil {
		t.Error()
		t.FailNow()
	}
	conf := Models.Config{}
	err = Config.GetInstance().Load(&conf, path)
	if err != nil {
		t.Error()
		t.Fail()
	}
	t.Log(conf)
}

func TestGetEngine(t *testing.T) {
	engine, err := database.GetEngine(SQLite)
	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}
	err = engine.Ping()
	if err != nil {
		t.Error(err.Error())
		t.Fail()
	} else {
		t.Log("sqlite数据库已经连接")
	}
}

func TestEngineSingleton(t *testing.T) {
	engine1, err := database.GetEngine(SQLite)
	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}
	engine2, err := database.GetEngine(SQLite)
	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}
	if engine1 == engine2 {
		t.Log("两个引擎是相同的")
	} else {
		t.Error("两个引擎不同")
		t.Fail()
	}
}

func genConfigFile(name string, config string) (string, error) {
	path := filepath.Join(".", name)
	io, err := os.Create(path)
	if err != nil {
		return "", err
	}
	if _, err := io.WriteString(config); err != nil {
		return "", err
	}
	return path, nil
}
