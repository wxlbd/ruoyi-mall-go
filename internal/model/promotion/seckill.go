package promotion

import (
	"backend-go/internal/model"
	"time"

)

// PromotionSeckillActivity 秒杀活动 DO
type PromotionSeckillActivity struct {
	ID               int64          `gorm:"primaryKey;autoIncrement;comment:编号" json:"id"`
	SpuID            int64          `gorm:"column:spu_id;not null;comment:秒杀活动商品" json:"spuId"`
	Name             string         `gorm:"size:255;not null;comment:秒杀活动名称" json:"name"`
	Status           int            `gorm:"default:0;not null;comment:状态" json:"status"`
	Remark           string         `gorm:"size:255;default:'';comment:备注" json:"remark"`
	StartTime        time.Time      `gorm:"column:start_time;not null;comment:活动开始时间" json:"startTime"`
	EndTime          time.Time      `gorm:"column:end_time;not null;comment:活动结束时间" json:"endTime"`
	Sort             int            `gorm:"default:0;not null;comment:排序" json:"sort"`
	ConfigIds        []int64        `gorm:"serializer:json;type:varchar(255);comment:秒杀时段 id" json:"configIds"`
	TotalLimitCount  int            `gorm:"default:0;comment:总限购数量" json:"totalLimitCount"`
	SingleLimitCount int            `gorm:"default:0;comment:单次限购数量" json:"singleLimitCount"`
	Stock            int            `gorm:"default:0;comment:秒杀库存(剩余库存秒杀时扣减)" json:"stock"`
	TotalStock       int            `gorm:"default:0;comment:秒杀总库存" json:"totalStock"`
	Creator          string         `gorm:"size:64;default:'';comment:创建者" json:"creator"`
	Updater          string         `gorm:"size:64;default:'';comment:更新者" json:"updater"`
	CreatedAt        time.Time      `gorm:"column:create_time;autoCreateTime;comment:创建时间" json:"createTime"`
	UpdatedAt        time.Time      `gorm:"column:update_time;autoUpdateTime;comment:更新时间" json:"updateTime"`
	Deleted          model.BitBool  `gorm:"column:deleted;softDelete:flag;default:0;comment:是否删除" json:"deleted"`
	TenantID         int64          `gorm:"column:tenant_id;default:0;comment:租户编号" json:"tenantId"`
}

func (PromotionSeckillActivity) TableName() string {
	return "promotion_seckill_activity"
}

// PromotionSeckillProduct 秒杀参与商品 DO
type PromotionSeckillProduct struct {
	ID                int64          `gorm:"primaryKey;autoIncrement;comment:编号" json:"id"`
	ActivityID        int64          `gorm:"column:activity_id;not null;comment:秒杀活动 id" json:"activityId"`
	ConfigIds         []int64        `gorm:"serializer:json;type:varchar(255);comment:秒杀时段 id" json:"configIds"`
	SpuID             int64          `gorm:"column:spu_id;not null;comment:商品 SPU 编号" json:"spuId"`
	SkuID             int64          `gorm:"column:sku_id;not null;comment:商品 SKU 编号" json:"skuId"`
	SeckillPrice      int            `gorm:"default:0;not null;comment:秒杀金额，单位：分" json:"seckillPrice"`
	Stock             int            `gorm:"default:0;not null;comment:秒杀库存" json:"stock"`
	ActivityStatus    int            `gorm:"default:0;not null;comment:秒杀商品状态" json:"activityStatus"`
	ActivityStartTime time.Time      `gorm:"column:activity_start_time;not null;comment:活动开始时间点" json:"activityStartTime"`
	ActivityEndTime   time.Time      `gorm:"column:activity_end_time;not null;comment:活动结束时间点" json:"activityEndTime"`
	Creator           string         `gorm:"size:64;default:'';comment:创建者" json:"creator"`
	Updater           string         `gorm:"size:64;default:'';comment:更新者" json:"updater"`
	CreatedAt         time.Time      `gorm:"column:create_time;autoCreateTime;comment:创建时间" json:"createTime"`
	UpdatedAt         time.Time      `gorm:"column:update_time;autoUpdateTime;comment:更新时间" json:"updateTime"`
	Deleted           model.BitBool  `gorm:"column:deleted;softDelete:flag;default:0;comment:是否删除" json:"deleted"`
	TenantID          int64          `gorm:"column:tenant_id;default:0;comment:租户编号" json:"tenantId"`
}

func (PromotionSeckillProduct) TableName() string {
	return "promotion_seckill_product"
}

// PromotionSeckillConfig 秒杀时段 DO
type PromotionSeckillConfig struct {
	ID            int64          `gorm:"primaryKey;autoIncrement;comment:编号" json:"id"`
	Name          string         `gorm:"size:255;not null;comment:秒杀时段名称" json:"name"`
	StartTime     string         `gorm:"size:10;not null;comment:开始时间点" json:"startTime"`
	EndTime       string         `gorm:"size:10;not null;comment:结束时间点" json:"endTime"`
	SliderPicUrls []string       `gorm:"serializer:json;type:varchar(2000);comment:秒杀轮播图" json:"sliderPicUrls"`
	Status        int            `gorm:"default:0;not null;comment:状态" json:"status"`
	Creator       string         `gorm:"size:64;default:'';comment:创建者" json:"creator"`
	Updater       string         `gorm:"size:64;default:'';comment:更新者" json:"updater"`
	CreatedAt     time.Time      `gorm:"column:create_time;autoCreateTime;comment:创建时间" json:"createTime"`
	UpdatedAt     time.Time      `gorm:"column:update_time;autoUpdateTime;comment:更新时间" json:"updateTime"`
	Deleted       model.BitBool  `gorm:"column:deleted;softDelete:flag;default:0;comment:是否删除" json:"deleted"`
	TenantID      int64          `gorm:"column:tenant_id;default:0;comment:租户编号" json:"tenantId"`
}

func (PromotionSeckillConfig) TableName() string {
	return "promotion_seckill_config"
}
