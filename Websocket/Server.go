package Websocket

import (
	"github.com/gorilla/websocket"
	"time"
)

const (
	// 1分钟没有收到客户端的返回消息，需要释放连接
	pongWait = 60 * time.Second
	// 服务端发起ping的时间间隔
	pingPeriod = 15 * time.Second
	// 消息最大长度，字节
	maxMessageSize = 512
)

var (
	instance = GetInstance()
)

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
		//// 处理pong的闭包函数，刷新deadline
		//handlePong := func(string) error {
		//	c.Conn.SetReadDeadline(time.Now().Add(pongWait));
		//	return nil
		//}
		// 最小间隔15s，至少会有一个pong从浏览器发送过来，否则认为已经断联
		// c.Conn.SetPongHandler(handlePong)
		defer func() {
			instance.UnregisterCh <- c
			c.Conn.Close()
		}()

		for {
			// 目前是出错了退出
			msgType, msg, err := c.Conn.ReadMessage()
			// going away就意味着客户端已经释放了连接，服务器端因此也需要释放
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
				// 写入出错时直接退出，说明SendCh已经被manager关闭了
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
