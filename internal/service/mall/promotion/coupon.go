package promotion

import (
	"context"
	"strings"
	"time"

	"github.com/samber/lo"
	promotion3 "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/mall/promotion"
	promotion2 "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/app/mall/promotion"
	"github.com/wxlbd/ruoyi-mall-go/internal/consts"
	member_model "github.com/wxlbd/ruoyi-mall-go/internal/model/member"
	"github.com/wxlbd/ruoyi-mall-go/internal/model/promotion"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/member"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

type CouponService struct {
	q           *query.Query
	userService *member.MemberUserService
}

func NewCouponService(q *query.Query, userService *member.MemberUserService) *CouponService {
	return &CouponService{
		q:           q,
		userService: userService,
	}
}

// CreateCouponTemplate 创建优惠券模板 (Admin)
func (s *CouponService) CreateCouponTemplate(ctx context.Context, req *promotion3.CouponTemplateCreateReq) (int64, error) {
	t := &promotion.PromotionCouponTemplate{
		Name:               req.Name,
		Description:        req.Description,
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
		DiscountLimitPrice: req.DiscountLimitPrice,
	}
	err := s.q.PromotionCouponTemplate.WithContext(ctx).Create(t)
	return t.ID, err
}

// UpdateCouponTemplate 更新优惠券模板 (Admin)
func (s *CouponService) UpdateCouponTemplate(ctx context.Context, req *promotion3.CouponTemplateUpdateReq) error {
	_, err := s.q.PromotionCouponTemplate.WithContext(ctx).Where(s.q.PromotionCouponTemplate.ID.Eq(req.ID)).Updates(promotion.PromotionCouponTemplate{
		Name:               req.Name,
		Description:        req.Description,
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
		DiscountLimitPrice: req.DiscountLimitPrice,
	})
	return err
}

// GetCouponTemplatePage 获得优惠券模板分页 (Admin)
func (s *CouponService) GetCouponTemplatePage(ctx context.Context, req *promotion3.CouponTemplatePageReq) (*pagination.PageResult[*promotion3.CouponTemplateResp], error) {
	q := s.q.PromotionCouponTemplate.WithContext(ctx)

	// 名称筛选 (去除前后空格)
	if req.Name != "" {
		trimmedName := strings.TrimSpace(req.Name)
		if trimmedName != "" {
			q = q.Where(s.q.PromotionCouponTemplate.Name.Like("%" + trimmedName + "%"))
		}
	}

	// 状态筛选
	if req.Status != nil {
		q = q.Where(s.q.PromotionCouponTemplate.Status.Eq(int(*req.Status)))
	}

	// 优惠类型筛选
	if req.DiscountType != nil {
		q = q.Where(s.q.PromotionCouponTemplate.DiscountType.Eq(*req.DiscountType))
	}

	// 领取类型筛选
	if len(req.CanTakeTypes) > 0 {
		q = q.Where(s.q.PromotionCouponTemplate.TakeType.In(req.CanTakeTypes...))
	}

	// 商品范围筛选
	if req.ProductScope != nil {
		q = q.Where(s.q.PromotionCouponTemplate.ProductScope.Eq(*req.ProductScope))
	}

	if req.ProductScopeValue != nil {
		// JSON_CONTAINS 需要原生 SQL，暂时跳过（前端未使用此筛选）
		// TODO: 如需支持，可用 Clauses 或 Scopes
	}

	// 创建时间范围筛选
	if len(req.CreateTime) == 2 && req.CreateTime[0] != nil && req.CreateTime[1] != nil {
		q = q.Where(s.q.PromotionCouponTemplate.CreateTime.Between(*req.CreateTime[0], *req.CreateTime[1]))
	}

	result, count, err := q.FindByPage(int((req.PageNo-1)*req.PageSize), int(req.PageSize))
	if err != nil {
		return nil, err
	}

	// 收集模板 IDs 用于批量查询统计
	templateIDs := make([]int64, len(result))
	for i, tmpl := range result {
		templateIDs[i] = tmpl.ID
	}

	// 批量查询每个模板的使用数量 (status = 2 表示已使用)
	useCountMap := make(map[int64]int)
	if len(templateIDs) > 0 {
		type UseCountResult struct {
			TemplateID int64 `gorm:"column:template_id"`
			UseCount   int   `gorm:"column:use_count"`
		}
		var useCountResults []UseCountResult
		c := s.q.PromotionCoupon
		err = c.WithContext(ctx).
			Select(c.TemplateID, c.TemplateID.Count().As("use_count")).
			Where(c.TemplateID.In(templateIDs...), c.Status.Eq(2)). // 2 = USED
			Group(c.TemplateID).
			Scan(&useCountResults)
		if err != nil {
			return nil, err
		}
		for _, r := range useCountResults {
			useCountMap[r.TemplateID] = r.UseCount
		}
	}

	// 转换为 Response DTO
	list := make([]*promotion3.CouponTemplateResp, len(result))
	for i, tmpl := range result {
		list[i] = s.convertTemplateToResp0(tmpl, useCountMap[tmpl.ID])
	}

	return &pagination.PageResult[*promotion3.CouponTemplateResp]{
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
		return errors.NewBizError(1006001000, "优惠券模板不存在")
	}
	if template == nil {
		return errors.NewBizError(1006001000, "优惠券模板不存在")
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
		return errors.NewBizError(1006001000, "优惠券模板不存在")
	}
	if template == nil {
		return errors.NewBizError(1006001000, "优惠券模板不存在")
	}

	// 删除
	_, err = t.WithContext(ctx).Where(t.ID.Eq(id)).Delete()
	return err
}

// GetCouponTemplate 获取优惠券模板详情 (Admin)
// 对应 Java: CouponTemplateService.getCouponTemplate
func (s *CouponService) GetCouponTemplate(ctx context.Context, id int64) (*promotion3.CouponTemplateResp, error) {
	t := s.q.PromotionCouponTemplate
	tmpl, err := t.WithContext(ctx).Where(t.ID.Eq(id)).First()
	if err != nil {
		return nil, err
	}
	if tmpl == nil {
		return nil, nil
	}
	return s.convertTemplateToResp(ctx, tmpl)
}

// GetCouponTemplateList 获取优惠券模板列表 (Admin)
// 对应 Java: CouponTemplateService.getCouponTemplateList(ids)
func (s *CouponService) GetCouponTemplateList(ctx context.Context, ids []int64) ([]*promotion3.CouponTemplateResp, error) {
	if len(ids) == 0 {
		return []*promotion3.CouponTemplateResp{}, nil
	}
	t := s.q.PromotionCouponTemplate
	list, err := t.WithContext(ctx).Where(t.ID.In(ids...)).Find()
	if err != nil {
		return nil, err
	}

	res := make([]*promotion3.CouponTemplateResp, 0, len(list))
	for _, item := range list {
		resp, _ := s.convertTemplateToResp(ctx, item)
		res = append(res, resp)
	}
	return res, nil
}

func (s *CouponService) convertTemplateToResp(ctx context.Context, tmpl *promotion.PromotionCouponTemplate) (*promotion3.CouponTemplateResp, error) {
	c := s.q.PromotionCoupon
	useCount, _ := c.WithContext(ctx).Where(c.TemplateID.Eq(tmpl.ID), c.Status.Eq(2)).Count()
	return s.convertTemplateToResp0(tmpl, int(useCount)), nil
}

func (s *CouponService) convertTemplateToResp0(tmpl *promotion.PromotionCouponTemplate, useCount int) *promotion3.CouponTemplateResp {
	return &promotion3.CouponTemplateResp{
		ID:                 tmpl.ID,
		Name:               tmpl.Name,
		Description:        tmpl.Description,
		Status:             tmpl.Status,
		TotalCount:         tmpl.TotalCount,
		TakeLimitCount:     tmpl.TakeLimitCount,
		TakeType:           tmpl.TakeType,
		UsePrice:           tmpl.UsePriceMin,
		ProductScope:       tmpl.ProductScope,
		ProductScopeValues: tmpl.ProductScopeValues,
		ValidityType:       tmpl.ValidityType,
		ValidStartTime:     tmpl.ValidStartTime,
		ValidEndTime:       tmpl.ValidEndTime,
		FixedStartTerm:     &tmpl.FixedStartTerm,
		FixedEndTerm:       &tmpl.FixedEndTerm,
		DiscountType:       tmpl.DiscountType,
		DiscountPercent:    &tmpl.DiscountPercent,
		DiscountPrice:      &tmpl.DiscountPrice,
		DiscountLimitPrice: &tmpl.DiscountLimitPrice,
		TakeCount:          tmpl.TakeCount,
		UseCount:           useCount,
		Creator:            tmpl.Creator,
		Updater:            tmpl.Updater,
		CreateTime:         tmpl.CreateTime,
		UpdateTime:         tmpl.UpdateTime,
	}
}

// GetCouponPage 获得优惠券分页 (Admin)
func (s *CouponService) GetCouponPage(ctx context.Context, req *promotion3.CouponPageReq) (*pagination.PageResult[*promotion3.CouponPageResp], error) {
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

	// 收集用户 IDs 用于批量查询
	userIDs := make([]int64, 0, len(result))
	for _, coupon := range result {
		userIDs = append(userIDs, coupon.UserID)
	}

	// 批量查询用户信息
	userMap := make(map[int64]*member_model.MemberUser)
	if len(userIDs) > 0 {
		var err error
		userMap, err = s.userService.GetUserMap(ctx, userIDs)
		if err != nil {
			// 用户查询失败不影响主流程，只是 nickname 为空
		}
	}

	// 转换为 Response DTO
	list := make([]*promotion3.CouponPageResp, len(result))
	for i, coupon := range result {
		list[i] = &promotion3.CouponPageResp{
			// RespVO 字段
			ID:         coupon.ID,
			CreateTime: coupon.CreateTime,

			// BaseVO 字段 - 基本信息
			TemplateID: coupon.TemplateID,
			Name:       coupon.Name,
			Status:     coupon.Status,

			// BaseVO 字段 - 领取情况
			UserID:   coupon.UserID,
			TakeType: coupon.TakeType,

			// BaseVO 字段 - 使用规则
			UsePrice:           coupon.UsePrice,
			ValidStartTime:     coupon.ValidStartTime,
			ValidEndTime:       coupon.ValidEndTime,
			ProductScope:       coupon.ProductScope,
			ProductScopeValues: coupon.ProductScopeValues,

			// BaseVO 字段 - 使用效果
			DiscountType:       coupon.DiscountType,
			DiscountPercent:    &coupon.DiscountPercent,
			DiscountPrice:      &coupon.DiscountPrice,
			DiscountLimitPrice: &coupon.DiscountLimitPrice,

			// BaseVO 字段 - 使用情况
			UseOrderID: &coupon.UseOrderID,
			UseTime:    coupon.UseTime,

			// PageItemRespVO 字段 - 关联字段
			Nickname: "", // 默认空字符串
		}
		// 填充用户昵称
		if user, ok := userMap[coupon.UserID]; ok {
			list[i].Nickname = user.Nickname
		}
	}

	return &pagination.PageResult[*promotion3.CouponPageResp]{
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
		return errors.NewBizError(1006002000, "优惠券不存在")
	}
	if coupon == nil {
		return errors.NewBizError(1006002000, "优惠券不存在")
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
		return errors.NewBizError(1006001000, "优惠券模板不存在")
	}
	if template == nil {
		return errors.NewBizError(1006001000, "优惠券模板不存在")
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
			TemplateID:         templateId,
			Name:               template.Name,
			UserID:             userId,
			Status:             consts.CouponStatusUnused, // 未使用
			UsePrice:           template.UsePriceMin,
			ValidStartTime:     validStartTime,
			ValidEndTime:       validEndTime,
			DiscountType:       template.DiscountType,
			DiscountPrice:      template.DiscountPrice,
			DiscountPercent:    template.DiscountPercent,
			DiscountLimitPrice: template.DiscountLimitPrice,
		}
		coupons = append(coupons, coupon)
	}

	// 4. 批量创建
	c := s.q.PromotionCoupon
	return c.WithContext(ctx).Create(coupons...)
}

