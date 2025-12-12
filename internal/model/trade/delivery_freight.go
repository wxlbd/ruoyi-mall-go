package trade

import (
	"backend-go/internal/model"
	"time"

)

// TradeDeliveryFreightTemplate 运费模板 DO
type TradeDeliveryFreightTemplate struct {
	ID         int64          `gorm:"primaryKey;autoIncrement;comment:编号" json:"id"`
	Name       string         `gorm:"size:64;not null;comment:模板名称" json:"name"`
	Type       int            `gorm:"default:0;not null;comment:类型" json:"type"`         // 1-买家承担运费 2-卖家包邮
	ChargeMode int            `gorm:"default:0;not null;comment:计费方式" json:"chargeMode"` // 1-按件 2-按重量 3-按体积
	Sort       int            `gorm:"default:0;not null;comment:排序" json:"sort"`
	Status     int            `gorm:"default:0;not null;comment:状态" json:"status"`
	Remark     string         `gorm:"size:255;default:'';comment:备注" json:"remark"`
	Creator    string         `gorm:"size:64;default:'';comment:创建者" json:"creator"`
	Updater    string         `gorm:"size:64;default:'';comment:更新者" json:"updater"`
	CreatedAt  time.Time      `gorm:"column:create_time;autoCreateTime;comment:创建时间" json:"createTime"`
	UpdatedAt  time.Time      `gorm:"column:update_time;autoUpdateTime;comment:更新时间" json:"updateTime"`
	Deleted    model.BitBool  `gorm:"column:deleted;softDelete:flag;default:0;comment:是否删除" json:"deleted"`
	TenantID   int64          `gorm:"column:tenant_id;default:0;comment:租户编号" json:"tenantId"`
}

func (TradeDeliveryFreightTemplate) TableName() string {
	return "trade_delivery_freight_template"
}

// TradeDeliveryFreightTemplateCharge 运费模板计费规则 DO
type TradeDeliveryFreightTemplateCharge struct {
	ID         int64          `gorm:"primaryKey;autoIncrement;comment:编号" json:"id"`
	TemplateID int64          `gorm:"column:template_id;not null;comment:模板编号" json:"templateId"`
	AreaIDs    string         `gorm:"column:area_ids;type:text;not null;comment:区域编号列表" json:"areaIds"` // 逗号分隔
	StartCount float64        `gorm:"column:start_count;not null;comment:首件/首重/首体积" json:"startCount"`
	StartPrice int            `gorm:"column:start_price;not null;comment:首费(分)" json:"startPrice"`
	ExtraCount float64        `gorm:"column:extra_count;not null;comment:续件/续重/续体积" json:"extraCount"`
	ExtraPrice int            `gorm:"column:extra_price;not null;comment:续费(分)" json:"extraPrice"`
	Creator    string         `gorm:"size:64;default:'';comment:创建者" json:"creator"`
	Updater    string         `gorm:"size:64;default:'';comment:更新者" json:"updater"`
	CreatedAt  time.Time      `gorm:"column:create_time;autoCreateTime;comment:创建时间" json:"createTime"`
	UpdatedAt  time.Time      `gorm:"column:update_time;autoUpdateTime;comment:更新时间" json:"updateTime"`
	Deleted    model.BitBool  `gorm:"column:deleted;softDelete:flag;default:0;comment:是否删除" json:"deleted"`
}

func (TradeDeliveryFreightTemplateCharge) TableName() string {
	return "trade_delivery_freight_template_charge"
}

// TradeDeliveryFreightTemplateFree 运费模板包邮规则 DO
type TradeDeliveryFreightTemplateFree struct {
	ID         int64          `gorm:"primaryKey;autoIncrement;comment:编号" json:"id"`
	TemplateID int64          `gorm:"column:template_id;not null;comment:模板编号" json:"templateId"`
	AreaIDs    string         `gorm:"column:area_ids;type:text;not null;comment:区域编号列表" json:"areaIds"` // 逗号分隔
	FreePrice  int            `gorm:"column:free_price;not null;comment:包邮金额(分)" json:"freePrice"`
	FreeCount  float64        `gorm:"column:free_count;not null;comment:包邮件数/重量/体积" json:"freeCount"`
	Creator    string         `gorm:"size:64;default:'';comment:创建者" json:"creator"`
	Updater    string         `gorm:"size:64;default:'';comment:更新者" json:"updater"`
	CreatedAt  time.Time      `gorm:"column:create_time;autoCreateTime;comment:创建时间" json:"createTime"`
	UpdatedAt  time.Time      `gorm:"column:update_time;autoUpdateTime;comment:更新时间" json:"updateTime"`
	Deleted    model.BitBool  `gorm:"column:deleted;softDelete:flag;default:0;comment:是否删除" json:"deleted"`
}

func (TradeDeliveryFreightTemplateFree) TableName() string {
	return "trade_delivery_freight_template_free"
}
