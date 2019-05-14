package Websocket

// Message is an object for websocket message which is mapped to json type
type Message struct {
	// 发送者
	Sender string `json:"sender,omitempty"`
	// 接收者
	Receiver string `json:"receiver,omitempty"`
	// 内容
	Content string `json:"content,omitempty"`
}
