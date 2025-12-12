package trade

import (
	"backend-go/internal/model"
	"time"

)

// TradeDeliveryExpress 物流公司 DO
type TradeDeliveryExpress struct {
	ID        int64          `gorm:"primaryKey;autoIncrement;comment:编号" json:"id"`
	Code      string         `gorm:"size:64;not null;comment:物流编码" json:"code"`
	Name      string         `gorm:"size:64;not null;comment:物流名称" json:"name"`
	Logo      string         `gorm:"size:256;default:'';comment:物流Logo" json:"logo"`
	Sort      int            `gorm:"default:0;not null;comment:排序" json:"sort"`
	Status    int            `gorm:"default:0;not null;comment:状态" json:"status"`
	Creator   string         `gorm:"size:64;default:'';comment:创建者" json:"creator"`
	Updater   string         `gorm:"size:64;default:'';comment:更新者" json:"updater"`
	CreatedAt time.Time      `gorm:"column:create_time;autoCreateTime;comment:创建时间" json:"createTime"`
	UpdatedAt time.Time      `gorm:"column:update_time;autoUpdateTime;comment:更新时间" json:"updateTime"`
	Deleted   model.BitBool  `gorm:"column:deleted;softDelete:flag;default:0;comment:是否删除" json:"deleted"`
	TenantID  int64          `gorm:"column:tenant_id;default:0;comment:租户编号" json:"tenantId"`
}

func (TradeDeliveryExpress) TableName() string {
	return "trade_delivery_express"
}

// TradeDeliveryPickUpStore 自提门店 DO
type TradeDeliveryPickUpStore struct {
	ID            int64          `gorm:"primaryKey;autoIncrement;comment:编号" json:"id"`
	Name          string         `gorm:"size:64;not null;comment:门店名称" json:"name"`
	Introduction  string         `gorm:"size:256;default:'';comment:门店简介" json:"introduction"`
	Phone         string         `gorm:"size:11;not null;comment:联系电话" json:"phone"`
	AreaID        int            `gorm:"column:area_id;not null;comment:区域编号" json:"areaId"`
	DetailAddress string         `gorm:"size:256;not null;comment:详细地址" json:"detailAddress"`
	Logo          string         `gorm:"size:256;not null;comment:门店Logo" json:"logo"`
	Latitude      float64        `gorm:"type:decimal(10,6);comment:纬度" json:"latitude"`
	Longitude     float64        `gorm:"type:decimal(10,6);comment:经度" json:"longitude"`
	Status        int            `gorm:"default:0;not null;comment:状态" json:"status"`
	Sort          int            `gorm:"default:0;not null;comment:排序" json:"sort"`
	Creator       string         `gorm:"size:64;default:'';comment:创建者" json:"creator"`
	Updater       string         `gorm:"size:64;default:'';comment:更新者" json:"updater"`
	CreatedAt     time.Time      `gorm:"column:create_time;autoCreateTime;comment:创建时间" json:"createTime"`
	UpdatedAt     time.Time      `gorm:"column:update_time;autoUpdateTime;comment:更新时间" json:"updateTime"`
	Deleted       model.BitBool  `gorm:"column:deleted;softDelete:flag;default:0;comment:是否删除" json:"deleted"`
	TenantID      int64          `gorm:"column:tenant_id;default:0;comment:租户编号" json:"tenantId"`
}

func (TradeDeliveryPickUpStore) TableName() string {
	return "trade_delivery_pick_up_store"
}
