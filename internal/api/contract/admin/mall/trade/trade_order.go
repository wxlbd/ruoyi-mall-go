package trade

import (
	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/mall/product"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
	"github.com/wxlbd/ruoyi-mall-go/pkg/types"
)

// ========== Admin Response DTOs ==========

// TradeOrderBase 订单基础响应
type TradeOrderBase struct {
	ID                    int64      `json:"id"`
	No                    string     `json:"no"`
	CreateTime            time.Time  `json:"createTime"`
	Type                  int        `json:"type"`
	Terminal              int        `json:"terminal"`
	UserID                int64      `json:"userId"`
	UserIP                string     `json:"userIp"`
	UserRemark            string     `json:"userRemark"`
	Status                int        `json:"status"`
	ProductCount          int        `json:"productCount"`
	FinishTime            *time.Time `json:"finishTime"`
	CancelTime            *time.Time `json:"cancelTime"`
	CancelType            int        `json:"cancelType"`
	Remark                string     `json:"remark"`
	PayOrderID            int64      `json:"payOrderId"`
	PayStatus             bool       `json:"payStatus"`
	PayTime               *time.Time `json:"payTime"`
	PayChannelCode        string     `json:"payChannelCode"`
	TotalPrice            int        `json:"totalPrice"`
	DiscountPrice         int        `json:"discountPrice"`
	DeliveryPrice         int        `json:"deliveryPrice"`
	AdjustPrice           int        `json:"adjustPrice"`
	PayPrice              int        `json:"payPrice"`
	DeliveryType          int        `json:"deliveryType"`
	LogisticsID           int64      `json:"logisticsId"`
	LogisticsNo           string     `json:"logisticsNo"`
	DeliveryTime          *time.Time `json:"deliveryTime"`
	ReceiveTime           *time.Time `json:"receiveTime"`
	ReceiverName          string     `json:"receiverName"`
	ReceiverMobile        string     `json:"receiverMobile"`
	ReceiverAreaID        int32      `json:"receiverAreaId"`
	ReceiverDetailAddress string     `json:"receiverDetailAddress"`
	RefundPrice           int        `json:"refundPrice"`
	CouponID              int64      `json:"couponId"`
	CouponPrice           int        `json:"couponPrice"`
	PointPrice            int        `json:"pointPrice"`
	VipPrice              int        `json:"vipPrice"`
}

// TradeOrderItemBase 订单项基础响应
type TradeOrderItemBase struct {
	ID         int64                            `json:"id"`
	UserID     int64                            `json:"userId"`
	OrderID    int64                            `json:"orderId"`
	SpuID      int64                            `json:"spuId"`
	SpuName    string                           `json:"spuName"`
	SkuID      int64                            `json:"skuId"`
	PicURL     string                           `json:"picUrl"`
	Count      int                              `json:"count"`
	Price      int                              `json:"price"`
	PayPrice   int                              `json:"payPrice"`
	Properties []product.ProductSkuPropertyResp `json:"properties"`
}

// TradeOrderPageItemResp 订单分页项
type TradeOrderPageItemResp struct {
	TradeOrderBase
	Items            []TradeOrderItemBase `json:"items"`
	User             interface{}          `json:"user"`          // 后续对齐 MemberUserResp
	BrokerageUser    interface{}          `json:"brokerageUser"` // 后续对齐 MemberUserResp
	ReceiverAreaName string               `json:"receiverAreaName"`
}

// TradeOrderLogResp 订单日志
type TradeOrderLogResp struct {
	Content    string    `json:"content"`
	CreateTime time.Time `json:"createTime"`
	UserType   int       `json:"userType"`
}

// TradeOrderDetailResp 订单详情
type TradeOrderDetailResp struct {
	TradeOrderBase
	Items            []TradeOrderItemBase `json:"items"`
	Logs             []TradeOrderLogResp  `json:"logs"`
	User             interface{}          `json:"user"`
	BrokerageUser    interface{}          `json:"brokerageUser"`
	ReceiverAreaName string               `json:"receiverAreaName"`
}

// AppTradeOrderSettlementReq 交易订单结算请求
type AppTradeOrderSettlementReq struct {
	Items                 []AppTradeOrderSettlementItem `json:"items" form:"items" binding:"dive"`
	CouponID              *int64                        `json:"couponId" form:"couponId"`
	PointStatus           *bool                         `json:"pointStatus" form:"pointStatus" binding:"required"`
	DeliveryType          int                           `json:"deliveryType" form:"deliveryType" binding:"required"` // 1: 快递, 2: 自提
	AddressID             *int64                        `json:"addressId" form:"addressId"`
	PickUpStoreID         *int64                        `json:"pickUpStoreId" form:"pickUpStoreId"`
	ReceiverName          string                        `json:"receiverName" form:"receiverName"`
	ReceiverMobile        string                        `json:"receiverMobile" form:"receiverMobile"`
	SeckillActivityID     *int64                        `json:"seckillActivityId" form:"seckillActivityId"`
	CombinationActivityID *int64                        `json:"combinationActivityId" form:"combinationActivityId"`
	CombinationHeadID     *int64                        `json:"combinationHeadId" form:"combinationHeadId"`
	BargainRecordID       *int64                        `json:"bargainRecordId" form:"bargainRecordId"`
	PointActivityID       *int64                        `json:"pointActivityId" form:"pointActivityId"`
}

