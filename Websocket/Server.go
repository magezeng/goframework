package Websocket

import (
	"github.com/gorilla/websocket"
	"time"
)

const (
	pongWait       = 60 * time.Second
	pingPeriod     = 15 * time.Second
	maxMessageSize = 512
)

var (
	instance = GetInstance()
)

type ServerManager struct {
	// 现在保持的连接,key为token
	Servers map[string]*Server
	// 增加客户端的channel
	RegisterCh chan *Server
	// 断开客户端的channel
	UnregisterCh chan *Server
}

type Server struct {
	// 用户的token，一个token只对应一个连接
	Token string
	// 发送数据的连接
	Conn *websocket.Conn
	// 发送数据的通道
	SendCh chan []byte
}

// start 启动ws服务端
func (c *Server) start() {
	c.read()
	c.write()
}

// read 从ws连接中读取数据
func (c *Server) read() {
	go func() {
		// 设置读取的deadline
		c.Conn.SetReadDeadline(time.Now().Add(pongWait));
		c.Conn.SetReadLimit(maxMessageSize)
		// 此API目前浏览器不支持，需要采用判断消息字符串的方式进行
		//// 处理ping的闭包函数，设置deadline并且写入Pong
		//handlePong := func(string) error {
		//	c.Conn.SetReadDeadline(time.Now().Add(pongWait));
		//	return nil
		//}
		// 最小间隔15s，至少会有一个ping发送过来
		// c.Conn.SetPongHandler(handlePong)
		defer func() {
			// 不进行读取后，进行释放
			instance.UnregisterCh <- c
			c.Conn.Close()
		}()

		for {
			// 目前是出错了退出
			msgType, msg, err := c.Conn.ReadMessage()
			// TODO: 出错的展示
			if err != nil {
				logger.Warn("读取消息出错，退出!", err)
				break
			}

			switch msgType {
			case websocket.TextMessage:
				// TODO: 对收到消息的处理
				logger.Info("收到了浏览器发来的消息: ", string(msg))
				if string(msg) == "pong" {
					c.Conn.SetReadDeadline(time.Now().Add(pongWait));
				}
			default:
				logger.Warn("不支持的消息类型，忽略!", msgType)
			}
		}
	}()
}

// 向连接中写入数据
func (c *Server) write() {
	go func() {
		ticker := time.NewTicker(pingPeriod)
		defer func() {
			ticker.Stop()
			c.Conn.Close()
		}()

		for {
			select {
			case message, ok := <-c.SendCh:
				// 写入出错时直接退出，说明SendCh已经被关闭了
				if !ok {
					c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
					logger.Info("发送通道已经关闭: ", c.Token)
					return
				}

				w, err := c.Conn.NextWriter(websocket.TextMessage)
				if err != nil {
					return
				}
				_, err = w.Write(message)
				if err != nil {
					logger.Info("写入数据失败: ", c.Token)
					return
				}
				if err := w.Close(); err != nil {
					logger.Info("关闭writer失败: ", c.Token)
					return
				}
			case <-ticker.C:
				// 按周期发送ping到浏览器
				err := c.Conn.WriteMessage(websocket.TextMessage, []byte("ping"))
				if err != nil {
					logger.Error("发送ping失败，可能是浏览器端已经离线: ", c.Token)
					return
				}
				logger.Info("发送ping到浏览器", c.Token)
			}
		}
	}()
}
