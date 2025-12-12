package promotion

import (
	"context"
	"time"

	"backend-go/internal/api/req"
	"backend-go/internal/api/resp"
	"backend-go/internal/model/promotion"
	"backend-go/internal/pkg/core"
	"backend-go/internal/repo/query"
	"backend-go/internal/service/member"
	prodSvc "backend-go/internal/service/product"
)

type CombinationRecordService interface {
	// App
	GetCombinationRecordSummary(ctx context.Context, activityID int64) (*resp.AppCombinationRecordSummaryRespVO, error)
	GetCombinationRecordPage(ctx context.Context, userID int64, req req.AppCombinationRecordPageReq) (*core.PageResult[*resp.AppCombinationRecordRespVO], error)
	GetCombinationRecordDetail(ctx context.Context, userID int64, id int64) (*resp.AppCombinationRecordDetailRespVO, error)
	GetLatestCombinationRecordList(ctx context.Context, activityID int64, count int) ([]*promotion.PromotionCombinationRecord, error)

	// Internal (for Order)
	ValidateCombinationRecord(ctx context.Context, userID int64, activityID int64, headID int64, skuID int64, count int) (*promotion.PromotionCombinationActivity, *promotion.PromotionCombinationProduct, error)
	CreateCombinationRecord(ctx context.Context, record *promotion.PromotionCombinationRecord) (int64, error)
	GetCombinationRecordPageAdmin(ctx context.Context, req *req.CombinationRecordPageReq) (*core.PageResult[*promotion.PromotionCombinationRecord], error)
}

type combinationRecordService struct {
	q           *query.Query
	activitySvc CombinationActivityService
	userSvc     *member.MemberUserService
	spuSvc      *prodSvc.ProductSpuService
	skuSvc      *prodSvc.ProductSkuService
}

func NewCombinationRecordService(
	q *query.Query,
	activitySvc CombinationActivityService,
	userSvc *member.MemberUserService,
	spuSvc *prodSvc.ProductSpuService,
	skuSvc *prodSvc.ProductSkuService,
) CombinationRecordService {
	return &combinationRecordService{
		q:           q,
		activitySvc: activitySvc,
		userSvc:     userSvc,
		spuSvc:      spuSvc,
		skuSvc:      skuSvc,
	}
}

func (s *combinationRecordService) GetCombinationRecordSummary(ctx context.Context, activityID int64) (*resp.AppCombinationRecordSummaryRespVO, error) {
	q := s.q.PromotionCombinationRecord

	count, err := q.WithContext(ctx).Distinct(q.UserID).Count()
	if err != nil {
		return nil, err
	}

	records, err := q.WithContext(ctx).
		// Where(status success). MVP: All records?
		Order(q.CreatedAt.Desc()).
		Limit(7).
		Find()
	if err != nil {
		return nil, err
	}

	avatars := make([]string, 0)
	for _, r := range records {
		if r.Avatar != "" {
			avatars = append(avatars, r.Avatar)
		}
	}

	return &resp.AppCombinationRecordSummaryRespVO{
		UserCount: count,
		Avatars:   avatars,
	}, nil
}

func (s *combinationRecordService) GetCombinationRecordPage(ctx context.Context, userID int64, req req.AppCombinationRecordPageReq) (*core.PageResult[*resp.AppCombinationRecordRespVO], error) {
	q := s.q.PromotionCombinationRecord
	do := q.WithContext(ctx).Where(q.UserID.Eq(userID))
	if req.Status != 0 {
		do = do.Where(q.Status.Eq(req.Status))
	}
	list, total, err := do.Order(q.CreatedAt.Desc()).FindByPage(req.GetOffset(), req.GetLimit())
	if err != nil {
		return nil, err
	}

	result := make([]*resp.AppCombinationRecordRespVO, len(list))
	for i, item := range list {
		result[i] = &resp.AppCombinationRecordRespVO{
			ID:               item.ID,
			ActivityID:       item.ActivityID,
			Nickname:         item.Nickname,
			Avatar:           item.Avatar,
			ExpireTime:       item.ExpireTime,
			UserSize:         item.UserSize,
			UserCount:        item.UserCount,
			Status:           item.Status,
			OrderID:          item.OrderID,
			SpuName:          item.SpuName,
			PicUrl:           item.PicUrl,
			Count:            item.Count,
			CombinationPrice: item.CombinationPrice,
		}
	}
	return &core.PageResult[*resp.AppCombinationRecordRespVO]{List: result, Total: total}, nil
}

func (s *combinationRecordService) GetCombinationRecordDetail(ctx context.Context, userID int64, id int64) (*resp.AppCombinationRecordDetailRespVO, error) {
	// MVP Not Implemented, required by Interface
	return nil, nil
}

func (s *combinationRecordService) GetLatestCombinationRecordList(ctx context.Context, activityID int64, count int) ([]*promotion.PromotionCombinationRecord, error) {
	q := s.q.PromotionCombinationRecord
	return q.WithContext(ctx).Where(q.ActivityID.Eq(activityID), q.Status.Eq(1)).Order(q.CreatedAt.Desc()).Limit(count).Find()
}

