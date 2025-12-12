package promotion

import (
	"backend-go/internal/model"
	"time"

)

// PromotionCombinationActivity 拼团活动 DO
// Maps to promotion_combination_activity
type PromotionCombinationActivity struct {
	ID               int64     `gorm:"primaryKey;autoIncrement;column:id;comment:活动编号" json:"id"`
	Name             string    `gorm:"column:name;type:varchar(64);not null;comment:拼团名称" json:"name"`
	SpuID            int64     `gorm:"column:spu_id;not null;comment:商品SPU编号" json:"spuId"`
	TotalLimitCount  int       `gorm:"column:total_limit_count;not null;comment:总限购数量" json:"totalLimitCount"`
	SingleLimitCount int       `gorm:"column:single_limit_count;not null;comment:单次限购数量" json:"singleLimitCount"`
	StartTime        time.Time `gorm:"column:start_time;not null;comment:开始时间" json:"startTime"`
	EndTime          time.Time `gorm:"column:end_time;not null;comment:结束时间" json:"endTime"`
	UserSize         int       `gorm:"column:user_size;not null;comment:几人团" json:"userSize"`
	VirtualGroup     bool      `gorm:"column:virtual_group;not null;default:false;comment:虚拟成团" json:"virtualGroup"`
	Status           int       `gorm:"column:status;not null;comment:活动状态" json:"status"`
	LimitDuration    int       `gorm:"column:limit_duration;not null;comment:限制时长(小时)" json:"limitDuration"`

	Creator   string         `gorm:"size:64;default:'';comment:创建者" json:"creator"`
	Updater   string         `gorm:"size:64;default:'';comment:更新者" json:"updater"`
	CreatedAt time.Time      `gorm:"column:create_time;autoCreateTime;comment:创建时间" json:"createTime"`
	UpdatedAt time.Time      `gorm:"column:update_time;autoUpdateTime;comment:更新时间" json:"updateTime"`
	Deleted   model.BitBool  `gorm:"column:deleted;softDelete:flag;default:0;comment:是否删除" json:"deleted"`
	TenantID  int64          `gorm:"column:tenant_id;default:0;comment:租户编号" json:"tenantId"`
}

// TableName 表名
func (PromotionCombinationActivity) TableName() string {
	return "promotion_combination_activity"
}

// PromotionCombinationProduct 拼团商品 DO
// Maps to promotion_combination_product
type PromotionCombinationProduct struct {
	ID                int64     `gorm:"primaryKey;autoIncrement;column:id;comment:编号" json:"id"`
	ActivityID        int64     `gorm:"column:activity_id;not null;comment:活动编号" json:"activityId"`
	SpuID             int64     `gorm:"column:spu_id;not null;comment:商品SPU编号" json:"spuId"`
	SkuID             int64     `gorm:"column:sku_id;not null;comment:商品SKU编号" json:"skuId"`
	CombinationPrice  int       `gorm:"column:combination_price;not null;comment:拼团价格" json:"combinationPrice"`
	ActivityStatus    int       `gorm:"column:activity_status;not null;comment:活动状态" json:"activityStatus"`
	ActivityStartTime time.Time `gorm:"column:activity_start_time;comment:活动开始时间" json:"activityStartTime"`
	ActivityEndTime   time.Time `gorm:"column:activity_end_time;comment:活动结束时间" json:"activityEndTime"`

	Creator   string         `gorm:"size:64;default:'';comment:创建者" json:"creator"`
	Updater   string         `gorm:"size:64;default:'';comment:更新者" json:"updater"`
	CreatedAt time.Time      `gorm:"column:create_time;autoCreateTime;comment:创建时间" json:"createTime"`
	UpdatedAt time.Time      `gorm:"column:update_time;autoUpdateTime;comment:更新时间" json:"updateTime"`
	Deleted   model.BitBool  `gorm:"column:deleted;softDelete:flag;default:0;comment:是否删除" json:"deleted"`
	TenantID  int64          `gorm:"column:tenant_id;default:0;comment:租户编号" json:"tenantId"`
}

// TableName 表名
func (PromotionCombinationProduct) TableName() string {
	return "promotion_combination_product"
}

// PromotionCombinationRecord 拼团记录 DO
// Maps to promotion_combination_record
type PromotionCombinationRecord struct {
	ID               int64     `gorm:"primaryKey;autoIncrement;column:id;comment:编号" json:"id"`
	ActivityID       int64     `gorm:"column:activity_id;not null;comment:活动编号" json:"activityId"`
	CombinationPrice int       `gorm:"column:combination_price;not null;comment:拼团单价" json:"combinationPrice"`
	SpuID            int64     `gorm:"column:spu_id;not null;comment:SPU编号" json:"spuId"`
	SpuName          string    `gorm:"column:spu_name;type:varchar(255);comment:商品名字" json:"spuName"`
	PicUrl           string    `gorm:"column:pic_url;type:varchar(255);comment:商品图片" json:"picUrl"`
	SkuID            int64     `gorm:"column:sku_id;not null;comment:SKU编号" json:"skuId"`
	Count            int       `gorm:"column:count;not null;comment:购买数量" json:"count"`
	UserID           int64     `gorm:"column:user_id;not null;comment:用户编号" json:"userId"`
	Nickname         string    `gorm:"column:nickname;type:varchar(64);comment:用户昵称" json:"nickname"`
	Avatar           string    `gorm:"column:avatar;type:varchar(255);comment:用户头像" json:"avatar"`
	HeadID           int64     `gorm:"column:head_id;not null;comment:团长编号" json:"headId"`
	Status           int       `gorm:"column:status;not null;comment:拼团状态" json:"status"`
	OrderID          int64     `gorm:"column:order_id;not null;comment:订单编号" json:"orderId"`
	UserSize         int       `gorm:"column:user_size;not null;comment:成团人数" json:"userSize"`
	UserCount        int       `gorm:"column:user_count;not null;comment:已入团人数" json:"userCount"`
	VirtualGroup     bool      `gorm:"column:virtual_group;not null;default:false;comment:是否虚拟成团" json:"virtualGroup"`
	ExpireTime       time.Time `gorm:"column:expire_time;comment:过期时间" json:"expireTime"`
	StartTime        time.Time `gorm:"column:start_time;comment:开始时间" json:"startTime"`
	EndTime          time.Time `gorm:"column:end_time;comment:结束时间" json:"endTime"`

	Creator   string         `gorm:"size:64;default:'';comment:创建者" json:"creator"`
	Updater   string         `gorm:"size:64;default:'';comment:更新者" json:"updater"`
	CreatedAt time.Time      `gorm:"column:create_time;autoCreateTime;comment:创建时间" json:"createTime"`
	UpdatedAt time.Time      `gorm:"column:update_time;autoUpdateTime;comment:更新时间" json:"updateTime"`
	Deleted   model.BitBool  `gorm:"column:deleted;softDelete:flag;default:0;comment:是否删除" json:"deleted"`
	TenantID  int64          `gorm:"column:tenant_id;default:0;comment:租户编号" json:"tenantId"`
}

// TableName 表名
func (PromotionCombinationRecord) TableName() string {
	return "promotion_combination_record"
}
