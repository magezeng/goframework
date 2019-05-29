package Websocket

import (
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"time"
	"tipu.com/go-framework/Logger"
)

var (
	managerInstance *ServerManager
	once            sync.Once
	logger          *Logger.Logger
)

type ServerManager struct {
	// 现在保持的连接,key为token
	Servers map[string]*Server
	// 增加客户端的channel
	RegisterCh chan *Server
	// 断开客户端的channel
	UnregisterCh chan *Server
}

func GetInstance() *ServerManager {
	once.Do(func() {
		manager := &ServerManager{
			RegisterCh:   make(chan *Server),
			UnregisterCh: make(chan *Server),
			// 可能以后有广播的需求，所以需要一个map
			Servers: make(map[string]*Server),
		}
		manager.startManager()
		logger = Logger.NewLogger().SetFileWriter("ws.log")
		logger.SetPrefix("Websocket")
		managerInstance = manager
	})
	return managerInstance
}

// 运行一个新的服务端
func (m *ServerManager) StartServer(token string, w http.ResponseWriter, r *http.Request, h http.Header) (err error) {
	conn, err := (&websocket.Upgrader{
		EnableCompression: true,
		ReadBufferSize:    1024 * 1024 * 1,
		WriteBufferSize:   1024 * 1024 * 1,
		// 设置跨域
		CheckOrigin: func(r *http.Request) bool { return true },
		// 设置握手超时的时间
		HandshakeTimeout: 1 * time.Minute,
		Subprotocols:     []string{r.Header.Get("Sec-WebSocket-Protocol")},
	}).Upgrade(w, r, h)
	if err != nil {
		return
	}
	// 设置缓冲区的大小，短token取最后30位
	server := &Server{Token: token, ShortToken: token[len(token)-20:], Conn: conn, BufferCh: make(chan string, bufferSize), lock: new(sync.Mutex)}
	// 注册到server集合中
	managerInstance.RegisterCh <- server

	return
}

// WriteMessage 写入数据到指定的server连接中
func (m *ServerManager) WriteMessage(token string, msg string) (err error) {
	if server, ok := m.Servers[token]; ok {
		server.BufferCh <- msg
	}
	return
}

// startManager 运行连接管理器
func (m *ServerManager) startManager() {
	go func() {
		for {
			select {
			// 连接信号
			case server := <-m.RegisterCh:
				m.Servers[server.Token] = server
				server.start()
				logger.Info("连接创建成功: ", server.ShortToken)
			// 断开信号
			case server := <-m.UnregisterCh:
				if _, ok := m.Servers[server.Token]; ok {
					// 此时不需要发送信息
					close(server.BufferCh)
					delete(m.Servers, server.Token)
					logger.Info("连接关闭成功: ", server.ShortToken)
				}
			}
		}
	}()
}
