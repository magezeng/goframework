package TipuQuantCenter

type ControlTerminal struct {
	IP        string  `json:"ip"`
	Version   string  `json:"version"`
	CPUUsage  float32 `json:"cpu_usage"`
	MemUsage  float32 `json:"mem_usage"`
	DiskUsage float32 `json:"disk_usage"`
}
