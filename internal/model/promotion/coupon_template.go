package promotion

import (
	"time"

	"gorm.io/gorm"
)

// PromotionCouponTemplate 优惠券模板
type PromotionCouponTemplate struct {
	ID                 int64          `gorm:"primaryKey;autoIncrement;comment:模板编号"`
	Name               string         `gorm:"column:name;type:varchar(64);not null;comment:优惠券名称"`
	Status             int            `gorm:"column:status;type:int;not null;comment:状态"` // 0: Disable, 1: Enable
	TotalCount         int            `gorm:"column:total_count;type:int;not null;comment:发放总量"`
	TakeCount          int            `gorm:"column:take_count;type:int;not null;default:0;comment:已领取数量"`
	TakeLimitCount     int            `gorm:"column:take_limit_count;type:int;not null;comment:每人限领数量"`
	TakeType           int            `gorm:"column:take_type;type:int;not null;comment:领取方式"` // 1: Manually, 2: Register, 3: Admin
	UsePriceMin        int            `gorm:"column:use_price_min;type:int;not null;comment:最低消费金额"`
	ProductScope       int            `gorm:"column:product_scope;type:int;not null;comment:商品范围"` // 1: All, 2: Category, 3: Spu
	ProductScopeValues []int64        `gorm:"column:product_scope_values;type:json;serializer:json;comment:商品范围值"`
	ValidityType       int            `gorm:"column:validity_type;type:int;not null;comment:有效期类型"` // 1: Date, 2: Term
	ValidStartTime     *time.Time     `gorm:"column:valid_start_time;comment:固定日期-开始"`
	ValidEndTime       *time.Time     `gorm:"column:valid_end_time;comment:固定日期-结束"`
	FixedStartTerm     int            `gorm:"column:fixed_start_term;type:int;comment:领取后生效天数"`
	FixedEndTerm       int            `gorm:"column:fixed_end_term;type:int;comment:有效期天数"`
	DiscountType       int            `gorm:"column:discount_type;type:int;not null;comment:优惠类型"` // 1: Price, 2: Percent
	DiscountPrice      int            `gorm:"column:discount_price;type:int;comment:优惠金额"`
	DiscountPercent    int            `gorm:"column:discount_percent;type:int;comment:折扣百分比"`
	DiscountLimit      int            `gorm:"column:discount_limit;type:int;comment:最多优惠金额"`
	Creator            string         `gorm:"column:creator;size:64;default:'';comment:创建者"`
	Updater            string         `gorm:"column:updater;size:64;default:'';comment:更新者"`
	CreatedAt          time.Time      `gorm:"column:create_time;autoCreateTime;comment:创建时间"`
	UpdatedAt          time.Time      `gorm:"column:update_time;autoUpdateTime;comment:更新时间"`
	DeletedAt          gorm.DeletedAt `gorm:"column:deleted;index;comment:删除时间"`
	Deleted            bool           `gorm:"column:deleted;type:tinyint(1);not null;default:0;comment:是否删除"`
}

func (PromotionCouponTemplate) TableName() string {
	return "promotion_coupon_template"
}
