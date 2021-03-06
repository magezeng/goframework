package MasterNodeModel

// 心跳数据，实际上就是serverStatus + masterNodeInfo
type MasterNodeServerSituation struct {
	// 没有时不进行设置
	IP        string  `json:"ip"`
	Version   string  `json:"version"`
	CPUUsage  float32 `json:"cpu_usage"`
	MemUsage  float32 `json:"mem_usage"`
	DiskUsage float32 `json:"disk_usage"`
}
