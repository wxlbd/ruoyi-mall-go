package member

import (
	"backend-go/internal/model"
	"time"

)

// MemberAddress 用户收件地址
type MemberAddress struct {
	ID            int64         `gorm:"primaryKey;autoIncrement;comment:收件地址编号" json:"id"`
	UserID        int64         `gorm:"column:user_id;not null;comment:用户编号" json:"userId"`
	Name          string        `gorm:"column:name;size:50;not null;comment:收件人名称" json:"name"`
	Mobile        string        `gorm:"column:mobile;size:20;not null;comment:手机号" json:"mobile"`
	AreaID        int64         `gorm:"column:area_id;not null;comment:地区编号" json:"areaId"`
	DetailAddress string        `gorm:"column:detail_address;size:255;not null;comment:收件详细地址" json:"detailAddress"`
	DefaultStatus model.BitBool `gorm:"column:default_status;not null;default:0;comment:是否默认" json:"defaultStatus"`

	Creator   string         `gorm:"size:64;default:'';comment:创建者" json:"creator"`
	Updater   string         `gorm:"size:64;default:'';comment:更新者" json:"updater"`
	CreatedAt time.Time      `gorm:"column:create_time;autoCreateTime;comment:创建时间" json:"createTime"`
	UpdatedAt time.Time      `gorm:"column:update_time;autoUpdateTime;comment:更新时间" json:"updateTime"`
	Deleted   model.BitBool  `gorm:"column:deleted;softDelete:flag;default:0;comment:是否删除" json:"deleted"`
}

func (MemberAddress) TableName() string {
	return "member_address"
}
