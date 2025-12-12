package promotion

import (
	"backend-go/internal/api/req"
	"backend-go/internal/api/resp"
	"backend-go/internal/model/promotion"
	"backend-go/internal/pkg/core"
	productSvc "backend-go/internal/service/product"
	promotionSvc "backend-go/internal/service/promotion"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

type PointActivityHandler struct {
	svc    *promotionSvc.PointActivityService
	spuSvc *productSvc.ProductSpuService
}

func NewPointActivityHandler(svc *promotionSvc.PointActivityService, spuSvc *productSvc.ProductSpuService) *PointActivityHandler {
	return &PointActivityHandler{
		svc:    svc,
		spuSvc: spuSvc,
	}
}

// CreatePointActivity 创建积分商城活动
func (h *PointActivityHandler) CreatePointActivity(c *gin.Context) {
	var r req.PointActivityCreateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteError(c, 400, err.Error())
		return
	}
	id, err := h.svc.CreatePointActivity(c, &r)
	if err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}
	core.WriteSuccess(c, id)
}

// UpdatePointActivity 更新积分商城活动
func (h *PointActivityHandler) UpdatePointActivity(c *gin.Context) {
	var r req.PointActivityUpdateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteError(c, 400, err.Error())
		return
	}
	if err := h.svc.UpdatePointActivity(c, &r); err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}
	core.WriteSuccess(c, true)
}

// ClosePointActivity 关闭积分商城活动
func (h *PointActivityHandler) ClosePointActivity(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		core.WriteError(c, 400, "参数错误")
		return
	}
	if err := h.svc.ClosePointActivity(c, id); err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}
	core.WriteSuccess(c, true)
}

// DeletePointActivity 删除积分商城活动
func (h *PointActivityHandler) DeletePointActivity(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		core.WriteError(c, 400, "参数错误")
		return
	}
	if err := h.svc.DeletePointActivity(c, id); err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}
	core.WriteSuccess(c, true)
}

// GetPointActivity 获得积分商城活动
func (h *PointActivityHandler) GetPointActivity(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		core.WriteError(c, 400, "参数错误")
		return
	}
	activity, products, err := h.svc.GetPointActivity(c, id)
	if err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}
	if activity == nil {
		core.WriteSuccess(c, nil)
		return
	}

	// 组装 VO
	vo := resp.PointActivityRespVO{
		ID:         activity.ID,
		SpuID:      activity.SpuID,
		Status:     activity.Status,
		Remark:     activity.Remark,
		Sort:       activity.Sort,
		Stock:      activity.Stock,
		TotalStock: activity.TotalStock,
		CreateTime: activity.CreatedAt,
	}

	productVOs := make([]resp.PointProductRespVO, len(products))
	for i, p := range products {
		productVOs[i] = resp.PointProductRespVO{
			ID:             p.ID,
			ActivityID:     p.ActivityID,
			SpuID:          p.SpuID,
			SkuID:          p.SkuID,
			Count:          p.Count,
			Point:          p.Point,
			Price:          p.Price,
			Stock:          p.Stock,
			ActivityStatus: p.ActivityStatus,
		}
	}
	vo.Products = productVOs

	core.WriteSuccess(c, vo)
}

// GetPointActivityPage 获得积分商城活动分页
func (h *PointActivityHandler) GetPointActivityPage(c *gin.Context) {
	var r req.PointActivityPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		core.WriteError(c, 400, err.Error())
		return
	}
	pageResult, err := h.svc.GetPointActivityPage(c, &r)
	if err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}

	// 拼接数据
	list := make([]*promotion.PromotionPointActivity, len(pageResult.List))
	for i := range pageResult.List {
		list[i] = &pageResult.List[i]
	}
	resultList, err := h.buildPointActivityRespVOList(c, list)
	if err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}

	core.WriteSuccess(c, core.PageResult[resp.PointActivityRespVO]{
		List:  resultList,
		Total: pageResult.Total,
	})
}

// GetPointActivityListByIds 获得积分商城活动列表
func (h *PointActivityHandler) GetPointActivityListByIds(c *gin.Context) {
	idsStr := c.QueryArray("ids")
	ids := make([]int64, 0, len(idsStr))
	for _, s := range idsStr {
		if id, err := strconv.ParseInt(s, 10, 64); err == nil {
			ids = append(ids, id)
		}
	}
	list, err := h.svc.GetPointActivityListByIds(c, ids)
	if err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}

	result, err := h.buildPointActivityRespVOList(c, list)
	if err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}
	core.WriteSuccess(c, result)
}

func (h *PointActivityHandler) buildPointActivityRespVOList(c *gin.Context, activityList []*promotion.PromotionPointActivity) ([]resp.PointActivityRespVO, error) {
	if len(activityList) == 0 {
		return []resp.PointActivityRespVO{}, nil
	}

	// 1. 获取活动商品列表
	activityIds := lo.Map(activityList, func(item *promotion.PromotionPointActivity, _ int) int64 {
		return item.ID
	})
	products, err := h.svc.GetPointProductListByActivityIds(c, activityIds)
	if err != nil {
		return nil, err
	}
	productsMap := lo.GroupBy(products, func(item *promotion.PromotionPointProduct) int64 {
		return item.ActivityID
	})

	// 2. 获取 SPU 信息
	spuIds := lo.Map(activityList, func(item *promotion.PromotionPointActivity, _ int) int64 {
		return item.SpuID
	})
	spuList, err := h.spuSvc.GetSpuList(c, spuIds)
	if err != nil {
		return nil, err
	}
	spuMap := lo.KeyBy(spuList, func(item *resp.ProductSpuResp) int64 {
		return item.ID
	})

	// 3. 组装结果
	result := make([]resp.PointActivityRespVO, len(activityList))
	for i, activity := range activityList {
		vo := resp.PointActivityRespVO{
			ID:         activity.ID,
			SpuID:      activity.SpuID,
			Status:     activity.Status,
			Remark:     activity.Remark,
			Sort:       activity.Sort,
			Stock:      activity.Stock,
			TotalStock: activity.TotalStock,
			CreateTime: activity.CreatedAt,
		}

		// 设置 Product 信息 (Min Point/Price)
		if actProducts, ok := productsMap[activity.ID]; ok && len(actProducts) > 0 {
			minProduct := lo.MinBy(actProducts, func(a, b *promotion.PromotionPointProduct) bool {
				return a.Point < b.Point
			})
			if minProduct != nil {
				vo.Point = minProduct.Point
				vo.Price = minProduct.Price
			}
		}

		// 设置 SPU 信息
		if spu, ok := spuMap[activity.SpuID]; ok {
			vo.SpuName = spu.Name
			vo.PicUrl = spu.PicURL
			vo.MarketPrice = spu.MarketPrice
		}

		result[i] = vo
	}
	return result, nil
}
