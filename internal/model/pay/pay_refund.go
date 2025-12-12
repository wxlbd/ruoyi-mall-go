package pay

import (
	"time"

	"gorm.io/gorm"
)

// PayRefund 支付退款单
// TableName: pay_refund
type PayRefund struct {
	ID                int64      `gorm:"column:id;primaryKey;autoIncrement;comment:退款单编号" json:"id"`
	No                string     `gorm:"column:no;comment:外部退款号" json:"no"`
	AppID             int64      `gorm:"column:app_id;comment:应用编号" json:"appId"`
	ChannelID         int64      `gorm:"column:channel_id;comment:渠道编号" json:"channelId"`
	ChannelCode       string     `gorm:"column:channel_code;comment:渠道编码" json:"channelCode"`
	OrderID           int64      `gorm:"column:order_id;comment:订单编号" json:"orderId"`
	OrderNo           string     `gorm:"column:order_no;comment:支付订单编号" json:"orderNo"`
	UserID            int64      `gorm:"column:user_id;comment:用户编号" json:"userId"`
	UserType          int        `gorm:"column:user_type;comment:用户类型" json:"userType"`
	MerchantOrderId   string     `gorm:"column:merchant_order_id;comment:商户订单编号" json:"merchantOrderId"`
	MerchantRefundId  string     `gorm:"column:merchant_refund_id;comment:商户退款订单号" json:"merchantRefundId"`
	NotifyURL         string     `gorm:"column:notify_url;comment:异步通知地址" json:"notifyUrl"`
	Status            int        `gorm:"column:status;comment:退款状态" json:"status"`
	PayPrice          int        `gorm:"column:pay_price;comment:支付金额" json:"payPrice"`
	RefundPrice       int        `gorm:"column:refund_price;comment:退款金额" json:"refundPrice"`
	Reason            string     `gorm:"column:reason;comment:退款原因" json:"reason"`
	UserIP            string     `gorm:"column:user_ip;comment:用户 IP" json:"userIp"`
	ChannelOrderNo    string     `gorm:"column:channel_order_no;comment:渠道订单号" json:"channelOrderNo"`
	ChannelRefundNo   string     `gorm:"column:channel_refund_no;comment:渠道退款单号" json:"channelRefundNo"`
	SuccessTime       *time.Time `gorm:"column:success_time;comment:退款成功时间" json:"successTime"`
	ChannelErrorCode  string     `gorm:"column:channel_error_code;comment:调用渠道的错误码" json:"channelErrorCode"`
	ChannelErrorMsg   string     `gorm:"column:channel_error_msg;comment:调用渠道的错误提示" json:"channelErrorMsg"`
	ChannelNotifyData string     `gorm:"column:channel_notify_data;comment:支付渠道的同步/异步通知的内容" json:"channelNotifyData"`

	CreatedAt time.Time      `gorm:"column:create_time;autoCreateTime;comment:创建时间" json:"createTime"`
	UpdatedAt time.Time      `gorm:"column:update_time;autoUpdateTime;comment:更新时间" json:"updateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted;index;comment:删除时间" json:"deletedTime"`
	Deleted   bool           `gorm:"column:deleted;default:0;comment:是否删除" json:"deleted"`
	Creator   string         `gorm:"column:creator;default:'';comment:创建者" json:"creator"`
	Updater   string         `gorm:"column:updater;default:'';comment:更新者" json:"updater"`
}

func (PayRefund) TableName() string {
	return "pay_refund"
}
