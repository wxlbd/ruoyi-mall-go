package infra

import (
	"strconv"
	"strings"

	system2 "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/system"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/system"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/excel"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

	"github.com/gin-gonic/gin"
)

type ConfigHandler struct {
	configSvc *system.ConfigService
}

func NewConfigHandler(configSvc *system.ConfigService) *ConfigHandler {
	return &ConfigHandler{
		configSvc: configSvc,
	}
}

// GetConfigPage 获得参数配置分页
func (h *ConfigHandler) GetConfigPage(c *gin.Context) {
	var req system2.ConfigPageReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	res, err := h.configSvc.GetConfigPage(c, &req)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, res)
}

// GetConfig 获得参数配置
func (h *ConfigHandler) GetConfig(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if id == 0 {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	res, err := h.configSvc.GetConfig(c, id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, res)
}

// GetConfigKey 根据参数键名查询参数值
func (h *ConfigHandler) GetConfigKey(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	config, err := h.configSvc.GetConfigByKey(c, key)
	if err != nil || config == nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	if !config.Visible {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	response.WriteSuccess(c, config.Value)
}

// CreateConfig 创建参数配置
func (h *ConfigHandler) CreateConfig(c *gin.Context) {
	var req system2.ConfigSaveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	id, err := h.configSvc.CreateConfig(c, &req)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, id)
}

// UpdateConfig 更新参数配置
func (h *ConfigHandler) UpdateConfig(c *gin.Context) {
	var req system2.ConfigSaveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	if err := h.configSvc.UpdateConfig(c, &req); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// DeleteConfig 删除参数配置
func (h *ConfigHandler) DeleteConfig(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if id == 0 {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	if err := h.configSvc.DeleteConfig(c, id); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// DeleteConfigList 批量删除参数配置
func (h *ConfigHandler) DeleteConfigList(c *gin.Context) {
	idsStr := c.Query("ids")
	if idsStr == "" {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	parts := strings.Split(idsStr, ",")
	ids := make([]int64, 0, len(parts))
	for _, part := range parts {
		id, _ := strconv.ParseInt(strings.TrimSpace(part), 10, 64)
		if id > 0 {
			ids = append(ids, id)
		}
	}
	if len(ids) == 0 {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	if err := h.configSvc.DeleteConfigList(c, ids); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// ExportConfigExcel 导出参数配置
func (h *ConfigHandler) ExportConfigExcel(c *gin.Context) {
	var req system2.ConfigPageReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	req.PageNo = 1
	req.PageSize = 1000000
	res, err := h.configSvc.GetConfigPage(c, &req)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	if err := excel.WriteExcel(c, "参数配置.xls", "数据", res.List); err != nil {
		response.WriteBizError(c, err)
	}
}
