package WebsocketsBaseCollect

import (
	"../Logger"
	"../Retry"
	"../Utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"net/url"
	"sync"
	"time"
)

type BaseDataCollect struct {
	url  string
	path string
	// 当前只保持一个websocket连接
	connect *websocket.Conn
	// 连接使用互斥锁，go的websockets链接只能同时发一个数据  用这个保证不会同时两个数据在发送
	connectMutex   *sync.Mutex
	receiveChannel chan interface{}
	// 数据来源，数据处理的代理对象
	aspectDelegate CollectAspectInterface
	// 使用通知来保证所有的线程可以第一时间退出
	ctx            context.Context
	cancel         context.CancelFunc
	closeWaitGroup *sync.WaitGroup
	// 对每个实例都是单独的
	logger *Logger.Logger
}

func (collect *BaseDataCollect) Init(url string, path string, apiKey string, aspectDelegate CollectAspectInterface) {
	collect.url = url
	collect.path = path
	collect.connectMutex = new(sync.Mutex)
	collect.aspectDelegate = aspectDelegate
	// 初始化链接退出上下文及wait group
	collect.ctx, collect.cancel = context.WithCancel(context.TODO())
	collect.closeWaitGroup = new(sync.WaitGroup)
	// 根据apiKey初始化日志
	collect.setCollectLogger(apiKey)
}

func (collect *BaseDataCollect) ConnectToService() (err error) {
	if collect.url == "" || collect.path == "" {
		err = errors.New("url 或 path 为空")
		return
	}

	// 连接前进行的操作
	err = collect.aspectDelegate.PreConnectToService()
	if err != nil {
		return
	}
	// 初始化一个ws连接
	tempURL := url.URL{Scheme: "wss", Host: collect.url, Path: collect.path, RawQuery: "compress=true"}
	collect.logger.Info("发起链接: ", tempURL.String())
	// 尝试三次连接，间隔2s，失败时返回错误
	result, err1 := Retry.Retry(3, 2*time.Second, func() (result interface{}, err error) {
		// 超时改为20秒
		websocket.DefaultDialer.HandshakeTimeout = time.Second * 20
		websocket.DefaultDialer.EnableCompression = true
		websocket.DefaultDialer.WriteBufferSize = 1024
		connect, _, err := websocket.DefaultDialer.Dial(tempURL.String(), nil)
		if err != nil {
			return
		}
		result = connect
		return
	})
	if err1 != nil {
		err = err1
		return
	}
	// 连接后进行的操作
	collect.connect = result.(*websocket.Conn)
	collect.receiveChannel = make(chan interface{})
	// 启动三个异步线程
	collect.Palpitate()
	collect.CollectData()
	collect.handleData()

	return collect.aspectDelegate.AfterConnectToService()
}

// DisConnect 断开数据采集器连接
func (collect *BaseDataCollect) DisConnect() {
	// 保证此处只能发出一次取消命令
	if collect.ctx.Err() == nil {
		// 通知上下文取消
		collect.cancel()
		collect.logger.Warn("发出了关闭所有线程和连接的命令")
		// 等待线程终止
		collect.closeWaitGroup.Wait()
		// 关闭websocket连接
		if collect.connect != nil {
			_ = collect.connect.Close()
		}
		// TODO: graceful close of a channel
		close(collect.receiveChannel)
		collect.logger.Warn("数据采集器释放所有线程和连接完成")
	}
}

// SendData 写入数据到websocket
func (collect *BaseDataCollect) SendData(data interface{}) (err error) {
	var jsBytes []byte
	if stringData, isString := data.(string); isString {
		jsBytes = []byte(stringData)
	} else {
		jsBytes, err = json.Marshal(data)
		if err != nil {
			return
		}
	}

	collect.connectMutex.Lock()
	defer collect.connectMutex.Unlock()

	// 最多进行3次
	maxRetry := collect.aspectDelegate.GetWebsocketsSendDataMaxRetry()
	_, err = Retry.Retry(maxRetry, 1, func() (result interface{}, err error) {
		err = collect.connect.WriteMessage(websocket.TextMessage, jsBytes)
		if err != nil {
			collect.logger.Error("websocket发送数据失败:", err)
			return
		}
		collect.logger.Info("websocket发送了: ", string(jsBytes))
		return
	})
	return
}

// Palpitate 心跳线程，相同间隔发送一个ping
// 不涉及资源竞争，接收到cancel就直接退出
func (collect *BaseDataCollect) Palpitate() {
	go func() {
		defer collect.logger.Warn("Palpitate 呼吸发送线程已经退出")
		for {
			select {
			case <-time.After(time.Second * collect.aspectDelegate.GetWebsocketsBreatheSendIntermit()):
				if collect.connect == nil {
					// 等待连接建立
					continue
				} else {
					// 呼吸发送发生错误时直接忽略错误，因为连续护送发不出去本来就是错误，可以允许等待呼吸接受不到数据时退出本次连接
					_ = collect.SendData(collect.aspectDelegate.GetPingString())
				}
			case <-collect.ctx.Done():
				collect.logger.Warn("Palpitate 呼吸发送线程收到立即退出的通知")
				return
			}
		}
	}()
}

