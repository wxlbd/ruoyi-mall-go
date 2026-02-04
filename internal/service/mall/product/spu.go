package product

import (
	"context"

	product2 "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/mall/product"
	product3 "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/app/mall/product"
	"github.com/wxlbd/ruoyi-mall-go/internal/consts"
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/internal/model/product"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"

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

func mergeVirtualSalesCount(spu *product.ProductSpu) {
	if spu == nil {
		return
	}
	spu.SalesCount += spu.VirtualSalesCount
}

// CreateSpu 创建 SPU
func (s *ProductSpuService) CreateSpu(ctx context.Context, req *product2.ProductSpuSaveReq) (int64, error) {
	// 校验分类
	if err := s.categorySvc.ValidateCategory(ctx, req.CategoryID); err != nil {
		return 0, err
	}
	// 校验分类层级（只允许二级分类）
	if err := s.categorySvc.ValidateCategoryLevel(ctx, req.CategoryID); err != nil {
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
		Status:             1, // ✅ 对齐 Java: 默认上架 (ENABLE=1)
		SalesCount:         0, // ✅ 对齐 Java: 默认销量为0
		BrowseCount:        0, // ✅ 对齐 Java: 默认浏览量为0
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
func (s *ProductSpuService) UpdateSpu(ctx context.Context, req *product2.ProductSpuSaveReq) error {
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
		return product.ErrSpuNotRecycle // 使用商品模块错误码
	}

	return s.q.Transaction(func(tx *query.Query) error {
		if _, err := tx.ProductSpu.WithContext(ctx).Where(tx.ProductSpu.ID.Eq(id)).Delete(); err != nil {
			return err
		}
		return s.skuSvc.DeleteSkuBySpuId(ctx, id)
	})
}

// UpdateSpuStatus 更新 SPU 状态
func (s *ProductSpuService) UpdateSpuStatus(ctx context.Context, req *product2.ProductSpuUpdateStatusReq) error {
	// 校验存在
	if _, err := s.validateSpuExists(ctx, int64(req.ID)); err != nil {
		return err
	}
	return s.q.Transaction(func(tx *query.Query) error {
		_, err := tx.ProductSpu.WithContext(ctx).Where(tx.ProductSpu.ID.Eq(int64(req.ID))).Update(tx.ProductSpu.Status, int32(*req.Status))
		return err
	})
}

// UpdateBrowseCount 更新浏览量
func (s *ProductSpuService) UpdateBrowseCount(ctx context.Context, id int64, count int) error {
	_, err := s.q.ProductSpu.WithContext(ctx).Where(s.q.ProductSpu.ID.Eq(id)).
		Update(s.q.ProductSpu.BrowseCount, s.q.ProductSpu.BrowseCount.Add(count))
	return err
}

// GetSpuDetail 获得 SPU 详情 - 返回模型数据，VO组装在Handler层完成
func (s *ProductSpuService) GetSpuDetail(ctx context.Context, id int64) (*product.ProductSpu, []*product.ProductSku, error) {
	// 获得商品 SPU
	spu, err := s.q.ProductSpu.WithContext(ctx).Where(s.q.ProductSpu.ID.Eq(id)).First()
	if err != nil {
		return nil, nil, product.ErrSpuNotExists // 使用商品模块错误码
	}

	// 检查商品状态，对齐Java版本的ProductSpuStatusEnum.isEnable检查
	if spu.Status != consts.ProductSpuStatusEnable { // Status=1表示上架状态
		return nil, nil, product.ErrSpuNotEnable // 使用商品模块错误码
	}

	// 获得商品 SKU
	skus, err := s.skuSvc.GetSkuListBySpuId(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	// 合并销量：实际销量 + 虚拟销量，对齐Java版本逻辑
	mergeVirtualSalesCount(spu)

	return spu, skus, nil
}

// GetSpuPage 获得 SPU 分页
func (s *ProductSpuService) GetSpuPage(ctx context.Context, req *product2.ProductSpuPageReq) (*pagination.PageResult[*product2.ProductSpuResp], error) {
	u := s.q.ProductSpu
	q := u.WithContext(ctx)

	if req.TabType != nil {
		switch *req.TabType {
		case 0:
			// 出售中 (Status = 1，上架状态)
			q = q.Where(u.Status.Eq(consts.ProductSpuStatusEnable))
		case 1:
			// 仓库中 (Status = 0，下架状态)
			q = q.Where(u.Status.Eq(consts.ProductSpuStatusDisable))
		case 2:
			// 已售空 (Stock = 0)
			q = q.Where(u.Stock.Eq(0))
		case 3:
			// 警戒库存 (Stock <= 10)
			q = q.Where(u.Stock.Lte(10))
		case 4:
			// 回收站 (Status = -1)
			q = q.Where(u.Status.Eq(consts.ProductSpuStatusRecycle))
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

	resList := lo.Map(list, func(item *product.ProductSpu, _ int) *product2.ProductSpuResp {
		return s.convertResp(item, nil)
	})

	return &pagination.PageResult[*product2.ProductSpuResp]{
		List:  resList,
		Total: total,
	}, nil
}

// GetTabsCount 获得 SPU Tab 统计
func (s *ProductSpuService) GetTabsCount(ctx context.Context) (map[int]int64, error) {
	u := s.q.ProductSpu
	// 0: For Sale (Status = 1，上架状态)
	count0, _ := u.WithContext(ctx).Where(u.Status.Eq(consts.ProductSpuStatusEnable)).Count()
	// 1: In Warehouse (Status = 0，下架状态)
	count1, _ := u.WithContext(ctx).Where(u.Status.Eq(consts.ProductSpuStatusDisable)).Count()
	// 2: Sold Out (Stock = 0)
	count2, _ := u.WithContext(ctx).Where(u.Stock.Eq(0)).Count()
	// 3: Alert Stock (Stock <= 10，对齐 Java)
	count3, _ := u.WithContext(ctx).Where(u.Stock.Lte(10)).Count()
	// 4: Recycle (Status = -1)
	count4, _ := u.WithContext(ctx).Where(u.Status.Eq(consts.ProductSpuStatusRecycle)).Count()

	return map[int]int64{
		0: count0,
		1: count1,
		2: count2,
		3: count3,
		4: count4,
	}, nil
}

// GetSpuPageForApp 获得商品 SPU 分页 (App)
func (s *ProductSpuService) GetSpuPageForApp(ctx context.Context, req *product3.AppProductSpuPageReq) (*pagination.PageResult[*product.ProductSpu], error) {
	u := s.q.ProductSpu
	q := u.WithContext(ctx).Where(u.Status.Eq(1)) // 上架状态 Status=1

	// 处理分类查询：如果指定了分类ID，则包含其子分类
	if req.CategoryID != nil && *req.CategoryID > 0 {
		categoryIds, err := s.categorySvc.GetCategoryAndChildrenIds(ctx, *req.CategoryID)
		if err != nil {
			return nil, err
		}
		if len(categoryIds) > 0 {
			q = q.Where(u.CategoryID.In(categoryIds...))
		}
	}
	if req.Keyword != nil && *req.Keyword != "" {
		q = q.Where(u.Name.Like("%" + *req.Keyword + "%"))
	}
	// Sort
	if req.SortField != nil && *req.SortField == "sales_count" {
		if req.SortAsc != nil && *req.SortAsc {
			q = q.Order(u.SalesCount.Asc())
		} else {
			q = q.Order(u.SalesCount.Desc())
		}
	} else if req.SortField != nil && *req.SortField == "price" {
		if req.SortAsc != nil && *req.SortAsc {
			q = q.Order(u.Price.Asc())
		} else {
			q = q.Order(u.Price.Desc())
		}
	} else {
		q = q.Order(u.Sort.Desc(), u.ID.Desc())
	}

	list, total, err := q.FindByPage((req.PageNo-1)*req.PageSize, req.PageSize)
	if err != nil {
		return nil, err
	}
	for _, spu := range list {
		mergeVirtualSalesCount(spu)
	}
	return &pagination.PageResult[*product.ProductSpu]{
		List:  list,
		Total: total,
	}, nil
}

// GetSpuList 获得 SPU 列表 (Simple) - 对齐Java版本逻辑
func (s *ProductSpuService) GetSpuList(ctx context.Context, ids []int64) ([]*product2.ProductSpuResp, error) {
	if len(ids) == 0 {
		return []*product2.ProductSpuResp{}, nil
	}

	// 查询SPU数据
	list, err := s.q.ProductSpu.WithContext(ctx).Where(s.q.ProductSpu.ID.In(ids...)).Find()
	if err != nil {
		return nil, err
	}

	// 创建ID到SPU的映射，用于保持返回顺序
	spuMap := make(map[int64]*product.ProductSpu)
	for _, spu := range list {
		spuMap[spu.ID] = spu
	}

	// 按照输入ID顺序构建结果，对齐Java版本的convertList(ids, spuMap::get)逻辑
	result := make([]*product2.ProductSpuResp, 0, len(ids))
	for _, id := range ids {
		if spu, exists := spuMap[id]; exists {
			// 合并销量：实际销量 + 虚拟销量，对齐Java版本逻辑
			mergeVirtualSalesCount(spu)
			result = append(result, s.convertResp(spu, nil))
		}
	}

	return result, nil
}

// GetSpuSimpleList 获得 SPU 精简列表
func (s *ProductSpuService) GetSpuSimpleList(ctx context.Context) ([]*product2.ProductSpuSimpleResp, error) {
	list, err := s.q.ProductSpu.WithContext(ctx).Where(s.q.ProductSpu.Status.Eq(0)).Order(s.q.ProductSpu.Sort.Desc()).Find()
	if err != nil {
		return nil, err
	}
	return lo.Map(list, func(item *product.ProductSpu, _ int) *product2.ProductSpuSimpleResp {
		return &product2.ProductSpuSimpleResp{
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
			Update(s.q.ProductSpu.Stock, s.q.ProductSpu.Stock.Add(int(incr)))
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
		return nil, product.ErrSpuNotExists // 使用商品模块错误码
	}
	return spu, nil
}

// initSpuFromSkus 计算 SPU 价格库存
func (s *ProductSpuService) initSpuFromSkus(spu *product.ProductSpu, skus []*product2.ProductSkuSaveReq) {
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

func (s *ProductSpuService) convertResp(spu *product.ProductSpu, skus []*product.ProductSku) *product2.ProductSpuResp {
	skuResps := make([]*product2.ProductSkuResp, 0)
	if len(skus) > 0 {
		skuResps = lo.Map(skus, func(item *product.ProductSku, _ int) *product2.ProductSkuResp {
			return s.skuSvc.convertSkuResp(item)
		})
	}

	// 确保数组字段返回[]而不是null，对齐Java版本处理逻辑
	sliderPicURLs := spu.SliderPicURLs
	if sliderPicURLs == nil {
		sliderPicURLs = []string{}
	}

	deliveryTypes := spu.DeliveryTypes
	if deliveryTypes == nil {
		deliveryTypes = []int{}
	}

	return &product2.ProductSpuResp{
		ID:                 spu.ID,
		Name:               spu.Name,
		Keyword:            spu.Keyword,
		Introduction:       spu.Introduction, // 确保返回完整介绍文本而不是空字符串
		Description:        spu.Description,
		CategoryID:         spu.CategoryID, // 确保返回正确分类ID而不是0
		BrandID:            spu.BrandID,
		PicURL:             spu.PicURL,
		SliderPicURLs:      sliderPicURLs, // 确保返回[]而不是null
		Sort:               spu.Sort,
		Status:             spu.Status,
		SpecType:           bool(spu.SpecType),
		Price:              spu.Price,
		MarketPrice:        spu.MarketPrice,
		CostPrice:          spu.CostPrice,
		Stock:              spu.Stock,
		DeliveryTypes:      deliveryTypes, // 确保返回[]而不是null
		DeliveryTemplateID: spu.DeliveryTemplateID,
		GiveIntegral:       spu.GiveIntegral,
		SubCommissionType:  bool(spu.SubCommissionType),
		SalesCount:         spu.SalesCount, // 注意：在GetSpuList中已合并虚拟销量
		VirtualSalesCount:  spu.VirtualSalesCount,
		BrowseCount:        spu.BrowseCount,
		CreateTime:         spu.CreateTime,
		Skus:               skuResps,
	}
}
