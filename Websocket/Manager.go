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
func (m *ServerManager) StartNewServer(key string, w http.ResponseWriter, r *http.Request, h http.Header) (err error) {
	conn, err := (
		&websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			// 设置跨域
			CheckOrigin: func(r *http.Request) bool { return true },
			// 设置握手超时的时间
			HandshakeTimeout: 1 * time.Minute,
			Subprotocols:     []string{r.Header.Get("Sec-WebSocket-Protocol")},
		}).Upgrade(w, r, h)
	if err != nil {
		return
	}
	//// 对于同一个用户来说，也有可能开了多个连接，所以此处不能用用户id来作为key
	//uuidModel, _ := uuid.NewV4()
	//id = uuidModel.String()

	server := &Server{Token: key, Conn: conn, SendCh: make(chan []byte)}
	// 下面都是异步操作
	server.start()
	// 注册到server集合中
	instance.RegisterCh <- server
	return
}

// 写入数据到指定的server连接中
func (m *ServerManager) WriteMessage(token string, msg string) (err error) {
	if server, ok := m.Servers[token]; ok {
		server.SendCh <- []byte(msg)
		logger.Info("发送数据: ", server.Token)
	}
	return
}

// 运行客户端管理器
func (m *ServerManager) startManager() {
	go func() {
		for {
			select {
			// 连接信号
			case server := <-m.RegisterCh:
				m.Servers[server.Token] = server
				logger.Info("连接创建成功: ", server.Token)
			// 断开信号
			case server := <-m.UnregisterCh:
				if _, ok := m.Servers[server.Token]; ok {
					// 此时不需要发送信息
					close(server.SendCh)
					delete(m.Servers, server.Token)
					logger.Info("连接关闭成功: ", server.Token)
				}
			}
		}
	}()
}