func (s *combinationRecordService) ValidateCombinationRecord(ctx context.Context, userID int64, activityID int64, headID int64, skuID int64, count int) (*promotion.PromotionCombinationActivity, *promotion.PromotionCombinationProduct, error) {
	activity, err := s.activitySvc.ValidateCombinationActivityCanJoin(ctx, activityID)
	if err != nil {
		return nil, nil, err
	}

	// 1.3 校验是否超出单次限购数量
	if count > activity.SingleLimitCount {
		return nil, nil, core.NewBizError(1001006012, "单次限购数量超出")
	}

	prod, err := s.q.PromotionCombinationProduct.WithContext(ctx).Where(
		s.q.PromotionCombinationProduct.ActivityID.Eq(activityID),
		s.q.PromotionCombinationProduct.SkuID.Eq(skuID),
	).First()
	if err != nil {
		return nil, nil, core.NewBizError(1001006004, "拼团活动商品不存在")
	}

	if headID > 0 {
		head, err := s.q.PromotionCombinationRecord.WithContext(ctx).Where(s.q.PromotionCombinationRecord.ID.Eq(headID)).First()
		if err != nil {
			return nil, nil, core.NewBizError(1001006005, "拼团不存在")
		}
		if head.Status != 0 { // 0: InProgress
			return nil, nil, core.NewBizError(1001006006, "拼团已结束")
		}
		if head.UserCount >= head.UserSize {
			return nil, nil, core.NewBizError(1001006007, "拼团人数已满")
		}
	}

	// 6.1 校验是否有拼团记录 (Already IN_PROGRESS) & Total Limit
	// Status!=2 (Failed) means InProgress(0) or Success(1)
	records, err := s.q.PromotionCombinationRecord.WithContext(ctx).Where(
		s.q.PromotionCombinationRecord.UserID.Eq(userID),
		s.q.PromotionCombinationRecord.ActivityID.Eq(activityID),
		s.q.PromotionCombinationRecord.Status.Neq(2),
	).Find()
	if err != nil {
		return nil, nil, err
	}

	totalCount := 0
	for _, r := range records {
		if r.Status == 0 { // InProgress
			return nil, nil, core.NewBizError(1001006013, "您已有该活动的拼团记录")
		}
		totalCount += r.Count
	}
	if totalCount+count > activity.TotalLimitCount {
		return nil, nil, core.NewBizError(1001006014, "总限购数量超出")
	}

	return activity, prod, nil
}

func (s *combinationRecordService) CreateCombinationRecord(ctx context.Context, record *promotion.PromotionCombinationRecord) (int64, error) {
	err := s.q.Transaction(func(tx *query.Query) error {
		if err := tx.PromotionCombinationRecord.WithContext(ctx).Create(record); err != nil {
			return err
		}
		// Update Head Status if joining
		if record.HeadID > 0 {
			return s.updateCombinationRecordWhenCreate(ctx, tx, record.HeadID)
		}
		return nil
	})
	return record.ID, err
}

// updateCombinationRecordWhenCreate 更新拼团记录状态
func (s *combinationRecordService) updateCombinationRecordWhenCreate(ctx context.Context, tx *query.Query, headID int64) error {
	// 1. Get Head
	head, err := tx.PromotionCombinationRecord.WithContext(ctx).Where(tx.PromotionCombinationRecord.ID.Eq(headID)).First()
	if err != nil {
		return err
	}
	// 2. Get Members
	members, err := tx.PromotionCombinationRecord.WithContext(ctx).Where(tx.PromotionCombinationRecord.HeadID.Eq(headID)).Find()
	if err != nil {
		return err
	}

	// 3. Get Activity for UserSize
	activity, err := s.activitySvc.GetCombinationActivity(ctx, head.ActivityID)
	if err != nil {
		return err
	}

	// 4. Update
	totalCount := 1 + len(members) // Head + Members
	isFull := totalCount >= activity.UserSize

	updates := make([]*promotion.PromotionCombinationRecord, 0, totalCount)
	updates = append(updates, head)
	updates = append(updates, members...)

	now := time.Now()
	for _, r := range updates {
		r.UserCount = totalCount
		if isFull {
			r.Status = 1 // Success
			r.EndTime = now
		}
		if _, err := tx.PromotionCombinationRecord.WithContext(ctx).Where(tx.PromotionCombinationRecord.ID.Eq(r.ID)).Updates(r); err != nil {
			return err
		}
	}
	// TODO: Send Success Notification if isFull
	return nil
}

// GetCombinationRecordPageAdmin 获得拼团记录分页 (Admin)
func (s *combinationRecordService) GetCombinationRecordPageAdmin(ctx context.Context, req *req.CombinationRecordPageReq) (*core.PageResult[*promotion.PromotionCombinationRecord], error) {
	q := s.q.PromotionCombinationRecord
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
	return &core.PageResult[*promotion.PromotionCombinationRecord]{List: list, Total: total}, nil
}
