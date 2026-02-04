package infra

import (
	"context"
	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/infra"
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

type ApiErrorLogService struct {
	q *query.Query
}

func NewApiErrorLogService(q *query.Query) *ApiErrorLogService {
	return &ApiErrorLogService{q: q}
}

// GetApiErrorLog 获取API错误日志详情
func (s *ApiErrorLogService) GetApiErrorLog(ctx context.Context, id int64) (*model.InfraApiErrorLog, error) {
	return s.q.InfraApiErrorLog.WithContext(ctx).Where(s.q.InfraApiErrorLog.ID.Eq(id)).First()
}

// GetApiErrorLogPage 获取API错误日志分页
func (s *ApiErrorLogService) GetApiErrorLogPage(ctx context.Context, r *infra.ApiErrorLogPageReq) (*pagination.PageResult[*model.InfraApiErrorLog], error) {
	q := s.q.InfraApiErrorLog.WithContext(ctx)

	if r.UserID != nil {
		q = q.Where(s.q.InfraApiErrorLog.UserID.Eq(*r.UserID))
	}
	if r.UserType != nil {
		q = q.Where(s.q.InfraApiErrorLog.UserType.Eq(*r.UserType))
	}
	if r.ApplicationName != "" {
		q = q.Where(s.q.InfraApiErrorLog.ApplicationName.Eq(r.ApplicationName))
	}
	if r.RequestURL != "" {
		q = q.Where(s.q.InfraApiErrorLog.RequestURL.Like("%" + r.RequestURL + "%"))
	}
	if len(r.ExceptionTime) == 2 {
		q = q.Where(s.q.InfraApiErrorLog.ExceptionTime.Between(r.ExceptionTime[0], r.ExceptionTime[1]))
	}
	if r.ProcessStatus != nil {
		q = q.Where(s.q.InfraApiErrorLog.ProcessStatus.Eq(*r.ProcessStatus))
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

	list, err := q.Order(s.q.InfraApiErrorLog.ID.Desc()).Offset(offset).Limit(pageSize).Find()
	if err != nil {
		return nil, err
	}

	return &pagination.PageResult[*model.InfraApiErrorLog]{
		List:  list,
		Total: total,
	}, nil
}

// UpdateApiErrorLogProcess 更新API错误日志处理状态
func (s *ApiErrorLogService) UpdateApiErrorLogProcess(ctx context.Context, id int64, processStatus int, processUserID int64) error {
	now := time.Now()
	_, err := s.q.InfraApiErrorLog.WithContext(ctx).Where(s.q.InfraApiErrorLog.ID.Eq(id)).Updates(map[string]interface{}{
		"process_status":  processStatus,
		"process_time":    now,
		"process_user_id": processUserID,
	})
	return err
}

// CleanErrorLog 物理删除 N 天前的错误日志
// 对齐 Java: ApiErrorLogServiceImpl#cleanErrorLog
func (s *ApiErrorLogService) CleanErrorLog(ctx context.Context, exceedDay int, deleteLimit int) (int, error) {
	if exceedDay <= 0 || deleteLimit <= 0 {
		return 0, nil
	}

	expireTime := time.Now().AddDate(0, 0, -exceedDay)
	count := 0
	for i := 0; i < int(^uint16(0)); i++ {
		logs, err := s.q.InfraApiErrorLog.WithContext(ctx).
			Where(s.q.InfraApiErrorLog.CreateTime.Lt(expireTime)).
			Order(s.q.InfraApiErrorLog.ID.Asc()).
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
		if _, err = s.q.InfraApiErrorLog.WithContext(ctx).Where(s.q.InfraApiErrorLog.ID.In(ids...)).Delete(); err != nil {
			return count, err
		}
		count += len(ids)
		if len(ids) < deleteLimit {
			break
		}
	}
	return count, nil
}
