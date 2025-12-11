package trade

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"backend-go/internal/api/req"
	"backend-go/internal/api/resp"
	"backend-go/internal/model/trade"
	"backend-go/internal/pkg/core"
	"backend-go/internal/repo/query"
)

type DeliveryFreightTemplateService struct {
	q *query.Query
}

func NewDeliveryFreightTemplateService(q *query.Query) *DeliveryFreightTemplateService {
	return &DeliveryFreightTemplateService{q: q}
}

// CreateDeliveryFreightTemplate 创建运费模板
func (s *DeliveryFreightTemplateService) CreateDeliveryFreightTemplate(ctx context.Context, r *req.DeliveryFreightTemplateSaveReq) (int64, error) {
	template := &trade.TradeDeliveryFreightTemplate{
		Name:       r.Name,
		Type:       r.Type,
		ChargeMode: r.ChargeMode,
		Sort:       r.Sort,
		Status:     r.Status,
		Remark:     r.Remark,
	}

	err := s.q.Transaction(func(tx *query.Query) error {
		// 1. 保存模板
		if err := tx.TradeDeliveryFreightTemplate.WithContext(ctx).Create(template); err != nil {
			return err
		}

		// 2. 保存计费规则
		if len(r.Charges) > 0 {
			var charges []*trade.TradeDeliveryFreightTemplateCharge
			for _, chargeReq := range r.Charges {
				areaIDs := s.convertAreaIDsToString(chargeReq.AreaIDs)
				charges = append(charges, &trade.TradeDeliveryFreightTemplateCharge{
					TemplateID: template.ID,
					AreaIDs:    areaIDs,
					StartCount: chargeReq.StartCount,
					StartPrice: chargeReq.StartPrice,
					ExtraCount: chargeReq.ExtraCount,
					ExtraPrice: chargeReq.ExtraPrice,
				})
			}
			if err := tx.TradeDeliveryFreightTemplateCharge.WithContext(ctx).Create(charges...); err != nil {
				return err
			}
		}

		// 3. 保存包邮规则
		if len(r.Frees) > 0 {
			var frees []*trade.TradeDeliveryFreightTemplateFree
			for _, freeReq := range r.Frees {
				areaIDs := s.convertAreaIDsToString(freeReq.AreaIDs)
				frees = append(frees, &trade.TradeDeliveryFreightTemplateFree{
					TemplateID: template.ID,
					AreaIDs:    areaIDs,
					FreePrice:  freeReq.FreePrice,
					FreeCount:  freeReq.FreeCount,
				})
			}
			if err := tx.TradeDeliveryFreightTemplateFree.WithContext(ctx).Create(frees...); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return 0, err
	}
	return template.ID, nil
}

// UpdateDeliveryFreightTemplate 更新运费模板
func (s *DeliveryFreightTemplateService) UpdateDeliveryFreightTemplate(ctx context.Context, r *req.DeliveryFreightTemplateSaveReq) error {
	return s.q.Transaction(func(tx *query.Query) error {
		// 1. 更新模板
		if _, err := tx.TradeDeliveryFreightTemplate.WithContext(ctx).Where(tx.TradeDeliveryFreightTemplate.ID.Eq(r.ID)).Updates(map[string]interface{}{
			"name":        r.Name,
			"type":        r.Type,
			"charge_mode": r.ChargeMode,
			"sort":        r.Sort,
			"status":      r.Status,
			"remark":      r.Remark,
		}); err != nil {
			return err
		}

		// 2. 删除旧的计费规则
		if _, err := tx.TradeDeliveryFreightTemplateCharge.WithContext(ctx).Where(tx.TradeDeliveryFreightTemplateCharge.TemplateID.Eq(r.ID)).Delete(); err != nil {
			return err
		}

		// 3. 保存新的计费规则
		if len(r.Charges) > 0 {
			var charges []*trade.TradeDeliveryFreightTemplateCharge
			for _, chargeReq := range r.Charges {
				areaIDs := s.convertAreaIDsToString(chargeReq.AreaIDs)
				charges = append(charges, &trade.TradeDeliveryFreightTemplateCharge{
					TemplateID: r.ID,
					AreaIDs:    areaIDs,
					StartCount: chargeReq.StartCount,
					StartPrice: chargeReq.StartPrice,
					ExtraCount: chargeReq.ExtraCount,
					ExtraPrice: chargeReq.ExtraPrice,
				})
			}
			if err := tx.TradeDeliveryFreightTemplateCharge.WithContext(ctx).Create(charges...); err != nil {
				return err
			}
		}

		// 4. 删除旧的包邮规则
		if _, err := tx.TradeDeliveryFreightTemplateFree.WithContext(ctx).Where(tx.TradeDeliveryFreightTemplateFree.TemplateID.Eq(r.ID)).Delete(); err != nil {
			return err
		}

		// 5. 保存新的包邮规则
		if len(r.Frees) > 0 {
			var frees []*trade.TradeDeliveryFreightTemplateFree
			for _, freeReq := range r.Frees {
				areaIDs := s.convertAreaIDsToString(freeReq.AreaIDs)
				frees = append(frees, &trade.TradeDeliveryFreightTemplateFree{
					TemplateID: r.ID,
					AreaIDs:    areaIDs,
					FreePrice:  freeReq.FreePrice,
					FreeCount:  freeReq.FreeCount,
				})
			}
			if err := tx.TradeDeliveryFreightTemplateFree.WithContext(ctx).Create(frees...); err != nil {
				return err
			}
		}

		return nil
	})
}

// DeleteDeliveryFreightTemplate 删除运费模板
func (s *DeliveryFreightTemplateService) DeleteDeliveryFreightTemplate(ctx context.Context, id int64) error {
	return s.q.Transaction(func(tx *query.Query) error {
		// 删除模板
		if _, err := tx.TradeDeliveryFreightTemplate.WithContext(ctx).Where(tx.TradeDeliveryFreightTemplate.ID.Eq(id)).Delete(); err != nil {
			return err
		}
		// 删除计费规则
		if _, err := tx.TradeDeliveryFreightTemplateCharge.WithContext(ctx).Where(tx.TradeDeliveryFreightTemplateCharge.TemplateID.Eq(id)).Delete(); err != nil {
			return err
		}
		// 删除包邮规则
		if _, err := tx.TradeDeliveryFreightTemplateFree.WithContext(ctx).Where(tx.TradeDeliveryFreightTemplateFree.TemplateID.Eq(id)).Delete(); err != nil {
			return err
		}
		return nil
	})
}

// GetDeliveryFreightTemplate 获取运费模板详情
func (s *DeliveryFreightTemplateService) GetDeliveryFreightTemplate(ctx context.Context, id int64) (*resp.DeliveryFreightTemplateResp, error) {
	// 1. 获取模板
	template, err := s.q.TradeDeliveryFreightTemplate.WithContext(ctx).Where(s.q.TradeDeliveryFreightTemplate.ID.Eq(id)).First()
	if err != nil {
		return nil, err
	}

	result := &resp.DeliveryFreightTemplateResp{
		ID:         template.ID,
		Name:       template.Name,
		Type:       template.Type,
		ChargeMode: template.ChargeMode,
		Sort:       template.Sort,
		Status:     template.Status,
		Remark:     template.Remark,
		CreateTime: template.CreatedAt,
	}

	// 2. 获取计费规则
	charges, err := s.q.TradeDeliveryFreightTemplateCharge.WithContext(ctx).Where(s.q.TradeDeliveryFreightTemplateCharge.TemplateID.Eq(id)).Find()
	if err != nil {
		return nil, err
	}
	for _, charge := range charges {
		areaIDs := s.convertAreaIDsToIntSlice(charge.AreaIDs)
		result.Charges = append(result.Charges, resp.DeliveryFreightTemplateChargeResp{
			AreaIDs:    areaIDs,
			StartCount: charge.StartCount,
			StartPrice: charge.StartPrice,
			ExtraCount: charge.ExtraCount,
			ExtraPrice: charge.ExtraPrice,
		})
	}

	// 3. 获取包邮规则
	frees, err := s.q.TradeDeliveryFreightTemplateFree.WithContext(ctx).Where(s.q.TradeDeliveryFreightTemplateFree.TemplateID.Eq(id)).Find()
	if err != nil {
		return nil, err
	}
	for _, free := range frees {
		areaIDs := s.convertAreaIDsToIntSlice(free.AreaIDs)
		result.Frees = append(result.Frees, resp.DeliveryFreightTemplateFreeResp{
			AreaIDs:   areaIDs,
			FreePrice: free.FreePrice,
			FreeCount: free.FreeCount,
		})
	}

	return result, nil
}

// GetDeliveryFreightTemplatePage 获取运费模板分页
func (s *DeliveryFreightTemplateService) GetDeliveryFreightTemplatePage(ctx context.Context, r *req.DeliveryFreightTemplatePageReq) (*core.PageResult[*trade.TradeDeliveryFreightTemplate], error) {
	q := s.q.TradeDeliveryFreightTemplate.WithContext(ctx)
	if r.Name != "" {
		q = q.Where(s.q.TradeDeliveryFreightTemplate.Name.Like("%" + r.Name + "%"))
	}
	if r.Status != nil {
		q = q.Where(s.q.TradeDeliveryFreightTemplate.Status.Eq(*r.Status))
	}

	pageNo := r.PageNo
	pageSize := r.PageSize
	if pageNo <= 0 {
		pageNo = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (pageNo - 1) * pageSize

	total, err := q.Count()
	if err != nil {
		return nil, err
	}

	list, err := q.Order(s.q.TradeDeliveryFreightTemplate.Sort.Asc()).Offset(offset).Limit(pageSize).Find()
	if err != nil {
		return nil, err
	}

	return &core.PageResult[*trade.TradeDeliveryFreightTemplate]{
		List:  list,
		Total: total,
	}, nil
}

// GetSimpleDeliveryFreightTemplateList 获取所有运费模板精简列表
func (s *DeliveryFreightTemplateService) GetSimpleDeliveryFreightTemplateList(ctx context.Context) ([]*resp.SimpleDeliveryFreightTemplateResp, error) {
	list, err := s.q.TradeDeliveryFreightTemplate.WithContext(ctx).Order(s.q.TradeDeliveryFreightTemplate.Sort.Asc()).Find()
	if err != nil {
		return nil, err
	}
	var res []*resp.SimpleDeliveryFreightTemplateResp
	for _, item := range list {
		res = append(res, &resp.SimpleDeliveryFreightTemplateResp{
			ID:   item.ID,
			Name: item.Name,
		})
	}
	return res, nil
}

// 辅助方法: 转换 AreaIDs 数组为逗号分隔字符串
func (s *DeliveryFreightTemplateService) convertAreaIDsToString(ids []int) string {
	if len(ids) == 0 {
		return ""
	}
	// json marshal works, or string join
	// Standard RuoYi uses comma separated string often, or uses JSON if specified.
	// My model says comment "逗号分隔".
	var strIDs []string
	for _, id := range ids {
		strIDs = append(strIDs, strconv.Itoa(id))
	}
	return strings.Join(strIDs, ",")
}

// 辅助方法: 转换逗号分隔字符串为 AreaIDs 数组
func (s *DeliveryFreightTemplateService) convertAreaIDsToIntSlice(str string) []int {
	if str == "" {
		return []int{}
	}
	// Try parsing as JSON first in case it is stored as "[]"? No, let's assume comma separated.
	// Check if it starts with [
	if strings.HasPrefix(str, "[") {
		var ids []int
		json.Unmarshal([]byte(str), &ids)
		return ids
	}

	parts := strings.Split(str, ",")
	var ids []int
	for _, p := range parts {
		if id, err := strconv.Atoi(p); err == nil {
			ids = append(ids, id)
		}
	}
	return ids
}

// CalculateFreight 计算运费
func (s *DeliveryFreightTemplateService) CalculateFreight(ctx context.Context, templateID int64, areaID int, count int) (int, error) {
	if templateID == 0 {
		return 0, nil
	}
	template, err := s.q.TradeDeliveryFreightTemplate.WithContext(ctx).Where(s.q.TradeDeliveryFreightTemplate.ID.Eq(templateID)).First()
	if err != nil {
		return 0, err
	}
	if template == nil {
		return 0, nil
	}

	// 1. Check Free Shipping Rules
	frees, err := s.q.TradeDeliveryFreightTemplateFree.WithContext(ctx).Where(s.q.TradeDeliveryFreightTemplateFree.TemplateID.Eq(templateID)).Find()
	if err != nil {
		return 0, err
	}
	for _, free := range frees {
		ids := s.convertAreaIDsToIntSlice(free.AreaIDs)
		if core.IntSliceContains(ids, areaID) {
			if float64(count) >= free.FreeCount || (free.FreePrice > 0 && 0 >= free.FreePrice) { // Logic simplified for count only now
				return 0, nil
			}
		}
	}

	// 2. Check Charge Rules
	charges, err := s.q.TradeDeliveryFreightTemplateCharge.WithContext(ctx).Where(s.q.TradeDeliveryFreightTemplateCharge.TemplateID.Eq(templateID)).Find()
	if err != nil {
		return 0, err
	}
	// Find matching region rule, otherwise use default
	var matchCharge *trade.TradeDeliveryFreightTemplateCharge
	for _, charge := range charges {
		ids := s.convertAreaIDsToIntSlice(charge.AreaIDs)
		if core.IntSliceContains(ids, areaID) {
			matchCharge = charge
			break
		}
	}
	// If no specific match, try to find default rule (usually 1 record with empty areaIDs)
	if matchCharge == nil {
		for _, charge := range charges {
			if charge.AreaIDs == "" {
				matchCharge = charge
				break
			}
		}
	}

	if matchCharge == nil {
		return 0, nil
	}

	// Calculate
	price := matchCharge.StartPrice
	if float64(count) > matchCharge.StartCount {
		extraCount := float64(count) - matchCharge.StartCount
		// Ceiling division for extra units
		if matchCharge.ExtraCount == 0 {
			matchCharge.ExtraCount = 1 // Avoid div by zero
		}

		div := extraCount / matchCharge.ExtraCount
		units := int(div)
		if float64(units)*matchCharge.ExtraCount < extraCount {
			units++
		}
		price += units * matchCharge.ExtraPrice
	}
	return price, nil
	return price, nil
}
