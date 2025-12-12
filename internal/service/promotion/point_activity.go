package promotion

import (
	"context"
	"time"

	"backend-go/internal/api/req"
	"backend-go/internal/model/product"
	"backend-go/internal/model/promotion"
	"backend-go/internal/pkg/core"
	"backend-go/internal/repo/query"
	productSvc "backend-go/internal/service/product"

	"github.com/samber/lo"
)

type PointActivityService struct {
	q      *query.Query
	spuSvc *productSvc.ProductSpuService
	skuSvc *productSvc.ProductSkuService
}

func NewPointActivityService(spuSvc *productSvc.ProductSpuService, skuSvc *productSvc.ProductSkuService) *PointActivityService {
	return &PointActivityService{
		q:      query.Q,
		spuSvc: spuSvc,
		skuSvc: skuSvc,
	}
}

// CreatePointActivity 创建积分商城活动
// 对应 Java: PointActivityServiceImpl.createPointActivity
func (s *PointActivityService) CreatePointActivity(ctx context.Context, req *req.PointActivityCreateReq) (int64, error) {
	// 1.1 校验商品是否存在
	if err := s.validateProductExists(ctx, req.SpuID, req.Products); err != nil {
		return 0, err
	}
	// 1.2 校验商品是否已经参加别的活动
	if err := s.validatePointActivityProductConflicts(ctx, 0, req.SpuID); err != nil {
		return 0, err
	}

	t := &promotion.PromotionPointActivity{
		SpuID:      req.SpuID,
		Status:     req.Status,
		Remark:     req.Remark,
		Sort:       req.Sort,
		Stock:      req.Stock,
		TotalStock: req.Stock, // 初始化时总库存等于库存
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	err := s.q.Transaction(func(tx *query.Query) error {
		// 1. 创建活动
		if err := tx.PromotionPointActivity.WithContext(ctx).Create(t); err != nil {
			return err
		}

		// 2. 创建活动商品
		products := make([]*promotion.PromotionPointProduct, len(req.Products))
		for i, p := range req.Products {
			products[i] = &promotion.PromotionPointProduct{
				ActivityID:     t.ID,
				SpuID:          t.SpuID,
				SkuID:          p.SkuID,
				Count:          p.Count,
				Point:          p.Point,
				Price:          p.Price,
				Stock:          p.Stock,
				ActivityStatus: t.Status,
			}
		}
		return tx.PromotionPointProduct.WithContext(ctx).Create(products...)
	})
	return t.ID, err
}

// UpdatePointActivity 更新积分商城活动
// 对应 Java: PointActivityServiceImpl.updatePointActivity
func (s *PointActivityService) UpdatePointActivity(ctx context.Context, req *req.PointActivityUpdateReq) error {
	// 1.1 校验存在
	activity, err := s.validatePointActivityExists(ctx, req.ID)
	if err != nil {
		return err
	}
	if activity.Status == 0 { // DISABLE
		return core.NewBizError(1006003001, "积分商城活动已关闭") // POINT_ACTIVITY_UPDATE_FAIL_STATUS_CLOSED
	}

	// 1.2 校验商品是否存在
	if err := s.validateProductExists(ctx, req.SpuID, req.Products); err != nil {
		return err
	}
	// 1.3 校验商品是否已经参加别的活动
	if err := s.validatePointActivityProductConflicts(ctx, req.ID, req.SpuID); err != nil {
		return err
	}

	return s.q.Transaction(func(tx *query.Query) error {
		// 2.1 更新活动
		updateObj := map[string]interface{}{
			"spu_id": req.SpuID,
			"status": req.Status,
			"remark": req.Remark,
			"sort":   req.Sort,
			"stock":  req.Stock,
		}
		if req.Stock > activity.TotalStock {
			updateObj["total_stock"] = req.Stock
		}

		if _, err := tx.PromotionPointActivity.WithContext(ctx).Where(tx.PromotionPointActivity.ID.Eq(req.ID)).Updates(updateObj); err != nil {
			return err
		}

		// 2.2 更新活动商品 (Diff Logic)
		return s.updatePointProduct(ctx, tx, req.ID, req.SpuID, req.Status, req.Products)
	})
}

// updatePointProduct 更新积分商品 (Diff Logic)
func (s *PointActivityService) updatePointProduct(ctx context.Context, tx *query.Query, activityID int64, spuID int64, activityStatus int, products []req.PointProductSaveReq) error {
	// 1. 查询旧的商品列表
	oldList, err := tx.PromotionPointProduct.WithContext(ctx).Where(tx.PromotionPointProduct.ActivityID.Eq(activityID)).Find()
	if err != nil {
		return err
	}
	oldMap := lo.KeyBy(oldList, func(item *promotion.PromotionPointProduct) int64 {
		return item.SkuID
	})

	toInsert := make([]*promotion.PromotionPointProduct, 0)
	toUpdate := make([]*promotion.PromotionPointProduct, 0)

	// Process New Items
	for _, p := range products {
		if oldItem, ok := oldMap[p.SkuID]; ok {
			// Update: Keep ID, update data
			updateItem := &promotion.PromotionPointProduct{
				ID:             oldItem.ID,
				ActivityID:     activityID,
				SpuID:          spuID, // Use the updated SpuID from the activity
				SkuID:          p.SkuID,
				Count:          p.Count,
				Point:          p.Point,
				Price:          p.Price,
				Stock:          p.Stock,
				ActivityStatus: activityStatus,
			}
			toUpdate = append(toUpdate, updateItem)
			delete(oldMap, p.SkuID) // Mark as processed
		} else {
			// Insert
			insertItem := &promotion.PromotionPointProduct{
				ActivityID:     activityID,
				SpuID:          spuID, // Use the updated SpuID from the activity
				SkuID:          p.SkuID,
				Count:          p.Count,
				Point:          p.Point,
				Price:          p.Price,
				Stock:          p.Stock,
				ActivityStatus: activityStatus,
			}
			toInsert = append(toInsert, insertItem)
		}
	}

	// Whatever is left in oldMap is Delete
	toDeleteIDs := lo.Map(lo.Values(oldMap), func(item *promotion.PromotionPointProduct, _ int) int64 {
		return item.ID
	})

	// Execute Batch Ops
	if len(toInsert) > 0 {
		if err := tx.PromotionPointProduct.WithContext(ctx).Create(toInsert...); err != nil {
			return err
		}
	}
	if len(toUpdate) > 0 {
		// GORM's Save method updates by primary key if it exists, or Updates(item)
		for _, item := range toUpdate {
			if _, err := tx.PromotionPointProduct.WithContext(ctx).Where(tx.PromotionPointProduct.ID.Eq(item.ID)).Updates(item); err != nil {
				return err
			}
		}
	}
	if len(toDeleteIDs) > 0 {
		if _, err := tx.PromotionPointProduct.WithContext(ctx).Where(tx.PromotionPointProduct.ID.In(toDeleteIDs...)).Delete(); err != nil {
			return err
		}
	}

	return nil
}

// ClosePointActivity 关闭积分商城活动
// 对应 Java: PointActivityServiceImpl.closePointActivity
func (s *PointActivityService) ClosePointActivity(ctx context.Context, id int64) error {
	// 校验存在
	activity, err := s.validatePointActivityExists(ctx, id)
	if err != nil {
		return err
	}
	if activity.Status == 0 { // already closed
		return core.NewBizError(1006003002, "积分商城活动已关闭") // POINT_ACTIVITY_CLOSE_FAIL_STATUS_CLOSED
	}

	return s.q.Transaction(func(tx *query.Query) error {
		// 1. 更新活动状态
		if _, err := tx.PromotionPointActivity.WithContext(ctx).Where(tx.PromotionPointActivity.ID.Eq(id)).Update(tx.PromotionPointActivity.Status, 0); err != nil {
			return err
		}
		// 2. 更新商品状态
		p := tx.PromotionPointProduct
		_, err := p.WithContext(ctx).Where(p.ActivityID.Eq(id)).Update(p.ActivityStatus, 0)
		return err
	})
}

// DeletePointActivity 删除积分商城活动
// 对应 Java: PointActivityServiceImpl.deletePointActivity
func (s *PointActivityService) DeletePointActivity(ctx context.Context, id int64) error {
	activity, err := s.validatePointActivityExists(ctx, id)
	if err != nil {
		return err
	}
	if activity.Status == 1 { // ENABLE
		return core.NewBizError(1006003003, "活动未关闭或未结束，不能删除") // POINT_ACTIVITY_DELETE_FAIL_STATUS_NOT_CLOSED_OR_END
	}

	// 逻辑删除
	return s.q.Transaction(func(tx *query.Query) error {
		if _, err := tx.PromotionPointActivity.WithContext(ctx).Where(tx.PromotionPointActivity.ID.Eq(id)).Delete(); err != nil {
			return err
		}
		p := tx.PromotionPointProduct
		_, err := p.WithContext(ctx).Where(p.ActivityID.Eq(id)).Delete()
		return err
	})
}

// GetPointActivity 获得积分商城活动
// 对应 Java: PointActivityServiceImpl.getPointActivity
func (s *PointActivityService) GetPointActivity(ctx context.Context, id int64) (*promotion.PromotionPointActivity, []*promotion.PromotionPointProduct, error) {
	t := s.q.PromotionPointActivity
	activity, err := t.WithContext(ctx).Where(t.ID.Eq(id)).First()
	if err != nil {
		return nil, nil, err
	}

	p := s.q.PromotionPointProduct
	products, err := p.WithContext(ctx).Where(p.ActivityID.Eq(id)).Find()
	if err != nil {
		return nil, nil, err
	}

	return activity, products, nil
}

// GetPointActivityPage 获得积分商城活动分页
// 对应 Java: PointActivityServiceImpl.getPointActivityPage
func (s *PointActivityService) GetPointActivityPage(ctx context.Context, req *req.PointActivityPageReq) (*core.PageResult[promotion.PromotionPointActivity], error) {
	q := s.q.PromotionPointActivity.WithContext(ctx)
	if req.Status != nil {
		q = q.Where(s.q.PromotionPointActivity.Status.Eq(int32(*req.Status)))
	}
	// TODO: 支持其他搜索条件 (Java只支持Status)

	result, count, err := q.FindByPage(int((req.PageNo-1)*req.PageSize), int(req.PageSize))
	if err != nil {
		return nil, err
	}

	list := make([]promotion.PromotionPointActivity, len(result))
	for i, v := range result {
		list[i] = *v
	}

	return &core.PageResult[promotion.PromotionPointActivity]{
		List:  list,
		Total: count,
	}, nil
}

// GetPointActivityListByIds 获得积分商城活动列表
func (s *PointActivityService) GetPointActivityListByIds(ctx context.Context, ids []int64) ([]*promotion.PromotionPointActivity, error) {
	if len(ids) == 0 {
		return []*promotion.PromotionPointActivity{}, nil
	}
	t := s.q.PromotionPointActivity
	return t.WithContext(ctx).Where(t.ID.In(ids...)).Find()
}

// GetPointProductListByActivityIds 获得积分商城活动商品列表
func (s *PointActivityService) GetPointProductListByActivityIds(ctx context.Context, activityIds []int64) ([]*promotion.PromotionPointProduct, error) {
	if len(activityIds) == 0 {
		return []*promotion.PromotionPointProduct{}, nil
	}
	p := s.q.PromotionPointProduct
	return p.WithContext(ctx).Where(p.ActivityID.In(activityIds...)).Find()
}

// validatePointActivityExists 校验积分商城活动是否存在
func (s *PointActivityService) validatePointActivityExists(ctx context.Context, id int64) (*promotion.PromotionPointActivity, error) {
	t := s.q.PromotionPointActivity
	activity, err := t.WithContext(ctx).Where(t.ID.Eq(id)).First()
	if err != nil {
		return nil, core.NewBizError(1006003000, "积分商城活动不存在") // POINT_ACTIVITY_NOT_EXISTS
	}
	if activity == nil {
		return nil, core.NewBizError(1006003000, "积分商城活动不存在")
	}
	return activity, nil
}

// validateProductExists 校验商品是否存在
func (s *PointActivityService) validateProductExists(ctx context.Context, spuID int64, products []req.PointProductSaveReq) error {
	// 1. 校验商品 spu 是否存在
	spu, err := s.spuSvc.GetSpu(ctx, spuID)
	if err != nil {
		return err
	}
	if spu == nil {
		return core.NewBizError(1006000002, "商品不存在") // SPU_NOT_EXISTS
	}

	// 2. 校验商品 sku 都存在
	skus, err := s.skuSvc.GetSkuListBySpuId(ctx, spuID)
	if err != nil {
		return err
	}
	skuMap := lo.KeyBy(skus, func(sku *product.ProductSku) int64 {
		return sku.ID
	})

	for _, p := range products {
		if _, ok := skuMap[p.SkuID]; !ok {
			return core.NewBizError(1006002002, "商品 SKU 不存在") // SKU_NOT_EXISTS
		}
	}
	return nil
}

// validatePointActivityProductConflicts 校验商品是否冲突
// 校验当前 SPU 是否已经参加了其他开启的积分商城活动
func (s *PointActivityService) validatePointActivityProductConflicts(ctx context.Context, id int64, spuID int64) error {
	t := s.q.PromotionPointActivity
	q := t.WithContext(ctx).Where(t.Status.Eq(1), t.SpuID.Eq(spuID)) // ENABLE and same SpuID
	if id > 0 {
		q = q.Where(t.ID.Neq(id)) // Exclude current activity for update operations
	}
	count, err := q.Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return core.NewBizError(1006003004, "该商品已经参加了其他积分活动") // POINT_ACTIVITY_PRODUCT_CONFLICTS
	}
	return nil
}

