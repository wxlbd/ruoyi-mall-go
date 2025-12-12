package product

import (
	"backend-go/internal/model"
	"time"

)

// ProductComment 商品评价
type ProductComment struct {
	ID                int64                `gorm:"primaryKey;autoIncrement;comment:评价编号"`
	UserID            int64                `gorm:"comment:评价人编号"`
	UserNickname      string               `gorm:"size:64;comment:评价人名称"`
	UserAvatar        string               `gorm:"size:255;comment:评价人头像"`
	Anonymous         bool                 `gorm:"comment:是否匿名"`
	OrderID           int64                `gorm:"comment:交易订单编号"`
	OrderItemID       int64                `gorm:"comment:交易订单项编号"`
	SpuID             int64                `gorm:"comment:商品SPU编号"`
	SpuName           string               `gorm:"size:255;comment:商品SPU名称"`
	SkuID             int64                `gorm:"comment:商品SKU编号"`
	SkuPicURL         string               `gorm:"size:255;comment:商品SKU图片地址"`
	SkuProperties     []ProductSkuProperty `gorm:"serializer:json;comment:属性数组"`
	Visible           bool                 `gorm:"comment:是否可见"`
	Scores            int                  `gorm:"comment:评分星级"`
	DescriptionScores int                  `gorm:"comment:描述星级"`
	BenefitScores     int                  `gorm:"comment:服务星级"`
	Content           string               `gorm:"type:text;comment:评论内容"`
	PicURLs           []string             `gorm:"serializer:json;comment:评论图片地址数组"`
	ReplyStatus       bool                 `gorm:"comment:商家是否回复"`
	ReplyUserID       int64                `gorm:"comment:回复管理员编号"`
	ReplyContent      string               `gorm:"type:text;comment:商家回复内容"`
	ReplyTime         *time.Time           `gorm:"comment:商家回复时间"`

	Creator   string         `gorm:"size:64;default:'';comment:创建者"`
	Updater   string         `gorm:"size:64;default:'';comment:更新者"`
	CreatedAt time.Time      `gorm:"column:create_time;autoCreateTime;comment:创建时间"`
	UpdatedAt time.Time      `gorm:"column:update_time;autoUpdateTime;comment:更新时间"`
	Deleted   model.BitBool  `gorm:"column:deleted;softDelete:flag;default:0;comment:是否删除"`
}

func (ProductComment) TableName() string {
	return "product_comment"
}
