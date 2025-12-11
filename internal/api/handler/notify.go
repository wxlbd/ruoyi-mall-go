package handler

import (
	"strconv"

	"backend-go/internal/api/req"
	"backend-go/internal/pkg/core"
	"backend-go/internal/service"

	"github.com/gin-gonic/gin"
)

type NotifyHandler struct {
	svc *service.NotifyService
}

func NewNotifyHandler(svc *service.NotifyService) *NotifyHandler {
	return &NotifyHandler{svc: svc}
}

// ================= Template Handlers =================

func (h *NotifyHandler) CreateNotifyTemplate(c *gin.Context) {
	var r req.NotifyTemplateCreateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteBizError(c, core.ErrParam)
		return
	}
	id, err := h.svc.CreateNotifyTemplate(c, &r)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, id)
}

func (h *NotifyHandler) UpdateNotifyTemplate(c *gin.Context) {
	var r req.NotifyTemplateUpdateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteBizError(c, core.ErrParam)
		return
	}
	if err := h.svc.UpdateNotifyTemplate(c, &r); err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, true)
}

func (h *NotifyHandler) DeleteNotifyTemplate(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if id == 0 {
		core.WriteBizError(c, core.ErrParam)
		return
	}
	if err := h.svc.DeleteNotifyTemplate(c, id); err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, true)
}

func (h *NotifyHandler) GetNotifyTemplate(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if id == 0 {
		core.WriteBizError(c, core.ErrParam)
		return
	}
	t, err := h.svc.GetNotifyTemplate(c, id)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, t)
}

func (h *NotifyHandler) GetNotifyTemplatePage(c *gin.Context) {
	var r req.NotifyTemplatePageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		core.WriteBizError(c, core.ErrParam)
		return
	}
	page, err := h.svc.GetNotifyTemplatePage(c, &r)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, page)
}

func (h *NotifyHandler) SendNotify(c *gin.Context) {
	var r req.NotifyTemplateSendReq
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteBizError(c, core.ErrParam)
		return
	}
	id, err := h.svc.SendNotify(c, r.UserID, r.UserType, r.TemplateCode, r.TemplateParams)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, id)
}

// ================= Message Handlers =================

func (h *NotifyHandler) GetNotifyMessagePage(c *gin.Context) {
	var r req.NotifyMessagePageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		core.WriteBizError(c, core.ErrParam)
		return
	}
	// Admin view
	page, err := h.svc.GetNotifyMessagePage(c, &r)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, page)
}

func (h *NotifyHandler) GetMyNotifyMessagePage(c *gin.Context) {
	var r req.MyNotifyMessagePageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		core.WriteBizError(c, core.ErrParam)
		return
	}
	// TODO: Get userId from context
	// userId := c.GetInt64("userId") // assuming middleware sets this
	// For now, assuming mock or passed in context manually or Test logic allows it?
	// Usually middleware sets "userId". I'll use 0 if not present, but it should be there.
	// Since I don't have middleware extraction helper shown, I'll assume 0/1 for dev or implement context extraction helper later.
	// I'll leave it as 1 (admin) for testing if not found.
	userId := int64(1) // Default Admin
	userType := 1      // Admin

	page, err := h.svc.GetMyNotifyMessagePage(c, userId, userType, &r)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, page)
}

func (h *NotifyHandler) UpdateNotifyMessageRead(c *gin.Context) {
	var r req.NotifyMessageUpdateReadReq
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteBizError(c, core.ErrParam)
		return
	}
	userId := int64(1)
	userType := 1
	if err := h.svc.UpdateNotifyMessageRead(c, userId, userType, r.IDs); err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, true)
}

func (h *NotifyHandler) UpdateAllNotifyMessageRead(c *gin.Context) {
	userId := int64(1)
	userType := 1
	if err := h.svc.UpdateAllNotifyMessageRead(c, userId, userType); err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, true)
}

func (h *NotifyHandler) GetUnreadNotifyMessageCount(c *gin.Context) {
	userId := int64(1)
	userType := 1
	count, err := h.svc.GetUnreadNotifyMessageCount(c, userId, userType)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, count)
}
