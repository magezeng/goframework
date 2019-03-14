package Model

// Result 返回结果时的通用结构
type Result struct {
	Code int16       `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

// ResultList 列表型结果的结构
type ResultList struct {
	Code  int16         `json:"code"`
	Rows  []interface{} `json:"rows"`
	Total int           `json:"total"`
	Msg   string        `json:"msg"`
}
