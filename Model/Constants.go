package Model

const (
	CHANGE_DEPTH = iota
	CHANGE_KLINE
	CHANGE_POSITION
	CHANGE_ACCOUNT
	CHANGE_ORDERS
)

const (
	DATA_TYPE_DEPTH    = "depth"
	DATA_TYPE_POSITION = "position"
	DATA_TYPE_ACCOUNT  = "account"
	DATA_TYPE_KLINE    = "kline"
	DATA_TYPE_ORDER    = "order"
)
