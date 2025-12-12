package promotion

import (
	"time"

	"gorm.io/gorm"
)

// PromotionPointProduct 积分商城商品
// 对应 Java: PointProductDO
type PromotionPointProduct struct {
	ID             int64          `gorm:"primaryKey;autoIncrement;comment:商品编号"`
	ActivityID     int64          `gorm:"column:activity_id;type:bigint;not null;comment:活动编号"`
	SpuID          int64          `gorm:"column:spu_id;type:bigint;not null;comment:商品SPU编号"`
	SkuID          int64          `gorm:"column:sku_id;type:bigint;not null;comment:商品SKU编号"`
	Count          int            `gorm:"column:count;type:int;not null;default:0;comment:可兑换次数"`
	Point          int            `gorm:"column:point;type:int;not null;default:0;comment:所需兑换积分"`
	Price          int            `gorm:"column:price;type:int;not null;default:0;comment:所需兑换金额"` // 单位：分
	Stock          int            `gorm:"column:stock;type:int;not null;default:0;comment:积分商城商品库存"`
	ActivityStatus int            `gorm:"column:activity_status;type:int;not null;comment:活动状态"`
	Creator        string         `gorm:"column:creator;size:64;default:'';comment:创建者"`
	Updater        string         `gorm:"column:updater;size:64;default:'';comment:更新者"`
	CreatedAt      time.Time      `gorm:"column:create_time;autoCreateTime;comment:创建时间"`
	UpdatedAt      time.Time      `gorm:"column:update_time;autoUpdateTime;comment:更新时间"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted;index;comment:删除时间"`
	Deleted        bool           `gorm:"column:deleted;type:tinyint(1);not null;default:0;comment:是否删除"`
}

func (PromotionPointProduct) TableName() string {
	return "promotion_point_product"
}
