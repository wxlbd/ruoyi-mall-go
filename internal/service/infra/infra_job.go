package infra

import (
	"context"
	"errors"
	"fmt"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/infra"
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

// JobStatus 任务状态
const (
	JobStatusInit   = 0 // 初始化
	JobStatusNormal = 1 // 开启
	JobStatusStop   = 2 // 暂停
)

type JobService struct {
	q         *query.Query
	scheduler *Scheduler
}

func NewJobService(q *query.Query, scheduler *Scheduler) *JobService {
	return &JobService{q: q, scheduler: scheduler}
}

// CreateJob 创建定时任务
func (s *JobService) CreateJob(ctx context.Context, r *infra.JobSaveReq) (int64, error) {
	// 1. 校验 Cron 表达式
	if err := s.scheduler.ValidateCronExpression(r.CronExpression); err != nil {
		return 0, err
	}

	// 2. 校验 Handler 存在
	if err := s.validateJobHandlerExists(r.HandlerName); err != nil {
		return 0, err
	}

	var jobID int64
	err := s.q.Transaction(func(tx *query.Query) error {
		// 3. 校验 handlerName 唯一性
		if err := s.validateHandlerNameUnique(ctx, r.HandlerName, nil); err != nil {
			return err
		}

		// 4. 为 MonitorTimeout 提供默认值
		monitorTimeout := r.MonitorTimeout
		if monitorTimeout == nil {
			defaultTimeout := 0
			monitorTimeout = &defaultTimeout
		}

		// 5. 创建任务记录（初始状态）
		job := &model.InfraJob{
			Name:           r.Name,
			Status:         JobStatusInit,
			HandlerName:    r.HandlerName,
			HandlerParam:   r.HandlerParam,
			CronExpression: r.CronExpression,
			RetryCount:     r.RetryCount,
			RetryInterval:  r.RetryInterval,
			MonitorTimeout: monitorTimeout,
		}
		if err := tx.InfraJob.WithContext(ctx).Create(job); err != nil {
			return err
		}
		jobID = job.ID

		// 6. 添加任务到调度器并更新状态为正常
		if s.scheduler != nil {
			if err := s.scheduler.AddJob(ctx, job.ID); err != nil {
				return err
			}
			// 更新状态为正常
			_, err := tx.InfraJob.WithContext(ctx).Where(tx.InfraJob.ID.Eq(job.ID)).Update(tx.InfraJob.Status, JobStatusNormal)
			return err
		}
		return nil
	})

	if err != nil {
		return 0, err
	}

	return jobID, nil
}

// UpdateJob 更新定时任务
func (s *JobService) UpdateJob(ctx context.Context, r *infra.JobSaveReq) error {
	if r.ID == nil {
		return errors.New("任务 ID 不能为空")
	}

	// 1. 校验 Cron 表达式
	if err := s.scheduler.ValidateCronExpression(r.CronExpression); err != nil {
		return err
	}

	// 2. 校验 Handler 存在
	if err := s.validateJobHandlerExists(r.HandlerName); err != nil {
		return err
	}

	return s.q.Transaction(func(tx *query.Query) error {
		// 3. 校验任务是否存在
		job, err := s.GetJob(ctx, *r.ID)
		if err != nil {
			return err
		}
		if job == nil {
			return errors.New("任务不存在")
		}
		// 只有开启状态，才可以修改.原因是，如果出暂停状态，修改调度时，会导致任务又开始执行
		if job.Status != JobStatusNormal {
			return errors.New("只有开启状态的任务才可以修改")
		}

		// 4. 校验 handlerName 唯一性
		if err := s.validateHandlerNameUnique(ctx, r.HandlerName, r.ID); err != nil {
			return err
		}

		// 5. 为 MonitorTimeout 提供默认值
		monitorTimeout := r.MonitorTimeout
		if monitorTimeout == nil {
			defaultTimeout := 0
			monitorTimeout = &defaultTimeout
		}

		_, err = tx.InfraJob.WithContext(ctx).Where(tx.InfraJob.ID.Eq(*r.ID)).Updates(map[string]interface{}{
			"name":            r.Name,
			"handler_name":    r.HandlerName,
			"handler_param":   r.HandlerParam,
			"cron_expression": r.CronExpression,
			"retry_count":     r.RetryCount,
			"retry_interval":  r.RetryInterval,
			"monitor_timeout": monitorTimeout,
		})
		if err != nil {
			return err
		}

		// 6. 重新调度
		if s.scheduler != nil {
			_ = s.scheduler.RemoveJob(*r.ID)
			return s.scheduler.AddJob(ctx, *r.ID)
		}
		return nil
	})
}

// validateJobHandlerExists 校验 Handler 是否已注册
func (s *JobService) validateJobHandlerExists(handlerName string) error {
	if s.scheduler == nil {
		return errors.New("调度器未初始化")
	}
	if !s.scheduler.HasHandler(handlerName) {
		return fmt.Errorf("定时任务处理器 '%s' 不存在，请检查并确保已在 Scheduler 注册。当前可用 Handler: %v",
			handlerName, s.scheduler.GetRegisteredHandlers())
	}
	return nil
}

// validateHandlerNameUnique 校验 Handler 唯一性
func (s *JobService) validateHandlerNameUnique(ctx context.Context, handlerName string, excludeID *int64) error {
	query := s.q.InfraJob.WithContext(ctx).Where(s.q.InfraJob.HandlerName.Eq(handlerName))
	if excludeID != nil {
		query = query.Where(s.q.InfraJob.ID.Neq(*excludeID))
	}
	count, err := query.Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("该处理器名称已被使用")
	}
	return nil
}

