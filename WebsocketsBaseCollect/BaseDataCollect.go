package WebsocketsBaseCollect

import (
	"bytes"
	"compress/flate"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net/url"
	"runtime/debug"
	"sync"
	"time"
	"tipu.com/go-framework/Broadcast"
)

type BaseDataCollect struct {
	url     string
	path    string
	connect *websocket.Conn

	IsPrintProcess bool

	//连接使用互斥锁，go的websockets链接只能同时发一个数据  用这个保证不会同时两个数据在发送
	connectMutex   *sync.Mutex
	receiveChannel chan interface{}
	//数据来源，数据处理的代理对象
	aspectDelegate CollectAspectInterface
	//使用通知来保证所有的线程可以第一时间退出
	closeBroadcast *Broadcast.Broadcast
	closeWaitGroup sync.WaitGroup
}

func (collect *BaseDataCollect) Init(url string, path string, aspectDelegate CollectAspectInterface) {
	fmt.Println("WebsocketsBaseCollect aspectDelegate.HandleData ", aspectDelegate.HandleData)
	fmt.Println("WebsocketsBaseCollect collect.HandleData ", collect.handleData)
	collect.url = url
	collect.path = path
	collect.connectMutex = new(sync.Mutex)
	collect.aspectDelegate = aspectDelegate

	collect.closeBroadcast = Broadcast.NewBroadcast()
	collect.closeWaitGroup = sync.WaitGroup{}

}

func (collect *BaseDataCollect) ConnectToService() {

	if collect.url == "" || collect.path == "" {
		panic("url 或 path 为空")
	}

	collect.aspectDelegate.PreConnectToService()

	defer func() {
		err := recover()

		if err == nil {
			collect.receiveChannel = make(chan interface{})
			collect.handleData()
			collect.CollectData()
			collect.Palpitate()
		} else {
			if collect.IsPrintProcess {
				fmt.Println(err)
			}
		}
		collect.aspectDelegate.AfterConnectToService(err == nil)
	}()

	tempURL := url.URL{Scheme: "wss", Host: collect.url, Path: collect.path, RawQuery: "compress=true"}
	log.Printf("发起链接 %s", tempURL.String())

	connect, _, err := websocket.DefaultDialer.Dial(tempURL.String(), nil)
	collect.connect = connect
	if err != nil {
		panic("链接错误:" + err.Error())
	}
	return
}

func (collect *BaseDataCollect) DisConnect() {
	collect.closeBroadcast.PostMessage("")
	if collect.connect != nil {
		_ = collect.connect.Close()
	}
	if collect.receiveChannel != nil {
		close(collect.receiveChannel)
	}
	collect.closeWaitGroup.Wait()
}

func (collect *BaseDataCollect) SendData(data interface{}) (err error) {
	var jsBytes []byte
	if stringData, isString := data.(string); isString {
		jsBytes = []byte(stringData)
	} else {
		jsBytes, err = json.Marshal(data)
	}
	if collect.IsPrintProcess {
		fmt.Println("发送了", string(jsBytes))
	}
	if err != nil {
		return
	}
	collect.connectMutex.Lock()
	defer collect.connectMutex.Unlock()

	//最多进行3次
	for i := 0; i < collect.aspectDelegate.GetWebsocketsSendDataMaxRetry(); i++ {
		err = collect.connect.WriteMessage(websocket.TextMessage, jsBytes)
		if err != nil {
			fmt.Println("error3:", err)
		} else {
			break
		}
	}
	return
}

func (collect *BaseDataCollect) Palpitate() {
	go func(collect *BaseDataCollect) {
		//waitGroup增加和释放，让外层在断开连接的时候可以等到所有的线程都退出之后才进行下一次开启
		collect.closeWaitGroup.Add(1)
		defer collect.closeWaitGroup.Done()
		//呼吸发送线程因为不会涉及到显示的资源竞争，可以不用增加资源判断，直接调用即可，在SendData函数内部会进行判断
		tempCloseReceiver := collect.closeBroadcast.AddReceiver()
		defer collect.closeBroadcast.RemoveReceiver(tempCloseReceiver)
		for {
			select {
			case <-time.After(time.Second * collect.aspectDelegate.GetWebsocketsBreatheSendIntermit()):
				if collect.connect == nil {
					continue
				} else {
					//呼吸发送发生错误时直接忽略错误，因为连续护送发不出去本来就是错误，可以允许等待呼吸接受不到数据时退出本次连接
					_ = collect.SendData(collect.aspectDelegate.GetPingString())
				}
			case <-tempCloseReceiver.ReveiceChannel:
				fmt.Println("收到断开连接的通知，呼吸发送线程退出")
				return
			}
		}
	}(collect)
}

