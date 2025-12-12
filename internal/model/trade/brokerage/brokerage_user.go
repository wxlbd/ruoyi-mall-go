package brokerage

import (
	"time"

	"gorm.io/gorm"
)

// BrokerageUser 分销用户
type BrokerageUser struct {
	ID               int64          `gorm:"primaryKey;autoIncrement;comment:用户编号"`
	BindUserID       int64          `gorm:"column:bind_user_id;default:0;comment:推广员编号"`
	BindUserTime     *time.Time     `gorm:"column:bind_user_time;comment:推广员绑定时间"`
	BrokerageEnabled bool           `gorm:"column:brokerage_enabled;type:tinyint(1);default:0;comment:是否有分销资格"`
	BrokerageTime    *time.Time     `gorm:"column:brokerage_time;comment:成为分销员时间"`
	BrokeragePrice   int            `gorm:"column:brokerage_price;default:0;comment:可用佣金"`
	FrozenPrice      int            `gorm:"column:frozen_price;default:0;comment:冻结佣金"`
	Creator          string         `gorm:"column:creator;size:64;default:'';comment:创建者"`
	Updater          string         `gorm:"column:updater;size:64;default:'';comment:更新者"`
	CreatedAt        time.Time      `gorm:"column:create_time;autoCreateTime;comment:创建时间"`
	UpdatedAt        time.Time      `gorm:"column:update_time;autoUpdateTime;comment:更新时间"`
	DeletedAt        gorm.DeletedAt `gorm:"column:deleted;index;comment:删除时间"`
	Deleted          bool           `gorm:"column:deleted;type:tinyint(1);not null;default:0;comment:是否删除"`
}

func (BrokerageUser) TableName() string {
	return "trade_brokerage_user"
}
