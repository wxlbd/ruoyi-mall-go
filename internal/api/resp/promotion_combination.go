package resp

import (
	"time"
)

// CombinationProductRespVO 拼团商品 Response VO
type CombinationProductRespVO struct {
	SpuID             int64     `json:"spuId"`
	SkuID             int64     `json:"skuId"`
	CombinationPrice  int       `json:"combinationPrice"`
	ActivityStatus    int       `json:"activityStatus"`
	ActivityStartTime time.Time `json:"activityStartTime"`
	ActivityEndTime   time.Time `json:"activityEndTime"`
}

// CombinationActivityRespVO 拼团活动 Response VO (Admin)
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
	Products         []CombinationProductRespVO `json:"products"`
	CreateTime       time.Time                  `json:"createTime"`
}

// AppCombinationActivityRespVO (Simple list item)
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

// AppCombinationActivityDetailRespVO contains detailed info including products
type AppCombinationActivityDetailRespVO struct {
	AppCombinationActivityRespVO
	Products []CombinationProductRespVO `json:"products"`
}

// AppCombinationRecordRespVO
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

// AppCombinationRecordDetailRespVO
type AppCombinationRecordDetailRespVO struct {
	HeadRecord    AppCombinationRecordRespVO   `json:"headRecord"`
	MemberRecords []AppCombinationRecordRespVO `json:"memberRecords"`
}

// AppCombinationRecordSummaryRespVO
type AppCombinationRecordSummaryRespVO struct {
	UserCount int64    `json:"userCount"`
	Avatars   []string `json:"avatars"`
}

// CombinationRecordPageItemRespVO 拼团记录 Admin 分页 VO
type CombinationRecordPageItemRespVO struct {
	ID               int64     `json:"id"`
	ActivityID       int64     `json:"activityId"`
	ActivityName     string    `json:"activityName"`
	SpuID            int64     `json:"spuId"`
	SpuName          string    `json:"spuName"`
	PicUrl           string    `json:"picUrl"`
	UserID           int64     `json:"userId"`
	Nickname         string    `json:"nickname"`
	Avatar           string    `json:"avatar"`
	UserCount        int       `json:"userCount"`
	UserSize         int       `json:"userSize"`
	Status           int       `json:"status"`
	CombinationPrice int       `json:"combinationPrice"`
	HeadID           int64     `json:"headId"`
	OrderID          int64     `json:"orderId"`
	VirtualGroup     bool      `json:"virtualGroup"`
	ExpireTime       time.Time `json:"expireTime"`
	StartTime        time.Time `json:"startTime"`
	EndTime          time.Time `json:"endTime"`
	CreateTime       time.Time `json:"createTime"`
}
