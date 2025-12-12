package pay

import (
	"backend-go/internal/model"
	"time"

)

// PayApp 支付应用 DO
type PayApp struct {
	ID                int64          `gorm:"primaryKey;autoIncrement;comment:应用编号" json:"id"`
	AppKey            string         `gorm:"size:64;not null;comment:应用标识" json:"appKey"`
	Name              string         `gorm:"size:64;not null;comment:应用名" json:"name"`
	Status            int            `gorm:"default:0;not null;comment:状态" json:"status"` // 参见 CommonStatusEnum
	Remark            string         `gorm:"size:255;default:'';comment:备注" json:"remark"`
	OrderNotifyURL    string         `gorm:"column:order_notify_url;size:1024;not null;comment:支付结果的回调地址" json:"orderNotifyUrl"`
	RefundNotifyURL   string         `gorm:"column:refund_notify_url;size:1024;not null;comment:退款结果的回调地址" json:"refundNotifyUrl"`
	TransferNotifyURL string         `gorm:"column:transfer_notify_url;size:1024;default:'';comment:转账结果的回调地址" json:"transferNotifyUrl"`
	Creator           string         `gorm:"size:64;default:'';comment:创建者" json:"creator"`
	Updater           string         `gorm:"size:64;default:'';comment:更新者" json:"updater"`
	CreatedAt         time.Time      `gorm:"column:create_time;autoCreateTime;comment:创建时间" json:"createTime"`
	UpdatedAt         time.Time      `gorm:"column:update_time;autoUpdateTime;comment:更新时间" json:"updateTime"`
	Deleted           model.BitBool  `gorm:"column:deleted;softDelete:flag;default:0;comment:是否删除" json:"deleted"`
}

func (PayApp) TableName() string {
	return "pay_app"
}
