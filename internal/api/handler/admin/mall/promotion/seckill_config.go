package promotion

import (
	"strconv"

	promotion2 "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/mall/promotion"
	"github.com/wxlbd/ruoyi-mall-go/internal/consts"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/mall/promotion"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

	"github.com/gin-gonic/gin"
)

type SeckillConfigHandler struct {
	svc *promotion.SeckillConfigService
}

func NewSeckillConfigHandler(svc *promotion.SeckillConfigService) *SeckillConfigHandler {
	return &SeckillConfigHandler{svc: svc}
}

// CreateSeckillConfig 创建
func (h *SeckillConfigHandler) CreateSeckillConfig(c *gin.Context) {
	var r promotion2.SeckillConfigCreateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteError(c, 400, err.Error()) // HTTP 400 Bad Request
		return
	}
	id, err := h.svc.CreateSeckillConfig(c.Request.Context(), &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, id)
}

// UpdateSeckillConfig 更新
func (h *SeckillConfigHandler) UpdateSeckillConfig(c *gin.Context) {
	var r promotion2.SeckillConfigUpdateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteError(c, 400, err.Error()) // HTTP 400 Bad Request
		return
	}
	if err := h.svc.UpdateSeckillConfig(c.Request.Context(), &r); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// UpdateSeckillConfigStatus 更新状态
func (h *SeckillConfigHandler) UpdateSeckillConfigStatus(c *gin.Context) {
	var r promotion2.SeckillConfigUpdateStatusReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteError(c, 400, err.Error())
		return
	}
	if err := h.svc.UpdateSeckillConfigStatus(c.Request.Context(), r.ID, *r.Status); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// DeleteSeckillConfig 删除
func (h *SeckillConfigHandler) DeleteSeckillConfig(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if err := h.svc.DeleteSeckillConfig(c.Request.Context(), id); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// GetSeckillConfig 获得
func (h *SeckillConfigHandler) GetSeckillConfig(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	res, err := h.svc.GetSeckillConfig(c.Request.Context(), id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, res)
}

// GetSeckillConfigPage 分页
func (h *SeckillConfigHandler) GetSeckillConfigPage(c *gin.Context) {
	var r promotion2.SeckillConfigPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteError(c, 400, err.Error())
		return
	}

	r.Status = response.NormalizeIntPtr(c, "status", r.Status)

	res, err := h.svc.GetSeckillConfigPage(c.Request.Context(), &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, res)
}

// GetSeckillConfigList 获得列表
func (h *SeckillConfigHandler) GetSeckillConfigList(c *gin.Context) {
	res, err := h.svc.GetSeckillConfigList(c.Request.Context())
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	var respList []promotion2.SeckillConfigResp
	for _, v := range res {
		respList = append(respList, promotion2.SeckillConfigResp{
			ID:            v.ID,
			Name:          v.Name,
			StartTime:     v.StartTime,
			EndTime:       v.EndTime,
			SliderPicUrls: v.SliderPicUrls,
			Status:        v.Status,
			CreateTime:    v.CreateTime,
		})
	}
	response.WriteSuccess(c, respList)
}

// GetSeckillConfigSimpleList 精简列表
func (h *SeckillConfigHandler) GetSeckillConfigSimpleList(c *gin.Context) {
	res, err := h.svc.GetSeckillConfigListByStatus(c.Request.Context(), consts.CommonStatusEnable) // 使用启用状态常量替代魔法数字 1
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	var respList []promotion2.SeckillConfigSimpleResp
	for _, v := range res {
		respList = append(respList, promotion2.SeckillConfigSimpleResp{
			ID:        v.ID,
			Name:      v.Name,
			StartTime: v.StartTime,
			EndTime:   v.EndTime,
		})
	}
	response.WriteSuccess(c, respList)
}
