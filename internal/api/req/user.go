package req

import (
	"time"

	"backend-go/internal/pkg/core"
)

type UserPageReq struct {
	core.PageParam
	Username     string     `form:"username"`
	Mobile       string     `form:"mobile"`
	Status       *int       `form:"status"`
	DeptID       int64      `form:"deptId"`
	CreateTimeGe *time.Time `form:"createTime[0]"` // Helper for time range
	CreateTimeLe *time.Time `form:"createTime[1]"`
}

type UserSaveReq struct {
	ID       int64   `json:"id"`
	Username string  `json:"username" binding:"required"`
	Nickname string  `json:"nickname" binding:"required"`
	Email    string  `json:"email"`
	Mobile   string  `json:"mobile"`
	Sex      int32   `json:"sex"`
	Avatar   string  `json:"avatar"`
	DeptID   int64   `json:"deptId"`
	PostIDs  []int64 `json:"postIds"`
	RoleIDs  []int64 `json:"roleIds"`
	Status   int     `json:"status"`
	Remark   string  `json:"remark"`
	Password string  `json:"password"` // Required for Create, Optional for Update (usually separate API)
}

type UserUpdateStatusReq struct {
	ID     int64 `json:"id" binding:"required"`
	Status int   `json:"status" binding:"required"`
}

type UserUpdatePasswordReq struct {
	ID       int64  `json:"id" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserResetPasswordReq struct {
	ID       int64  `json:"id" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserExportReq struct {
	Username     string     `form:"username"`
	Mobile       string     `form:"mobile"`
	Status       *int       `form:"status"`
	DeptID       int64      `form:"deptId"`
	CreateTimeGe *time.Time `form:"createTime[0]"`
	CreateTimeLe *time.Time `form:"createTime[1]"`
}
