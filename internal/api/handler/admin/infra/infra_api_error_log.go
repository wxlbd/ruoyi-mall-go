package infra

import (
	infra2 "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/infra"
	pkgContext "github.com/wxlbd/ruoyi-mall-go/pkg/context"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/infra"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/excel"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"
	"github.com/wxlbd/ruoyi-mall-go/pkg/utils"

	"github.com/gin-gonic/gin"
)

type ApiErrorLogHandler struct {
	svc *infra.ApiErrorLogService
}

func NewApiErrorLogHandler(svc *infra.ApiErrorLogService) *ApiErrorLogHandler {
	return &ApiErrorLogHandler{svc: svc}
}

// GetApiErrorLogPage 获取API错误日志分页
func (h *ApiErrorLogHandler) GetApiErrorLogPage(c *gin.Context) {
	var r infra2.ApiErrorLogPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}

	pageResult, err := h.svc.GetApiErrorLogPage(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	list := make([]infra2.ApiErrorLogResp, len(pageResult.List))
	for i, log := range pageResult.List {
		list[i] = infra2.ApiErrorLogResp{
			ID:                        log.ID,
			TraceID:                   log.TraceID,
			UserID:                    log.UserID,
			UserType:                  log.UserType,
			ApplicationName:           log.ApplicationName,
			RequestMethod:             log.RequestMethod,
			RequestURL:                log.RequestURL,
			RequestParams:             log.RequestParams,
			UserIP:                    log.UserIP,
			UserAgent:                 log.UserAgent,
			ExceptionTime:             log.ExceptionTime,
			ExceptionName:             log.ExceptionName,
			ExceptionMessage:          log.ExceptionMessage,
			ExceptionRootCauseMessage: log.ExceptionRootCauseMessage,
			ExceptionStackTrace:       log.ExceptionStackTrace,
			ExceptionClassName:        log.ExceptionClassName,
			ExceptionFileName:         log.ExceptionFileName,
			ExceptionMethodName:       log.ExceptionMethodName,
			ExceptionLineNumber:       log.ExceptionLineNumber,
			ProcessStatus:             log.ProcessStatus,
			ProcessTime:               log.ProcessTime,
			ProcessUserID:             log.ProcessUserID,
			CreateTime:                log.CreateTime,
		}
	}

	response.WriteSuccess(c, pagination.PageResult[infra2.ApiErrorLogResp]{
		List:  list,
		Total: pageResult.Total,
	})
}

// UpdateApiErrorLogProcess 更新API错误日志处理状态
func (h *ApiErrorLogHandler) UpdateApiErrorLogProcess(c *gin.Context) {
	id := utils.ParseInt64(c.Query("id"))
	processStatus := int(utils.ParseInt64(c.Query("processStatus")))
	if id == 0 {
		response.WriteBizError(c, errors.ErrParam)
		return
	}

	processUserID := pkgContext.GetLoginUserID(c)
	if processUserID == 0 {
		response.WriteBizError(c, errors.ErrUnauthorized)
		return
	}

	if err := h.svc.UpdateApiErrorLogProcess(c, id, processStatus, processUserID); err != nil {
		response.WriteError(c, 500, err.Error())
		return
	}

	response.WriteSuccess(c, true)
}

// GetApiErrorLog 获取API错误日志详情
func (h *ApiErrorLogHandler) GetApiErrorLog(c *gin.Context) {
	id := utils.ParseInt64(c.Query("id"))
	if id == 0 {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	log, err := h.svc.GetApiErrorLog(c, id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, infra2.ApiErrorLogResp{
		ID:                        log.ID,
		TraceID:                   log.TraceID,
		UserID:                    log.UserID,
		UserType:                  log.UserType,
		ApplicationName:           log.ApplicationName,
		RequestMethod:             log.RequestMethod,
		RequestURL:                log.RequestURL,
		RequestParams:             log.RequestParams,
		UserIP:                    log.UserIP,
		UserAgent:                 log.UserAgent,
		ExceptionTime:             log.ExceptionTime,
		ExceptionName:             log.ExceptionName,
		ExceptionMessage:          log.ExceptionMessage,
		ExceptionRootCauseMessage: log.ExceptionRootCauseMessage,
		ExceptionStackTrace:       log.ExceptionStackTrace,
		ExceptionClassName:        log.ExceptionClassName,
		ExceptionFileName:         log.ExceptionFileName,
		ExceptionMethodName:       log.ExceptionMethodName,
		ExceptionLineNumber:       log.ExceptionLineNumber,
		ProcessStatus:             log.ProcessStatus,
		ProcessTime:               log.ProcessTime,
		ProcessUserID:             log.ProcessUserID,
		CreateTime:                log.CreateTime,
	})
}

// ExportApiErrorLogExcel 导出API错误日志
func (h *ApiErrorLogHandler) ExportApiErrorLogExcel(c *gin.Context) {
	var r infra2.ApiErrorLogPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	r.PageNo = 1
	r.PageSize = 1000000
	pageResult, err := h.svc.GetApiErrorLogPage(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	list := make([]infra2.ApiErrorLogResp, len(pageResult.List))
	for i, log := range pageResult.List {
		list[i] = infra2.ApiErrorLogResp{
			ID:                        log.ID,
			TraceID:                   log.TraceID,
			UserID:                    log.UserID,
			UserType:                  log.UserType,
			ApplicationName:           log.ApplicationName,
			RequestMethod:             log.RequestMethod,
			RequestURL:                log.RequestURL,
			RequestParams:             log.RequestParams,
			UserIP:                    log.UserIP,
			UserAgent:                 log.UserAgent,
			ExceptionTime:             log.ExceptionTime,
			ExceptionName:             log.ExceptionName,
			ExceptionMessage:          log.ExceptionMessage,
			ExceptionRootCauseMessage: log.ExceptionRootCauseMessage,
			ExceptionStackTrace:       log.ExceptionStackTrace,
			ExceptionClassName:        log.ExceptionClassName,
			ExceptionFileName:         log.ExceptionFileName,
			ExceptionMethodName:       log.ExceptionMethodName,
			ExceptionLineNumber:       log.ExceptionLineNumber,
			ProcessStatus:             log.ProcessStatus,
			ProcessTime:               log.ProcessTime,
			ProcessUserID:             log.ProcessUserID,
			CreateTime:                log.CreateTime,
		}
	}

	if err := excel.WriteExcel(c, "API 错误日志.xls", "数据", list); err != nil {
		response.WriteBizError(c, err)
	}
}