// ========== App 端方法 ==========

// GetCouponTemplateForApp 获取单个优惠券模板 (App 端)
// 对齐 Java: AppCouponTemplateController.getCouponTemplate
func (s *CouponService) GetCouponTemplateForApp(ctx context.Context, id int64, userId int64) (*promotion2.AppCouponTemplateResp, error) {
	t := s.q.PromotionCouponTemplate
	template, err := t.WithContext(ctx).Where(t.ID.Eq(id)).First()
	if err != nil {
		return nil, nil // 返回 null 对齐 Java
	}

	// 处理是否可领取 (对齐 Java: couponService.getUserCanCanTakeMap)
	canTakeMap, err := s.GetUserCanTakeMap(ctx, userId, []int64{template.ID})
	if err != nil {
		return nil, err
	}

	return s.convertToAppCouponTemplateResp(template, canTakeMap[template.ID]), nil
}

// GetCouponTemplateListForApp 获取优惠券模板列表 (App 端)
// 对齐 Java: AppCouponTemplateController.getCouponTemplateList (带条件)
func (s *CouponService) GetCouponTemplateListForApp(ctx context.Context, spuId *int64, productScope *int, count int, userId int64) ([]*promotion2.AppCouponTemplateResp, error) {
	t := s.q.PromotionCouponTemplate
	q := t.WithContext(ctx).Where(t.Status.Eq(consts.CommonStatusEnable)) // 只查询启用状态
	q = q.Where(t.TakeType.Eq(consts.CouponTakeTypeUser))                 // 领取方式 = 直接领取 (CouponTakeTypeEnum.USER)

	if productScope != nil {
		q = q.Where(t.ProductScope.Eq(*productScope))
	}

	q = q.Order(t.ID.Desc()) // 按 ID 降序排列，对齐 Java 实现

	if count <= 0 {
		count = 10
	}
	q = q.Limit(count)

	templates, err := q.Find()
	if err != nil {
		return nil, err
	}

	// 获取用户是否可领取
	templateIds := make([]int64, len(templates))
	for i, tmpl := range templates {
		templateIds[i] = tmpl.ID
	}
	canTakeMap, err := s.GetUserCanTakeMap(ctx, userId, templateIds)
	if err != nil {
		return nil, err
	}

	// 转换响应
	result := make([]*promotion2.AppCouponTemplateResp, len(templates))
	for i, tmpl := range templates {
		result[i] = s.convertToAppCouponTemplateResp(tmpl, canTakeMap[tmpl.ID])
	}
	return result, nil
}

