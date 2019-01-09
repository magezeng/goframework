package Broadcast

type BroadcastReceiver struct {
	id             string
	ReveiceChannel chan BroadcastContent
}
