package member

import (
	"backend-go/internal/model"
	"time"

)

// MemberSignInConfig 签到规则
type MemberSignInConfig struct {
	ID         int64          `gorm:"primaryKey;autoIncrement;comment:规则自增主键" json:"id"`
	Day        int            `gorm:"comment:签到第 x 天" json:"day"`
	Point      int            `gorm:"comment:奖励积分" json:"point"`
	Experience int            `gorm:"comment:奖励经验" json:"experience"`
	Status     int            `gorm:"default:0;comment:状态" json:"status"` // 参见 CommonStatusEnum
	Creator    string         `gorm:"size:64;default:'';comment:创建者" json:"creator"`
	Updater    string         `gorm:"size:64;default:'';comment:更新者" json:"updater"`
	CreatedAt  time.Time      `gorm:"column:create_time;autoCreateTime;comment:创建时间" json:"createTime"`
	UpdatedAt  time.Time      `gorm:"column:update_time;autoUpdateTime;comment:更新时间" json:"updateTime"`
	Deleted    model.BitBool  `gorm:"column:deleted;softDelete:flag;default:0;comment:是否删除" json:"deleted"`
}

func (MemberSignInConfig) TableName() string {
	return "member_sign_in_config"
}

// MemberSignInRecord 签到记录
type MemberSignInRecord struct {
	ID         int64          `gorm:"primaryKey;autoIncrement;comment:编号" json:"id"`
	UserID     int64          `gorm:"column:user_id;comment:签到用户" json:"userId"`
	Day        int            `gorm:"comment:第几天签到" json:"day"`
	Point      int            `gorm:"comment:签到的积分" json:"point"`
	Experience int            `gorm:"comment:签到的经验" json:"experience"`
	Creator    string         `gorm:"size:64;default:'';comment:创建者" json:"creator"`
	Updater    string         `gorm:"size:64;default:'';comment:更新者" json:"updater"`
	CreatedAt  time.Time      `gorm:"column:create_time;autoCreateTime;comment:创建时间" json:"createTime"`
	UpdatedAt  time.Time      `gorm:"column:update_time;autoUpdateTime;comment:更新时间" json:"updateTime"`
	Deleted    model.BitBool  `gorm:"column:deleted;softDelete:flag;default:0;comment:是否删除" json:"deleted"`
}

func (MemberSignInRecord) TableName() string {
	return "member_sign_in_record"
}
