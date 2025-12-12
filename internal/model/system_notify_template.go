package model

import (
	"time"

)

// SystemNotifyTemplate 站内信模版
type SystemNotifyTemplate struct {
	ID       int64  `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name     string `gorm:"column:name;not null;comment:模板名称" json:"name"`
	Code     string `gorm:"column:code;not null;comment:模板编码" json:"code"`
	Nickname string `gorm:"column:nickname;comment:发送人名称" json:"nickname"`
	Content  string `gorm:"column:content;not null;comment:模板内容" json:"content"`
	Type     int    `gorm:"column:type;not null;comment:类型" json:"type"` // 1: System, 2: ?
	Params   string `gorm:"column:params;comment:参数数组" json:"params"`    // JSON array
	Status   int    `gorm:"column:status;not null;default:0;comment:状态" json:"status"`
	Remark   string `gorm:"column:remark;comment:备注" json:"remark"`

	Creator   string         `gorm:"column:creator;size:64;default:'';comment:创建者" json:"creator"`
	Updater   string         `gorm:"column:updater;size:64;default:'';comment:更新者" json:"updater"`
	CreatedAt time.Time      `gorm:"column:create_time;autoCreateTime;comment:创建时间" json:"createTime"`
	UpdatedAt time.Time      `gorm:"column:update_time;autoUpdateTime;comment:更新时间" json:"updateTime"`
	Deleted   BitBool        `gorm:"column:deleted;softDelete:flag;default:0;comment:是否删除" json:"-"`
}

func (SystemNotifyTemplate) TableName() string {
	return "system_notify_template"
}
