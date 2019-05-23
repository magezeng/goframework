package Websocket

import (
	"github.com/gorilla/websocket"
	"strings"
	"sync"
	"time"
)

const (
	// 1.5分钟没有收到客户端的返回消息，需要释放连接
	pongWait = 90 * time.Second
	// 服务端发起ping的时间间隔
	pingPeriod = 15 * time.Second
	// 消息读取最大长度，字节
	maxMessageSize = 2048
	// 打包数据的发送频率
	messageWritePeriod = 500 * time.Millisecond
	// 最大缓冲的消息数量
	bufferSize = 500
)

var (
	instance = GetInstance()
)

type Server struct {
	// 用户的token，一个token只对应一个连接
	Token string
	// 发送数据的连接
	Conn *websocket.Conn
	// 缓冲通道
	BufferCh chan string
	// websocket的写锁，防止ping和发送普通消息引起concurrent write错误
	lock *sync.Mutex
}

// start 启动ws服务端
func (s *Server) start() {
	s.read()
	s.write()
}

// read 从ws连接中读取数据
func (s *Server) read() {
	go func() {
		// 设置读取的deadline
		s.Conn.SetReadDeadline(time.Now().Add(pongWait));
		s.Conn.SetReadLimit(maxMessageSize)
		// 此API目前浏览器不支持，需要采用判断消息字符串的方式进行
		//// 处理pong的闭包函数，刷新deadline
		//handlePong := func(string) error {
		//	c.Conn.SetReadDeadline(time.Now().Add(pongWait));
		//	return nil
		//}
		// 最小间隔15s，至少会有一个pong从浏览器发送过来，否则认为已经断联
		// c.Conn.SetPongHandler(handlePong)
		defer func() {
			instance.UnregisterCh <- s
			s.Conn.Close()
		}()

		for {
			// 目前是出错了退出
			msgType, msg, err := s.Conn.ReadMessage()
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
					s.Conn.SetReadDeadline(time.Now().Add(pongWait));
				}
			default:
				logger.Warn("不支持的消息类型，忽略!", msgType)
			}
		}
	}()
}

// 向连接中写入数据
func (s *Server) write() {
	go func() {
		ticker := time.NewTicker(pingPeriod)
		messageTicker := time.NewTicker(messageWritePeriod)
		defer func() {
			ticker.Stop()
			messageTicker.Stop()
			s.Conn.Close()
		}()

		for {
			select {
			case <-ticker.C:
				err := s.writeMessageSync(websocket.TextMessage, "ping")
				if err != nil {
					logger.Error("发送ping失败，可能是浏览器端已经离线: ", s.Token)
					return
				}
				logger.Info("发送ping到浏览器", s.Token)
			case <-messageTicker.C:
				s.packageMessage()
			}
		}
	}()
}

func (s *Server) packageMessage() {
	if len(s.BufferCh) > 0 {
		go func() {
			builder := strings.Builder{}
			chLength := len(s.BufferCh)
			// 读取当前内部的数量的一批数据，默认是10个
			for i := 0; i < chLength; i++ {
				c := <-s.BufferCh
				builder.WriteString(c + "\n")
			}
			if builder.String() == "" {
				return
			}
			// 如果有消息才写入，否则什么都不做
			err := s.writeMessageSync(websocket.TextMessage, builder.String())
			if err != nil {
				logger.Info("写入数据失败: ", s.Token)
				return
			}

			return
		}()
	}
}

// 同步写入数据到websocket
func (s *Server) writeMessageSync(messageType int, message string) (err error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	err = s.Conn.WriteMessage(messageType, []byte(message))
	logger.Info("写入了数据: ", message)
	return
}
