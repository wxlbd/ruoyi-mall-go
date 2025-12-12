package promotion

import (
	"context"

	"backend-go/internal/api/req"
	"backend-go/internal/model/promotion"
	"backend-go/internal/pkg/core"
	"backend-go/internal/repo/query"
)

type BargainRecordService struct {
	q *query.Query
}

func NewBargainRecordService(q *query.Query) *BargainRecordService {
	return &BargainRecordService{q: q}
}

// GetBargainRecordUserCountMap 获得砍价活动的用户参与数量 Map
func (s *BargainRecordService) GetBargainRecordUserCountMap(ctx context.Context, activityIds []int64, status *int) (map[int64]int, error) {
	if len(activityIds) == 0 {
		return make(map[int64]int), nil
	}
	q := s.q.PromotionBargainRecord
	do := q.WithContext(ctx).Select(q.ActivityID, q.UserID.Count()).Where(q.ActivityID.In(activityIds...))
	// Note: CountDistinct Not supported directly by .Count() usually?
	// Use Group logic.
	if status != nil {
		do = do.Where(q.Status.Eq(*status))
	}
	// Group by ActivityID
	rows, err := do.Group(q.ActivityID).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int64]int)
	for rows.Next() {
		var activityID int64
		var count int
		if err := rows.Scan(&activityID, &count); err == nil {
			result[activityID] = count
		}
	}
	return result, nil
}

// GetBargainRecordUserCount 获得砍价记录数量 (按活动和状态)
func (s *BargainRecordService) GetBargainRecordUserCount(ctx context.Context, activityID int64, status int) (int64, error) {
	q := s.q.PromotionBargainRecord
	do := q.WithContext(ctx).Where(q.Status.Eq(status))
	if activityID > 0 {
		do = do.Where(q.ActivityID.Eq(activityID))
	}
	return do.Count()
}

// GetBargainRecordList 获得砍价记录列表 (按状态，带限制)
func (s *BargainRecordService) GetBargainRecordList(ctx context.Context, status int, limit int) ([]*promotion.PromotionBargainRecord, error) {
	q := s.q.PromotionBargainRecord
	return q.WithContext(ctx).Where(q.Status.Eq(status)).Order(q.EndTime.Desc()).Limit(limit).Find()
}

// GetBargainRecord 获得砍价记录
func (s *BargainRecordService) GetBargainRecord(ctx context.Context, id int64) (*promotion.PromotionBargainRecord, error) {
	q := s.q.PromotionBargainRecord
	return q.WithContext(ctx).Where(q.ID.Eq(id)).First()
}

// GetLastBargainRecord 获得用户在活动中的最近一次记录
func (s *BargainRecordService) GetLastBargainRecord(ctx context.Context, userID int64, activityID int64) (*promotion.PromotionBargainRecord, error) {
	q := s.q.PromotionBargainRecord
	return q.WithContext(ctx).Where(q.UserID.Eq(userID), q.ActivityID.Eq(activityID)).Order(q.ID.Desc()).First()
}

// CreateBargainRecord 创建砍价记录
func (s *BargainRecordService) CreateBargainRecord(ctx context.Context, userID int64, req *req.AppBargainRecordCreateReq) (int64, error) {
	// TODO: Implement Logic
	return 0, nil
}

// GetBargainRecordPage 获得砍价记录分页
func (s *BargainRecordService) GetBargainRecordPage(ctx context.Context, userID int64, p *core.PageParam) (*core.PageResult[*promotion.PromotionBargainRecord], error) {
	q := s.q.PromotionBargainRecord
	// List Logic: Where(UserID = userID).Order(CreateTime Desc)
	list, total, err := q.WithContext(ctx).Where(q.UserID.Eq(userID)).Order(q.CreatedAt.Desc()).FindByPage(p.GetOffset(), p.PageSize)
	if err != nil {
		return nil, err
	}
	return &core.PageResult[*promotion.PromotionBargainRecord]{List: list, Total: total}, nil
}

// GetBargainRecordPageAdmin 获得砍价记录分页 (Admin)
func (s *BargainRecordService) GetBargainRecordPageAdmin(ctx context.Context, req *req.BargainRecordPageReq) (*core.PageResult[*promotion.PromotionBargainRecord], error) {
	q := s.q.PromotionBargainRecord
	do := q.WithContext(ctx)

	if req.Status != nil {
		do = do.Where(q.Status.Eq(*req.Status))
	}
	if len(req.DateRange) == 2 {
		do = do.Where(q.CreatedAt.Between(req.DateRange[0], req.DateRange[1]))
	}

	list, total, err := do.Order(q.CreatedAt.Desc()).FindByPage(int((req.PageNo-1)*req.PageSize), int(req.PageSize))
	if err != nil {
		return nil, err
	}
	return &core.PageResult[*promotion.PromotionBargainRecord]{List: list, Total: total}, nil
}
