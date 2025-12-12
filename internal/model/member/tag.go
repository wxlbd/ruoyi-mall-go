package member

import (
	"time"

	"gorm.io/gorm"
)

// MemberTag 会员标签
// Table: member_tag
type MemberTag struct {
	ID        int64          `gorm:"primaryKey;autoIncrement;comment:编号" json:"id"`
	Name      string         `gorm:"column:name;type:varchar(30);not null;default:'';comment:标签名称" json:"name"`
	Status    int            `gorm:"column:status;type:int;not null;default:0;comment:状态" json:"status"`
	Remark    string         `gorm:"column:remark;type:varchar(500);default:'';comment:备注" json:"remark"`
	Creator   string         `gorm:"column:creator;size:64;default:'';comment:创建者"`
	Updater   string         `gorm:"column:updater;size:64;default:'';comment:更新者"`
	CreatedAt time.Time      `gorm:"column:create_time;autoCreateTime;comment:创建时间"`
	UpdatedAt time.Time      `gorm:"column:update_time;autoUpdateTime;comment:更新时间"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted;index;comment:删除时间"`
	Deleted   bool           `gorm:"column:deleted;type:tinyint(1);not null;default:0;comment:是否删除"`
}

func (MemberTag) TableName() string {
	return "member_tag"
}
