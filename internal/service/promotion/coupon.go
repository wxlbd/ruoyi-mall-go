package promotion

import (
	"backend-go/internal/api/req"
	"backend-go/internal/model/promotion"
	"backend-go/internal/pkg/core"
	"backend-go/internal/repo/query"
	"context"
	"time"
)

type CouponService struct {
	q *query.Query
}

func NewCouponService() *CouponService {
	return &CouponService{
		q: query.Q,
	}
}

// CreateCouponTemplate 创建优惠券模板 (Admin)
func (s *CouponService) CreateCouponTemplate(ctx context.Context, req *req.CouponTemplateCreateReq) (int64, error) {
	t := &promotion.PromotionCouponTemplate{
		Name:               req.Name,
		Status:             req.Status,
		TotalCount:         req.TotalCount,
		TakeLimitCount:     req.TakeLimitCount,
		TakeType:           req.TakeType,
		UsePriceMin:        req.UsePriceMin,
		ProductScope:       req.ProductScope,
		ProductScopeValues: req.ProductScopeValues,
		ValidityType:       req.ValidityType,
		ValidStartTime:     req.ValidStartTime,
		ValidEndTime:       req.ValidEndTime,
		FixedStartTerm:     req.FixedStartTerm,
		FixedEndTerm:       req.FixedEndTerm,
		DiscountType:       req.DiscountType,
		DiscountPrice:      req.DiscountPrice,
		DiscountPercent:    req.DiscountPercent,
		DiscountLimit:      req.DiscountLimit,
	}
	err := s.q.PromotionCouponTemplate.WithContext(ctx).Create(t)
	return t.ID, err
}

// UpdateCouponTemplate 更新优惠券模板 (Admin)
func (s *CouponService) UpdateCouponTemplate(ctx context.Context, req *req.CouponTemplateUpdateReq) error {
	_, err := s.q.PromotionCouponTemplate.WithContext(ctx).Where(s.q.PromotionCouponTemplate.ID.Eq(req.ID)).Updates(promotion.PromotionCouponTemplate{
		Name:               req.Name,
		Status:             req.Status,
		TotalCount:         req.TotalCount,
		TakeLimitCount:     req.TakeLimitCount,
		TakeType:           req.TakeType,
		UsePriceMin:        req.UsePriceMin,
		ProductScope:       req.ProductScope,
		ProductScopeValues: req.ProductScopeValues,
		ValidityType:       req.ValidityType,
		ValidStartTime:     req.ValidStartTime,
		ValidEndTime:       req.ValidEndTime,
		FixedStartTerm:     req.FixedStartTerm,
		FixedEndTerm:       req.FixedEndTerm,
		DiscountType:       req.DiscountType,
		DiscountPrice:      req.DiscountPrice,
		DiscountPercent:    req.DiscountPercent,
		DiscountLimit:      req.DiscountLimit,
	})
	return err
}

// GetCouponTemplatePage 获得优惠券模板分页 (Admin)
func (s *CouponService) GetCouponTemplatePage(ctx context.Context, req *req.CouponTemplatePageReq) (*core.PageResult[promotion.PromotionCouponTemplate], error) {
	q := s.q.PromotionCouponTemplate.WithContext(ctx)
	if req.Name != "" {
		q = q.Where(s.q.PromotionCouponTemplate.Name.Like("%" + req.Name + "%"))
	}
	if req.Status != nil {
		q = q.Where(s.q.PromotionCouponTemplate.Status.Eq(*req.Status))
	}

	result, count, err := q.FindByPage(int((req.PageNo-1)*req.PageSize), int(req.PageSize))
	if err != nil {
		return nil, err
	}

	list := make([]promotion.PromotionCouponTemplate, len(result))
	for i, v := range result {
		list[i] = *v
	}

	return &core.PageResult[promotion.PromotionCouponTemplate]{
		List:  list,
		Total: count,
	}, nil
}

// UpdateCouponTemplateStatus 更新优惠券模板状态 (Admin)
// 对应 Java: CouponTemplateService.updateCouponTemplateStatus
func (s *CouponService) UpdateCouponTemplateStatus(ctx context.Context, id int64, status int32) error {
	t := s.q.PromotionCouponTemplate
	// 校验存在
	template, err := t.WithContext(ctx).Where(t.ID.Eq(id)).First()
	if err != nil {
		return core.NewBizError(1006001000, "优惠券模板不存在")
	}
	if template == nil {
		return core.NewBizError(1006001000, "优惠券模板不存在")
	}

	// 更新状态
	_, err = t.WithContext(ctx).Where(t.ID.Eq(id)).Update(t.Status, status)
	return err
}

