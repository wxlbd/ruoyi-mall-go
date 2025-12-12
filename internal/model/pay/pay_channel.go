package pay

import (
	"backend-go/internal/model"
	"time"

)

// PayChannel 支付渠道 DO
type PayChannel struct {
	ID        int64          `gorm:"primaryKey;autoIncrement;comment:渠道编号" json:"id"`
	Code      string         `gorm:"size:32;not null;comment:渠道编码" json:"code"`
	Status    int            `gorm:"default:0;not null;comment:状态" json:"status"` // 参见 CommonStatusEnum
	FeeRate   float64        `gorm:"default:0;comment:渠道费率" json:"feeRate"`       // 单位：百分比
	Remark    string         `gorm:"size:255;default:'';comment:备注" json:"remark"`
	AppID     int64          `gorm:"column:app_id;not null;comment:应用编号" json:"appId"`
	Config    string         `gorm:"type:json;serializer:json;comment:支付渠道配置" json:"config"` // JSON Configuration
	TenantID  int64          `gorm:"column:tenant_id;default:0;comment:租户编号" json:"tenantId"`
	Creator   string         `gorm:"size:64;default:'';comment:创建者" json:"creator"`
	Updater   string         `gorm:"size:64;default:'';comment:更新者" json:"updater"`
	CreatedAt time.Time      `gorm:"column:create_time;autoCreateTime;comment:创建时间" json:"createTime"`
	UpdatedAt time.Time      `gorm:"column:update_time;autoUpdateTime;comment:更新时间" json:"updateTime"`
	Deleted   model.BitBool  `gorm:"column:deleted;softDelete:flag;default:0;comment:是否删除" json:"deleted"`
}

func (PayChannel) TableName() string {
	return "pay_channel"
}
