package product

import (
	"backend-go/internal/model"
	"time"

)

// ProductBrand 商品品牌
type ProductBrand struct {
	ID          int64          `gorm:"primaryKey;autoIncrement;comment:主键" json:"id"`
	Name        string         `gorm:"size:255;not null;comment:品牌名称" json:"name"`
	PicURL      string         `gorm:"column:pic_url;size:255;default:'';comment:品牌图片" json:"picUrl"`
	Sort        int            `gorm:"default:0;comment:品牌排序" json:"sort"`
	Description string         `gorm:"size:1024;default:'';comment:品牌描述" json:"description"`
	Status      int            `gorm:"default:0;comment:状态" json:"status"` // 0: 开启, 1: 关闭
	Creator     string         `gorm:"size:64;default:'';comment:创建者" json:"creator"`
	Updater     string         `gorm:"size:64;default:'';comment:更新者" json:"updater"`
	CreatedAt   time.Time      `gorm:"column:create_time;autoCreateTime;comment:创建时间" json:"createTime"`
	UpdatedAt   time.Time      `gorm:"column:update_time;autoUpdateTime;comment:更新时间" json:"updateTime"`
	Deleted     model.BitBool  `gorm:"column:deleted;softDelete:flag;default:0;comment:是否删除" json:"deleted"`
}

func (ProductBrand) TableName() string {
	return "product_brand"
}
