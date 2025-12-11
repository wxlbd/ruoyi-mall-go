package handler

import (
	"strconv"

	"backend-go/internal/api/req"
	"backend-go/internal/pkg/core"
	"backend-go/internal/service"

	"github.com/gin-gonic/gin"
)

type SensitiveWordHandler struct {
	svc *service.SensitiveWordService
}

func NewSensitiveWordHandler(svc *service.SensitiveWordService) *SensitiveWordHandler {
	return &SensitiveWordHandler{svc: svc}
}

// CreateSensitiveWord 创建敏感词
func (h *SensitiveWordHandler) CreateSensitiveWord(c *gin.Context) {
	var r req.SensitiveWordCreateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteBizError(c, core.ErrParam)
		return
	}
	id, err := h.svc.CreateSensitiveWord(c, &r)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, id)
}

// UpdateSensitiveWord 更新敏感词
func (h *SensitiveWordHandler) UpdateSensitiveWord(c *gin.Context) {
	var r req.SensitiveWordUpdateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteBizError(c, core.ErrParam)
		return
	}
	if err := h.svc.UpdateSensitiveWord(c, &r); err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, true)
}

// DeleteSensitiveWord 删除敏感词
func (h *SensitiveWordHandler) DeleteSensitiveWord(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if id == 0 {
		core.WriteBizError(c, core.ErrParam)
		return
	}
	if err := h.svc.DeleteSensitiveWord(c, id); err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, true)
}

// GetSensitiveWord 获得敏感词
func (h *SensitiveWordHandler) GetSensitiveWord(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if id == 0 {
		core.WriteBizError(c, core.ErrParam)
		return
	}
	word, err := h.svc.GetSensitiveWord(c, id)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, word)
}

// GetSensitiveWordPage 获得敏感词分页
func (h *SensitiveWordHandler) GetSensitiveWordPage(c *gin.Context) {
	var r req.SensitiveWordPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		core.WriteBizError(c, core.ErrParam)
		return
	}
	page, err := h.svc.GetSensitiveWordPage(c, &r)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, page)
}

// ValidateSensitiveWord 验证敏感词
func (h *SensitiveWordHandler) ValidateSensitiveWord(c *gin.Context) {
	text := c.Query("text")
	tag := c.Query("tag") // Single tag for simple test, or array
	var tags []string
	if tag != "" {
		tags = append(tags, tag)
	}

	words := h.svc.ValidateSensitiveWord(c, text, tags)
	core.WriteSuccess(c, words)
}

// ExportSensitiveWord 导出敏感词
func (h *SensitiveWordHandler) ExportSensitiveWord(c *gin.Context) {
	// TODO: Implement Export (Can limit only Page for now or reuse GetPage logic with large size)
	core.WriteSuccess(c, nil)
}
