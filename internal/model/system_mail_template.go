package model

import (
	"time"

)

// SystemMailTemplate 邮件模版
type SystemMailTemplate struct {
	ID        int64          `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name      string         `gorm:"column:name;not null;comment:模板名称" json:"name"`
	Code      string         `gorm:"column:code;not null;comment:模板编码" json:"code"`
	AccountID int64          `gorm:"column:account_id;not null;comment:发送的邮箱账号编号" json:"accountId"`
	Nickname  string         `gorm:"column:nickname;comment:发送人名称" json:"nickname"`
	Title     string         `gorm:"column:title;not null;comment:模板标题" json:"title"`
	Content   string         `gorm:"column:content;not null;comment:模板内容" json:"content"`
	Params    string         `gorm:"column:params;comment:参数数组" json:"params"` // JSON array of param names
	Status    int            `gorm:"column:status;not null;default:0;comment:状态" json:"status"`
	Remark    string         `gorm:"column:remark;comment:备注" json:"remark"`
	Creator   string         `gorm:"column:creator;size:64;default:'';comment:创建者" json:"creator"`
	Updater   string         `gorm:"column:updater;size:64;default:'';comment:更新者" json:"updater"`
	CreatedAt time.Time      `gorm:"column:create_time;autoCreateTime;comment:创建时间" json:"createTime"`
	UpdatedAt time.Time      `gorm:"column:update_time;autoUpdateTime;comment:更新时间" json:"updateTime"`
	Deleted   BitBool        `gorm:"column:deleted;softDelete:flag;default:0;comment:是否删除" json:"-"`
}

func (SystemMailTemplate) TableName() string {
	return "system_mail_template"
}
