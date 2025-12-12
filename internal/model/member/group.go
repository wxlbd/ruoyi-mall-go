package member

import (
	"time"

	"gorm.io/gorm"
)

// MemberGroup 会员分组
type MemberGroup struct {
	ID        int64          `gorm:"primaryKey;autoIncrement;comment:编号" json:"id"`
	Name      string         `gorm:"column:name;comment:名称" json:"name"`     // 名称
	Remark    string         `gorm:"column:remark;comment:备注" json:"remark"` // 备注
	Status    int            `gorm:"column:status;comment:状态" json:"status"` // 状态
	Creator   string         `gorm:"column:creator;size:64;default:'';comment:创建者"`
	Updater   string         `gorm:"column:updater;size:64;default:'';comment:更新者"`
	CreatedAt time.Time      `gorm:"column:create_time;autoCreateTime;comment:创建时间"`
	UpdatedAt time.Time      `gorm:"column:update_time;autoUpdateTime;comment:更新时间"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted;index;comment:删除时间"`
	Deleted   bool           `gorm:"column:deleted;type:tinyint(1);not null;default:0;comment:是否删除"`
}

// TableName 表名
func (MemberGroup) TableName() string {
	return "member_group"
}
