package product

import (
	"context"

	product2 "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/mall/product"
	"github.com/wxlbd/ruoyi-mall-go/internal/model/product"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"

	"github.com/samber/lo"
)

type ProductBrandService struct {
	q *query.Query
}

func NewProductBrandService(q *query.Query) *ProductBrandService {
	return &ProductBrandService{q: q}
}

// CreateBrand 创建品牌
func (s *ProductBrandService) CreateBrand(ctx context.Context, req *product2.ProductBrandCreateReq) (int64, error) {
	// 校验名称唯一
	if err := s.validateBrandNameUnique(ctx, 0, req.Name); err != nil {
		return 0, err
	}

	brand := &product.ProductBrand{
		Name:        req.Name,
		PicURL:      req.PicURL,
		Sort:        req.Sort,
		Description: req.Description,
		Status:      req.Status,
	}
	err := s.q.ProductBrand.WithContext(ctx).Create(brand)
	return brand.ID, err
}

// UpdateBrand 更新品牌
func (s *ProductBrandService) UpdateBrand(ctx context.Context, req *product2.ProductBrandUpdateReq) error {
	// 校验存在
	if err := s.ValidateProductBrand(ctx, req.ID); err != nil {
		return err
	}
	// 校验名称唯一
	if err := s.validateBrandNameUnique(ctx, req.ID, req.Name); err != nil {
		return err
	}
	_, err := s.q.ProductBrand.WithContext(ctx).
		Where(s.q.ProductBrand.ID.Eq(req.ID)).
		UpdateColumns(map[string]any{
			"name":        req.Name,
			"pic_url":     req.PicURL,
			"sort":        req.Sort,
			"description": req.Description,
			"status":      req.Status,
		})
	return err
}

// DeleteBrand 删除品牌
func (s *ProductBrandService) DeleteBrand(ctx context.Context, id int64) error {
	// 校验存在
	if err := s.ValidateProductBrand(ctx, id); err != nil {
		return err
	}
	_, err := s.q.ProductBrand.WithContext(ctx).Where(s.q.ProductBrand.ID.Eq(id)).Delete()
	return err
}

// GetBrand 获得品牌
func (s *ProductBrandService) GetBrand(ctx context.Context, id int64) (*product2.ProductBrandResp, error) {
	u := s.q.ProductBrand
	brand, err := u.WithContext(ctx).Where(u.ID.Eq(id)).First()
	if err != nil {
		return nil, nil // Or error if strict
	}
	return s.convertResp(brand), nil
}

// GetBrandPage 获得品牌分页
func (s *ProductBrandService) GetBrandPage(ctx context.Context, req *product2.ProductBrandPageReq) (*pagination.PageResult[*product2.ProductBrandResp], error) {
	u := s.q.ProductBrand
	q := u.WithContext(ctx)
	if req.Name != "" {
		q = q.Where(u.Name.Like("%" + req.Name + "%"))
	}
	if req.Status != nil {
		q = q.Where(u.Status.Eq(*req.Status))
	}

	list, total, err := q.Order(u.Sort.Asc()).FindByPage((req.PageNo-1)*req.PageSize, req.PageSize)
	if err != nil {
		return nil, err
	}

	resList := lo.Map(list, func(item *product.ProductBrand, _ int) *product2.ProductBrandResp {
		return s.convertResp(item)
	})
	return &pagination.PageResult[*product2.ProductBrandResp]{
		List:  resList,
		Total: total,
	}, nil
}

// GetBrandList 获得品牌列表
func (s *ProductBrandService) GetBrandList(ctx context.Context, req *product2.ProductBrandListReq) ([]*product2.ProductBrandResp, error) {
	u := s.q.ProductBrand
	q := u.WithContext(ctx)
	if req.Name != "" {
		q = q.Where(u.Name.Like("%" + req.Name + "%"))
	}
	list, err := q.Order(u.Sort.Asc()).Find()
	if err != nil {
		return nil, err
	}
	return lo.Map(list, func(item *product.ProductBrand, _ int) *product2.ProductBrandResp {
		return s.convertResp(item)
	}), nil
}

// ValidateProductBrand 校验品牌是否存在
func (s *ProductBrandService) ValidateProductBrand(ctx context.Context, id int64) error {
	u := s.q.ProductBrand
	count, err := u.WithContext(ctx).Where(u.ID.Eq(id)).Count()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.NewBizError(1006001000, "品牌不存在") // BRAND_NOT_EXISTS
	}
	return nil
}

// validateBrandNameUnique 校验品牌名称是否唯一
func (s *ProductBrandService) validateBrandNameUnique(ctx context.Context, id int64, name string) error {
	u := s.q.ProductBrand
	brand, err := u.WithContext(ctx).Where(u.Name.Eq(name), u.ID.Neq(id)).First()
	if err == nil && brand != nil {
		if id == 0 || brand.ID != id {
			return errors.NewBizError(1006001001, "品牌名称已存在") // BRAND_NAME_EXISTS
		}
	}
	return nil
}

func (s *ProductBrandService) convertResp(item *product.ProductBrand) *product2.ProductBrandResp {
	return &product2.ProductBrandResp{
		ID:          item.ID,
		Name:        item.Name,
		PicURL:      item.PicURL,
		Sort:        item.Sort,
		Description: item.Description,
		Status:      item.Status,
		CreateTime:  item.CreateTime,
	}
}