type AppTradeOrderSettlementItem struct {
	SkuID  int64 `json:"skuId" form:"skuId"`
	Count  int   `json:"count" form:"count"`
	CartID int64 `json:"cartId" form:"cartId"`
}

// AppTradeOrderCreateReq 交易订单创建请求
type AppTradeOrderCreateReq struct {
	AppTradeOrderSettlementReq
	Remark string `json:"remark"`
}

// AppTradeOrderPageReq 交易订单分页请求
type AppTradeOrderPageReq struct {
	pagination.PageParam
	Status        *int  `form:"status"`
	CommentStatus *bool `form:"commentStatus"`
}

// TradeOrderPageReq 管理后台 - 交易订单分页请求
type TradeOrderPageReq struct {
	pagination.PageParam
	No               string   `form:"no"`
	UserID           *int64   `form:"userId"`
	UserNickname     string   `form:"userNickname"`
	UserMobile       string   `form:"userMobile"`
	DeliveryType     *int     `form:"deliveryType"`
	LogisticsID      *int64   `form:"logisticsId"`
	PickUpStoreIDs   []int64  `form:"pickUpStoreIds"`
	PickUpVerifyCode string   `form:"pickUpVerifyCode"`
	Type             *int     `form:"type"`
	Status           *int     `form:"status"`
	PayChannelCode   string   `form:"payChannelCode"`
	CreateTime       []string `form:"createTime[]"` // time range
	Terminal         *int     `form:"terminal"`
	CommentStatus    *bool    `form:"commentStatus"`
}

// TradeOrderDeliveryReq 订单发货请求
type TradeOrderDeliveryReq struct {
	ID int64 `json:"id" binding:"required"`
	// 说明：对齐 Java 版本（TradeOrderDeliveryReqVO）
	// - logisticsId 仅要求非空（可为 0 表示无需发货）
	// - logisticsNo 可为空（无需发货时为 ""）
	LogisticsID *int64 `json:"logisticsId" binding:"required"`
	LogisticsNo string `json:"logisticsNo"`
}

// TradeOrderUpdateAddressReq 更新订单地址请求
type TradeOrderUpdateAddressReq struct {
	ID                    int64  `json:"id" binding:"required"`
	ReceiverName          string `json:"receiverName" binding:"required"`
	ReceiverMobile        string `json:"receiverMobile" binding:"required"`
	ReceiverAreaID        int    `json:"receiverAreaId" binding:"required"`
	ReceiverDetailAddress string `json:"receiverDetailAddress" binding:"required"`
}

// TradeOrderUpdatePriceReq 更新订单价格请求
type TradeOrderUpdatePriceReq struct {
	ID          int64 `json:"id" binding:"required"`
	AdjustPrice int   `json:"adjustPrice" binding:"required"` // 单位：分
}

// TradeOrderRemarkReq 订单备注请求
type TradeOrderRemarkReq struct {
	ID     int64  `json:"id" binding:"required"`
	Remark string `json:"remark" binding:"required"`
}

// AppTradeOrderItemCommentCreateReq 用户 App - 商品评价创建请求
type AppTradeOrderItemCommentCreateReq struct {
	Anonymous         bool     `json:"anonymous"`
	OrderItemID       int64    `json:"orderItemId" binding:"required"`
	DescriptionScores int      `json:"descriptionScores" binding:"required,min=1,max=5"`
	BenefitScores     int      `json:"benefitScores" binding:"required,min=1,max=5"`
	Content           string   `json:"content" binding:"required"`
	PicUrls           []string `json:"picUrls" binding:"max=9"`
}

// AppTradeOrderSettlementProductReq 获得订单结算的商品信息请求
type AppTradeOrderSettlementProductReq struct {
	SkuID int64 `form:"skuId" binding:"required"`
	Count int   `form:"count" binding:"required,min=1"`
}