// CollectData 将接收到的数据转换格式后发送到receiveChannel中,让handleData去处理
// 这个方法是异步的，如果出错或者其它线程出错，则本线程中止
func (collect *BaseDataCollect) CollectData() {
	go func() {
		collect.closeWaitGroup.Add(1)
		defer func() {
			// 此处group的Done 一定要放在调用外部异常之前 不然会造成外层在等内层Done  内层在等外层的所有执行完  的循环
			collect.closeWaitGroup.Done()
			if ee := recover(); ee != nil {
				if err, isError := ee.(error); isError {
					// 断开重连
					collect.ThrowAbnormal(err)
				}
			}
		}()
		// 读取信息的线程和外部通信的结构
		type Message struct {
			messageType int
			message     []byte
			err         error
		}

		for {
			// 读取信息的线程和外部通信的通道
			tempChannel := make(chan Message)
			go func() {
				defer close(tempChannel)
				tempMessage := Message{}
				tempMessage.messageType, tempMessage.message, tempMessage.err = collect.connect.ReadMessage()
				tempChannel <- tempMessage
			}()
			select {
			case tempMessage := <-tempChannel:
				if tempMessage.err != nil {
					panic(NewTraceWithMsg("CollectData 数据采集线程出错: " + tempMessage.err.Error()))
				}

				var tempText []byte
				switch tempMessage.messageType {
				case websocket.TextMessage:
					// 不需要解压
					tempText = tempMessage.message
				case websocket.BinaryMessage:
					var err error
					// 解压缩失败的情况下，报错
					tempText, err = Utils.GzipDecode2(tempMessage.message)
					if err != nil {
						panic(NewTraceWithMsg("CollectData 数据采集线程出错: " + err.Error()))
					}
				}

				if tempText != nil {
					collect.logger.Info("CollectData 接收到: ", string(tempText))
					collect.receiveChannel <- tempText
				}
			case <-collect.ctx.Done():
				collect.logger.Info("CollectData 数据采集线程收到立即退出的通知，立即退出")
				return
			}
		}
	}()
}

func (collect *BaseDataCollect) ThrowAbnormal(tempError error) {
	// 假如调用此函数的地方发生了异常   则调用代理的异常处理
	collect.logger.Error("ThrowAbnormal 异常详细信息:", tempError.Error())
	collect.aspectDelegate.OnAbnormal()
}

// handleData 接收到数据时进行转换处理，并发送到外层的HandleData进行进一步的处理
// 这个方法是异步的
func (collect *BaseDataCollect) handleData() {
	go func() {
		collect.closeWaitGroup.Add(1)
		defer func() {
			//此处group的Done 一定要放在调用外部异常之前 不然会造成外层在等内层Done  内层在等外层的所有执行完  的循环
			collect.closeWaitGroup.Done()
			if ee := recover(); ee != nil {
				if err, isError := ee.(error); isError {
					// 断开重连
					collect.ThrowAbnormal(err)
				}
			}
		}()
		for {
			select {
			case responseBytes := <-collect.receiveChannel:
				// 将interface转为String
				responseString := string(responseBytes.([]byte))
				//过滤掉呼吸返回
				if collect.aspectDelegate.IsPong(responseString) {
					continue
				}
				//将业务数据转换为Map
				obj, err := collect.aspectDelegate.ChangeResponseToStruct(responseBytes.([]byte))
				if err != nil {
					panic(NewTraceWithMsg("handleData 业务数据转换为map失败: " + err.Error()))
				}
				// 将业务Map数据通过代理传递到外层
				// 异步是为了防止循环等待
				collect.aspectDelegate.HandleData(obj)

			case <-time.After(time.Second * collect.aspectDelegate.GetWebsocketsBreatheReciveTimeOut()):
				// 发送呼吸数据是5秒发送一次  10秒未收到呼吸返回则表示当前链接出现了故障 需要进行链接死亡后的操作
				timeout := collect.aspectDelegate.GetWebsocketsBreatheReciveTimeOut()
				panic(NewTraceWithMsg(fmt.Sprintf("handleData %ds呼吸超时", timeout)))

			case <-collect.ctx.Done():
				// 当接收到外层的断开连接的通知之后   退出循环  结束该子线程
				collect.logger.Warn("handleData 数据处理线程收到立即退出的通知，立即退出")
				return
			}
		}
	}()
}

// setCollectLogger 单独设置每一个用户的log，以防冲突
func (collect *BaseDataCollect) setCollectLogger(apiKey string) {
	collect.logger = Logger.NewLogger()
	// 有apiKey说明是用户请求
	if apiKey != "" {
		collect.logger.SetFileWriter("base-" + apiKey[0:5] + ".log")
		collect.logger.SetPrefix("[Base][用户:" + apiKey[0:5] + "]")
	} else {
		collect.logger.SetFileWriter("base-public.log")
		collect.logger.SetPrefix("[Base]")
	}
}