// GetCouponTemplateListByIdsForApp 按 ID 获取优惠券模板列表 (App 端)
// 对齐 Java: AppCouponTemplateController.getCouponTemplateList (按ids)
func (s *CouponService) GetCouponTemplateListByIdsForApp(ctx context.Context, ids []int64, userId int64) ([]*promotion2.AppCouponTemplateResp, error) {
	if len(ids) == 0 {
		return []*promotion2.AppCouponTemplateResp{}, nil
	}
	t := s.q.PromotionCouponTemplate
	templates, err := t.WithContext(ctx).Where(t.ID.In(ids...)).Find()
	if err != nil {
		return nil, err
	}

	// 获取用户是否可领取
	canTakeMap, err := s.GetUserCanTakeMap(ctx, userId, ids)
	if err != nil {
		return nil, err
	}

	// 转换响应
	result := make([]*promotion2.AppCouponTemplateResp, len(templates))
	for i, tmpl := range templates {
		result[i] = s.convertToAppCouponTemplateResp(tmpl, canTakeMap[tmpl.ID])
	}
	return result, nil
}

// GetCouponTemplatePageForApp 获取优惠券模板分页 (App 端)
// 对齐 Java: AppCouponTemplateController.getCouponTemplatePage
func (s *CouponService) GetCouponTemplatePageForApp(ctx context.Context, r *promotion2.AppCouponTemplatePageReq, userId int64) (*pagination.PageResult[*promotion2.AppCouponTemplateResp], error) {
	t := s.q.PromotionCouponTemplate
	q := t.WithContext(ctx).Where(t.Status.Eq(consts.CommonStatusEnable)) // 只查询启用状态
	q = q.Where(t.TakeType.Eq(consts.CouponTakeTypeUser))                 // 领取方式 = 直接领取

	if r.ProductScope != nil {
		q = q.Where(t.ProductScope.Eq(*r.ProductScope))
	}

	q = q.Order(t.ID.Desc()) // 按 ID 降序排列，对齐 Java 实现

	templates, count, err := q.FindByPage((r.PageNo-1)*r.PageSize, r.PageSize)
	if err != nil {
		return nil, err
	}

	// 获取用户是否可领取
	templateIds := lo.Map(templates, func(tmpl *promotion.PromotionCouponTemplate, _ int) int64 {
		return tmpl.ID
	})
	canTakeMap, err := s.GetUserCanTakeMap(ctx, userId, templateIds)
	if err != nil {
		return nil, err
	}

	// 转换响应
	result := lo.Map(templates, func(tmpl *promotion.PromotionCouponTemplate, _ int) *promotion2.AppCouponTemplateResp {
		return s.convertToAppCouponTemplateResp(tmpl, canTakeMap[tmpl.ID])
	})

	return &pagination.PageResult[*promotion2.AppCouponTemplateResp]{
		List:  result,
		Total: count,
	}, nil
}

