package req

import (
	"time"

	"backend-go/internal/pkg/core"
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
	core.PageParam
	Name   string `json:"name"`
	Status int    `json:"status"`
}

// AppCombinationRecordPageReq 拼团记录分页 Request VO
type AppCombinationRecordPageReq struct {
	core.PageParam
	Status int `json:"status"` // 0-进行中 1-成功 2-失败
}

// CombinationRecordPageReq 拼团记录分页 Request VO (Admin)
type CombinationRecordPageReq struct {
	core.PageParam
	Status    *int        `json:"status" form:"status"`
	Name      string      `json:"name" form:"name"` // User Nickname?
	DateRange []time.Time `json:"dateRange" form:"dateRange" time_format:"2006-01-02 15:04:05"`
}
