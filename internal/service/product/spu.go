package product

import (
	"backend-go/internal/api/req"
	"backend-go/internal/api/resp"
	"backend-go/internal/model"
	"backend-go/internal/model/product"
	"backend-go/internal/pkg/core"
	"backend-go/internal/repo/query"
	"context"

	"github.com/samber/lo"
)

type ProductSpuService struct {
	q           *query.Query
	skuSvc      *ProductSkuService
	brandSvc    *ProductBrandService
	categorySvc *ProductCategoryService
}

func NewProductSpuService(q *query.Query, skuSvc *ProductSkuService, brandSvc *ProductBrandService, categorySvc *ProductCategoryService) *ProductSpuService {
	s := &ProductSpuService{
		q:           q,
		skuSvc:      skuSvc,
		brandSvc:    brandSvc,
		categorySvc: categorySvc,
	}
	s.skuSvc.SetSpuService(s)
	return s
}

// CreateSpu 创建 SPU
func (s *ProductSpuService) CreateSpu(ctx context.Context, req *req.ProductSpuSaveReq) (int64, error) {
	// 校验分类
	if err := s.categorySvc.ValidateCategory(ctx, req.CategoryID); err != nil {
		return 0, err
	}
	// 校验品牌
	if err := s.brandSvc.ValidateProductBrand(ctx, req.BrandID); err != nil {
		return 0, err
	}
	// 校验 SKU
	if err := s.skuSvc.ValidateSkuList(ctx, req.Skus, *req.SpecType); err != nil {
		return 0, err
	}

	spu := &product.ProductSpu{
		Name:               req.Name,
		Keyword:            req.Keyword,
		Introduction:       req.Introduction,
		Description:        req.Description,
		CategoryID:         req.CategoryID,
		BrandID:            req.BrandID,
		PicURL:             req.PicURL,
		SliderPicURLs:      req.SliderPicURLs,
		Sort:               req.Sort,
		SpecType:           model.BitBool(*req.SpecType),
		DeliveryTypes:      req.DeliveryTypes,
		DeliveryTemplateID: req.DeliveryTemplateID,
		GiveIntegral:       req.GiveIntegral,
		SubCommissionType:  model.BitBool(*req.SubCommissionType),
		VirtualSalesCount:  req.VirtualSalesCount,
		Status:             0, // Default to 0? Or from req? Java defaults to ENABLE if not set, logic is in initSpuFromSkus
	}

	// 初始化 SPU 信息 (价格、库存等)
	s.initSpuFromSkus(spu, req.Skus)

	// 事务执行
	err := s.q.Transaction(func(tx *query.Query) error {
		if err := tx.ProductSpu.WithContext(ctx).Create(spu); err != nil {
			return err
		}
		if err := s.skuSvc.CreateSkuList(ctx, spu.ID, req.Skus); err != nil {
			return err
		}
		return nil
	})
	return spu.ID, err
}

// UpdateSpu 更新 SPU
func (s *ProductSpuService) UpdateSpu(ctx context.Context, req *req.ProductSpuSaveReq) error {
	// 校验存在
	spu, err := s.validateSpuExists(ctx, req.ID)
	if err != nil {
		return err
	}
	// 校验分类、品牌
	if err := s.categorySvc.ValidateCategory(ctx, req.CategoryID); err != nil {
		return err
	}
	if err := s.brandSvc.ValidateProductBrand(ctx, req.BrandID); err != nil {
		return err
	}
	// 校验 SKU
	if err := s.skuSvc.ValidateSkuList(ctx, req.Skus, *req.SpecType); err != nil {
		return err
	}

	updateSpu := &product.ProductSpu{
		ID:                 req.ID,
		Name:               req.Name,
		Keyword:            req.Keyword,
		Introduction:       req.Introduction,
		Description:        req.Description,
		CategoryID:         req.CategoryID,
		BrandID:            req.BrandID,
		PicURL:             req.PicURL,
		SliderPicURLs:      req.SliderPicURLs,
		Sort:               req.Sort,
		SpecType:           model.BitBool(*req.SpecType),
		DeliveryTypes:      req.DeliveryTypes,
		DeliveryTemplateID: req.DeliveryTemplateID,
		GiveIntegral:       req.GiveIntegral,
		SubCommissionType:  model.BitBool(*req.SubCommissionType),
		VirtualSalesCount:  req.VirtualSalesCount,
		Status:             spu.Status, // Keep status
	}
	s.initSpuFromSkus(updateSpu, req.Skus)

	return s.q.Transaction(func(tx *query.Query) error {
		if _, err := tx.ProductSpu.WithContext(ctx).Where(tx.ProductSpu.ID.Eq(req.ID)).Updates(updateSpu); err != nil {
			return err
		}
		return s.skuSvc.UpdateSkuList(ctx, req.ID, req.Skus)
	})
}