// UpdatePointStockDecr 扣减积分商城活动库存
// 对应 Java: PointActivityServiceImpl.updatePointStockDecr
func (s *PointActivityService) UpdatePointStockDecr(ctx context.Context, id int64, skuID int64, count int) error {
	// 1. 校验活动是否存在
	activity, products, err := s.GetPointActivity(ctx, id)
	if err != nil {
		return err
	}
	if activity == nil {
		return core.NewBizError(1006003000, "积分商城活动不存在")
	}
	if activity.Status != 1 { // ENABLE
		return core.NewBizError(1006003002, "积分商城活动已关闭")
	}

	// 2. 校验商品是否存在
	product, found := lo.Find(products, func(item *promotion.PromotionPointProduct) bool {
		return item.SkuID == skuID
	})
	if !found {
		return core.NewBizError(1006002002, "商品 SKU 不存在")
	}

	// 3. 校验库存是否充足
	if product.Stock < count {
		return core.NewBizError(1006003005, "积分商品库存不足") // POINT_ACTIVITY_STOCK_NOT_ENOUGH
	}

	// 4. 扣减库存
	return s.q.Transaction(func(tx *query.Query) error {
		// 4.1 扣减活动商品库存
		if _, err := tx.PromotionPointProduct.WithContext(ctx).
			Where(tx.PromotionPointProduct.ID.Eq(product.ID), tx.PromotionPointProduct.Stock.Gte(int32(count))).
			Update(tx.PromotionPointProduct.Stock, tx.PromotionPointProduct.Stock.Sub(int32(count))); err != nil {
			return err
		}
		// 4.2 扣减活动总库存
		if _, err := tx.PromotionPointActivity.WithContext(ctx).
			Where(tx.PromotionPointActivity.ID.Eq(activity.ID), tx.PromotionPointActivity.Stock.Gte(int32(count))).
			Update(tx.PromotionPointActivity.Stock, tx.PromotionPointActivity.Stock.Sub(int32(count))); err != nil {
			return err
		}
		return nil
	})
}

