package WebsocketsBaseCollect

import (
	"time"
)

type CollectAspectInterface interface {
	PreConnectToService()
	AfterConnectToService(apiKey string) error
	OnAbnormal()
	HandleData(interface{})
	GetPingString() string
	IsPong(response string) bool
	GetWebsocketsSendDataMaxRetry() int
	GetWebsocketsBreatheSendIntermit() time.Duration
	GetWebsocketsBreatheReciveTimeOut() time.Duration
	ChangeResponseToStruct([]byte) (interface{}, error)
}
