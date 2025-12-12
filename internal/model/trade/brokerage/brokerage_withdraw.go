package brokerage

import (
	"time"

	"gorm.io/gorm"
)

// BrokerageWithdraw 佣金提现
type BrokerageWithdraw struct {
	ID                  int64          `gorm:"primaryKey;autoIncrement;comment:编号"`
	UserID              int64          `gorm:"column:user_id;not null;comment:用户编号"`
	Price               int            `gorm:"column:price;not null;comment:提现金额"`
	FeePrice            int            `gorm:"column:fee_price;default:0;comment:提现手续费"`
	TotalPrice          int            `gorm:"column:total_price;not null;comment:当前总佣金"`
	Type                int            `gorm:"column:type;not null;comment:提现类型"`
	UserName            string         `gorm:"column:user_name;size:64;default:'';comment:提现姓名"`
	UserAccount         string         `gorm:"column:user_account;size:64;default:'';comment:提现账号"`
	QrCodeURL           string         `gorm:"column:qr_code_url;size:255;default:'';comment:收款码"`
	BankName            string         `gorm:"column:bank_name;size:100;default:'';comment:银行名称"`
	BankAddress         string         `gorm:"column:bank_address;size:200;default:'';comment:开户地址"`
	Status              int            `gorm:"column:status;not null;comment:状态"`
	AuditReason         string         `gorm:"column:audit_reason;size:255;default:'';comment:审核驳回原因"`
	AuditTime           *time.Time     `gorm:"column:audit_time;comment:审核时间"`
	Remark              string         `gorm:"column:remark;size:255;default:'';comment:备注"`
	PayTransferID       int64          `gorm:"column:pay_transfer_id;default:0;comment:转账单编号"`
	TransferChannelCode string         `gorm:"column:transfer_channel_code;size:16;default:'';comment:转账渠道"`
	TransferTime        *time.Time     `gorm:"column:transfer_time;comment:转账成功时间"`
	TransferErrorMsg    string         `gorm:"column:transfer_error_msg;size:255;default:'';comment:转账错误提示"`
	Creator             string         `gorm:"column:creator;size:64;default:'';comment:创建者"`
	Updater             string         `gorm:"column:updater;size:64;default:'';comment:更新者"`
	CreatedAt           time.Time      `gorm:"column:create_time;autoCreateTime;comment:创建时间"`
	UpdatedAt           time.Time      `gorm:"column:update_time;autoUpdateTime;comment:更新时间"`
	DeletedAt           gorm.DeletedAt `gorm:"column:deleted;index;comment:删除时间"`
	Deleted             bool           `gorm:"column:deleted;type:tinyint(1);not null;default:0;comment:是否删除"`
}

func (BrokerageWithdraw) TableName() string {
	return "trade_brokerage_withdraw"
}
