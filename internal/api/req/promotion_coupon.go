package req

import "time"

// CouponTemplateCreateReq 创建优惠券模板 Request
type CouponTemplateCreateReq struct {
	Name               string     `json:"name"`
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
	DiscountLimit      int        `json:"discountLimit"`
}

// CouponTemplateUpdateReq 更新优惠券模板 Request
type CouponTemplateUpdateReq struct {
	ID                 int64      `json:"id"`
	Name               string     `json:"name"`
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
	DiscountLimit      int        `json:"discountLimit"`
}

// CouponTemplatePageReq 优惠券模板分页 Request
type CouponTemplatePageReq struct {
	PageNo   int    `form:"pageNo,default=1"`
	PageSize int    `form:"pageSize,default=10"`
	Name     string `form:"name"`
	Status   *int   `form:"status"`
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
