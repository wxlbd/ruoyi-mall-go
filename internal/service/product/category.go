package product

import (
	"backend-go/internal/api/req"
	"backend-go/internal/api/resp"
	"backend-go/internal/model/product"
	"backend-go/internal/pkg/core"
	"backend-go/internal/repo/query"
	"context"

	"github.com/samber/lo"
)

type ProductCategoryService struct {
	q *query.Query
}

func NewProductCategoryService(q *query.Query) *ProductCategoryService {
	return &ProductCategoryService{q: q}
}

// CreateCategory 创建商品分类
func (s *ProductCategoryService) CreateCategory(ctx context.Context, req *req.ProductCategoryCreateReq) (int64, error) {
	// 校验父分类
	if err := s.validateParentCategory(ctx, req.ParentID); err != nil {
		return 0, err
	}

	category := &product.ProductCategory{
		ParentID: req.ParentID,
		Name:     req.Name,
		PicURL:   req.PicURL,
		Sort:     req.Sort,
		Status:   req.Status,
	}
	err := s.q.ProductCategory.WithContext(ctx).Create(category)
	return category.ID, err
}

// UpdateCategory 更新商品分类
func (s *ProductCategoryService) UpdateCategory(ctx context.Context, req *req.ProductCategoryUpdateReq) error {
	// 校验存在
	if err := s.ValidateCategory(ctx, req.ID); err != nil {
		return err
	}
	// 校验父分类
	if err := s.validateParentCategory(ctx, req.ParentID); err != nil {
		return err
	}
	// 校验不能设置自己为父分类
	if req.ID == req.ParentID {
		return core.NewBizError(1006001004, "不能设置自己为父分类")
	}

	u := s.q.ProductCategory
	_, err := u.WithContext(ctx).Where(u.ID.Eq(req.ID)).Updates(&product.ProductCategory{
		ParentID: req.ParentID,
		Name:     req.Name,
		PicURL:   req.PicURL,
		Sort:     req.Sort,
		Status:   req.Status,
	})
	return err
}

// DeleteCategory 删除商品分类
func (s *ProductCategoryService) DeleteCategory(ctx context.Context, id int64) error {
	// 校验存在
	if err := s.ValidateCategory(ctx, id); err != nil {
		return err
	}
	// 校验是否有子分类
	count, err := s.q.ProductCategory.WithContext(ctx).Where(s.q.ProductCategory.ParentID.Eq(id)).Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return core.NewBizError(1006001001, "存在子分类，无法删除")
	}
	// 校验是否绑定了 SPU
	spuCount, err := s.q.ProductSpu.WithContext(ctx).Where(s.q.ProductSpu.CategoryID.Eq(id)).Count()
	if err != nil {
		return err
	}
	if spuCount > 0 {
		return core.NewBizError(1006001004, "存在商品绑定，无法删除")
	}

	_, err = s.q.ProductCategory.WithContext(ctx).Where(s.q.ProductCategory.ID.Eq(id)).Delete()
	return err
}

// GetCategory 获得商品分类
func (s *ProductCategoryService) GetCategory(ctx context.Context, id int64) (*resp.ProductCategoryResp, error) {
	u := s.q.ProductCategory
	category, err := u.WithContext(ctx).Where(u.ID.Eq(id)).First()
	if err != nil {
		return nil, nil
	}
	return s.convertResp(category), nil
}

// GetCategoryList 获得商品分类列表
func (s *ProductCategoryService) GetCategoryList(ctx context.Context, req *req.ProductCategoryListReq) ([]*resp.ProductCategoryResp, error) {
	u := s.q.ProductCategory
	q := u.WithContext(ctx)
	if req.Name != "" {
		q = q.Where(u.Name.Like("%" + req.Name + "%"))
	}
	if req.ParentID != nil {
		q = q.Where(u.ParentID.Eq(*req.ParentID))
	}
	if req.Status != nil {
		q = q.Where(u.Status.Eq(*req.Status))
	}
	list, err := q.Order(u.Sort.Desc(), u.ID.Asc()).Find() // Sort desc, ID asc
	if err != nil {
		return nil, err
	}
	return lo.Map(list, func(item *product.ProductCategory, _ int) *resp.ProductCategoryResp {
		return s.convertResp(item)
	}), nil
}

func (s *ProductCategoryService) validateParentCategory(ctx context.Context, parentId int64) error {
	if parentId == 0 {
		return nil
	}
	// 父分类必须存在
	u := s.q.ProductCategory
	parent, err := u.WithContext(ctx).Where(u.ID.Eq(parentId)).First()
	if err != nil {
		return core.NewBizError(1006001002, "父分类不存在")
	}
	// 父分类不能是二级分类 (即 parentId != 0) -> 意味着只能创建二级分类 (parentId 指向一级), 不能创建三级
	// Logic: If parent's ParentID is NOT 0, it means parent is ALREADY a child (Level 2).
	// So we cannot add a child to it.
	if parent.ParentID != 0 {
		return core.NewBizError(1006001003, "父分类不能是二级分类")
	}
	return nil
}

func (s *ProductCategoryService) ValidateCategory(ctx context.Context, id int64) error {
	u := s.q.ProductCategory
	count, err := u.WithContext(ctx).Where(u.ID.Eq(id)).Count()
	if err != nil {
		return err
	}
	if count == 0 {
		return core.NewBizError(1006001000, "分类不存在")
	}
	return nil
}

func (s *ProductCategoryService) convertResp(item *product.ProductCategory) *resp.ProductCategoryResp {
	return &resp.ProductCategoryResp{
		ID:        item.ID,
		ParentID:  item.ParentID,
		Name:      item.Name,
		PicURL:    item.PicURL,
		Sort:      item.Sort,
		Status:    item.Status,
		CreatedAt: item.CreatedAt,
	}
}
