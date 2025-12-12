package member

import (
	"time"

	"backend-go/internal/model"

	"gorm.io/gorm"
)

// MemberConfig 会员配置
// Table: member_config
type MemberConfig struct {
	ID                        int64          `gorm:"primaryKey;autoIncrement;comment:编号" json:"id"`
	PointTradeDeductEnable    int            `gorm:"column:point_trade_deduct_enable;type:tinyint;default:0;comment:积分抵扣开关" json:"pointTradeDeductEnable"` // 1-开启 0-关闭
	PointTradeDeductUnitPrice int            `gorm:"column:point_trade_deduct_unit_price;default:0;comment:积分抵扣单位价格" json:"pointTradeDeductUnitPrice"`     // 积分抵扣，单位：分
	PointTradeDeductMaxPrice  int            `gorm:"column:point_trade_deduct_max_price;default:0;comment:积分抵扣最大值" json:"pointTradeDeductMaxPrice"`
	PointTradeGivePoint       int            `gorm:"column:point_trade_give_point;default:0;comment:1 元赠送多少分" json:"pointTradeGivePoint"`
	Creator                   string         `gorm:"column:creator;size:64;default:'';comment:创建者"`
	Updater                   string         `gorm:"column:updater;size:64;default:'';comment:更新者"`
	CreatedAt                 time.Time      `gorm:"column:create_time;autoCreateTime;comment:创建时间"`
	UpdatedAt                 time.Time      `gorm:"column:update_time;autoUpdateTime;comment:更新时间"`
	DeletedAt                 gorm.DeletedAt `gorm:"column:deleted;index;comment:删除时间"`
	Deleted                   model.BitBool  `gorm:"column:deleted;type:tinyint(1);not null;default:0;comment:是否删除"`
}

func (MemberConfig) TableName() string {
	return "member_config"
}
