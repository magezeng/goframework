package WeChat

import "testing"

func TestSendMessageToWeChat(t *testing.T) {
	message := "go-framework中测试用的字符串"
	err := SendMessageToWeChat(message)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("发送消息成功!")
}
