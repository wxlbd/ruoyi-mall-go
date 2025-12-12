package promotion

import (
	"strconv"

	"backend-go/internal/api/req"
	"backend-go/internal/pkg/core"
	"backend-go/internal/service/promotion"

	"github.com/gin-gonic/gin"
)

type CouponHandler struct {
	svc *promotion.CouponService
}

func NewCouponHandler(svc *promotion.CouponService) *CouponHandler {
	return &CouponHandler{svc: svc}
}

// CreateCouponTemplate 创建模板
func (h *CouponHandler) CreateCouponTemplate(c *gin.Context) {
	var r req.CouponTemplateCreateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteError(c, 400, err.Error())
		return
	}
	id, err := h.svc.CreateCouponTemplate(c, &r)
	if err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}
	core.WriteSuccess(c, id)
}

// UpdateCouponTemplate 更新模板
func (h *CouponHandler) UpdateCouponTemplate(c *gin.Context) {
	var r req.CouponTemplateUpdateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteError(c, 400, err.Error())
		return
	}
	if err := h.svc.UpdateCouponTemplate(c, &r); err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}
	core.WriteSuccess(c, true)
}

// GetCouponTemplatePage 模板分页
func (h *CouponHandler) GetCouponTemplatePage(c *gin.Context) {
	var r req.CouponTemplatePageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		core.WriteError(c, 400, err.Error())
		return
	}
	list, err := h.svc.GetCouponTemplatePage(c, &r)
	if err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}
	core.WriteSuccess(c, list)
}

// GetCouponPage 发放记录
func (h *CouponHandler) GetCouponPage(c *gin.Context) {
	var r req.CouponPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		core.WriteError(c, 400, err.Error())
		return
	}
	list, err := h.svc.GetCouponPage(c, &r)
	if err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}
	core.WriteSuccess(c, list)
}

// UpdateCouponTemplateStatus 更新模板状态
// 对应 Java: CouponTemplateController.updateCouponTemplateStatus
func (h *CouponHandler) UpdateCouponTemplateStatus(c *gin.Context) {
	var r req.CouponTemplateUpdateStatusReq
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteError(c, 400, err.Error())
		return
	}
	if err := h.svc.UpdateCouponTemplateStatus(c, r.ID, r.Status); err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}
	core.WriteSuccess(c, true)
}

// DeleteCouponTemplate 删除模板
// 对应 Java: CouponTemplateController.deleteCouponTemplate
func (h *CouponHandler) DeleteCouponTemplate(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		core.WriteError(c, 400, "参数错误")
		return
	}
	if err := h.svc.DeleteCouponTemplate(c, id); err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}
	core.WriteSuccess(c, true)
}

// GetCouponTemplate 获取模板详情
// 对应 Java: CouponTemplateController.getCouponTemplate
func (h *CouponHandler) GetCouponTemplate(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		core.WriteError(c, 400, "参数错误")
		return
	}
	template, err := h.svc.GetCouponTemplate(c, id)
	if err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}
	core.WriteSuccess(c, template)
}

// GetCouponTemplateList 获取模板列表
// 对应 Java: CouponTemplateController.getCouponTemplateList
func (h *CouponHandler) GetCouponTemplateList(c *gin.Context) {
	idsStr := c.QueryArray("ids")
	ids := make([]int64, 0, len(idsStr))
	for _, s := range idsStr {
		if id, err := strconv.ParseInt(s, 10, 64); err == nil {
			ids = append(ids, id)
		}
	}
	list, err := h.svc.GetCouponTemplateList(c, ids)
	if err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}
	core.WriteSuccess(c, list)
}

// DeleteCoupon 删除/回收优惠券
// 对应 Java: CouponController.deleteCoupon
func (h *CouponHandler) DeleteCoupon(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		core.WriteError(c, 400, "参数错误")
		return
	}
	if err := h.svc.DeleteCoupon(c, id); err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}
	core.WriteSuccess(c, true)
}

// SendCoupon 发送优惠券
// 对应 Java: CouponController.sendCoupon
func (h *CouponHandler) SendCoupon(c *gin.Context) {
	var r req.CouponSendReq
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteError(c, 400, err.Error())
		return
	}
	if err := h.svc.TakeCouponByAdmin(c, r.TemplateID, r.UserIDs); err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}
	core.WriteSuccess(c, true)
}
