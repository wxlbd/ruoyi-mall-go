package infra

import (
	"github.com/gin-gonic/gin"
	infra2 "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/infra"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/infra"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/excel"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"
	"github.com/wxlbd/ruoyi-mall-go/pkg/utils"
)

type JobHandler struct {
	svc *infra.JobService
}

func NewJobHandler(svc *infra.JobService) *JobHandler {
	return &JobHandler{svc: svc}
}

// CreateJob 创建定时任务
func (h *JobHandler) CreateJob(c *gin.Context) {
	var r infra2.JobSaveReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	id, err := h.svc.CreateJob(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, id)
}

// UpdateJob 更新定时任务
func (h *JobHandler) UpdateJob(c *gin.Context) {
	var r infra2.JobSaveReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	if err := h.svc.UpdateJob(c, &r); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// UpdateJobStatus 更新定时任务状态
func (h *JobHandler) UpdateJobStatus(c *gin.Context) {
	id := utils.ParseInt64(c.Query("id"))
	status := int(utils.ParseInt64(c.Query("status")))
	if err := h.svc.UpdateJobStatus(c, id, status); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// DeleteJob 删除定时任务
func (h *JobHandler) DeleteJob(c *gin.Context) {
	id := utils.ParseInt64(c.Query("id"))
	if err := h.svc.DeleteJob(c, id); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// DeleteJobList 批量删除定时任务
func (h *JobHandler) DeleteJobList(c *gin.Context) {
	ids := utils.ParseIDs(c.QueryArray("ids"))
	if len(ids) == 0 {
		ids = utils.ParseIDs([]string{c.Query("ids")})
	}
	if len(ids) == 0 {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	if err := h.svc.DeleteJobList(c, ids); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// GetJob 获取定时任务
func (h *JobHandler) GetJob(c *gin.Context) {
	id := utils.ParseInt64(c.Query("id"))
	job, err := h.svc.GetJob(c, id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	if job == nil {
		response.WriteBizError(c, errors.ErrNotFound)
		return
	}
	response.WriteSuccess(c, infra2.JobResp{
		ID:             job.ID,
		Name:           job.Name,
		Status:         job.Status,
		HandlerName:    job.HandlerName,
		HandlerParam:   job.HandlerParam,
		CronExpression: job.CronExpression,
		RetryCount:     job.RetryCount,
		RetryInterval:  job.RetryInterval,
		MonitorTimeout: job.MonitorTimeout,
		CreateTime:     job.CreateTime,
	})
}

// GetJobPage 获取定时任务分页
func (h *JobHandler) GetJobPage(c *gin.Context) {
	var r infra2.JobPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	pageResult, err := h.svc.GetJobPage(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	list := make([]infra2.JobResp, len(pageResult.List))
	for i, job := range pageResult.List {
		list[i] = infra2.JobResp{
			ID:             job.ID,
			Name:           job.Name,
			Status:         job.Status,
			HandlerName:    job.HandlerName,
			HandlerParam:   job.HandlerParam,
			CronExpression: job.CronExpression,
			RetryCount:     job.RetryCount,
			RetryInterval:  job.RetryInterval,
			MonitorTimeout: job.MonitorTimeout,
			CreateTime:     job.CreateTime,
		}
	}

	response.WriteSuccess(c, pagination.PageResult[infra2.JobResp]{
		List:  list,
		Total: pageResult.Total,
	})
}

// TriggerJob 触发定时任务
func (h *JobHandler) TriggerJob(c *gin.Context) {
	id := utils.ParseInt64(c.Query("id"))
	if err := h.svc.TriggerJob(c, id); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// SyncJob 同步定时任务
func (h *JobHandler) SyncJob(c *gin.Context) {
	if err := h.svc.SyncJob(c); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// ExportJobExcel 导出定时任务 Excel
func (h *JobHandler) ExportJobExcel(c *gin.Context) {
	var r infra2.JobPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	// 导出全量（受当前过滤条件约束）
	r.PageNo = 1
	r.PageSize = 1000000
	pageResult, err := h.svc.GetJobPage(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	list := make([]infra2.JobResp, len(pageResult.List))
	for i, job := range pageResult.List {
		list[i] = infra2.JobResp{
			ID:             job.ID,
			Name:           job.Name,
			Status:         job.Status,
			HandlerName:    job.HandlerName,
			HandlerParam:   job.HandlerParam,
			CronExpression: job.CronExpression,
			RetryCount:     job.RetryCount,
			RetryInterval:  job.RetryInterval,
			MonitorTimeout: job.MonitorTimeout,
			CreateTime:     job.CreateTime,
		}
	}

	if err := excel.WriteExcel(c, "定时任务.xls", "数据", list); err != nil {
		response.WriteBizError(c, err)
	}
}

// GetJobNextTimes 获取定时任务的下 n 次执行时间
func (h *JobHandler) GetJobNextTimes(c *gin.Context) {
	id := utils.ParseInt64(c.Query("id"))
	count := int(utils.ParseInt64(c.Query("count")))
	if count <= 0 {
		count = 5 // 默认 5 次
	}
	times, err := h.svc.GetJobNextTimes(c, id, count)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, times)
}
