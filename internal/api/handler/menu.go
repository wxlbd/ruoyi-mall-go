package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"backend-go/internal/api/req"
	"backend-go/internal/pkg/core"
	"backend-go/internal/service"
)

type MenuHandler struct {
	svc *service.MenuService
}

func NewMenuHandler(svc *service.MenuService) *MenuHandler {
	return &MenuHandler{
		svc: svc,
	}
}

// CreateMenu 创建菜单
// @Router /system/menu/create [post]
func (h *MenuHandler) CreateMenu(c *gin.Context) {
	var r req.MenuCreateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(200, core.ErrParam)
		return
	}
	id, err := h.svc.CreateMenu(c.Request.Context(), &r)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(id))
}

// UpdateMenu 更新菜单
// @Router /system/menu/update [put]
func (h *MenuHandler) UpdateMenu(c *gin.Context) {
	var r req.MenuUpdateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(200, core.ErrParam)
		return
	}
	if err := h.svc.UpdateMenu(c.Request.Context(), &r); err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(true))
}

// DeleteMenu 删除菜单
// @Router /system/menu/delete [delete]
func (h *MenuHandler) DeleteMenu(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	if err := h.svc.DeleteMenu(c.Request.Context(), id); err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(true))
}

// GetMenuList 获取菜单列表
// @Router /system/menu/list [get]
func (h *MenuHandler) GetMenuList(c *gin.Context) {
	var r req.MenuListReq
	if err := c.ShouldBindQuery(&r); err != nil {
		c.JSON(200, core.ErrParam)
		return
	}
	list, err := h.svc.GetMenuList(c.Request.Context(), &r)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(list))
}

// GetMenu 获取菜单详情
// @Router /system/menu/get [get]
func (h *MenuHandler) GetMenu(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	item, err := h.svc.GetMenu(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(item))
}

// GetSimpleMenuList 获取精简菜单列表
// @Router /system/menu/simple-list [get]
func (h *MenuHandler) GetSimpleMenuList(c *gin.Context) {
	list, err := h.svc.GetSimpleMenuList(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(list))
}
