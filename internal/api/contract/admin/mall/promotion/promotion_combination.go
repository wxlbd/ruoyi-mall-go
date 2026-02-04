package promotion

import (
	"time"

	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

// CombinationActivityBaseVO 拼团活动 Base VO
type CombinationActivityBaseVO struct {
	Name             string    `json:"name"`
	SpuID            int64     `json:"spuId"`
	TotalLimitCount  int       `json:"totalLimitCount"`
	SingleLimitCount int       `json:"singleLimitCount"`
	StartTime        time.Time `json:"startTime"`
	EndTime          time.Time `json:"endTime"`
	UserSize         int       `json:"userSize"`
	VirtualGroup     bool      `json:"virtualGroup"`
	LimitDuration    int       `json:"limitDuration"`
}

// CombinationProductBaseVO 拼团商品 Base VO
type CombinationProductBaseVO struct {
	SpuID            int64 `json:"spuId"`
	SkuID            int64 `json:"skuId"`
	CombinationPrice int   `json:"combinationPrice"`
}

// CombinationActivityCreateReq 拼团活动创建 Request VO
type CombinationActivityCreateReq struct {
	CombinationActivityBaseVO
	Products []CombinationProductBaseVO `json:"products"`
}

// CombinationActivityUpdateReq 拼团活动更新 Request VO
type CombinationActivityUpdateReq struct {
	ID int64 `json:"id"`
	CombinationActivityCreateReq
}

// CombinationActivityPageReq 拼团活动分页 Request VO
type CombinationActivityPageReq struct {
	pagination.PageParam
	Name       string       `form:"name"`
	Status     *int         `form:"status"`
	CreateTime []*time.Time `form:"createTime"`
}

// AppCombinationRecordPageReq 拼团记录分页 Request VO
type AppCombinationRecordPageReq struct {
	pagination.PageParam
	Status int `json:"status"` // 0-进行中 1-成功 2-失败
}

// CombinationRecordPageReq 拼团记录分页 Request VO (Admin)
type CombinationRecordPageReq struct {
	pagination.PageParam
	Status    *int        `json:"status" form:"status"`
	Name      string      `json:"name" form:"name"` // User Nickname?
	DateRange []time.Time `json:"dateRange" form:"dateRange" time_format:"2006-01-02 15:04:05"`
}

// AppCombinationRecordSummaryRespVO App 拼团记录摘要 Response VO
type AppCombinationRecordSummaryRespVO struct {
	UserCount int64    `json:"userCount"`
	Avatars   []string `json:"avatars"`
}

// AppCombinationRecordRespVO App 拼团记录 Response VO
type AppCombinationRecordRespVO struct {
	ID               int64     `json:"id"`
	ActivityID       int64     `json:"activityId"`
	Nickname         string    `json:"nickname"`
	Avatar           string    `json:"avatar"`
	ExpireTime       time.Time `json:"expireTime"`
	UserSize         int       `json:"userSize"`
	UserCount        int       `json:"userCount"`
	Status           int       `json:"status"`
	OrderID          int64     `json:"orderId"`
	SpuName          string    `json:"spuName"`
	PicUrl           string    `json:"picUrl"`
	Count            int       `json:"count"`
	CombinationPrice int       `json:"combinationPrice"`
}

// AppCombinationRecordDetailRespVO App 拼团记录详情 Response VO
type AppCombinationRecordDetailRespVO struct {
	HeadRecord    AppCombinationRecordRespVO   `json:"headRecord"`
	MemberRecords []AppCombinationRecordRespVO `json:"memberRecords"`
	OrderID       int64                        `json:"orderId"`
}

// CombinationRecordSummaryVO Admin 拼团记录摘要 Response VO
type CombinationRecordSummaryVO struct {
	UserCount         int64 `json:"userCount"`
	SuccessCount      int64 `json:"successCount"`
	VirtualGroupCount int64 `json:"virtualGroupCount"`
}

// CombinationProductRespVO 拼团商品 Response VO
type CombinationProductRespVO struct {
	SpuID             int64     `json:"spuId"`
	SkuID             int64     `json:"skuId"`
	CombinationPrice  int       `json:"combinationPrice"`
	ActivityStatus    int       `json:"activityStatus"`
	ActivityStartTime time.Time `json:"activityStartTime"`
	ActivityEndTime   time.Time `json:"activityEndTime"`
}

// CombinationActivityRespVO 拼团活动 Response VO
type CombinationActivityRespVO struct {
	ID               int64                      `json:"id"`
	Name             string                     `json:"name"`
	SpuID            int64                      `json:"spuId"`
	TotalLimitCount  int                        `json:"totalLimitCount"`
	SingleLimitCount int                        `json:"singleLimitCount"`
	StartTime        time.Time                  `json:"startTime"`
	EndTime          time.Time                  `json:"endTime"`
	UserSize         int                        `json:"userSize"`
	VirtualGroup     bool                       `json:"virtualGroup"`
	LimitDuration    int                        `json:"limitDuration"`
	Status           int                        `json:"status"`
	CreateTime       time.Time                  `json:"createTime"`
	Products         []CombinationProductRespVO `json:"products"`
	// Derived Fields
	SpuName          string `json:"spuName"`
	PicUrl           string `json:"picUrl"`
	MarketPrice      int    `json:"marketPrice"`
	CombinationPrice int    `json:"combinationPrice"`
}

// CombinationActivityPageItemRespVO 拼团活动分页项 Response VO
type CombinationActivityPageItemRespVO struct {
	CombinationActivityRespVO
	GroupCount        int `json:"groupCount"`
	GroupSuccessCount int `json:"groupSuccessCount"`
	RecordCount       int `json:"recordCount"`
}

// AppCombinationActivityRespVO App 拼团活动 Response VO
type AppCombinationActivityRespVO struct {
	ID               int64  `json:"id"`
	Name             string `json:"name"`
	UserSize         int    `json:"userSize"`
	SpuID            int64  `json:"spuId"`
	SpuName          string `json:"spuName"`
	PicUrl           string `json:"picUrl"`
	MarketPrice      int    `json:"marketPrice"`
	CombinationPrice int    `json:"combinationPrice"`
}

// AppCombinationActivityDetailProduct App 拼团详情关联商品 VO
type AppCombinationActivityDetailProduct struct {
	SkuID            int64 `json:"skuId"`
	CombinationPrice int   `json:"combinationPrice"`
}

// AppCombinationActivityDetailRespVO App 拼团活动详情 Response VO
type AppCombinationActivityDetailRespVO struct {
	ID               int64                                 `json:"id"`
	Name             string                                `json:"name"`
	Status           int                                   `json:"status"`
	StartTime        *time.Time                            `json:"startTime"`
	EndTime          *time.Time                            `json:"endTime"`
	UserSize         int                                   `json:"userSize"`
	SuccessCount     int                                   `json:"successCount"`
	SpuID            int64                                 `json:"spuId"`
	TotalLimitCount  int                                   `json:"totalLimitCount"`
	SingleLimitCount int                                   `json:"singleLimitCount"`
	Products         []AppCombinationActivityDetailProduct `json:"products"`
}

// CombinationRecordPageItemRespVO 拼团记录分页项 Response VO
type CombinationRecordPageItemRespVO struct {
	ID               int64     `json:"id"`
	ActivityID       int64     `json:"activityId"`
	ActivityName     string    `json:"activityName"`
	UserID           int64     `json:"userId"`
	Nickname         string    `json:"nickname"`
	Avatar           string    `json:"avatar"`
	StartTime        time.Time `json:"startTime"`
	EndTime          time.Time `json:"endTime"`
	ExpireTime       time.Time `json:"expireTime"`
	UserSize         int       `json:"userSize"`
	UserCount        int       `json:"userCount"`
	Status           int       `json:"status"`
	OrderID          int64     `json:"orderId"`
	HeadID           int64     `json:"headId"`
	VirtualGroup     bool      `json:"virtualGroup"`
	SpuID            int64     `json:"spuId"`
	SpuName          string    `json:"spuName"`
	PicUrl           string    `json:"picUrl"`
	Count            int       `json:"count"`
	CombinationPrice int       `json:"combinationPrice"`
	CreateTime       time.Time `json:"createTime"`
}