// AppTradeOrderSettlementQueryReq 订单结算查询请求 (用于GET请求)
type AppTradeOrderSettlementQueryReq struct {
	SkuIds                types.ListFromCSV[int64] `form:"skuIds"`  // SKU ID列表，逗号分隔
	Counts                types.ListFromCSV[int]   `form:"counts"`  // 数量列表，逗号分隔
	CartIds               types.ListFromCSV[int64] `form:"cartIds"` // 购物车ID列表，逗号分隔
	CouponID              *int64                   `form:"couponId"`
	PointStatus           *bool                    `form:"pointStatus" binding:"required"`
	DeliveryType          int                      `form:"deliveryType" binding:"required"` // 1: 快递, 2: 自提
	AddressID             *int64                   `form:"addressId"`
	PickUpStoreID         *int64                   `form:"pickUpStoreId"`
	ReceiverName          string                   `form:"receiverName"`
	ReceiverMobile        string                   `form:"receiverMobile"`
	SeckillActivityID     *int64                   `form:"seckillActivityId"`
	CombinationActivityID *int64                   `form:"combinationActivityId"`
	CombinationHeadID     *int64                   `form:"combinationHeadId"`
	BargainRecordID       *int64                   `form:"bargainRecordId"`
	PointActivityID       *int64                   `form:"pointActivityId"`
}

func (q *AppTradeOrderSettlementQueryReq) ToSettlementItems() []AppTradeOrderSettlementItem {
	if len(q.SkuIds) == 0 {
		return nil
	}
	items := make([]AppTradeOrderSettlementItem, len(q.SkuIds))
	for i, skuId := range q.SkuIds {
		items[i].SkuID = skuId
		if i < len(q.Counts) {
			items[i].Count = q.Counts[i]
		}
		if i < len(q.CartIds) {
			items[i].CartID = q.CartIds[i]
		}
	}
	return items
}

// AppTradeOrderUpdatePaidReq 更新订单为已支付请求
type AppTradeOrderUpdatePaidReq struct {
	ID int64 `json:"id" binding:"required"`
}

// AppTradeOrderReceiveReq 确认收货请求
type AppTradeOrderReceiveReq struct {
	ID int64 `json:"id" binding:"required"`
}

// AppTradeOrderSettlementResp 订单结算响应
type AppTradeOrderSettlementResp struct {
	Type       int                                `json:"type"`
	Price      AppTradeOrderSettlementPrice       `json:"price"`
	Items      []AppTradeOrderSettlementItemResp  `json:"items"`
	Coupons    []AppTradeOrderSettlementCoupon    `json:"coupons"`
	Promotions []AppTradeOrderSettlementPromotion `json:"promotions"`
	Address    *AppTradeOrderSettlementAddress    `json:"address"`
	UsePoint   int                                `json:"usePoint"`
	TotalPoint int                                `json:"totalPoint"`
}

type AppTradeOrderSettlementPrice struct {
	TotalPrice    int `json:"totalPrice"`
	DiscountPrice int `json:"discountPrice"`
	DeliveryPrice int `json:"deliveryPrice"`
	CouponPrice   int `json:"couponPrice"`
	PointPrice    int `json:"pointPrice"`
	VipPrice      int `json:"vipPrice"`
	PayPrice      int `json:"payPrice"`
}

type AppTradeOrderSettlementItemResp struct {
	CategoryID int64                            `json:"categoryId"`
	SpuID      int64                            `json:"spuId"`
	SpuName    string                           `json:"spuName"`
	SkuID      int64                            `json:"skuId"`
	PicURL     string                           `json:"picUrl"`
	Price      int                              `json:"price"`
	Count      int                              `json:"count"`
	CartID     *int64                           `json:"cartId"`
	Properties []product.ProductSkuPropertyResp `json:"properties"`
}

type AppTradeOrderSettlementCoupon struct {
	ID                 int64   `json:"id"`
	Name               string  `json:"name"`
	UsePrice           int     `json:"usePrice"`
	ValidStartTime     int64   `json:"validStartTime"`
	ValidEndTime       int64   `json:"validEndTime"`
	DiscountType       int     `json:"discountType"`
	DiscountPercent    int     `json:"discountPercent"`
	DiscountPrice      int     `json:"discountPrice"`
	DiscountLimitPrice int     `json:"discountLimitPrice"`
	Match              bool    `json:"match"`
	MismatchReason     *string `json:"mismatchReason"`
}

type AppTradeOrderSettlementPromotion struct {
	ID            int64                                  `json:"id"`
	Name          string                                 `json:"name"`
	Type          int                                    `json:"type"`
	TotalPrice    int                                    `json:"totalPrice"`
	DiscountPrice int                                    `json:"discountPrice"`
	Match         bool                                   `json:"match"`
	Description   string                                 `json:"description"`
	Items         []AppTradeOrderSettlementPromotionItem `json:"items"`
}

type AppTradeOrderSettlementPromotionItem struct {
	SkuID         int64 `json:"skuId"`
	TotalPrice    int   `json:"totalPrice"`
	DiscountPrice int   `json:"discountPrice"`
	PayPrice      int   `json:"payPrice"`
}

