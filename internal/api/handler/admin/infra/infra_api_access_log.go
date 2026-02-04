package infra

import (
	infra2 "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/infra"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/infra"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/excel"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"
	"github.com/wxlbd/ruoyi-mall-go/pkg/utils"

	"github.com/gin-gonic/gin"
)

type ApiAccessLogHandler struct {
	svc *infra.ApiAccessLogService
}

func NewApiAccessLogHandler(svc *infra.ApiAccessLogService) *ApiAccessLogHandler {
	return &ApiAccessLogHandler{svc: svc}
}

// GetApiAccessLogPage 获取API访问日志分页
func (h *ApiAccessLogHandler) GetApiAccessLogPage(c *gin.Context) {
	var r infra2.ApiAccessLogPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}

	pageResult, err := h.svc.GetApiAccessLogPage(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	list := make([]infra2.ApiAccessLogResp, len(pageResult.List))
	for i, log := range pageResult.List {
		list[i] = infra2.ApiAccessLogResp{
			ID:              log.ID,
			TraceID:         log.TraceID,
			UserID:          log.UserID,
			UserType:        log.UserType,
			ApplicationName: log.ApplicationName,
			RequestMethod:   log.RequestMethod,
			RequestURL:      log.RequestURL,
			RequestParams:   log.RequestParams,
			ResponseBody:    log.ResponseBody,
			UserIP:          log.UserIP,
			UserAgent:       log.UserAgent,
			OperateModule:   log.OperateModule,
			OperateName:     log.OperateName,
			OperateType:     log.OperateType,
			BeginTime:       log.BeginTime,
			EndTime:         log.EndTime,
			Duration:        log.Duration,
			ResultCode:      log.ResultCode,
			ResultMsg:       log.ResultMsg,
			CreateTime:      log.CreateTime,
		}
	}

	response.WriteSuccess(c, pagination.PageResult[infra2.ApiAccessLogResp]{
		List:  list,
		Total: pageResult.Total,
	})
}

// GetApiAccessLog 获取API访问日志详情
func (h *ApiAccessLogHandler) GetApiAccessLog(c *gin.Context) {
	id := utils.ParseInt64(c.Query("id"))
	if id == 0 {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	log, err := h.svc.GetApiAccessLog(c, id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, infra2.ApiAccessLogResp{
		ID:              log.ID,
		TraceID:         log.TraceID,
		UserID:          log.UserID,
		UserType:        log.UserType,
		ApplicationName: log.ApplicationName,
		RequestMethod:   log.RequestMethod,
		RequestURL:      log.RequestURL,
		RequestParams:   log.RequestParams,
		ResponseBody:    log.ResponseBody,
		UserIP:          log.UserIP,
		UserAgent:       log.UserAgent,
		OperateModule:   log.OperateModule,
		OperateName:     log.OperateName,
		OperateType:     log.OperateType,
		BeginTime:       log.BeginTime,
		EndTime:         log.EndTime,
		Duration:        log.Duration,
		ResultCode:      log.ResultCode,
		ResultMsg:       log.ResultMsg,
		CreateTime:      log.CreateTime,
	})
}

// ExportApiAccessLogExcel 导出API访问日志
func (h *ApiAccessLogHandler) ExportApiAccessLogExcel(c *gin.Context) {
	var r infra2.ApiAccessLogPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	r.PageNo = 1
	r.PageSize = 1000000
	pageResult, err := h.svc.GetApiAccessLogPage(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	list := make([]infra2.ApiAccessLogResp, len(pageResult.List))
	for i, log := range pageResult.List {
		list[i] = infra2.ApiAccessLogResp{
			ID:              log.ID,
			TraceID:         log.TraceID,
			UserID:          log.UserID,
			UserType:        log.UserType,
			ApplicationName: log.ApplicationName,
			RequestMethod:   log.RequestMethod,
			RequestURL:      log.RequestURL,
			RequestParams:   log.RequestParams,
			ResponseBody:    log.ResponseBody,
			UserIP:          log.UserIP,
			UserAgent:       log.UserAgent,
			OperateModule:   log.OperateModule,
			OperateName:     log.OperateName,
			OperateType:     log.OperateType,
			BeginTime:       log.BeginTime,
			EndTime:         log.EndTime,
			Duration:        log.Duration,
			ResultCode:      log.ResultCode,
			ResultMsg:       log.ResultMsg,
			CreateTime:      log.CreateTime,
		}
	}

	if err := excel.WriteExcel(c, "API 访问日志.xls", "数据", list); err != nil {
		response.WriteBizError(c, err)
	}
}
