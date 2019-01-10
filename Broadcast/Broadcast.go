package Broadcast

import (
	"github.com/satori/go.uuid"
	"sync"
)

type BroadcastContent struct {
	Message string
	Data    interface{}
}

type Broadcast struct {
	messageWaitChannel chan BroadcastContent
	receivers          map[string]BroadcastReceiver
	receiversRWMutex   *sync.RWMutex
}

func NewBroadcast() *Broadcast {
	broadcast := &Broadcast{make(chan BroadcastContent), map[string]BroadcastReceiver{}, new(sync.RWMutex)}
	broadcast.waitMessage()
	return broadcast
}

func (broadcast *Broadcast) waitMessage() {
	go func(broadcast *Broadcast) {
		for true {
			content := <-broadcast.messageWaitChannel
			broadcast.receiversRWMutex.RLock()
			for _, element := range broadcast.receivers {
				element.ReveiceChannel <- content
			}
			broadcast.receiversRWMutex.RUnlock()
		}
	}(broadcast)
}

func (broadcast *Broadcast) PostMessage(message string) {
	broadcast.PostMessageAndData(message, "")
}

func (broadcast *Broadcast) PostMessageAndData(message string, data interface{}) {
	broadcast.messageWaitChannel <- BroadcastContent{message, data}
}

func (broadcast *Broadcast) AddReceiver() BroadcastReceiver {

	UUID, _ := uuid.NewV4()
	tempUUID := UUID.String()
	receiver := BroadcastReceiver{tempUUID, make(chan BroadcastContent)}
	broadcast.receiversRWMutex.Lock()
	broadcast.receivers[tempUUID] = receiver
	broadcast.receiversRWMutex.Unlock()
	return receiver
}
func (broadcast *Broadcast) RemoveReceiver(receiver BroadcastReceiver) {
	delete(broadcast.receivers, receiver.id)
}
