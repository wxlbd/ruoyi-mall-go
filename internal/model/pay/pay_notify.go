package pay

import (
	"time"

	"gorm.io/gorm"
)

// PayNotifyTask 支付通知任务
// TableName: pay_notify_task
type PayNotifyTask struct {
	ID                 int64      `gorm:"column:id;primaryKey;autoIncrement;comment:任务编号" json:"id"`
	AppID              int64      `gorm:"column:app_id;comment:应用编号" json:"appId"`
	Type               int        `gorm:"column:type;comment:通知类型" json:"type"`
	DataID             int64      `gorm:"column:data_id;comment:数据编号" json:"dataId"`
	MerchantOrderId    string     `gorm:"column:merchant_order_id;comment:商户订单编号" json:"merchantOrderId"`
	MerchantRefundId   string     `gorm:"column:merchant_refund_id;comment:商户退款编号" json:"merchantRefundId"`
	MerchantTransferId string     `gorm:"column:merchant_transfer_id;comment:商户转账编号" json:"merchantTransferId"`
	Status             int        `gorm:"column:status;comment:通知状态" json:"status"`
	NextNotifyTime     *time.Time `gorm:"column:next_notify_time;comment:下一次通知时间" json:"nextNotifyTime"`
	LastExecuteTime    *time.Time `gorm:"column:last_execute_time;comment:最后一次执行时间" json:"lastExecuteTime"`
	NotifyTimes        int        `gorm:"column:notify_times;comment:当前通知次数" json:"notifyTimes"`
	MaxNotifyTimes     int        `gorm:"column:max_notify_times;comment:最大可通知次数" json:"maxNotifyTimes"`
	NotifyURL          string     `gorm:"column:notify_url;comment:通知地址" json:"notifyUrl"`

	CreatedAt time.Time      `gorm:"column:create_time;autoCreateTime;comment:创建时间" json:"createTime"`
	UpdatedAt time.Time      `gorm:"column:update_time;autoUpdateTime;comment:更新时间" json:"updateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted;index;comment:删除时间" json:"deletedTime"`
	Deleted   bool           `gorm:"column:deleted;default:0;comment:是否删除" json:"deleted"`
	Creator   string         `gorm:"column:creator;default:'';comment:创建者" json:"creator"`
	Updater   string         `gorm:"column:updater;default:'';comment:更新者" json:"updater"`
}

func (PayNotifyTask) TableName() string {
	return "pay_notify_task"
}

// PayNotifyLog 支付通知日志
// TableName: pay_notify_log
type PayNotifyLog struct {
	ID          int64  `gorm:"column:id;primaryKey;autoIncrement;comment:日志编号" json:"id"`
	TaskID      int64  `gorm:"column:task_id;comment:任务编号" json:"taskId"`
	NotifyTimes int    `gorm:"column:notify_times;comment:第几次被通知" json:"notifyTimes"`
	Response    string `gorm:"column:response;comment:HTTP 响应结果" json:"response"`
	Status      int    `gorm:"column:status;comment:支付通知状态" json:"status"`

	CreatedAt time.Time      `gorm:"column:create_time;autoCreateTime;comment:创建时间" json:"createTime"`
	UpdatedAt time.Time      `gorm:"column:update_time;autoUpdateTime;comment:更新时间" json:"updateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted;index;comment:删除时间" json:"deletedTime"`
	Deleted   bool           `gorm:"column:deleted;default:0;comment:是否删除" json:"deleted"`
	Creator   string         `gorm:"column:creator;default:'';comment:创建者" json:"creator"`
	Updater   string         `gorm:"column:updater;default:'';comment:更新者" json:"updater"`
}

func (PayNotifyLog) TableName() string {
	return "pay_notify_log"
}
