package MasterNodeModel

type CoinNodeSituation struct {
	// 币种在中心数据库内的编号
	CoinIndex uint
	// 币种名称
	CoinName string
	// 节点状态
	Situation map[string]string
}
