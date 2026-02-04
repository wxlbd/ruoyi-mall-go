package trade

import (
	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

// BrokerageUserCreateReq 创建分销用户 Request
type BrokerageUserCreateReq struct {
	UserID     model.FlexInt64 `json:"userId" binding:"required"`
	BindUserID model.FlexInt64 `json:"bindUserId"`
}

// BrokerageUserUpdateBindUserReq 修改推广员 Request
type BrokerageUserUpdateBindUserReq struct {
	ID         model.FlexInt64 `json:"id" binding:"required"`
	BindUserID model.FlexInt64 `json:"bindUserId" binding:"required"`
}

// BrokerageUserClearBindUserReq 清除推广员 Request
type BrokerageUserClearBindUserReq struct {
	ID model.FlexInt64 `json:"id" binding:"required"`
}

// BrokerageUserUpdateBrokerageEnabledReq 修改推广资格 Request
type BrokerageUserUpdateBrokerageEnabledReq struct {
	ID      model.FlexInt64 `json:"id" binding:"required"`
	Enabled bool            `json:"enabled"`
}

// BrokerageUserPageReq 分销用户分页 Request
type BrokerageUserPageReq struct {
	pagination.PageParam
	BindUserID       int64    `form:"bindUserId"`
	BrokerageEnabled *bool    `form:"brokerageEnabled"`
	CreateTime       []string `form:"createTime[]"`   // Range
	Level            int      `form:"level"`          // 推广层级过滤
	BindUserTime     []string `form:"bindUserTime[]"` // Range
}

// BrokerageRecordPageReq 分销记录分页 Request
type BrokerageRecordPageReq struct {
	pagination.PageParam
	UserID          *int64   `form:"userId"`          // 用户编号（改为指针类型）
	BizType         *int     `form:"bizType"`         // 业务类型
	Status          *int     `form:"status"`          // 状态
	CreateTime      []string `form:"createTime[]"`    // Range
	SourceUserLevel *int     `form:"sourceUserLevel"` // 来源用户等级
	BizID           string   `form:"bizId"`           // 业务编号
}

// BrokerageWithdrawPageReq 分销提现分页 Request
type BrokerageWithdrawPageReq struct {
	pagination.PageParam
	UserID      int64    `form:"userId"`
	Type        int      `form:"type"`
	Status      int      `form:"status"`
	UserName    string   `form:"userName"`
	UserAccount string   `form:"userAccount"`
	BankName    string   `form:"bankName"`
	CreateTime  []string `form:"createTime[]"` // Range
}

// BrokerageWithdrawRejectReq 分销提现驳回 Request
type BrokerageWithdrawRejectReq struct {
	ID          model.FlexInt64 `json:"id" binding:"required"`
	AuditReason string          `json:"auditReason" binding:"required"`
}

// AppBrokerageUserRankPageReq 分销用户排行分页 Request (App)
type AppBrokerageUserRankPageReq struct {
	pagination.PageParam
	Times []string `form:"times[]"` // 时间范围 [start, end]
}

// ========== Response DTOs ==========

// BrokerageUserResp 分销用户 Response
type BrokerageUserResp struct {
	ID               int64      `json:"id"`
	BindUserID       int64      `json:"bindUserId"`
	BindUserTime     *time.Time `json:"bindUserTime"`
	BrokerageEnabled bool       `json:"brokerageEnabled"`
	BrokerageTime    *time.Time `json:"brokerageTime"`
	Price            int        `json:"price"`
	FrozenPrice      int        `json:"frozenPrice"`
	CreateTime       time.Time  `json:"createTime"`
	// Member Info
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

// BrokerageRecordResp 分销记录 Response
type BrokerageRecordResp struct {
	ID              int64      `json:"id"`
	UserID          int64      `json:"userId"`
	BizType         int        `json:"bizType"`
	BizID           string     `json:"bizId"`
	SourceUserID    int64      `json:"sourceUserId"`
	SourceUserLevel int        `json:"sourceUserLevel"`
	Price           int        `json:"price"`
	Status          int        `json:"status"`
	FrozenDays      int        `json:"frozenDays"`
	UnfreezeTime    *time.Time `json:"unfreezeTime"`
	Title           string     `json:"title"`
	CreateTime      time.Time  `json:"createTime"`
	// Derived Fields
	UserNickname       string `json:"userNickname"`
	UserAvatar         string `json:"userAvatar"`
	SourceUserNickname string `json:"sourceUserNickname"`
	SourceUserAvatar   string `json:"sourceUserAvatar"`
}

// BrokerageWithdrawResp 分销提现 Response
type BrokerageWithdrawResp struct {
	ID               int64      `json:"id"`
	UserID           int64      `json:"userId"`
	Price            int        `json:"price"`
	FeePrice         int        `json:"feePrice"`
	TotalPrice       int        `json:"totalPrice"`
	Type             int        `json:"type"`
	UserName         string     `json:"userName"`
	UserAccount      string     `json:"userAccount"`
	QRCodeUrl        string     `json:"qrCodeUrl"`
	BankName         string     `json:"bankName"`
	BankAddress      string     `json:"bankAddress"`
	Status           int        `json:"status"`
	AuditReason      string     `json:"auditReason"`
	AuditTime        *time.Time `json:"auditTime"`
	Remark           string     `json:"remark"`
	PayTransferID    int64      `json:"payTransferId"`
	TransferErrorMsg string     `json:"transferErrorMsg"`
	CreateTime       time.Time  `json:"createTime"`
	// Member Info
	UserNickname string `json:"userNickname"`
}