func (collect *BaseDataCollect) CollectData() {
	go func(collect *BaseDataCollect) {
		//waitGroup增加和释放，让外层在断开连接的时候可以等到所有的线程都退出之后才进行下一次开启
		collect.closeWaitGroup.Add(1)
		defer func() {
			//此处group的Done 一定要放在调用外部异常之前 不然会造成外层在等内层Done  内层在等外层的所有执行完  的循环
			collect.closeWaitGroup.Done()
			if ee := recover(); ee != nil {
				debug.PrintStack()
				if err, isError := ee.(error); isError {
					fmt.Println("CollectData 抛出异常")
					collect.ThrowAbnormal(err)
				}
			}
		}()

		//读取信息的线程和外部通信的结构
		type Message struct {
			messageType int
			message     []byte
			err         error
		}

		tempCloseReceiver := collect.closeBroadcast.AddReceiver()
		defer collect.closeBroadcast.RemoveReceiver(tempCloseReceiver)
		for {
			//读取信息的线程和外部通信的通道
			tempChannel := make(chan Message)
			go func(messageChan chan Message) {
				defer close(tempChannel)
				tempMessage := Message{}
				tempMessage.messageType, tempMessage.message, tempMessage.err = collect.connect.ReadMessage()
				tempChannel <- tempMessage
			}(tempChannel)

			select {
			case tempMessage := <-tempChannel:
				if tempMessage.err != nil {
					panic(tempMessage.err)
					return
				}
				var tempText []byte
				switch tempMessage.messageType {
				case websocket.TextMessage:
					//不需要解压
					tempText = tempMessage.message
				case websocket.BinaryMessage:
					//需要解压
					var err error
					tempText, err = gzipDecode2(tempMessage.message)
					if err != nil {
						panic(err)
					}
				}
				if tempText != nil {
					if collect.IsPrintProcess {
						fmt.Println("接收到", string(tempText))
					}
					collect.receiveChannel <- tempText
				}
			case <-tempCloseReceiver.ReveiceChannel:
				fmt.Println("收到断开连接的通知，数据采集线程退出")
				return
			}
		}
	}(collect)
}

func (collect *BaseDataCollect) ThrowAbnormal(tempError error) {
	//假如调用此函数的地方发生了异常   则调用代理的异常处理
	if tempError != nil {
		fmt.Println("抛出异常打印:", tempError)
		collect.aspectDelegate.OnAbnormal()
	}
	fmt.Println(tempError)
}

func (collect *BaseDataCollect) handleData() {
	go func(collect *BaseDataCollect) {
		//waitGroup增加和释放，让外层在断开连接的时候可以等到所有的线程都退出之后才进行下一次开启
		collect.closeWaitGroup.Add(1)
		defer func() {
			//此处group的Done 一定要放在调用外部异常之前 不然会造成外层在等内层Done  内层在等外层的所有执行完  的循环
			collect.closeWaitGroup.Done()
			if ee := recover(); ee != nil {
				debug.PrintStack()
				if err, isError := ee.(error); isError {
					fmt.Println("handleData抛出异常")
					collect.ThrowAbnormal(err)
				}
			}
		}()
		tempCloseReceiver := collect.closeBroadcast.AddReceiver()
		defer collect.closeBroadcast.RemoveReceiver(tempCloseReceiver)
		for {
			select {
			case responseBytes := <-collect.receiveChannel:
				//将interface转为String
				responseString := string(responseBytes.([]byte))
				//过滤掉呼吸返回
				if collect.aspectDelegate.IsPong(responseString) {
					continue
				}
				//将业务数据转换为Map
				obj, err := collect.aspectDelegate.ChangeResponseToStruct(responseBytes.([]byte))
				if err != nil {
					log.Println("error1:", err)
					continue
				}
				//将业务Map数据通过代理传递到外层

				fmt.Printf("%p", collect.aspectDelegate)
				fmt.Println("-----------------1", collect.aspectDelegate)
				fmt.Println("-----------------2", collect.aspectDelegate.HandleData)
				fmt.Println("-----------------3", obj)
				collect.aspectDelegate.HandleData(obj)
				fmt.Println("已经调用了上层")

			case <-time.After(time.Second * collect.aspectDelegate.GetWebsocketsBreatheReciveTimeOut()):
				//case <-time.After(time.Millisecond * 50):
				//发送呼吸数据是5秒发送一次  10秒未收到呼吸返回则表示当前链接出现了故障 需要进行链接死亡后的操作
				panic(errors.New("2呼吸超时"))
				return
			case <-tempCloseReceiver.ReveiceChannel:
				//当接收到外层的断开连接的通知之后   退出循环  结束该子线程
				return
			}
		}
	}(collect)
}

func gzipDecode2(in []byte) ([]byte, error) {
	reader := flate.NewReader(bytes.NewReader(in))
	defer func() {
		_ = reader.Close()
	}()
	return ioutil.ReadAll(reader)
}
