package req

import (
	"backend-go/internal/pkg/core"
)

type SensitiveWordCreateReq struct {
	Name        string   `json:"name" binding:"required"`
	Tags        []string `json:"tags" binding:"required"`
	Status      int      `json:"status" binding:"required"`
	Description string   `json:"description"`
}

type SensitiveWordUpdateReq struct {
	ID          int64    `json:"id" binding:"required"`
	Name        string   `json:"name" binding:"required"`
	Tags        []string `json:"tags" binding:"required"`
	Status      int      `json:"status" binding:"required"`
	Description string   `json:"description"`
}

type SensitiveWordDeleteReq struct {
	ID int64 `form:"id" binding:"required"`
}

type SensitiveWordPageReq struct {
	core.PageParam
	Name      string `form:"name"`
	Tag       string `form:"tag"`
	Status    *int   `form:"status"`
	StartDate string `form:"startDate"` // YYYY-MM-DD
	EndDate   string `form:"endDate"`   // YYYY-MM-DD
}

type SensitiveWordExportReq struct {
	Name      string `form:"name"`
	Tag       string `form:"tag"`
	Status    *int   `form:"status"`
	StartDate string `form:"startDate"`
	EndDate   string `form:"endDate"`
}

type SensitiveWordValidateReq struct {
	Text string   `form:"text" binding:"required"`
	Tags []string `form:"tags"` // Optional filter
}
