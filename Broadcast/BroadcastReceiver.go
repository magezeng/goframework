package Broadcast

type BroadcastReceiver struct {
	id             string
	reveiceChannel chan BroadcastContent
}

func (receiver *BroadcastReceiver) ReceiveMessage() (content BroadcastContent) {
	return <-receiver.reveiceChannel
}
