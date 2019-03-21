package Model

const (
	EVENT_DEPTH_CHANGE = iota
	EVENT_KLINE_CHANGE
	EVENT_POSITION_CHANGE
	EVENT_ACCOUNT_CHANGE
	EVENT_ORDERS_CHANGE
)

const (
	DATA_TYPE_DEPTH    = "depth"
	DATA_TYPE_POSITION = "position"
	DATA_TYPE_ACCOUNT  = "account"
	DATA_TYPE_KLINE    = "kline"
	DATA_TYPE_ORDER    = "order"
)

func ConvertEventTypeToDataType(eventType int) string {
	switch eventType {
	case EVENT_DEPTH_CHANGE:
		return DATA_TYPE_DEPTH
	case EVENT_KLINE_CHANGE:
		return DATA_TYPE_KLINE
	case EVENT_POSITION_CHANGE:
		return DATA_TYPE_POSITION
	case EVENT_ACCOUNT_CHANGE:
		return DATA_TYPE_ACCOUNT
	case EVENT_ORDERS_CHANGE:
		return DATA_TYPE_ORDER
	default:
		return ""
	}
}

func ConvertDataTypeToEventType(dataType string) int {
	switch dataType {
	case DATA_TYPE_DEPTH:
		return EVENT_DEPTH_CHANGE
	case DATA_TYPE_KLINE:
		return EVENT_KLINE_CHANGE
	case DATA_TYPE_ACCOUNT:
		return EVENT_ACCOUNT_CHANGE
	case DATA_TYPE_POSITION:
		return EVENT_POSITION_CHANGE
	case DATA_TYPE_ORDER:
		return EVENT_ORDERS_CHANGE
	default:
		return -1
	}
}