// DeleteCouponTemplate 删除优惠券模板 (Admin)
// 对应 Java: CouponTemplateService.deleteCouponTemplate
func (s *CouponService) DeleteCouponTemplate(ctx context.Context, id int64) error {
	t := s.q.PromotionCouponTemplate
	// 校验存在
	template, err := t.WithContext(ctx).Where(t.ID.Eq(id)).First()
	if err != nil {
		return core.NewBizError(1006001000, "优惠券模板不存在")
	}
	if template == nil {
		return core.NewBizError(1006001000, "优惠券模板不存在")
	}

	// 删除
	_, err = t.WithContext(ctx).Where(t.ID.Eq(id)).Delete()
	return err
}

// GetCouponTemplate 获取优惠券模板详情 (Admin)
// 对应 Java: CouponTemplateService.getCouponTemplate
func (s *CouponService) GetCouponTemplate(ctx context.Context, id int64) (*promotion.PromotionCouponTemplate, error) {
	t := s.q.PromotionCouponTemplate
	return t.WithContext(ctx).Where(t.ID.Eq(id)).First()
}

// GetCouponTemplateList 获取优惠券模板列表 (Admin)
// 对应 Java: CouponTemplateService.getCouponTemplateList(ids)
func (s *CouponService) GetCouponTemplateList(ctx context.Context, ids []int64) ([]*promotion.PromotionCouponTemplate, error) {
	if len(ids) == 0 {
		return []*promotion.PromotionCouponTemplate{}, nil
	}
	t := s.q.PromotionCouponTemplate
	return t.WithContext(ctx).Where(t.ID.In(ids...)).Find()
}

// GetCouponPage 获得优惠券分页 (Admin)
func (s *CouponService) GetCouponPage(ctx context.Context, req *req.CouponPageReq) (*core.PageResult[promotion.PromotionCoupon], error) {
	q := s.q.PromotionCoupon.WithContext(ctx)
	if req.UserID != nil {
		q = q.Where(s.q.PromotionCoupon.UserID.Eq(*req.UserID))
	}
	if req.Status != nil {
		q = q.Where(s.q.PromotionCoupon.Status.Eq(*req.Status))
	}

	result, count, err := q.FindByPage(int((req.PageNo-1)*req.PageSize), int(req.PageSize))
	if err != nil {
		return nil, err
	}

	list := make([]promotion.PromotionCoupon, len(result))
	for i, v := range result {
		list[i] = *v
	}

	return &core.PageResult[promotion.PromotionCoupon]{
		List:  list,
		Total: count,
	}, nil
}

// DeleteCoupon 删除/回收优惠券 (Admin)
// 对应 Java: CouponService.deleteCoupon
func (s *CouponService) DeleteCoupon(ctx context.Context, id int64) error {
	c := s.q.PromotionCoupon
	// 校验存在
	coupon, err := c.WithContext(ctx).Where(c.ID.Eq(id)).First()
	if err != nil {
		return core.NewBizError(1006002000, "优惠券不存在")
	}
	if coupon == nil {
		return core.NewBizError(1006002000, "优惠券不存在")
	}

	// 删除
	_, err = c.WithContext(ctx).Where(c.ID.Eq(id)).Delete()
	return err
}

// TakeCouponByAdmin 管理员发送优惠券给用户 (Admin)
// 对应 Java: CouponService.takeCouponByAdmin
func (s *CouponService) TakeCouponByAdmin(ctx context.Context, templateId int64, userIds []int64) error {
	if len(userIds) == 0 {
		return nil
	}

	// 1. 获取优惠券模板
	t := s.q.PromotionCouponTemplate
	template, err := t.WithContext(ctx).Where(t.ID.Eq(templateId)).First()
	if err != nil {
		return core.NewBizError(1006001000, "优惠券模板不存在")
	}
	if template == nil {
		return core.NewBizError(1006001000, "优惠券模板不存在")
	}

	// 2. 计算有效期
	var validStartTime, validEndTime time.Time
	if template.ValidStartTime != nil {
		validStartTime = *template.ValidStartTime
	} else {
		validStartTime = time.Now()
	}
	if template.ValidEndTime != nil {
		validEndTime = *template.ValidEndTime
	} else {
		// 默认30天后过期
		validEndTime = time.Now().AddDate(0, 0, 30)
	}

	// 3. 为每个用户创建优惠券
	coupons := make([]*promotion.PromotionCoupon, 0, len(userIds))
	for _, userId := range userIds {
		coupon := &promotion.PromotionCoupon{
			TemplateID:      templateId,
			Name:            template.Name,
			UserID:          userId,
			Status:          1, // 未使用
			UsePriceMin:     template.UsePriceMin,
			ValidStartTime:  validStartTime,
			ValidEndTime:    validEndTime,
			DiscountType:    template.DiscountType,
			DiscountPrice:   template.DiscountPrice,
			DiscountPercent: template.DiscountPercent,
			DiscountLimit:   template.DiscountLimit,
		}
		coupons = append(coupons, coupon)
	}

	// 4. 批量创建
	c := s.q.PromotionCoupon
	return c.WithContext(ctx).Create(coupons...)
}
