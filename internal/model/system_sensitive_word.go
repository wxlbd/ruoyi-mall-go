package model

import (
	"time"

	"gorm.io/gorm"
)

// SystemSensitiveWord 敏感词
type SystemSensitiveWord struct {
	ID          int64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name        string   `gorm:"column:name;not null;comment:敏感词" json:"name"`
	Tags        []string `gorm:"column:tags;serializer:json;comment:标签" json:"tags"`
	Status      int      `gorm:"column:status;not null;default:0;comment:状态" json:"status"`
	Description string   `gorm:"column:description;comment:描述" json:"description"`

	Creator   string         `gorm:"column:creator;size:64;default:'';comment:创建者" json:"creator"`
	Updater   string         `gorm:"column:updater;size:64;default:'';comment:更新者" json:"updater"`
	CreatedAt time.Time      `gorm:"column:create_time;autoCreateTime;comment:创建时间" json:"createTime"`
	UpdatedAt time.Time      `gorm:"column:update_time;autoUpdateTime;comment:更新时间" json:"updateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_time;index;comment:删除时间" json:"-"`
	Deleted   BitBool        `gorm:"column:deleted;softDelete:flag;default:0;comment:是否删除" json:"-"`
}

func (SystemSensitiveWord) TableName() string {
	return "system_sensitive_word"
}