// convertToAppCouponTemplateResp 转换模板 Model 到 App 响应 VO
func (s *CouponService) convertToAppCouponTemplateResp(t *promotion.PromotionCouponTemplate, canTake bool) *promotion2.AppCouponTemplateResp {
	// 处理ProductScopeValues，如果为空则返回空数组而不是null
	productScopeValues := t.ProductScopeValues
	if len(productScopeValues) == 0 {
		productScopeValues = []int64{}
	}

	// 处理时间戳：转换为毫秒时间戳或 null
	var validStartTime, validEndTime interface{}
	if t.ValidStartTime != nil {
		validStartTime = t.ValidStartTime.UnixMilli()
	}
	if t.ValidEndTime != nil {
		validEndTime = t.ValidEndTime.UnixMilli()
	}

	// 处理 fixedStartTerm 和 fixedEndTerm：当 validityType=1 时返回 null
	var fixedStartTerm, fixedEndTerm *int
	if t.ValidityType == 2 { // 2 = 领取后N天
		fixedStartTerm = &t.FixedStartTerm
		fixedEndTerm = &t.FixedEndTerm
	}

	return &promotion2.AppCouponTemplateResp{
		ID:                 t.ID,
		Name:               t.Name,
		Description:        t.Description,
		TotalCount:         t.TotalCount,
		TakeLimitCount:     t.TakeLimitCount,
		UsePrice:           t.UsePriceMin,
		ProductScope:       int(t.ProductScope),
		ProductScopeValues: productScopeValues,
		ValidityType:       t.ValidityType,
		ValidStartTime:     validStartTime,
		ValidEndTime:       validEndTime,
		FixedStartTerm:     fixedStartTerm,
		FixedEndTerm:       fixedEndTerm,
		DiscountType:       t.DiscountType,
		DiscountPercent:    t.DiscountPercent,
		DiscountPrice:      t.DiscountPrice,
		DiscountLimitPrice: t.DiscountLimitPrice,
		TakeCount:          t.TakeCount,
		CanTake:            canTake,
	}
}

