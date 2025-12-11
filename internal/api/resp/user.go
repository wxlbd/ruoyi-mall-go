package resp

import "time"

type UserRespVO struct {
	ID         int64     `json:"id"`
	Username   string    `json:"username"`
	Nickname   string    `json:"nickname"`
	Remark     string    `json:"remark"`
	DeptID     int64     `json:"deptId"`
	PostIDs    []int64   `json:"postIds"`
	RoleIDs    []int64   `json:"roleIds"`
	Email      string    `json:"email"`
	Mobile     string    `json:"mobile"`
	Sex        int32     `json:"sex"`
	Avatar     string    `json:"avatar"`
	Status     int32     `json:"status"`
	LoginIP    string    `json:"loginIp"`
	LoginDate  time.Time `json:"loginDate"`
	CreateTime time.Time `json:"createTime"`
}

type UserProfileRespVO struct {
	*UserRespVO
	Roles []*RoleRespVO `json:"roles,omitempty"`
	Posts []*PostRespVO `json:"posts,omitempty"`
}

type UserSimpleRespVO struct {
	ID       int64  `json:"id"`
	Nickname string `json:"nickname"`
}

// UserImportRespVO generic import response if needed
// UserImportRespVO generic import response if needed
type UserImportRespVO struct {
	CreateUsernames  []string          `json:"createUsernames"`
	UpdateUsernames  []string          `json:"updateUsernames"`
	FailureUsernames map[string]string `json:"failureUsernames"` // username -> error
}

// UserImportExcelVO Excel Import/Export Struct
type UserImportExcelVO struct {
	Username string `excel:"登录名称"`
	Nickname string `excel:"用户名称"`
	Email    string `excel:"用户邮箱"`
	Mobile   string `excel:"手机号码"`
	Sex      string `excel:"用户性别"` // 0=男, 1=女, 2=未知 -> Converted to string for Excel
	Status   string `excel:"帐号状态"` // 0=正常, 1=停用
	DeptID   int64  `excel:"部门编号"`
}
