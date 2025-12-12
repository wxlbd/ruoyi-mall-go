package req

// ProductSpuSaveReq 创建/更新商品 SPU Request
type ProductSpuSaveReq struct {
	ID                 int64                `json:"id"` // 更新时必填
	Name               string               `json:"name" binding:"required"`
	Keyword            string               `json:"keyword"`
	Introduction       string               `json:"introduction"`
	Description        string               `json:"description"`
	CategoryID         int64                `json:"categoryId" binding:"required"`
	BrandID            int64                `json:"brandId" binding:"required"`
	PicURL             string               `json:"picUrl" binding:"required"`
	SliderPicURLs      []string             `json:"sliderPicUrls"`
	Sort               int                  `json:"sort" binding:"min=0"`
	SpecType           *bool                `json:"specType" binding:"required"` // false: 单规格, true: 多规格
	DeliveryTypes      []int                `json:"deliveryTypes" binding:"required"`
	DeliveryTemplateID int64                `json:"deliveryTemplateId" binding:"required"`
	GiveIntegral       int                  `json:"giveIntegral" binding:"min=0"`
	SubCommissionType  *bool                `json:"subCommissionType" binding:"required"`
	VirtualSalesCount  int                  `json:"virtualSalesCount" binding:"min=0"`
	Skus               []*ProductSkuSaveReq `json:"skus" binding:"required,dive"`
}

// ProductSkuSaveReq SKU 保存 Request
type ProductSkuSaveReq struct {
	ID                   int64                   `json:"id"` // 更新时可选
	Properties           []ProductSkuPropertyReq `json:"properties"`
	Price                int                     `json:"price" binding:"required,min=0"`
	MarketPrice          int                     `json:"marketPrice" binding:"required,min=0"`
	CostPrice            int                     `json:"costPrice" binding:"required,min=0"`
	BarCode              string                  `json:"barCode"`
	PicURL               string                  `json:"picUrl"`
	Stock                int                     `json:"stock" binding:"required,min=0"`
	Weight               float64                 `json:"weight"`
	Volume               float64                 `json:"volume"`
	FirstBrokeragePrice  int                     `json:"firstBrokeragePrice"`
	SecondBrokeragePrice int                     `json:"secondBrokeragePrice"`
}

type ProductSkuPropertyReq struct {
	PropertyID   int64  `json:"propertyId" binding:"required"`
	PropertyName string `json:"propertyName"` // 冗余，后端可重新查询填充
	ValueID      int64  `json:"valueId" binding:"required"`
	ValueName    string `json:"valueName"` // 冗余，后端可重新查询填充
}

// ProductSpuUpdateStatusReq 更新商品状态 Request
type ProductSpuUpdateStatusReq struct {
	ID     int64 `json:"id" binding:"required"`
	Status int   `json:"status" binding:"required,oneof=0 1"` // 0: 上架, 1: 下架
}

// ProductSpuPageReq 分页查询商品 Request
type ProductSpuPageReq struct {
	PageNo     int      `form:"pageNo" binding:"required,min=1"`
	PageSize   int      `form:"pageSize" binding:"required,min=1,max=100"`
	TabType    *int     `form:"tabType"` // 标签类型，见 ProductSpuPageReq.FOR_SALE 等常量
	Name       string   `form:"name"`
	CategoryID int64    `form:"categoryId"`
	CreateTime []string `form:"createTime[]"`
}

// ProductSkuUpdateStockReq SKU 库存更新 Request
type ProductSkuUpdateStockReq struct {
	Items []ProductSkuUpdateStockItemReq
}

type ProductSkuUpdateStockItemReq struct {
	ID        int64
	IncrCount int
}

const (
	SpuTabForSale     = 0 // 出售中
	SpuTabInWarehouse = 1 // 仓库中
	SpuTabSoldOut     = 2 // 已售空
	SpuTabAlertStock  = 3 // 警戒库存
	SpuTabRecycleBin  = 4 // 回收站
)

// ProductSpuListReq 根据 ID 列表查询 SPU Request
type ProductSpuListReq struct {
	SpuIDs []int64 `form:"spuIds" binding:"required"`
}