// UpdatePointStockIncr 增加积分商城活动库存
// 对应 Java: PointActivityServiceImpl.updatePointStockIncr
func (s *PointActivityService) UpdatePointStockIncr(ctx context.Context, id int64, skuID int64, count int) error {
	// 1. 校验活动是否存在
	activity, products, err := s.GetPointActivity(ctx, id)
	if err != nil {
		return err
	}
	if activity == nil {
		return core.NewBizError(1006003000, "积分商城活动不存在")
	}

	// 2. 校验商品是否存在
	product, found := lo.Find(products, func(item *promotion.PromotionPointProduct) bool {
		return item.SkuID == skuID
	})
	if !found {
		return core.NewBizError(1006002002, "商品 SKU 不存在")
	}

	// 3. 增加库存
	return s.q.Transaction(func(tx *query.Query) error {
		// 3.1 增加活动商品库存
		if _, err := tx.PromotionPointProduct.WithContext(ctx).
			Where(tx.PromotionPointProduct.ID.Eq(product.ID)).
			Update(tx.PromotionPointProduct.Stock, tx.PromotionPointProduct.Stock.Add(int32(count))); err != nil {
			return err
		}
		// 3.2 增加活动总库存
		if _, err := tx.PromotionPointActivity.WithContext(ctx).
			Where(tx.PromotionPointActivity.ID.Eq(activity.ID)).
			Update(tx.PromotionPointActivity.Stock, tx.PromotionPointActivity.Stock.Add(int32(count))); err != nil {
			return err
		}
		return nil
	})
}

