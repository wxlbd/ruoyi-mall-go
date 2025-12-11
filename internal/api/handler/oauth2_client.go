package handler

import (
	"strconv"

	"backend-go/internal/api/req"
	"backend-go/internal/pkg/core"
	"backend-go/internal/service"

	"github.com/gin-gonic/gin"
)

type OAuth2ClientHandler struct {
	svc *service.OAuth2ClientService
}

func NewOAuth2ClientHandler(svc *service.OAuth2ClientService) *OAuth2ClientHandler {
	return &OAuth2ClientHandler{svc: svc}
}

func (h *OAuth2ClientHandler) CreateOAuth2Client(c *gin.Context) {
	var r req.OAuth2ClientCreateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteBizError(c, core.ErrParam)
		return
	}
	id, err := h.svc.CreateOAuth2Client(c, &r)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, id)
}

func (h *OAuth2ClientHandler) UpdateOAuth2Client(c *gin.Context) {
	var r req.OAuth2ClientUpdateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteBizError(c, core.ErrParam)
		return
	}
	if err := h.svc.UpdateOAuth2Client(c, &r); err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, true)
}

func (h *OAuth2ClientHandler) DeleteOAuth2Client(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if id == 0 {
		core.WriteBizError(c, core.ErrParam)
		return
	}
	if err := h.svc.DeleteOAuth2Client(c, id); err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, true)
}

func (h *OAuth2ClientHandler) GetOAuth2Client(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if id == 0 {
		core.WriteBizError(c, core.ErrParam)
		return
	}
	client, err := h.svc.GetOAuth2Client(c, id)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, client)
}

func (h *OAuth2ClientHandler) GetOAuth2ClientPage(c *gin.Context) {
	var r req.OAuth2ClientPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		core.WriteBizError(c, core.ErrParam)
		return
	}
	page, err := h.svc.GetOAuth2ClientPage(c, &r)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, page)
}