// DeleteJob 删除定时任务
func (s *JobService) DeleteJob(ctx context.Context, id int64) error {
	if s.scheduler != nil {
		_ = s.scheduler.RemoveJob(id)
	}
	_, err := s.q.InfraJob.WithContext(ctx).Where(s.q.InfraJob.ID.Eq(id)).Delete()
	return err
}

// DeleteJobList 批量删除定时任务
func (s *JobService) DeleteJobList(ctx context.Context, ids []int64) error {
	for _, id := range ids {
		if id <= 0 {
			continue
		}
		if err := s.DeleteJob(ctx, id); err != nil {
			return err
		}
	}
	return nil
}

// GetJob 获取定时任务
func (s *JobService) GetJob(ctx context.Context, id int64) (*model.InfraJob, error) {
	return s.q.InfraJob.WithContext(ctx).Where(s.q.InfraJob.ID.Eq(id)).First()
}

// GetJobPage 获取定时任务分页
func (s *JobService) GetJobPage(ctx context.Context, r *infra.JobPageReq) (*pagination.PageResult[*model.InfraJob], error) {
	q := s.q.InfraJob.WithContext(ctx)

	if r.Name != "" {
		q = q.Where(s.q.InfraJob.Name.Like("%" + r.Name + "%"))
	}
	if r.HandlerName != "" {
		q = q.Where(s.q.InfraJob.HandlerName.Like("%" + r.HandlerName + "%"))
	}
	if r.Status != nil {
		q = q.Where(s.q.InfraJob.Status.Eq(*r.Status))
	}

	pageNo := r.PageNo
	pageSize := r.PageSize
	if pageNo <= 0 {
		pageNo = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (pageNo - 1) * pageSize

	total, err := q.Count()
	if err != nil {
		return nil, err
	}

	list, err := q.Order(s.q.InfraJob.ID.Desc()).Offset(offset).Limit(pageSize).Find()
	if err != nil {
		return nil, err
	}

	return &pagination.PageResult[*model.InfraJob]{
		List:  list,
		Total: total,
	}, nil
}

// UpdateJobStatus 更新定时任务状态
func (s *JobService) UpdateJobStatus(ctx context.Context, id int64, status int) error {
	_, err := s.q.InfraJob.WithContext(ctx).Where(s.q.InfraJob.ID.Eq(id)).Update(s.q.InfraJob.Status, status)
	if err != nil {
		return err
	}
	if s.scheduler != nil {
		return s.scheduler.UpdateJobStatus(ctx, id, status)
	}
	return nil
}

// TriggerJob 触发定时任务
func (s *JobService) TriggerJob(ctx context.Context, id int64) error {
	if s.scheduler != nil {
		return s.scheduler.TriggerJob(ctx, id)
	}
	return errors.New("调度器未初始化")
}

// SyncJob 同步定时任务 (从数据库重加载)
func (s *JobService) SyncJob(ctx context.Context) error {
	if s.scheduler == nil {
		return errors.New("调度器未初始化")
	}
	// 查询所有开启的任务
	jobs, err := s.q.InfraJob.WithContext(ctx).Where(s.q.InfraJob.Status.Eq(JobStatusNormal)).Find()
	if err != nil {
		return err
	}

	// 重新加载所有任务 (这里简单实现：清空再添加，或者是 Scheduler 内部处理)
	// 假设 Scheduler 有 Reload 方法，或者我们手动遍历
	// 简单起见，调用 Scheduler 的 Initialize (如果支持) 或逐个添加
	// 这里我们假设需要刷新整个调度器，但为了安全，我们只对状态正常的任务进行确保添加
	// 更好的做法是 Scheduler 提供 Sync 接口

	// 临时方案：遍历所有任务，确保它们在调度器中
	for _, job := range jobs {
		_ = s.scheduler.AddJob(ctx, job.ID)
	}
	return nil
}

// GetJobNextTimes 获取下几次执行时间
func (s *JobService) GetJobNextTimes(ctx context.Context, id int64, count int) ([]string, error) {
	job, err := s.GetJob(ctx, id)
	if err != nil {
		return nil, err
	}
	if job == nil {
		return nil, errors.New("任务不存在")
	}

	// 使用 Cron 库解析
	// 注意：这里需要引入 robo/cron 或类似的库解析 Cron 表达式
	// 为了简化，这里暂时返回空列表，后续补充具体 parsing 逻辑或调用 scheduler 方法
	// 如果 Scheduler 暴露了 Parse 逻辑最好

	return s.scheduler.GetNextTimes(job.CronExpression, count)
}
