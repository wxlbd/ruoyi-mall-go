package promotion

import (
	"time"

	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

// CouponTemplateCreateReq 创建优惠券模板 Request
type CouponTemplateCreateReq struct {
	Name               string     `json:"name"`
	Description        string     `json:"description"`
	Status             int        `json:"status"`
	TotalCount         int        `json:"totalCount"`
	TakeLimitCount     int        `json:"takeLimitCount"`
	TakeType           int        `json:"takeType"`
	UsePriceMin        int        `json:"usePriceMin"`
	ProductScope       int        `json:"productScope"`
	ProductScopeValues []int64    `json:"productScopeValues"`
	ValidityType       int        `json:"validityType"`
	ValidStartTime     *time.Time `json:"validStartTime"`
	ValidEndTime       *time.Time `json:"validEndTime"`
	FixedStartTerm     int        `json:"fixedStartTerm"`
	FixedEndTerm       int        `json:"fixedEndTerm"`
	DiscountType       int        `json:"discountType"`
	DiscountPrice      int        `json:"discountPrice"`
	DiscountPercent    int        `json:"discountPercent"`
	DiscountLimitPrice int        `json:"discountLimitPrice"`
}

// CouponTemplateUpdateReq 更新优惠券模板 Request
type CouponTemplateUpdateReq struct {
	ID                 int64      `json:"id"`
	Name               string     `json:"name"`
	Description        string     `json:"description"`
	Status             int        `json:"status"`
	TotalCount         int        `json:"totalCount"`
	TakeLimitCount     int        `json:"takeLimitCount"`
	TakeType           int        `json:"takeType"`
	UsePriceMin        int        `json:"usePriceMin"`
	ProductScope       int        `json:"productScope"`
	ProductScopeValues []int64    `json:"productScopeValues"`
	ValidityType       int        `json:"validityType"`
	ValidStartTime     *time.Time `json:"validStartTime"`
	ValidEndTime       *time.Time `json:"validEndTime"`
	FixedStartTerm     int        `json:"fixedStartTerm"`
	FixedEndTerm       int        `json:"fixedEndTerm"`
	DiscountType       int        `json:"discountType"`
	DiscountPrice      int        `json:"discountPrice"`
	DiscountPercent    int        `json:"discountPercent"`
	DiscountLimitPrice int        `json:"discountLimitPrice"`
}

// CouponTemplatePageReq 优惠券模板分页 Request
type CouponTemplatePageReq struct {
	pagination.PageParam
	Name              string       `form:"name"`
	Status            *int32       `form:"status"`
	DiscountType      *int         `form:"discountType"`
	CreateTime        []*time.Time `form:"createTime"`
	CanTakeTypes      []int        `form:"canTakeTypes"`
	ProductScope      *int         `form:"productScope"`
	ProductScopeValue *int64       `form:"productScopeValue"`
}

// CouponPageReq 优惠券分页 Request
type CouponPageReq struct {
	PageNo   int    `form:"pageNo,default=1"`
	PageSize int    `form:"pageSize,default=10"`
	UserID   *int64 `form:"userId"`
	Status   *int   `form:"status"`
}

// CouponTemplateUpdateStatusReq 更新优惠券模板状态 Request
// 对应 Java: CouponTemplateUpdateStatusReqVO
type CouponTemplateUpdateStatusReq struct {
	ID     int64 `json:"id"`
	Status int32 `json:"status"`
}

// CouponSendReq 发送优惠券 Request
// 对应 Java: CouponSendReqVO
type CouponSendReq struct {
	TemplateID int64   `json:"templateId"`
	UserIDs    []int64 `json:"userIds"`
}

// CouponTemplateResp 优惠券模板 Response VO
type CouponTemplateResp struct {
	ID                 int64      `json:"id"`
	Name               string     `json:"name"`
	Description        string     `json:"description"`
	Status             int        `json:"status"`
	TotalCount         int        `json:"totalCount"`
	TakeLimitCount     int        `json:"takeLimitCount"`
	TakeType           int        `json:"takeType"`
	UsePrice           int        `json:"usePrice"`
	ProductScope       int        `json:"productScope"`
	ProductScopeValues []int64    `json:"productScopeValues"`
	ValidityType       int        `json:"validityType"`
	ValidStartTime     *time.Time `json:"validStartTime"`
	ValidEndTime       *time.Time `json:"validEndTime"`
	FixedStartTerm     *int       `json:"fixedStartTerm"`
	FixedEndTerm       *int       `json:"fixedEndTerm"`
	DiscountType       int        `json:"discountType"`
	DiscountPercent    *int       `json:"discountPercent"`
	DiscountPrice      *int       `json:"discountPrice"`
	DiscountLimitPrice *int       `json:"discountLimitPrice"`
	TakeCount          int        `json:"takeCount"`
	UseCount           int        `json:"useCount"`
	Creator            string     `json:"creator"`
	Updater            string     `json:"updater"`
	CreateTime         time.Time  `json:"createTime"`
	UpdateTime         time.Time  `json:"updateTime"`
}

// CouponPageResp 优惠券 Response VO
type CouponPageResp struct {
	ID                 int64      `json:"id"`
	TemplateID         int64      `json:"templateId"`
	Name               string     `json:"name"`
	Status             int        `json:"status"`
	UserID             int64      `json:"userId"`
	Nickname           string     `json:"nickname"`
	TakeType           int        `json:"takeType"`
	UsePrice           int        `json:"usePrice"`
	ValidStartTime     time.Time  `json:"validStartTime"`
	ValidEndTime       time.Time  `json:"validEndTime"`
	ProductScope       int        `json:"productScope"`
	ProductScopeValues []int64    `json:"productScopeValues"`
	DiscountType       int        `json:"discountType"`
	DiscountPercent    *int       `json:"discountPercent"`
	DiscountPrice      *int       `json:"discountPrice"`
	DiscountLimitPrice *int       `json:"discountLimitPrice"`
	UseOrderID         *int64     `json:"useOrderId"`
	UseTime            *time.Time `json:"useTime"`
	CreateTime         time.Time  `json:"createTime"`
}
