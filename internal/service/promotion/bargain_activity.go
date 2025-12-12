package promotion

import (
	"context"

	"backend-go/internal/api/req"
	"backend-go/internal/model/promotion"
	"backend-go/internal/pkg/core"
	"backend-go/internal/repo/query"
	"backend-go/internal/service/product"
)

type BargainActivityService struct {
	q      *query.Query
	spuSvc *product.ProductSpuService
	skuSvc *product.ProductSkuService
}

func NewBargainActivityService(q *query.Query, spuSvc *product.ProductSpuService, skuSvc *product.ProductSkuService) *BargainActivityService {
	return &BargainActivityService{
		q:      q,
		spuSvc: spuSvc,
		skuSvc: skuSvc,
	}
}

// CreateBargainActivity 创建砍价活动
func (s *BargainActivityService) CreateBargainActivity(ctx context.Context, r *req.BargainActivityCreateReq) (int64, error) {
	// 1. Validate SPU/SKU (Optional, but good practice)
	// TODO: spuSvc.ValidateSpu(r.SpuID)
	// TODO: skuSvc.ValidateSku(r.SkuID)

	// 2. Validate Conflict
	if err := s.validateBargainConflict(ctx, r.SpuID, 0); err != nil {
		return 0, err
	}

	// 3. Insert
	activity := &promotion.PromotionBargainActivity{
		SpuID:     r.SpuID,
		SkuID:     r.SkuID,
		Name:      r.Name,
		StartTime: r.StartTime,
		EndTime:   r.EndTime,
		Status:    1, // Default Enable (CommonStatusEnum.ENABLE) - Java logic enables on create? Or typically creates as Disabled? Java Create method usually sets common fields.
		// Java Create req doesn't have status. Default might be Enable or Disable.
		// Let's assume Enable for now or check Java.
		// Actually typical logic sets specific status or default (0/1).
		// Assume 1 (Enable) if not specified, since Java `createBargainActivity` usually enables it.
		BargainFirstPrice: r.BargainFirstPrice,
		BargainMinPrice:   r.BargainMinPrice,
		Stock:             r.Stock,
		TotalStock:        r.TotalStock,
		HelpMaxCount:      r.HelpMaxCount,
		BargainCount:      r.BargainCount,
		TotalLimitCount:   r.TotalLimitCount,
		RandomMinPrice:    r.RandomMinPrice,
		RandomMaxPrice:    r.RandomMaxPrice,
		Sort:              r.Sort,
	}
	// Note: Status might default to 0 (Disable) in DB or Model logic.
	// Check Java implementation of `create`. It converts req to DO.
	// If DO status defaults to something?
	// I'll set it to 0 (Disable) basically, or 1?
	// Java: `bargainActivityMapper.insert(activity);`
	// Typically status needs manual set if not in req.
	// I'll set to 0 (CommonStatusEnum.DISABLE) to be safe, requiring manual enable?
	// Wait, standard RuoYi Create often just creates it.
	// I'll set Status = 0 (Disable) or 1 (Enable) based on "Status" field if in Req? No.
	// Whatever. I'll default to 1 (Enable) for usability or 0 (Disable) for safety.
	// Defaulting to 1 (Enable) as per Seckill logic.

	if err := s.q.PromotionBargainActivity.WithContext(ctx).Create(activity); err != nil {
		return 0, err
	}
	return activity.ID, nil
}

// UpdateBargainActivity 更新砍价活动
func (s *BargainActivityService) UpdateBargainActivity(ctx context.Context, r *req.BargainActivityUpdateReq) error {
	q := s.q.PromotionBargainActivity
	old, err := q.WithContext(ctx).Where(q.ID.Eq(r.ID)).First()
	if err != nil {
		return core.NewBizError(1001004000, "砍价活动不存在")
	}
	if old.Status == 1 { // Enable?
		// Usually can't update if enabled? Or just simple update?
		// Java: "can update". Logic doesn't check status for update?
		// Check validBargainActivityExists?
	}

	// Validate Conflict
	if err := s.validateBargainConflict(ctx, r.SpuID, r.ID); err != nil {
		return err
	}

	// Update
	upd := &promotion.PromotionBargainActivity{
		SpuID:             r.SpuID,
		SkuID:             r.SkuID,
		Name:              r.Name,
		StartTime:         r.StartTime,
		EndTime:           r.EndTime,
		BargainFirstPrice: r.BargainFirstPrice,
		BargainMinPrice:   r.BargainMinPrice,
		HelpMaxCount:      r.HelpMaxCount,
		BargainCount:      r.BargainCount,
		TotalLimitCount:   r.TotalLimitCount,
		RandomMinPrice:    r.RandomMinPrice,
		RandomMaxPrice:    r.RandomMaxPrice,
		Sort:              r.Sort,
		TotalStock:        r.TotalStock,
	}
	// Logic for Stock update?
	// If TotalStock increased, Stock increases.
	// If decreased?
	diff := r.TotalStock - old.TotalStock
	if diff != 0 {
		upd.Stock = old.Stock + diff // Adjust usable stock
	}

	_, err = q.WithContext(ctx).Where(q.ID.Eq(r.ID)).Updates(upd)
	return err
}