type AppTradeOrderSettlementAddress struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	Mobile        string `json:"mobile"`
	AreaID        int32  `json:"areaId"`
	AreaName      string `json:"areaName"`
	DetailAddress string `json:"detailAddress"`
	DefaultStatus bool   `json:"defaultStatus"`
}

// TradeOrderSummaryResp 交易订单统计响应
type TradeOrderSummaryResp struct {
	OrderCount     int64 `json:"orderCount"`
	OrderPayPrice  int64 `json:"orderPayPrice"`
	AfterSaleCount int64 `json:"afterSaleCount"`
	AfterSalePrice int64 `json:"afterSalePrice"`
}

// ========== App Response DTOs ==========

// AppTradeOrderCreateResp 订单创建响应
type AppTradeOrderCreateResp struct {
	ID         int64 `json:"id"`
	PayOrderID int64 `json:"payOrderId"`
}

// AppTradeOrderItemResp 订单项响应
type AppTradeOrderItemResp struct {
	ID              int64                            `json:"id"`
	OrderID         int64                            `json:"orderId"`
	SpuID           int64                            `json:"spuId"`
	SpuName         string                           `json:"spuName"`
	SkuID           int64                            `json:"skuId"`
	PicURL          string                           `json:"picUrl"`
	Count           int                              `json:"count"`
	CommentStatus   bool                             `json:"commentStatus"`
	Price           int                              `json:"price"`
	PayPrice        int                              `json:"payPrice"`
	AfterSaleID     int64                            `json:"afterSaleId"`
	AfterSaleStatus int                              `json:"afterSaleStatus"`
	Properties      []product.ProductSkuPropertyResp `json:"properties"`
}

// AppTradeOrderDetailResp 订单详情响应
type AppTradeOrderDetailResp struct {
	ID                    int64                   `json:"id"`
	No                    string                  `json:"no"`
	Type                  int                     `json:"type"`
	CreateTime            types.JsonDateTime      `json:"createTime"`
	UserRemark            string                  `json:"userRemark"`
	Status                int                     `json:"status"`
	ProductCount          int                     `json:"productCount"`
	FinishTime            *types.JsonDateTime     `json:"finishTime"`
	CancelTime            *types.JsonDateTime     `json:"cancelTime"`
	CommentStatus         bool                    `json:"commentStatus"`
	PayStatus             bool                    `json:"payStatus"`
	PayOrderID            int64                   `json:"payOrderId"`
	PayTime               *types.JsonDateTime     `json:"payTime"`
	PayChannelCode        string                  `json:"payChannelCode"`
	TotalPrice            int                     `json:"totalPrice"`
	DiscountPrice         int                     `json:"discountPrice"`
	DeliveryPrice         int                     `json:"deliveryPrice"`
	AdjustPrice           int                     `json:"adjustPrice"`
	PayPrice              int                     `json:"payPrice"`
	DeliveryType          int                     `json:"deliveryType"`
	LogisticsID           int64                   `json:"logisticsId"`
	LogisticsNo           string                  `json:"logisticsNo"`
	DeliveryTime          *types.JsonDateTime     `json:"deliveryTime"`
	ReceiveTime           *types.JsonDateTime     `json:"receiveTime"`
	ReceiverName          string                  `json:"receiverName"`
	ReceiverMobile        string                  `json:"receiverMobile"`
	ReceiverAreaID        int                     `json:"receiverAreaId"`
	ReceiverDetailAddress string                  `json:"receiverDetailAddress"`
	RefundStatus          int                     `json:"refundStatus"`
	RefundPrice           int                     `json:"refundPrice"`
	CouponID              int64                   `json:"couponId"`
	CouponPrice           int                     `json:"couponPrice"`
	PointPrice            int                     `json:"pointPrice"`
	VipPrice              int                     `json:"vipPrice"`
	CombinationRecordID   int64                   `json:"combinationRecordId"`
	Items                 []AppTradeOrderItemResp `json:"items"`
}

// AppTradeOrderPageItemResp 订单分页项响应
type AppTradeOrderPageItemResp struct {
	ID                  int64                   `json:"id"`
	No                  string                  `json:"no"`
	Type                int                     `json:"type"`
	Status              int                     `json:"status"`
	ProductCount        int                     `json:"productCount"`
	CommentStatus       bool                    `json:"commentStatus"`
	CreateTime          types.JsonDateTime      `json:"createTime"`
	PayOrderID          int64                   `json:"payOrderId"`
	PayPrice            int                     `json:"payPrice"`
	DeliveryType        int                     `json:"deliveryType"`
	Items               []AppTradeOrderItemResp `json:"items"`
	CombinationRecordID int64                   `json:"combinationRecordId"`
}
