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

type JobLogHandler struct {
	svc *infra.JobLogService
}

func NewJobLogHandler(svc *infra.JobLogService) *JobLogHandler {
	return &JobLogHandler{svc: svc}
}

// GetJobLog 获取定时任务日志
func (h *JobLogHandler) GetJobLog(c *gin.Context) {
	id := utils.ParseInt64(c.Query("id"))
	log, err := h.svc.GetJobLog(c, id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	if log == nil {
		response.WriteBizError(c, errors.ErrNotFound)
		return
	}
	response.WriteSuccess(c, infra2.JobLogResp{
		ID:           log.ID,
		JobID:        log.JobID,
		HandlerName:  log.HandlerName,
		HandlerParam: log.HandlerParam,
		ExecuteIndex: log.ExecuteIndex,
		BeginTime:    log.BeginTime,
		EndTime:      log.EndTime,
		Duration:     log.Duration,
		Status:       log.Status,
		Result:       log.Result,
		CreateTime:   log.CreateTime,
	})
}

// GetJobLogPage 获取定时任务日志分页
func (h *JobLogHandler) GetJobLogPage(c *gin.Context) {
	var r infra2.JobLogPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	pageResult, err := h.svc.GetJobLogPage(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	list := make([]infra2.JobLogResp, len(pageResult.List))
	for i, log := range pageResult.List {
		list[i] = infra2.JobLogResp{
			ID:           log.ID,
			JobID:        log.JobID,
			HandlerName:  log.HandlerName,
			HandlerParam: log.HandlerParam,
			ExecuteIndex: log.ExecuteIndex,
			BeginTime:    log.BeginTime,
			EndTime:      log.EndTime,
			Duration:     log.Duration,
			Status:       log.Status,
			Result:       log.Result,
			CreateTime:   log.CreateTime,
		}
	}

	response.WriteSuccess(c, pagination.PageResult[infra2.JobLogResp]{
		List:  list,
		Total: pageResult.Total,
	})
}

// ExportJobLogExcel 导出定时任务日志 Excel
func (h *JobLogHandler) ExportJobLogExcel(c *gin.Context) {
	var r infra2.JobLogPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	// 导出全量（受当前过滤条件约束）
	r.PageNo = 1
	r.PageSize = 1000000
	pageResult, err := h.svc.GetJobLogPage(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	list := make([]infra2.JobLogResp, len(pageResult.List))
	for i, log := range pageResult.List {
		list[i] = infra2.JobLogResp{
			ID:           log.ID,
			JobID:        log.JobID,
			HandlerName:  log.HandlerName,
			HandlerParam: log.HandlerParam,
			ExecuteIndex: log.ExecuteIndex,
			BeginTime:    log.BeginTime,
			EndTime:      log.EndTime,
			Duration:     log.Duration,
			Status:       log.Status,
			Result:       log.Result,
			CreateTime:   log.CreateTime,
		}
	}

	if err := excel.WriteExcel(c, "任务日志.xls", "数据", list); err != nil {
		response.WriteError(c, 500, err.Error())
	}
}
