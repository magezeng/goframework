package MasterNodeModel

import (
	"fmt"
	"sync"
	"time"
)

// 心跳数据，实际上就是serverStatus + masterNodeInfo
type MasterNodeServerSituation struct {
	// 没有时不进行设置
	IP        string                    `json:"ip"`
	CPUUsage  float32                   `json:"cpu_usage"`
	MemUsage  float32                   `json:"mem_usage"`
	DiskUsage float32                   `json:"disk_usage"`
	Nodes     map[string]CoinMasterNode `json:"nodes"`
	Lock      *sync.RWMutex
}

func (serverSituation MasterNodeServerSituation) String() string {
	return fmt.Sprintf("IP地址：%s, CPU使用: %f, 内存使用: %f, 硬盘使用: %f",
		serverSituation.IP, serverSituation.CPUUsage, serverSituation.MemUsage, serverSituation.DiskUsage)
}

func (serverSituation MasterNodeServerSituation) SetSyncTable(coinName string, info CoinMasterNode) {
	serverSituation.Lock.Lock()
	defer serverSituation.Lock.Unlock()
	serverSituation.Nodes[coinName] = info
}

func (serverSituation MasterNodeServerSituation) GetFromSyncTable(coinName string, id uint) CoinMasterNode {
	serverSituation.Lock.RLock()
	defer serverSituation.Lock.RUnlock()
	return serverSituation.Nodes[coinName]
}

func (serverSituation MasterNodeServerSituation) UpdateProcessStatusToSyncTable(coinName string, processStatus int8) {
	serverSituation.Lock.Lock()
	defer serverSituation.Lock.Unlock()
	masterNodeInfo := serverSituation.Nodes[coinName]
	masterNodeInfo.ProcessStatus = processStatus
	serverSituation.Nodes[coinName] = masterNodeInfo
}

func (serverSituation MasterNodeServerSituation) UpdateStartAtToSyncTable(coinName string, startAt time.Time) {
	serverSituation.Lock.Lock()
	defer serverSituation.Lock.Unlock()
	masterNodeInfo := serverSituation.Nodes[coinName]
	masterNodeInfo.StartAt = startAt
	serverSituation.Nodes[coinName] = masterNodeInfo
}