// ValidateJoinPointActivity 校验是否参加积分商城活动
// 对应 Java: PointActivityServiceImpl.validateJoinPointActivity
// return: PointProduct activity
func (s *PointActivityService) ValidateJoinPointActivity(ctx context.Context, activityID int64, spuID int64, skuID int64, count int) (*promotion.PromotionPointProduct, error) {
	// 1. 校验活动是否存在
	activity, products, err := s.GetPointActivity(ctx, activityID)
	if err != nil {
		return nil, err
	}
	if activity == nil {
		return nil, core.NewBizError(1006003000, "积分商城活动不存在")
	}
	if activity.Status != 1 { // ENABLE
		return nil, core.NewBizError(1006003002, "积分商城活动已关闭")
	}

	// 2. 校验商品是否存在
	product, found := lo.Find(products, func(item *promotion.PromotionPointProduct) bool {
		return item.SkuID == skuID
	})
	if !found {
		return nil, core.NewBizError(1006002002, "商品 SKU 不存在")
	}
	if product.SpuID != spuID {
		return nil, core.NewBizError(1006002002, "商品 SPU 不匹配") // Should not happen if data integrity is good
	}

	// 3. 校验库存是否充足
	if product.Stock < count {
		return nil, core.NewBizError(1006003005, "积分商品库存不足")
	}

	return product, nil
}
