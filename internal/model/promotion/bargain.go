package promotion

import (
	"backend-go/internal/model"
	"time"

)

// PromotionBargainActivity 砍价活动 DO
type PromotionBargainActivity struct {
	ID                int64     `gorm:"primaryKey;autoIncrement;comment:编号" json:"id"`
	Name              string    `gorm:"size:255;not null;comment:砍价活动名称" json:"name"`
	StartTime         time.Time `gorm:"column:start_time;not null;comment:活动开始时间" json:"startTime"`
	EndTime           time.Time `gorm:"column:end_time;not null;comment:活动结束时间" json:"endTime"`
	Status            int       `gorm:"default:0;not null;comment:状态" json:"status"`
	SpuID             int64     `gorm:"column:spu_id;not null;comment:商品 SPU 编号" json:"spuId"`
	SkuID             int64     `gorm:"column:sku_id;not null;comment:商品 SKU 编号" json:"skuId"`
	BargainFirstPrice int       `gorm:"column:bargain_first_price;comment:砍价起始价格，单位：分" json:"bargainFirstPrice"`
	BargainMinPrice   int       `gorm:"column:bargain_min_price;comment:砍价底价，单位：分" json:"bargainMinPrice"`
	Stock             int       `gorm:"default:0;comment:砍价库存" json:"stock"`
	TotalStock        int       `gorm:"default:0;comment:砍价总库存" json:"totalStock"`
	HelpMaxCount      int       `gorm:"column:help_max_count;comment:砍价人数" json:"helpMaxCount"`
	BargainCount      int       `gorm:"column:bargain_count;comment:帮砍次数" json:"bargainCount"`
	TotalLimitCount   int       `gorm:"column:total_limit_count;comment:总限购数量" json:"totalLimitCount"`
	RandomMinPrice    int       `gorm:"column:random_min_price;comment:用户每次砍价的最小金额" json:"randomMinPrice"`
	RandomMaxPrice    int       `gorm:"column:random_max_price;comment:用户每次砍价的最大金额" json:"randomMaxPrice"`
	Sort              int       `gorm:"default:0;comment:排序" json:"sort"`
	Remark            string    `gorm:"size:255;default:'';comment:备注" json:"remark"`

	Creator   string         `gorm:"size:64;default:'';comment:创建者" json:"creator"`
	Updater   string         `gorm:"size:64;default:'';comment:更新者" json:"updater"`
	CreatedAt time.Time      `gorm:"column:create_time;autoCreateTime;comment:创建时间" json:"createTime"`
	UpdatedAt time.Time      `gorm:"column:update_time;autoUpdateTime;comment:更新时间" json:"updateTime"`
	Deleted   model.BitBool  `gorm:"column:deleted;softDelete:flag;default:0;comment:是否删除" json:"deleted"`
	TenantID  int64          `gorm:"column:tenant_id;default:0;comment:租户编号" json:"tenantId"`
}

// TableName PromoBargainActivity's table name
func (*PromotionBargainActivity) TableName() string {
	return "promotion_bargain_activity"
}

// PromotionBargainRecord 砍价记录 DO
type PromotionBargainRecord struct {
	ID                int64     `gorm:"primaryKey;autoIncrement;comment:编号" json:"id"`
	UserID            int64     `gorm:"column:user_id;not null;comment:用户编号" json:"userId"`
	ActivityID        int64     `gorm:"column:activity_id;not null;comment:砍价活动编号" json:"activityId"`
	SpuID             int64     `gorm:"column:spu_id;not null;comment:商品 SPU 编号" json:"spuId"`
	SkuID             int64     `gorm:"column:sku_id;not null;comment:商品 SKU 编号" json:"skuId"`
	BargainFirstPrice int       `gorm:"column:bargain_first_price;comment:砍价起始价格，单位：分" json:"bargainFirstPrice"`
	BargainPrice      int       `gorm:"column:bargain_price;comment:当前砍价，单位：分" json:"bargainPrice"`
	Status            int       `gorm:"default:0;not null;comment:状态" json:"status"`
	EndTime           time.Time `gorm:"column:end_time;comment:结束时间" json:"endTime"`
	OrderID           int64     `gorm:"column:order_id;comment:订单编号" json:"orderId"`

	Creator   string         `gorm:"size:64;default:'';comment:创建者" json:"creator"`
	Updater   string         `gorm:"size:64;default:'';comment:更新者" json:"updater"`
	CreatedAt time.Time      `gorm:"column:create_time;autoCreateTime;comment:创建时间" json:"createTime"`
	UpdatedAt time.Time      `gorm:"column:update_time;autoUpdateTime;comment:更新时间" json:"updateTime"`
	Deleted   model.BitBool  `gorm:"column:deleted;softDelete:flag;default:0;comment:是否删除" json:"deleted"`
	TenantID  int64          `gorm:"column:tenant_id;default:0;comment:租户编号" json:"tenantId"`
}

// TableName PromoBargainRecord's table name
func (*PromotionBargainRecord) TableName() string {
	return "promotion_bargain_record"
}

// PromotionBargainHelp 砍价助力 DO
type PromotionBargainHelp struct {
	ID          int64 `gorm:"primaryKey;autoIncrement;comment:编号" json:"id"`
	ActivityID  int64 `gorm:"column:activity_id;not null;comment:砍价活动编号" json:"activityId"`
	RecordID    int64 `gorm:"column:record_id;not null;comment:砍价记录编号" json:"recordId"`
	UserID      int64 `gorm:"column:user_id;not null;comment:用户编号" json:"userId"`
	ReducePrice int   `gorm:"column:reduce_price;comment:减少价格，单位：分" json:"reducePrice"`

	Creator   string         `gorm:"size:64;default:'';comment:创建者" json:"creator"`
	Updater   string         `gorm:"size:64;default:'';comment:更新者" json:"updater"`
	CreatedAt time.Time      `gorm:"column:create_time;autoCreateTime;comment:创建时间" json:"createTime"`
	UpdatedAt time.Time      `gorm:"column:update_time;autoUpdateTime;comment:更新时间" json:"updateTime"`
	Deleted   model.BitBool  `gorm:"column:deleted;softDelete:flag;default:0;comment:是否删除" json:"deleted"`
	TenantID  int64          `gorm:"column:tenant_id;default:0;comment:租户编号" json:"tenantId"`
}

// TableName PromoBargainHelp's table name
func (*PromotionBargainHelp) TableName() string {
	return "promotion_bargain_help"
}
