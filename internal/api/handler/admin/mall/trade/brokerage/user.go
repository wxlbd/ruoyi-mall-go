package brokerage

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/mall/trade"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/mall/trade/brokerage"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/member"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"
	"github.com/wxlbd/ruoyi-mall-go/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type BrokerageUserHandler struct {
	svc       *brokerage.BrokerageUserService
	memberSvc *member.MemberUserService
	// recordSvc    *brokerage.BrokerageRecordService
	// withdrawSvc  *brokerage.BrokerageWithdrawService
	logger *zap.Logger
}

func NewBrokerageUserHandler(svc *brokerage.BrokerageUserService, memberSvc *member.MemberUserService, logger *zap.Logger) *BrokerageUserHandler {
	return &BrokerageUserHandler{
		svc:       svc,
		memberSvc: memberSvc,
		logger:    logger,
	}
}

// CreateBrokerageUser 创建分销用户
// @Router /admin-api/trade/brokerage-user/create [post]
func (h *BrokerageUserHandler) CreateBrokerageUser(c *gin.Context) {
	var r trade.BrokerageUserCreateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteError(c, 400, err.Error())
		return
	}

	id, err := h.svc.CreateBrokerageUser(c.Request.Context(), &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, id)
}

// UpdateBindUser 修改推广员
// @Router /admin-api/trade/brokerage-user/update-bind-user [put]
func (h *BrokerageUserHandler) UpdateBindUser(c *gin.Context) {
	var r trade.BrokerageUserUpdateBindUserReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteError(c, 400, err.Error())
		return
	}
	if err := h.svc.UpdateBrokerageUserId(c.Request.Context(), int64(r.ID), int64(r.BindUserID)); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// ClearBindUser 清除推广员
// @Router /admin-api/trade/brokerage-user/clear-bind-user [put]
func (h *BrokerageUserHandler) ClearBindUser(c *gin.Context) {
	var r trade.BrokerageUserClearBindUserReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteError(c, 400, err.Error())
		return
	}
	if err := h.svc.UpdateBrokerageUserId(c.Request.Context(), int64(r.ID), 0); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// UpdateBrokerageEnabled 修改推广资格
// @Router /admin-api/trade/brokerage-user/update-brokerage-enable [put]
func (h *BrokerageUserHandler) UpdateBrokerageEnabled(c *gin.Context) {
	var r trade.BrokerageUserUpdateBrokerageEnabledReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteError(c, 400, err.Error())
		return
	}
	if err := h.svc.UpdateBrokerageUserEnabled(c.Request.Context(), int64(r.ID), r.Enabled); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// GetBrokerageUser 获得分销用户
// @Router /admin-api/trade/brokerage-user/get [get]
func (h *BrokerageUserHandler) GetBrokerageUser(c *gin.Context) {
	id := utils.ParseInt64(c.Query("id"))
	if id == 0 {
		response.WriteError(c, 400, "id is required")
		return
	}
	user, err := h.svc.GetBrokerageUser(c.Request.Context(), id)
	if err != nil {
		response.WriteError(c, 500, err.Error())
		return
	}

	res := trade.BrokerageUserResp{
		ID:               user.ID,
		BindUserID:       user.BindUserID,
		BindUserTime:     user.BindUserTime,
		BrokerageEnabled: bool(user.BrokerageEnabled),
		BrokerageTime:    user.BrokerageTime,
		Price:            user.BrokeragePrice,
		FrozenPrice:      user.FrozenPrice,
		CreateTime:       user.CreateTime,
	}

	// Fill Member Info
	memberUser, err := h.memberSvc.GetUser(c.Request.Context(), id)
	if err == nil && memberUser != nil {
		res.Avatar = memberUser.Avatar
		res.Nickname = memberUser.Nickname
	}

	response.WriteSuccess(c, res)
}

// GetBrokerageUserPage 获得分销用户分页
// @Router /admin-api/trade/brokerage-user/page [get]
func (h *BrokerageUserHandler) GetBrokerageUserPage(c *gin.Context) {
	var r trade.BrokerageUserPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteError(c, 400, err.Error())
		return
	}

	pageResult, err := h.svc.GetBrokerageUserPage(c.Request.Context(), &r)
	if err != nil {
		response.WriteError(c, 500, err.Error())
		return
	}

	// Aggregate Data
	userIds := make([]int64, len(pageResult.List))
	for i, u := range pageResult.List {
		userIds[i] = u.ID
	}

	userMap, _ := h.memberSvc.GetUserMap(c.Request.Context(), userIds)

	list := make([]trade.BrokerageUserResp, len(pageResult.List))
	for i, u := range pageResult.List {
		res := trade.BrokerageUserResp{
			ID:               u.ID,
			BindUserID:       u.BindUserID,
			BindUserTime:     u.BindUserTime,
			BrokerageEnabled: bool(u.BrokerageEnabled),
			BrokerageTime:    u.BrokerageTime,
			Price:            u.BrokeragePrice,
			FrozenPrice:      u.FrozenPrice,
			CreateTime:       u.CreateTime,
		}
		if mu, ok := userMap[u.ID]; ok {
			res.Avatar = mu.Avatar
			res.Nickname = mu.Nickname
		}

		// TODO: Aggregate BrokerageRecord and Withdraw data

		list[i] = res
	}

	response.WriteSuccess(c, pagination.PageResult[trade.BrokerageUserResp]{
		List:  list,
		Total: pageResult.Total,
	})
}
