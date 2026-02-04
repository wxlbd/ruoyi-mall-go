package infra

import (
	"context"
	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/infra"
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

type ApiAccessLogService struct {
	q *query.Query
}

func NewApiAccessLogService(q *query.Query) *ApiAccessLogService {
	return &ApiAccessLogService{q: q}
}

// GetApiAccessLog 获取API访问日志详情
func (s *ApiAccessLogService) GetApiAccessLog(ctx context.Context, id int64) (*model.InfraApiAccessLog, error) {
	return s.q.InfraApiAccessLog.WithContext(ctx).Where(s.q.InfraApiAccessLog.ID.Eq(id)).First()
}

// GetApiAccessLogPage 获取API访问日志分页
func (s *ApiAccessLogService) GetApiAccessLogPage(ctx context.Context, r *infra.ApiAccessLogPageReq) (*pagination.PageResult[*model.InfraApiAccessLog], error) {
	q := s.q.InfraApiAccessLog.WithContext(ctx)

	if r.UserID != nil {
		q = q.Where(s.q.InfraApiAccessLog.UserID.Eq(*r.UserID))
	}
	if r.UserType != nil {
		q = q.Where(s.q.InfraApiAccessLog.UserType.Eq(*r.UserType))
	}
	if r.ApplicationName != "" {
		q = q.Where(s.q.InfraApiAccessLog.ApplicationName.Eq(r.ApplicationName))
	}
	if r.RequestURL != "" {
		q = q.Where(s.q.InfraApiAccessLog.RequestURL.Like("%" + r.RequestURL + "%"))
	}
	if len(r.BeginTime) == 2 {
		q = q.Where(s.q.InfraApiAccessLog.BeginTime.Between(r.BeginTime[0], r.BeginTime[1]))
	}
	if r.Duration != nil {
		q = q.Where(s.q.InfraApiAccessLog.Duration.Lte(*r.Duration))
	}
	if r.ResultCode != nil {
		q = q.Where(s.q.InfraApiAccessLog.ResultCode.Eq(*r.ResultCode))
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

	list, err := q.Order(s.q.InfraApiAccessLog.ID.Desc()).Offset(offset).Limit(pageSize).Find()
	if err != nil {
		return nil, err
	}

	return &pagination.PageResult[*model.InfraApiAccessLog]{
		List:  list,
		Total: total,
	}, nil
}

// CleanAccessLog 物理删除 N 天前的访问日志
// 对齐 Java: ApiAccessLogServiceImpl#cleanAccessLog
func (s *ApiAccessLogService) CleanAccessLog(ctx context.Context, exceedDay int, deleteLimit int) (int, error) {
	if exceedDay <= 0 || deleteLimit <= 0 {
		return 0, nil
	}

	expireTime := time.Now().AddDate(0, 0, -exceedDay)
	count := 0
	for i := 0; i < int(^uint16(0)); i++ {
		logs, err := s.q.InfraApiAccessLog.WithContext(ctx).
			Where(s.q.InfraApiAccessLog.CreateTime.Lt(expireTime)).
			Order(s.q.InfraApiAccessLog.ID.Asc()).
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
		if _, err = s.q.InfraApiAccessLog.WithContext(ctx).Where(s.q.InfraApiAccessLog.ID.In(ids...)).Delete(); err != nil {
			return count, err
		}
		count += len(ids)
		if len(ids) < deleteLimit {
			break
		}
	}
	return count, nil
}
