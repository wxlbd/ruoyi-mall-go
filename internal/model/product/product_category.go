package product

import (
	"backend-go/internal/model"
	"time"

	"gorm.io/gorm"
)

// ProductCategory 商品分类
type ProductCategory struct {
	ID          int64  `gorm:"primaryKey;autoIncrement;comment:分类编号" json:"id"`
	ParentID    int64  `gorm:"column:parent_id;not null;default:0;comment:父分类编号" json:"parentId"`
	Name        string `gorm:"size:255;not null;comment:分类名称" json:"name"`
	PicURL      string `gorm:"column:pic_url;size:255;default:'';comment:移动端分类图" json:"picUrl"`
	Sort        int32  `gorm:"default:0;comment:分类排序" json:"sort"`
	Status      int32  `gorm:"default:0;comment:开启状态" json:"status"` // 参见 CommonStatusEnum
	Description string `gorm:"size:512;default:'';comment:分类描述" json:"description"`

	Creator   string         `gorm:"size:64;default:'';comment:创建者" json:"creator"`
	Updater   string         `gorm:"size:64;default:'';comment:更新者" json:"updater"`
	CreatedAt time.Time      `gorm:"column:create_time;autoCreateTime;comment:创建时间" json:"createTime"`
	UpdatedAt time.Time      `gorm:"column:update_time;autoUpdateTime;comment:更新时间" json:"updateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_time;index;comment:删除时间" json:"-"`
	Deleted   model.BitBool  `gorm:"column:deleted;softDelete:flag;default:0;comment:是否删除" json:"deleted"`
}

func (ProductCategory) TableName() string {
	return "product_category"
}
