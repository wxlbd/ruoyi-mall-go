package promotion

import (
	"strconv"

	productContract "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/mall/product"
	promotion2 "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/mall/promotion"
	"github.com/wxlbd/ruoyi-mall-go/internal/consts"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/mall/product"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/mall/promotion"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	promotionModel "github.com/wxlbd/ruoyi-mall-go/internal/model/promotion"

	"github.com/gin-gonic/gin"
)

type SeckillActivityHandler struct {
	svc    *promotion.SeckillActivityService
	spuSvc *product.ProductSpuService // Needed for response composition
}

func NewSeckillActivityHandler(svc *promotion.SeckillActivityService, spuSvc *product.ProductSpuService) *SeckillActivityHandler {
	return &SeckillActivityHandler{svc: svc, spuSvc: spuSvc}
}

// CreateSeckillActivity 创建
func (h *SeckillActivityHandler) CreateSeckillActivity(c *gin.Context) {
	var r promotion2.SeckillActivityCreateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteError(c, 400, err.Error())
		return
	}
	id, err := h.svc.CreateSeckillActivity(c.Request.Context(), &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, id)
}

// UpdateSeckillActivity 更新
func (h *SeckillActivityHandler) UpdateSeckillActivity(c *gin.Context) {
	var r promotion2.SeckillActivityUpdateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteError(c, 400, err.Error())
		return
	}
	if err := h.svc.UpdateSeckillActivity(c.Request.Context(), &r); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// DeleteSeckillActivity 删除
