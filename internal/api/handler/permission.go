package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"backend-go/internal/api/req"
	"backend-go/internal/pkg/core"
	"backend-go/internal/service"
)

type PermissionHandler struct {
	svc       *service.PermissionService
	tenantSvc *service.TenantService
}

func NewPermissionHandler(svc *service.PermissionService, tenantSvc *service.TenantService) *PermissionHandler {
	return &PermissionHandler{
		svc:       svc,
		tenantSvc: tenantSvc,
	}
}

func (h *PermissionHandler) GetRoleMenuList(c *gin.Context) {
	roleIdStr := c.Query("roleId")
	roleId, _ := strconv.ParseInt(roleIdStr, 10, 64)
	list, err := h.svc.GetRoleMenuListByRoleId(c.Request.Context(), []int64{roleId})
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(list))
}

func (h *PermissionHandler) AssignRoleMenu(c *gin.Context) {
	var r req.PermissionAssignRoleMenuReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(200, core.ErrParam)
		return
	}
	// Filter menus by tenant
	err := h.tenantSvc.HandleTenantMenu(c, func(allowedMenuIds []int64) {
		if allowedMenuIds == nil {
			return
		}
		// Filter r.MenuIDs
		allowedSet := make(map[int64]bool)
		for _, id := range allowedMenuIds {
			allowedSet[id] = true
		}
		filtered := make([]int64, 0, len(r.MenuIDs))
		for _, id := range r.MenuIDs {
			if allowedSet[id] {
				filtered = append(filtered, id)
			}
		}
		r.MenuIDs = filtered
	})
	if err != nil {
		c.Error(err)
		return
	}

	if err := h.svc.AssignRoleMenu(c.Request.Context(), r.RoleID, r.MenuIDs); err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(true))
}

func (h *PermissionHandler) AssignRoleDataScope(c *gin.Context) {
	var r req.PermissionAssignRoleDataScopeReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(200, core.ErrParam)
		return
	}
	if err := h.svc.AssignRoleDataScope(c.Request.Context(), r.RoleID, r.DataScope, r.DataScopeDeptIDs); err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(true))
}

func (h *PermissionHandler) GetUserRoleList(c *gin.Context) {
	userIdStr := c.Query("userId")
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)
	list, err := h.svc.GetUserRoleIdListByUserId(c.Request.Context(), userId)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(list))
}

func (h *PermissionHandler) AssignUserRole(c *gin.Context) {
	var r req.PermissionAssignUserRoleReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(200, core.ErrParam)
		return
	}
	if err := h.svc.AssignUserRole(c.Request.Context(), r.UserID, r.RoleIDs); err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(true))
}
