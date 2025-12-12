package pay

import (
	"backend-go/internal/model"
	"time"

)

// PayOrder 支付订单
type PayOrder struct {
	ID              int64          `gorm:"column:id;primaryKey;autoIncrement"`
	AppID           int64          `gorm:"column:app_id"`
	ChannelID       int64          `gorm:"column:channel_id"`
	ChannelCode     string         `gorm:"column:channel_code"`
	UserID          int64          `gorm:"column:user_id"`
	UserType        int            `gorm:"column:user_type"`
	MerchantOrderId string         `gorm:"column:merchant_order_id"`
	Subject         string         `gorm:"column:subject"`
	Body            string         `gorm:"column:body"`
	NotifyURL       string         `gorm:"column:notify_url"`
	Price           int            `gorm:"column:price"` // Unit: fen
	ChannelFeeRate  float64        `gorm:"column:channel_fee_rate"`
	ChannelFeePrice int            `gorm:"column:channel_fee_price"` // Unit: fen
	Status          int            `gorm:"column:status"`
	UserIP          string         `gorm:"column:user_ip"`
	ExpireTime      time.Time      `gorm:"column:expire_time"`
	SuccessTime     *time.Time     `gorm:"column:success_time"`
	ExtensionID     int64          `gorm:"column:extension_id"`
	No              string         `gorm:"column:no"`
	RefundPrice     int            `gorm:"column:refund_price"` // Unit: fen
	ChannelUserID   string         `gorm:"column:channel_user_id"`
	ChannelOrderNo  string         `gorm:"column:channel_order_no"`
	Creator         string         `gorm:"size:64;default:'';comment:创建者" json:"creator"`
	Updater         string         `gorm:"size:64;default:'';comment:更新者" json:"updater"`
	CreatedAt       time.Time      `gorm:"column:create_time;autoCreateTime;comment:创建时间" json:"createTime"`
	UpdatedAt       time.Time      `gorm:"column:update_time;autoUpdateTime;comment:更新时间" json:"updateTime"`
	Deleted         model.BitBool  `gorm:"column:deleted;softDelete:flag;default:0;comment:是否删除" json:"deleted"`
}

func (PayOrder) TableName() string {
	return "pay_order"
}
