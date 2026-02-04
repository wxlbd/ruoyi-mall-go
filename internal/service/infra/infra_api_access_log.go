package infra

import (
	"context"

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