// GetUserCanTakeMap 获取用户是否可领取某模板的 Map
// 对齐 Java: CouponService.getUserCanCanTakeMap
func (s *CouponService) GetUserCanTakeMap(ctx context.Context, userId int64, templateIds []int64) (map[int64]bool, error) {
	result := make(map[int64]bool)
	if userId == 0 || len(templateIds) == 0 {
		for _, id := range templateIds {
			result[id] = true // 未登录用户默认可领取
		}
		return result, nil
	}

	// 查询模板的领取限制
	t := s.q.PromotionCouponTemplate
	templates, err := t.WithContext(ctx).Where(t.ID.In(templateIds...)).Find()
	if err != nil {
		return nil, err
	}

	// 初始化结果，默认都为可领取
	for _, template := range templates {
		result[template.ID] = true
	}

	// 过滤出需要检查限制的模板（有限领数量且不为-1）
	var limitedTemplateIds []int64
	for _, template := range templates {
		if template.TakeLimitCount > 0 { // 只检查有限领数量且大于0的模板
			limitedTemplateIds = append(limitedTemplateIds, template.ID)
		} else {
			// 无限制或-1表示不限制，直接设为可领取
			result[template.ID] = true
		}
	}

	if len(limitedTemplateIds) == 0 {
		return result, nil
	}

	// 查询用户对每个模板已领取的数量
	c := s.q.PromotionCoupon
	coupons, err := c.WithContext(ctx).
		Where(c.UserID.Eq(userId)).
		Where(c.TemplateID.In(limitedTemplateIds...)).
		Find()
	if err != nil {
		return nil, err
	}

	// 统计每个模板的领取数量
	takeCountMap := make(map[int64]int)
	for _, coupon := range coupons {
		takeCountMap[coupon.TemplateID]++
	}

	// 检查是否超过限制
	for _, template := range templates {
		if template.TakeLimitCount > 0 { // 只检查有限领数量的模板
			takeCount := takeCountMap[template.ID]
			result[template.ID] = takeCount < template.TakeLimitCount
		}
	}

	return result, nil
}