func (h *SeckillActivityHandler) DeleteSeckillActivity(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if err := h.svc.DeleteSeckillActivity(c.Request.Context(), id); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// CloseSeckillActivity 关闭
func (h *SeckillActivityHandler) CloseSeckillActivity(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if err := h.svc.CloseSeckillActivity(c.Request.Context(), id); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// GetSeckillActivity 获得详情
func (h *SeckillActivityHandler) GetSeckillActivity(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	act, err := h.svc.GetSeckillActivity(c.Request.Context(), id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	if act == nil {
		response.WriteSuccess(c, nil)
		return
	}
	products, err := h.svc.GetSeckillProductListByActivityID(c.Request.Context(), id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	detail := promotion2.SeckillActivityDetailResp{}
	detail.SeckillActivityResp = promotion2.SeckillActivityResp{
		ID:               act.ID,
		SpuID:            act.SpuID,
		Name:             act.Name,
		Status:           act.Status,
		Remark:           act.Remark,
		StartTime:        act.StartTime,
		EndTime:          act.EndTime,
		Sort:             act.Sort,
		ConfigIds:        act.ConfigIds,
		TotalLimitCount:  act.TotalLimitCount,
		SingleLimitCount: act.SingleLimitCount,
		Stock:            act.Stock,
		TotalStock:       act.TotalStock,
		CreateTime:       act.CreateTime,
	}
	for _, p := range products {
		detail.Products = append(detail.Products, promotion2.SeckillProductResp{
			ID:           p.ID,
			ActivityID:   p.ActivityID,
			SpuID:        p.SpuID,
			SkuID:        p.SkuID,
			SeckillPrice: p.SeckillPrice,
			Stock:        p.Stock,
		})
	}
	response.WriteSuccess(c, detail)
}

// GetSeckillActivityPage 分页
func (h *SeckillActivityHandler) GetSeckillActivityPage(c *gin.Context) {
	var r promotion2.SeckillActivityPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteError(c, 400, err.Error())
		return
	}

	r.Status = response.NormalizeIntPtr(c, "status", r.Status)

	res, err := h.svc.GetSeckillActivityPage(c.Request.Context(), &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	if len(res.List) == 0 {
		response.WriteSuccess(c, pagination.PageResult[promotion2.SeckillActivityResp]{
			List:  []promotion2.SeckillActivityResp{},
			Total: res.Total,
		})
		return
	}

	// 收集 ActivityIDs 和 SpuIDs
	activityIds := make([]int64, len(res.List))
	spuIds := make([]int64, len(res.List))
	for i, v := range res.List {
		activityIds[i] = v.ID
		spuIds[i] = v.SpuID
	}

	// 批量获取 Products
	products, _ := h.svc.GetSeckillProductListByActivityIds(c.Request.Context(), activityIds)
	productMap := make(map[int64][]*promotionModel.PromotionSeckillProduct)
	for _, p := range products {
		productMap[p.ActivityID] = append(productMap[p.ActivityID], p)
	}

	// 批量获取 SPU
	spuList, _ := h.spuSvc.GetSpuList(c.Request.Context(), spuIds)
	spuMap := make(map[int64]*productContract.ProductSpuResp)
	for _, spu := range spuList {
		spuMap[spu.ID] = spu
	}

	// 构建响应
	list := make([]promotion2.SeckillActivityResp, len(res.List))
	for i, v := range res.List {
		item := promotion2.SeckillActivityResp{
			ID:               v.ID,
			SpuID:            v.SpuID,
			Name:             v.Name,
			Status:           v.Status,
			Remark:           v.Remark,
			StartTime:        v.StartTime,
			EndTime:          v.EndTime,
			Sort:             v.Sort,
			ConfigIds:        v.ConfigIds,
			TotalLimitCount:  v.TotalLimitCount,
			SingleLimitCount: v.SingleLimitCount,
			Stock:            v.Stock,
			TotalStock:       v.TotalStock,
			CreateTime:       v.CreateTime,
		}

		// 拼接 Products
		if prods, ok := productMap[v.ID]; ok {
			item.Products = make([]promotion2.SeckillProductResp, len(prods))
			minPrice := 0
			for j, p := range prods {
				item.Products[j] = promotion2.SeckillProductResp{
					ID:           p.ID,
					ActivityID:   p.ActivityID,
					SpuID:        p.SpuID,
					SkuID:        p.SkuID,
					SeckillPrice: p.SeckillPrice,
					Stock:        p.Stock,
				}
				if j == 0 || p.SeckillPrice < minPrice {
					minPrice = p.SeckillPrice
				}
			}
			item.SeckillPrice = minPrice
		}

		// 拼接 SPU
		if spu, ok := spuMap[v.SpuID]; ok {
			item.SpuName = spu.Name
			item.PicUrl = spu.PicURL
			item.MarketPrice = spu.MarketPrice
		}

		list[i] = item
	}

	response.WriteSuccess(c, pagination.PageResult[promotion2.SeckillActivityResp]{
		List:  list,
		Total: res.Total,
	})
}

// GetSeckillActivityListByIds 获得秒杀活动列表
func (h *SeckillActivityHandler) GetSeckillActivityListByIds(c *gin.Context) {
	idsStr := c.Query("ids")
	var ids []int64
	var intList model.IntListFromCSV
	if err := intList.Scan(idsStr); err != nil {
		response.WriteError(c, 400, "参数错误")
		return
	}
	for _, id := range intList {
		ids = append(ids, int64(id))
	}

	activityList, err := h.svc.GetSeckillActivityListByIds(c.Request.Context(), ids)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	// 过滤禁用状态 (对齐 Java: CommonStatusEnum.isDisable)
	var activeList []*promotionModel.PromotionSeckillActivity
	for _, act := range activityList {
		if act.Status == consts.CommonStatusEnable { // 使用 CommonStatusEnable 常量替代魔法数字 1
			activeList = append(activeList, act)
		}
	}
	if len(activeList) == 0 {
		response.WriteSuccess(c, []promotion2.SeckillActivityResp{})
		return
	}

	// 获取活动 ID 和 SPU ID
	actIds := make([]int64, len(activeList))
	spuIds := make([]int64, len(activeList))
	for i, act := range activeList {
		actIds[i] = act.ID
		spuIds[i] = act.SpuID
	}

	// 批量获取秒杀商品
	products, err := h.svc.GetSeckillProductListByActivityIds(c.Request.Context(), actIds)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	// 批量获取 SPU 信息
	spuList, err := h.spuSvc.GetSpuList(c.Request.Context(), spuIds)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	spuMap := make(map[int64]string)
	for _, spu := range spuList {
		spuMap[spu.ID] = spu.Name
	}

	// 构建商品 Map (按活动 ID 分组)
	productMap := make(map[int64][]*promotionModel.PromotionSeckillProduct)
	for _, p := range products {
		productMap[p.ActivityID] = append(productMap[p.ActivityID], p)
	}

	// 构建响应
	result := make([]promotion2.SeckillActivityResp, len(activeList))
	for i, act := range activeList {
		result[i] = promotion2.SeckillActivityResp{
			ID:               act.ID,
			SpuID:            act.SpuID,
			Name:             act.Name,
			Status:           act.Status,
			Remark:           act.Remark,
			StartTime:        act.StartTime,
			EndTime:          act.EndTime,
			Sort:             act.Sort,
			ConfigIds:        act.ConfigIds,
			TotalLimitCount:  act.TotalLimitCount,
			SingleLimitCount: act.SingleLimitCount,
			Stock:            act.Stock,
			TotalStock:       act.TotalStock,
			CreateTime:       act.CreateTime,
		}
	}
	response.WriteSuccess(c, result)
}
