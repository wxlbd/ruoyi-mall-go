package product

import (
	"backend-go/internal/model"
	"time"
)

// ProductSpu 商品 SPU
type ProductSpu struct {
	ID                 int64                `gorm:"primaryKey;autoIncrement;comment:主键" json:"id"`
	Name               string               `gorm:"size:255;not null;comment:商品名称" json:"name"`
	Keyword            string               `gorm:"size:255;default:'';comment:关键字" json:"keyword"`
	Introduction       string               `gorm:"size:1024;default:'';comment:商品简介" json:"introduction"`
	Description        string               `gorm:"type:text;comment:商品详情" json:"description"`
	CategoryID         int64                `gorm:"column:category_id;not null;comment:商品分类编号" json:"categoryId"`
	BrandID            int64                `gorm:"column:brand_id;default:0;comment:商品品牌编号" json:"brandId"`
	PicURL             string               `gorm:"column:pic_url;size:255;not null;comment:商品封面图" json:"picUrl"`
	SliderPicURLs      []string             `gorm:"column:slider_pic_urls;type:json;serializer:json;comment:商品轮播图" json:"sliderPicUrls"`
	Sort               int                  `gorm:"default:0;comment:排序字段" json:"sort"`
	Status             int                  `gorm:"default:0;comment:商品状态" json:"status"`                    // 0: 上架, 1: 下架, -1: 回收站
	SpecType           model.BitBool        `gorm:"column:spec_type;default:0;comment:规格类型" json:"specType"` // false: 单规格, true: 多规格
	Price              int                  `gorm:"default:0;comment:商品价格" json:"price"`                     // 单位：分
	MarketPrice        int                  `gorm:"default:0;comment:市场价" json:"marketPrice"`                // 单位：分
	CostPrice          int                  `gorm:"default:0;comment:成本价" json:"costPrice"`                  // 单位：分
	Stock              int                  `gorm:"default:0;comment:库存" json:"stock"`
	DeliveryTypes      model.IntListFromCSV `gorm:"column:delivery_types;comment:配送方式数组" json:"deliveryTypes"`
	DeliveryTemplateID int64                `gorm:"column:delivery_template_id;default:0;comment:物流配置模板编号" json:"deliveryTemplateId"`
	GiveIntegral       int                  `gorm:"default:0;comment:赠送积分" json:"giveIntegral"`
	SubCommissionType  model.BitBool        `gorm:"column:sub_commission_type;default:0;comment:分销类型" json:"subCommissionType"` // false: 默认, true: 自行设置
	SalesCount         int                  `gorm:"default:0;comment:商品销量" json:"salesCount"`
	VirtualSalesCount  int                  `gorm:"default:0;comment:虚拟销量" json:"virtualSalesCount"`
	BrowseCount        int                  `gorm:"default:0;comment:浏览量" json:"browseCount"`
	Creator            string               `gorm:"size:64;default:'';comment:创建者" json:"creator"`
	Updater            string               `gorm:"size:64;default:'';comment:更新者" json:"updater"`
	CreatedAt          time.Time            `gorm:"column:create_time;autoCreateTime;comment:创建时间" json:"createTime"`
	UpdatedAt          time.Time            `gorm:"column:update_time;autoUpdateTime;comment:更新时间" json:"updateTime"`
	Deleted            model.BitBool        `gorm:"column:deleted;softDelete:flag;default:0;comment:是否删除" json:"deleted"`
	TenantID           int64                `gorm:"column:tenant_id;default:0;comment:租户编号" json:"tenantId"`
}

func (ProductSpu) TableName() string {
	return "product_spu"
}
