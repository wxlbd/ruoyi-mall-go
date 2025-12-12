package pay

import (
	"backend-go/internal/model"
	"time"

)

// PayOrderExtension 支付订单拓展
type PayOrderExtension struct {
	ID                int64          `gorm:"column:id;primaryKey;autoIncrement"`
	No                string         `gorm:"column:no"`
	OrderID           int64          `gorm:"column:order_id"`
	ChannelID         int64          `gorm:"column:channel_id"`
	ChannelCode       string         `gorm:"column:channel_code"`
	UserIP            string         `gorm:"column:user_ip"`
	Status            int            `gorm:"column:status"`
	ChannelExtras     string         `gorm:"column:channel_extras"` // JSON String
	ChannelErrorCode  string         `gorm:"column:channel_error_code"`
	ChannelErrorMsg   string         `gorm:"column:channel_error_msg"`
	ChannelNotifyData string         `gorm:"column:channel_notify_data"`
	Creator           string         `gorm:"size:64;default:'';comment:创建者" json:"creator"`
	Updater           string         `gorm:"size:64;default:'';comment:更新者" json:"updater"`
	CreatedAt         time.Time      `gorm:"column:create_time;autoCreateTime;comment:创建时间" json:"createTime"`
	UpdatedAt         time.Time      `gorm:"column:update_time;autoUpdateTime;comment:更新时间" json:"updateTime"`
	Deleted           model.BitBool  `gorm:"column:deleted;softDelete:flag;default:0;comment:是否删除" json:"deleted"`
}

func (PayOrderExtension) TableName() string {
	return "pay_order_extension"
}
