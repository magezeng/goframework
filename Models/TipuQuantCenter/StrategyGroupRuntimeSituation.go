package TipuQuantCenter

type StrategyGroupRuntime struct {
	StrategyGroupId uint    `json:"strategy_group_id"`
	Status          *bool   `json:"status"`
	CPUUsage        float32 `json:"cpu_usage"`
	MemUsage        float32 `json:"mem_usage"`
}
