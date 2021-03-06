package MasterNodeModel

type Command struct {
	Key  string
	Path string
	Args []string
}

type CoinNodeConcoct struct {
	// 节点在数据库内的编号
	MasterNodeId uint
	// 币种名称
	CoinName string
	// 节点安装文件URL
	InstallFileURL string
	// 文件处理命令，节点端不需要关心文件内部是什么样，只需要把文件下载之后，将文件按照以下的命令进行替换内容，需要替换的内容在参数内
	FileHandleCommand Command
	// 节点端不需要关心具体获取状态的方式是什么，只需要按照中心传过来的命令运行，将运行结果回传到中心即可
	SituationCommand map[string]Command
}
type MasterNodeSituation struct {
	// 节点在数据库内的编号
	MasterNodeId uint
	// 币种名称
	CoinName string
	// 币种状态结果
	SituationResult map[string]string
}

type CoinNodeDelete struct {
	// 节点在数据库内的编号
	MasterNodeId uint
	// 币种名称
	CoinName string
	//节点删除文件URL
	DeleteFileURL string
}
