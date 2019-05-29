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
	pingPeriod = 30 * time.Second
	// 消息读取最大长度，字节
	maxMessageSize = 2048
	// 打包数据的发送频率
	messageWritePeriod = 500 * time.Millisecond
	// 最大缓冲的消息数量
	bufferSize = 500
)

type Server struct {
	// 用来做记录的短token
	ShortToken string
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
	s.Conn.SetReadLimit(maxMessageSize)
	s.read()
	s.write()
}

// read 从ws连接中读取数据
func (s *Server) read() {
	go func() {
		defer func() {
			// 必须释放而不管中间关闭发生的错误
			_ = s.Conn.Close()
			GetInstance().UnregisterCh <- s
		}()

		// 设置读取的deadline
		err := s.Conn.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			logger.Error("首次设置ReadDeadline失败: ", err, s.ShortToken)
			return
		}
		// Pong Handler目前浏览器不支持，无法接收ping帧，但是服务端相互连接可以支持
		// 这里采用判断消息字符串的方式进行

		//// 处理pong的闭包函数，刷新deadline
		//handlePong := func(string) error {
		//	c.Conn.SetReadDeadline(time.Now().Add(pongWait));
		//	return nil
		//}
		// c.Conn.SetPongHandler(handlePong)

		for {
			// 目前是出了任何错退出，包括1001错误
			msgType, msg, err := s.Conn.ReadMessage()
			if err != nil {
				// websocket.IsUnexpectedCloseError(websocket.CloseGoingAway)
				logger.Warn("读取消息出错: ", err, s.ShortToken)
				return
			}

			switch msgType {
			case websocket.TextMessage:
				// 目前只有可能收到浏览器发来的pong而不会有其它消息
				logger.Info("收到了浏览器发来的消息: ", string(msg), s.ShortToken)
				if string(msg) == "pong" {
					err = s.Conn.SetReadDeadline(time.Now().Add(pongWait))
					if err != nil {
						logger.Error("刷新ReadDeadline失败: ", err, s.ShortToken)
						return
					}
				}
			default:
				logger.Warn("不支持的消息类型: ", msgType, s.ShortToken)
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
			_ = s.Conn.Close()
		}()

		for {
			select {
			case <-ticker.C:
				err := s.writeMessageSync(websocket.TextMessage, "ping")
				if err != nil {
					logger.Error("发送ping失败，可能是浏览器端已经离线: ", s.ShortToken)
					return
				}
				logger.Info("发送ping到浏览器", s.ShortToken)
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
			// 读取当前通道内部的数量的一批数据
			for i := 0; i < chLength; i++ {
				c := <-s.BufferCh
				builder.WriteString(c + "\n")
			}
			// 如果有消息才写入，否则什么都不做
			err := s.writeMessageSync(websocket.TextMessage, builder.String())
			if err != nil {
				logger.Info("写入数据失败: ", s.ShortToken)
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
	logger.Info("准备写入数据: ", message, s.ShortToken)
	return
}
