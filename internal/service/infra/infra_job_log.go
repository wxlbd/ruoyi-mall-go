package infra

import (
	"context"
	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/infra"
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

type JobLogService struct {
	q *query.Query
}

func NewJobLogService(q *query.Query) *JobLogService {
	return &JobLogService{q: q}
}

// GetJobLog 获取定时任务日志
func (s *JobLogService) GetJobLog(ctx context.Context, id int64) (*model.InfraJobLog, error) {
	return s.q.InfraJobLog.WithContext(ctx).Where(s.q.InfraJobLog.ID.Eq(id)).First()
}

// GetJobLogPage 获取定时任务日志分页
func (s *JobLogService) GetJobLogPage(ctx context.Context, r *infra.JobLogPageReq) (*pagination.PageResult[*model.InfraJobLog], error) {
	q := s.q.InfraJobLog.WithContext(ctx)

	if r.JobID != nil {
		q = q.Where(s.q.InfraJobLog.JobID.Eq(*r.JobID))
	}
	if r.HandlerName != "" {
		q = q.Where(s.q.InfraJobLog.HandlerName.Like("%" + r.HandlerName + "%"))
	}
	if r.Status != nil {
		q = q.Where(s.q.InfraJobLog.Status.Eq(*r.Status))
	}
	if len(r.BeginTime) == 2 {
		q = q.Where(s.q.InfraJobLog.BeginTime.Between(r.BeginTime[0], r.BeginTime[1]))
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

	list, err := q.Order(s.q.InfraJobLog.ID.Desc()).Offset(offset).Limit(pageSize).Find()
	if err != nil {
		return nil, err
	}

	return &pagination.PageResult[*model.InfraJobLog]{
		List:  list,
		Total: total,
	}, nil
}

// CleanJobLog 物理删除 N 天前的任务日志
// 对齐 Java: JobLogServiceImpl#cleanJobLog
func (s *JobLogService) CleanJobLog(ctx context.Context, exceedDay int, deleteLimit int) (int, error) {
	if exceedDay <= 0 || deleteLimit <= 0 {
		return 0, nil
	}

	expireTime := time.Now().AddDate(0, 0, -exceedDay)
	count := 0
	for i := 0; i < int(^uint16(0)); i++ {
		logs, err := s.q.InfraJobLog.WithContext(ctx).
			Where(s.q.InfraJobLog.CreateTime.Lt(expireTime)).
			Order(s.q.InfraJobLog.ID.Asc()).
			Limit(deleteLimit).
			Find()
		if err != nil {
			return count, err
		}
		if len(logs) == 0 {
			break
		}

		ids := make([]int64, 0, len(logs))
		for _, log := range logs {
			ids = append(ids, log.ID)
		}
		if _, err = s.q.InfraJobLog.WithContext(ctx).Where(s.q.InfraJobLog.ID.In(ids...)).Delete(); err != nil {
			return count, err
		}
		count += len(ids)
		if len(ids) < deleteLimit {
			break
		}
	}
	return count, nil
}
