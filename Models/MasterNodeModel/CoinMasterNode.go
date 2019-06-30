package MasterNodeModel

import (
	"github.com/jinzhu/gorm"
	"time"
)

type CoinMasterNode struct {
	gorm.Model
	// 用于指示这个节点位于哪个服务器上面
	CoinName      string    `gorm:"not null;type:varchar(20)" form:"coin_name" binding:"required,max=20" json:"coin_name"`
	StartAt       time.Time `gorm:"type:timestamp" form:"start_at" json:"start_at"`
	Status        string    `gorm:"type:varchar(200)" form:"status" binding:"required" json:"status"`
	PayAddress    string    `gorm:"type:varchar(100)" form:"pay_address" binding:"required,max=100" json:"pay_address"`
	ProcessStatus int8      `json:"process_status"`
	MasterNodeID  uint
}
