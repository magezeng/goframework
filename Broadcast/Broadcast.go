package Broadcast

import (
	"github.com/satori/go.uuid"
)

type BroadcastContent struct {
	Message string
	Data    interface{}
}

type Broadcast struct {
	messageWaitChannel chan BroadcastContent
	receivers          map[string]BroadcastReceiver
}

func NewBroadcast() *Broadcast {
	broadcast := &Broadcast{make(chan BroadcastContent), map[string]BroadcastReceiver{}}
	broadcast.waitMessage()
	return broadcast
}

func (broadcast *Broadcast) waitMessage() {
	go func(broadcast *Broadcast) {
		for true {
			content := <-broadcast.messageWaitChannel
			for _, element := range broadcast.receivers {
				element.reveiceChannel <- content
			}
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
	tempUUID := uuid.NewV1().String()
	receiver := BroadcastReceiver{tempUUID, make(chan BroadcastContent)}
	broadcast.receivers[tempUUID] = receiver
	return receiver
}
func (broadcast *Broadcast) RemoveReceiver(receiver BroadcastReceiver) {
	delete(broadcast.receivers, receiver.id)
}
