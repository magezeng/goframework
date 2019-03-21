package Broadcast

type BroadcastReceiver struct {
	Id             string
	ReveiceChannel chan BroadcastContent
}
