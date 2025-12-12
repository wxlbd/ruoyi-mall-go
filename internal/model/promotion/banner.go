package promotion

import (
	"time"

	"gorm.io/gorm"
)

// PromotionBanner 首页轮播图
// Table: promotion_banner
type PromotionBanner struct {
	ID        int64          `gorm:"primaryKey;autoIncrement;comment:编号" json:"id"`
	Title     string         `gorm:"column:title;type:varchar(64);not null;comment:标题" json:"title"`
	PicURL    string         `gorm:"column:pic_url;type:varchar(255);not null;comment:图片地址" json:"picUrl"`
	Url       string         `gorm:"column:url;type:varchar(255);not null;comment:跳转地址" json:"url"`
	Status    int            `gorm:"column:status;type:tinyint;not null;default:0;comment:状态" json:"status"` // 0: 开启, 1: 关闭
	Sort      int            `gorm:"column:sort;type:int;not null;default:0;comment:排序" json:"sort"`
	Position  int            `gorm:"column:position;type:tinyint;not null;default:1;comment:位置" json:"position"` // 1: 首页
	Memo      string         `gorm:"column:memo;type:varchar(255);comment:备注" json:"memo"`
	Creator   string         `gorm:"column:creator;size:64;default:'';comment:创建者"`
	Updater   string         `gorm:"column:updater;size:64;default:'';comment:更新者"`
	CreatedAt time.Time      `gorm:"column:create_time;autoCreateTime;comment:创建时间"`
	UpdatedAt time.Time      `gorm:"column:update_time;autoUpdateTime;comment:更新时间"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted;index;comment:删除时间"`
	Deleted   bool           `gorm:"column:deleted;type:tinyint(1);not null;default:0;comment:是否删除"`
}

func (PromotionBanner) TableName() string {
	return "promotion_banner"
}
