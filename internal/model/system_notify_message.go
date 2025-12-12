package model

import (
	"time"

)

// SystemNotifyMessage 站内信消息
type SystemNotifyMessage struct {
	ID               int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID           int64      `gorm:"column:user_id;not null;comment:用户编号" json:"userId"`
	UserType         int        `gorm:"column:user_type;not null;comment:用户类型" json:"userType"`
	TemplateID       int64      `gorm:"column:template_id;not null;comment:模板编号" json:"templateId"`
	TemplateCode     string     `gorm:"column:template_code;not null;comment:模板编码" json:"templateCode"`
	TemplateNickname string     `gorm:"column:template_nickname;comment:模版发送人名称" json:"templateNickname"`
	TemplateContent  string     `gorm:"column:template_content;not null;comment:模版内容" json:"templateContent"`
	TemplateType     int        `gorm:"column:template_type;not null;comment:模版类型" json:"templateType"`
	TemplateParams   string     `gorm:"column:template_params;not null;comment:模版参数" json:"templateParams"` // JSON map
	ReadStatus       bool       `gorm:"column:read_status;not null;default:0;comment:是否已读" json:"readStatus"`
	ReadTime         *time.Time `gorm:"column:read_time;comment:阅读时间" json:"readTime"`

	Creator   string         `gorm:"column:creator;size:64;default:'';comment:创建者" json:"creator"`
	Updater   string         `gorm:"column:updater;size:64;default:'';comment:更新者" json:"updater"`
	CreatedAt time.Time      `gorm:"column:create_time;autoCreateTime;comment:创建时间" json:"createTime"`
	UpdatedAt time.Time      `gorm:"column:update_time;autoUpdateTime;comment:更新时间" json:"updateTime"`
	Deleted   BitBool        `gorm:"column:deleted;softDelete:flag;default:0;comment:是否删除" json:"-"`
}

func (SystemNotifyMessage) TableName() string {
	return "system_notify_message"
}
