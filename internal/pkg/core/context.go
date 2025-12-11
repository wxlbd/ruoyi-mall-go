package core

import "github.com/gin-gonic/gin"

const (
	CtxUserIDKey    = "userID"
	CtxLoginUserKey = "loginUser"
)

// LoginUser 登录用户信息，与 Java 的 LoginUser 对齐
type LoginUser struct {
	UserID   int64  `json:"userId"`
	UserType int    `json:"userType"` // 0: Member, 1: Admin
	TenantID int64  `json:"tenantId"`
	Nickname string `json:"nickname"`
}

func GetLoginUserID(c *gin.Context) int64 {
	v, exists := c.Get(CtxUserIDKey)
	if !exists {
		return 0
	}
	if id, ok := v.(int64); ok {
		return id
	}
	return 0
}

func GetUserId(c *gin.Context) int64 {
	return GetLoginUserID(c)
}

// GetLoginUser 获取完整的登录用户信息
func GetLoginUser(c *gin.Context) *LoginUser {
	v, exists := c.Get(CtxLoginUserKey)
	if !exists {
		return nil
	}
	if user, ok := v.(*LoginUser); ok {
		return user
	}
	return nil
}

// SetLoginUser 设置登录用户信息到上下文
func SetLoginUser(c *gin.Context, user *LoginUser) {
	if user != nil {
		c.Set(CtxUserIDKey, user.UserID)
		c.Set(CtxLoginUserKey, user)
	}
}

// GetTenantId 获得租户编号
func GetTenantId(c *gin.Context) int64 {
	user := GetLoginUser(c)
	if user == nil {
		return 0
	}
	return user.TenantID
}