// DeleteSpu 删除 SPU
func (s *ProductSpuService) DeleteSpu(ctx context.Context, id int64) error {
	// 校验存在
	spu, err := s.validateSpuExists(ctx, id)
	if err != nil {
		return err
	}
	// 校验状态 (只有回收站可以删除)
	if spu.Status != -1 { // RECYCLE_BIN
		return core.NewBizError(1006000004, "商品必须是回收站状态才能删除") // SPU_NOT_RECYCLE
	}

	return s.q.Transaction(func(tx *query.Query) error {
		if _, err := tx.ProductSpu.WithContext(ctx).Where(tx.ProductSpu.ID.Eq(id)).Delete(); err != nil {
			return err
		}
		return s.skuSvc.DeleteSkuBySpuId(ctx, id)
	})
}

// UpdateSpuStatus 更新 SPU 状态
func (s *ProductSpuService) UpdateSpuStatus(ctx context.Context, req *req.ProductSpuUpdateStatusReq) error {
	if _, err := s.validateSpuExists(ctx, req.ID); err != nil {
		return err
	}
	_, err := s.q.ProductSpu.WithContext(ctx).Where(s.q.ProductSpu.ID.Eq(req.ID)).Update(s.q.ProductSpu.Status, req.Status)
	return err
}

// UpdateBrowseCount 更新浏览量
func (s *ProductSpuService) UpdateBrowseCount(ctx context.Context, id int64, count int) error {
	_, err := s.q.ProductSpu.WithContext(ctx).Where(s.q.ProductSpu.ID.Eq(id)).
		Update(s.q.ProductSpu.BrowseCount, s.q.ProductSpu.BrowseCount.Add(count))
	return err
}

