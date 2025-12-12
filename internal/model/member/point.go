package member

import (
	"time"

	"backend-go/internal/model"

	"gorm.io/gorm"
)

// MemberPointRecord 用户积分记录
// Table: member_point_record
type MemberPointRecord struct {
	ID          int64          `gorm:"primaryKey;autoIncrement;comment:自增主键" json:"id"`
	UserID      int64          `gorm:"column:user_id;not null;comment:用户编号" json:"userId"`
	BizID       string         `gorm:"column:biz_id;size:64;comment:业务编码" json:"bizId"`
	BizType     int            `gorm:"column:biz_type;not null;comment:业务类型" json:"bizType"` // MemberPointBizTypeEnum
	Title       string         `gorm:"column:title;size:64;not null;comment:积分标题" json:"title"`
	Description string         `gorm:"column:description;size:255;default:'';comment:积分描述" json:"description"`
	Point       int            `gorm:"column:point;not null;comment:变动积分" json:"point"`              // 1、正数表示获得积分 2、负数表示消耗积分
	TotalPoint  int            `gorm:"column:total_point;not null;comment:变动后的积分" json:"totalPoint"` // 变动后的积分
	Creator     string         `gorm:"column:creator;size:64;default:'';comment:创建者"`
	Updater     string         `gorm:"column:updater;size:64;default:'';comment:更新者"`
	CreatedAt   time.Time      `gorm:"column:create_time;autoCreateTime;comment:创建时间"`
	UpdatedAt   time.Time      `gorm:"column:update_time;autoUpdateTime;comment:更新时间"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted;index;comment:删除时间"`
	Deleted     model.BitBool  `gorm:"column:deleted;type:tinyint(1);not null;default:0;comment:是否删除"`
}

func (MemberPointRecord) TableName() string {
	return "member_point_record"
}
