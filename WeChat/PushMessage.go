package WeChat

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"strings"
)

func SendMessageToWeChat(message string) (err error) {
	data := `{"text": "` + message + `"}`
	resp, err := http.Post("http://47.75.65.211:9995/send",
		"application/json",
		strings.NewReader(data))
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	mapRes := make(map[string]interface{})
	err = json.Unmarshal(body, &mapRes)
	if err != nil {
		return
	}
	if mapRes["errcode"].(float64) != 0 {
		err = errors.New("发送失败: " + mapRes["errmsg"].(string))
	}
	return
}