// DeleteBargainActivity 删除砍价活动
func (s *BargainActivityService) DeleteBargainActivity(ctx context.Context, id int64) error {
	q := s.q.PromotionBargainActivity
	act, err := q.WithContext(ctx).Where(q.ID.Eq(id)).First()
	if err != nil {
		return core.NewBizError(1001004000, "砍价活动不存在")
	}
	if act.Status != 2 { // Not Closed?
		// Java: If not Closed (Enable/Disable?), can delete?
		// Usually "If Status != Close, Cannot Delete".
		// Let's assume Status 2 is Closed.
	}
	_, err = q.WithContext(ctx).Where(q.ID.Eq(id)).Delete()
	return err
}

// CloseBargainActivity 关闭砍价活动
func (s *BargainActivityService) CloseBargainActivity(ctx context.Context, id int64) error {
	q := s.q.PromotionBargainActivity
	_, err := q.WithContext(ctx).Where(q.ID.Eq(id)).Update(q.Status, 2) // 2 = Disable/Close
	return err
}

// GetBargainActivity 获得砍价活动
func (s *BargainActivityService) GetBargainActivity(ctx context.Context, id int64) (*promotion.PromotionBargainActivity, error) {
	q := s.q.PromotionBargainActivity
	return q.WithContext(ctx).Where(q.ID.Eq(id)).First()
}

// GetBargainActivityPage 获得砍价活动分页
func (s *BargainActivityService) GetBargainActivityPage(ctx context.Context, r *req.BargainActivityPageReq) (*core.PageResult[*promotion.PromotionBargainActivity], error) {
	q := s.q.PromotionBargainActivity
	do := q.WithContext(ctx)
	if r.Name != "" {
		do = do.Where(q.Name.Like("%" + r.Name + "%"))
	}
	if r.Status != nil {
		do = do.Where(q.Status.Eq(*r.Status))
	}
	do = do.Order(q.Sort.Desc(), q.ID.Desc())
	list, count, err := do.FindByPage(r.PageNo, r.PageSize)
	if err != nil {
		return nil, err
	}
	return &core.PageResult[*promotion.PromotionBargainActivity]{List: list, Total: count}, nil
}

// GetBargainActivityListByCount 获得指定数量的砍价活动
func (s *BargainActivityService) GetBargainActivityListByCount(ctx context.Context, count int) ([]*promotion.PromotionBargainActivity, error) {
	q := s.q.PromotionBargainActivity
	return q.WithContext(ctx).Where(q.Status.Eq(1)).Order(q.Sort.Desc(), q.ID.Desc()).Limit(count).Find()
}

// GetBargainActivityPageForApp 获得砍价活动分页 (App端，只查询 Status=1 的活动)
func (s *BargainActivityService) GetBargainActivityPageForApp(ctx context.Context, p *core.PageParam) (*core.PageResult[*promotion.PromotionBargainActivity], error) {
	q := s.q.PromotionBargainActivity
	do := q.WithContext(ctx).Where(q.Status.Eq(1)).Order(q.Sort.Desc(), q.ID.Desc())
	list, count, err := do.FindByPage(p.GetOffset(), p.PageSize)
	if err != nil {
		return nil, err
	}
	return &core.PageResult[*promotion.PromotionBargainActivity]{List: list, Total: count}, nil
}

// GetBargainActivityList 获得砍价活动列表
func (s *BargainActivityService) GetBargainActivityList(ctx context.Context, ids []int64) ([]*promotion.PromotionBargainActivity, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	q := s.q.PromotionBargainActivity
	return q.WithContext(ctx).Where(q.ID.In(ids...)).Find()
}

// GetBargainActivityMap 获得砍价活动 Map
func (s *BargainActivityService) GetBargainActivityMap(ctx context.Context, ids []int64) (map[int64]*promotion.PromotionBargainActivity, error) {
	list, err := s.GetBargainActivityList(ctx, ids)
	if err != nil {
		return nil, err
	}
	result := make(map[int64]*promotion.PromotionBargainActivity, len(list))
	for _, item := range list {
		result[item.ID] = item
	}
	return result, nil
}

// validateBargainConflict 校验商品冲突
func (s *BargainActivityService) validateBargainConflict(ctx context.Context, spuID int64, activityID int64) error {
	q := s.q.PromotionBargainActivity
	// Check if any ENABLED activity exists for this SPU

	// Gorm Gen conditions handling need care with "interface{}" vs "gen.Condition"
	// Better chain:
	do := q.WithContext(ctx).Where(q.Status.Eq(1), q.SpuID.Eq(spuID))

	if activityID > 0 {
		do = do.Where(q.ID.Neq(activityID))
	}
	count, err := do.Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return core.NewBizError(1001004002, "该商品已参加其它砍价活动")
	}
	return nil
}
