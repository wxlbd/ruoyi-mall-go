package member

import (
	"time"

	"backend-go/internal/model"

)

type MemberUser struct {
	ID               int64      `gorm:"primaryKey;autoIncrement;comment:用户ID" json:"id"`
	Mobile           string     `gorm:"size:11;comment:手机" json:"mobile"`
	Password         string     `gorm:"size:100;default:'';comment:密码" json:"-"`
	Status           int32      `gorm:"default:0;comment:状态" json:"status"` // 参见 CommonStatusEnum
	RegisterIP       string     `gorm:"column:register_ip;size:32;default:'';comment:注册IP" json:"registerIp"`
	RegisterTerminal int32      `gorm:"column:register_terminal;default:0;comment:注册终端" json:"registerTerminal"` // 参见 TerminalEnum
	LoginIP          string     `gorm:"column:login_ip;size:32;default:'';comment:最后登录IP" json:"loginIp"`
	LoginDate        *time.Time `gorm:"column:login_date;comment:最后登录时间" json:"loginDate"`

	Nickname string     `gorm:"size:30;default:'';comment:用户昵称" json:"nickname"`
	Avatar   string     `gorm:"size:255;default:'';comment:头像" json:"avatar"`
	Name     string     `gorm:"size:30;default:'';comment:真实姓名" json:"name"`
	Sex      int32      `gorm:"default:0;comment:性别" json:"sex"` // 参见 SexEnum
	Birthday *time.Time `gorm:"comment:出生日期" json:"birthday"`
	AreaID   int32      `gorm:"column:area_id;comment:所在地" json:"areaId"`
	Mark     string     `gorm:"size:255;default:'';comment:备注" json:"mark"`

	Point      int32   `gorm:"default:0;comment:积分" json:"point"`
	TagIds     []int64 `gorm:"type:json;serializer:json;comment:标签编号数组" json:"tagIds"`
	LevelID    int64   `gorm:"column:level_id;comment:等级编号" json:"levelId"`
	Experience int32   `gorm:"default:0;comment:经验" json:"experience"`
	GroupID    int64   `gorm:"column:group_id;comment:分组编号" json:"groupId"`

	TenantID         int64          `gorm:"column:tenant_id;default:0;comment:租户编号" json:"tenantId"`
	Creator          string         `gorm:"size:64;default:'';comment:创建者" json:"creator"`
	Updater          string         `gorm:"size:64;default:'';comment:更新者" json:"updater"`
	CreatedAt        time.Time      `gorm:"column:create_time;autoCreateTime;comment:创建时间" json:"createTime"`
	UpdatedAt        time.Time      `gorm:"column:update_time;autoUpdateTime;comment:更新时间" json:"updateTime"`
	Deleted          model.BitBool  `gorm:"column:deleted;softDelete:flag;default:0;comment:是否删除" json:"deleted"`
	BrokerageEnabled model.BitBool  `gorm:"column:brokerage_enabled;default:0;comment:是否成为推广员" json:"brokerageEnabled"`
}

func (MemberUser) TableName() string {
	return "member_user"
}
