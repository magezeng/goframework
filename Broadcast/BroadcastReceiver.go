package Broadcast

type BroadcastReceiver struct {
	Id             string
	ReceiveChannel chan BroadcastContent
}
