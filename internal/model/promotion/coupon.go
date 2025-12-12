package promotion

import (
	"time"

	"gorm.io/gorm"
)

// PromotionCoupon 优惠券
type PromotionCoupon struct {
	ID              int64          `gorm:"primaryKey;autoIncrement;comment:优惠券编号"`
	TemplateID      int64          `gorm:"column:template_id;type:bigint;not null;comment:模板编号"`
	Name            string         `gorm:"column:name;type:varchar(64);not null;comment:优惠券名称"`
	Status          int            `gorm:"column:status;type:int;not null;comment:状态"` // 1: Unused, 2: Used, 3: Expired
	UserID          int64          `gorm:"column:user_id;type:bigint;not null;comment:用户编号"`
	UseOrderID      int64          `gorm:"column:use_order_id;type:bigint;comment:使用订单编号"`
	UseTime         *time.Time     `gorm:"column:use_time;comment:使用时间"`
	ValidStartTime  time.Time      `gorm:"column:valid_start_time;not null;comment:有效期开始时间"`
	ValidEndTime    time.Time      `gorm:"column:valid_end_time;not null;comment:有效期结束时间"`
	DiscountType    int            `gorm:"column:discount_type;type:int;not null;comment:优惠类型"`
	DiscountPrice   int            `gorm:"column:discount_price;type:int;comment:优惠金额"`
	DiscountPercent int            `gorm:"column:discount_percent;type:int;comment:折扣百分比"`
	DiscountLimit   int            `gorm:"column:discount_limit;type:int;comment:最多优惠金额"`
	UsePriceMin     int            `gorm:"column:use_price_min;type:int;not null;comment:最低消费金额"`
	Creator         string         `gorm:"column:creator;size:64;default:'';comment:创建者"`
	Updater         string         `gorm:"column:updater;size:64;default:'';comment:更新者"`
	CreatedAt       time.Time      `gorm:"column:create_time;autoCreateTime;comment:创建时间"`
	UpdatedAt       time.Time      `gorm:"column:update_time;autoUpdateTime;comment:更新时间"`
	DeletedAt       gorm.DeletedAt `gorm:"column:deleted;index;comment:删除时间"`
	Deleted         bool           `gorm:"column:deleted;type:tinyint(1);not null;default:0;comment:是否删除"`
}

func (PromotionCoupon) TableName() string {
	return "promotion_coupon"
}
