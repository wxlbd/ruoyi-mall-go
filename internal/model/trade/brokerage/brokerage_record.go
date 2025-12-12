package brokerage

import (
	"time"

	"gorm.io/gorm"
)

// BrokerageRecord 佣金记录
type BrokerageRecord struct {
	ID              int64          `gorm:"primaryKey;autoIncrement;comment:编号"`
	UserID          int64          `gorm:"column:user_id;not null;comment:用户编号"`
	BizID           string         `gorm:"column:biz_id;size:64;not null;comment:业务编号"`
	BizType         int            `gorm:"column:biz_type;not null;comment:业务类型"`
	Title           string         `gorm:"column:title;size:64;not null;comment:标题"`
	Description     string         `gorm:"column:description;size:255;not null;comment:说明"`
	Price           int            `gorm:"column:price;not null;comment:金额"`
	TotalPrice      int            `gorm:"column:total_price;not null;comment:当前总佣金"`
	Status          int            `gorm:"column:status;not null;comment:状态"`
	FrozenDays      int            `gorm:"column:frozen_days;default:0;comment:冻结时间（天）"`
	UnfreezeTime    *time.Time     `gorm:"column:unfreeze_time;comment:解冻时间"`
	SourceUserLevel int            `gorm:"column:source_user_level;default:0;comment:来源用户等级"`
	SourceUserID    int64          `gorm:"column:source_user_id;default:0;comment:来源用户编号"`
	Creator         string         `gorm:"column:creator;size:64;default:'';comment:创建者"`
	Updater         string         `gorm:"column:updater;size:64;default:'';comment:更新者"`
	CreatedAt       time.Time      `gorm:"column:create_time;autoCreateTime;comment:创建时间"`
	UpdatedAt       time.Time      `gorm:"column:update_time;autoUpdateTime;comment:更新时间"`
	DeletedAt       gorm.DeletedAt `gorm:"column:deleted;index;comment:删除时间"`
	Deleted         bool           `gorm:"column:deleted;type:tinyint(1);not null;default:0;comment:是否删除"`
}

func (BrokerageRecord) TableName() string {
	return "trade_brokerage_record"
}