// GetSpuDetail 获得 SPU 详情
func (s *ProductSpuService) GetSpuDetail(ctx context.Context, id int64) (*resp.ProductSpuResp, error) {
	spu, err := s.q.ProductSpu.WithContext(ctx).Where(s.q.ProductSpu.ID.Eq(id)).First()
	if err != nil {
		return nil, nil // Return nil if not found, or error
	}

	skus, err := s.skuSvc.GetSkuListBySpuId(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.convertResp(spu, skus), nil
}

// GetSpuPage 获得 SPU 分页
func (s *ProductSpuService) GetSpuPage(ctx context.Context, req *req.ProductSpuPageReq) (*core.PageResult[*resp.ProductSpuResp], error) {
	u := s.q.ProductSpu
	q := u.WithContext(ctx)

	if req.TabType != nil {
		switch *req.TabType {
		case 0:
			// 出售中 (Status = 0)
			q = q.Where(u.Status.Eq(0))
		case 1:
			// 仓库中 (Status = 1)
			q = q.Where(u.Status.Eq(1))
		case 2:
			// 已售空 (Stock = 0)
			q = q.Where(u.Stock.Eq(0))
		case 3:
			// 警戒库存 (Stock <= 10)
			q = q.Where(u.Stock.Lte(10))
		case 4:
			// 回收站 (Status = -1)
			q = q.Where(u.Status.Eq(-1))
		}
	}

	if req.Name != "" {
		q = q.Where(u.Name.Like("%" + req.Name + "%"))
	}
	if req.CategoryID > 0 {
		q = q.Where(u.CategoryID.Eq(req.CategoryID))
	}

	list, total, err := q.Order(u.Sort.Desc(), u.ID.Desc()).FindByPage((req.PageNo-1)*req.PageSize, req.PageSize)
	if err != nil {
		return nil, err
	}

	resList := lo.Map(list, func(item *product.ProductSpu, _ int) *resp.ProductSpuResp {
		return s.convertResp(item, nil)
	})

	return &core.PageResult[*resp.ProductSpuResp]{
		List:  resList,
		Total: total,
	}, nil
}

// GetTabsCount 获得 SPU Tab 统计
func (s *ProductSpuService) GetTabsCount(ctx context.Context) (map[int]int64, error) {
	u := s.q.ProductSpu
	// Simple count implementation
	// 0: For Sale
	count0, _ := u.WithContext(ctx).Where(u.Status.Eq(0)).Count()
	// 1: In Warehouse
	count1, _ := u.WithContext(ctx).Where(u.Status.Eq(1)).Count()
	// 2: Sold Out
	count2, _ := u.WithContext(ctx).Where(u.Stock.Eq(0)).Count()
	// 3: Alert
	count3, _ := u.WithContext(ctx).Where(u.Stock.Lt(10)).Count() // Mock alert
	// 4: Recycle
	count4, _ := u.WithContext(ctx).Where(u.Status.Eq(-1)).Count()

	return map[int]int64{
		0: count0,
		1: count1,
		2: count2,
		3: count3,
		4: count4,
	}, nil
}

// GetSpuList 获得 SPU 列表 (Simple)
func (s *ProductSpuService) GetSpuList(ctx context.Context, ids []int64) ([]*resp.ProductSpuResp, error) {
	if len(ids) == 0 {
		return []*resp.ProductSpuResp{}, nil
	}
	list, err := s.q.ProductSpu.WithContext(ctx).Where(s.q.ProductSpu.ID.In(ids...)).Find()
	if err != nil {
		return nil, err
	}
	return lo.Map(list, func(item *product.ProductSpu, _ int) *resp.ProductSpuResp {
		return s.convertResp(item, nil)
	}), nil
}

// GetSpuListByStatus (Ref for Simple List)
func (s *ProductSpuService) GetSpuSimpleList(ctx context.Context) ([]*resp.ProductSpuSimpleResp, error) {
	list, err := s.q.ProductSpu.WithContext(ctx).Where(s.q.ProductSpu.Status.Eq(0)).Order(s.q.ProductSpu.Sort.Desc()).Find()
	if err != nil {
		return nil, err
	}
	return lo.Map(list, func(item *product.ProductSpu, _ int) *resp.ProductSpuSimpleResp {
		return &resp.ProductSpuSimpleResp{
			ID:          item.ID,
			Name:        item.Name,
			PicURL:      item.PicURL,
			Price:       item.Price,
			MarketPrice: item.MarketPrice,
			CostPrice:   item.CostPrice,
			Stock:       item.Stock,
		}
	}), nil
}

// UpdateSpuStock 更新 SPU 库存
func (s *ProductSpuService) UpdateSpuStock(ctx context.Context, stockIncr map[int64]int) error {
	for spuID, incr := range stockIncr {
		if incr == 0 {
			continue
		}
		// Update stock
		// Note: We don't strictly check SPU stock >= 0 here because it's an aggregate.
		// SKU level check is the authority.
		_, err := s.q.ProductSpu.WithContext(ctx).Where(s.q.ProductSpu.ID.Eq(spuID)).
			Update(s.q.ProductSpu.Stock, s.q.ProductSpu.Stock.Add(incr))
		if err != nil {
			return err
		}
	}
	return nil
}

// GetSpuCountByCategoryId 获得分类下的 SPU 数量
func (s *ProductSpuService) GetSpuCountByCategoryId(ctx context.Context, categoryID int64) (int64, error) {
	return s.q.ProductSpu.WithContext(ctx).Where(s.q.ProductSpu.CategoryID.Eq(categoryID)).Count()
}

// GetSpu 获得 SPU (Model)
func (s *ProductSpuService) GetSpu(ctx context.Context, id int64) (*product.ProductSpu, error) {
	return s.q.ProductSpu.WithContext(ctx).Where(s.q.ProductSpu.ID.Eq(id)).First()
}

// Internal Helpers

func (s *ProductSpuService) validateSpuExists(ctx context.Context, id int64) (*product.ProductSpu, error) {
	spu, err := s.q.ProductSpu.WithContext(ctx).Where(s.q.ProductSpu.ID.Eq(id)).First()
	if err != nil {
		return nil, core.NewBizError(1006000002, "商品不存在") // SPU_NOT_EXISTS
	}
	return spu, nil
}

// initSpuFromSkus 计算 SPU 价格库存
func (s *ProductSpuService) initSpuFromSkus(spu *product.ProductSpu, skus []*req.ProductSkuSaveReq) {
	if len(skus) == 0 {
		return
	}

	minPrice := skus[0].Price
	minMarketPrice := skus[0].MarketPrice
	minCostPrice := skus[0].CostPrice
	totalStock := 0

	for _, sku := range skus {
		if sku.Price < minPrice {
			minPrice = sku.Price
		}
		if sku.MarketPrice < minMarketPrice {
			minMarketPrice = sku.MarketPrice
		}
		if sku.CostPrice < minCostPrice {
			minCostPrice = sku.CostPrice
		}
		totalStock += sku.Stock
	}

	spu.Price = minPrice
	spu.MarketPrice = minMarketPrice
	spu.CostPrice = minCostPrice
	spu.Stock = totalStock
}

func (s *ProductSpuService) convertResp(spu *product.ProductSpu, skus []*product.ProductSku) *resp.ProductSpuResp {
	skuResps := make([]*resp.ProductSkuResp, 0)
	if len(skus) > 0 {
		skuResps = lo.Map(skus, func(item *product.ProductSku, _ int) *resp.ProductSkuResp {
			return s.convertSkuResp(item)
		})
	}

	return &resp.ProductSpuResp{
		ID:                 spu.ID,
		Name:               spu.Name,
		Keyword:            spu.Keyword,
		Introduction:       spu.Introduction,
		Description:        spu.Description,
		CategoryID:         spu.CategoryID,
		BrandID:            spu.BrandID,
		PicURL:             spu.PicURL,
		SliderPicURLs:      spu.SliderPicURLs,
		Sort:               spu.Sort,
		Status:             spu.Status,
		SpecType:           bool(spu.SpecType),
		Price:              spu.Price,
		MarketPrice:        spu.MarketPrice,
		CostPrice:          spu.CostPrice,
		Stock:              spu.Stock,
		DeliveryTypes:      spu.DeliveryTypes,
		DeliveryTemplateID: spu.DeliveryTemplateID,
		GiveIntegral:       spu.GiveIntegral,
		SubCommissionType:  bool(spu.SubCommissionType),
		SalesCount:         spu.SalesCount,
		VirtualSalesCount:  spu.VirtualSalesCount,
		BrowseCount:        spu.BrowseCount,
		CreatedAt:          spu.CreatedAt,
		Skus:               skuResps,
	}
}

func (s *ProductSpuService) convertSkuResp(sku *product.ProductSku) *resp.ProductSkuResp {
	properties := make([]resp.ProductSkuPropertyResp, len(sku.Properties))
	for i, p := range sku.Properties {
		properties[i] = resp.ProductSkuPropertyResp{
			PropertyID:   p.PropertyID,
			PropertyName: p.PropertyName,
			ValueID:      p.ValueID,
			ValueName:    p.ValueName,
		}
	}
	return &resp.ProductSkuResp{
		ID:                   sku.ID,
		SpuID:                sku.SpuID,
		Properties:           properties,
		Price:                sku.Price,
		MarketPrice:          sku.MarketPrice,
		CostPrice:            sku.CostPrice,
		BarCode:              sku.BarCode,
		PicURL:               sku.PicURL,
		Stock:                sku.Stock,
		Weight:               sku.Weight,
		Volume:               sku.Volume,
		FirstBrokeragePrice:  sku.FirstBrokeragePrice,
		SecondBrokeragePrice: sku.SecondBrokeragePrice,
		SalesCount:           sku.SalesCount,
	}
}
