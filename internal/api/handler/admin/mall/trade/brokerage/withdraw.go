package brokerage

import (
	"github.com/gin-gonic/gin"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/mall/trade"

	brokerageModel "github.com/wxlbd/ruoyi-mall-go/internal/model/trade/brokerage"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/mall/trade/brokerage"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/member"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"
	"github.com/wxlbd/ruoyi-mall-go/pkg/utils"
)

type BrokerageWithdrawHandler struct {
	withdrawSvc   *brokerage.BrokerageWithdrawService
	memberUserSvc *member.MemberUserService
}

func NewBrokerageWithdrawHandler(withdrawSvc *brokerage.BrokerageWithdrawService, memberUserSvc *member.MemberUserService) *BrokerageWithdrawHandler {
	return &BrokerageWithdrawHandler{
		withdrawSvc:   withdrawSvc,
		memberUserSvc: memberUserSvc,
	}
}

// ApproveBrokerageWithdraw 通过申请
func (h *BrokerageWithdrawHandler) ApproveBrokerageWithdraw(c *gin.Context) {
	id := utils.ParseInt64(c.Query("id"))
	if id == 0 {
		response.WriteError(c, 400, "参数错误")
		return
	}

	// 10: AUDIT_SUCCESS (See Enum in Java)
	if err := h.withdrawSvc.AuditBrokerageWithdraw(c, id, 10, ""); err != nil {
		response.WriteError(c, 500, err.Error())
		return
	}

	response.WriteSuccess(c, true)
}

// RejectBrokerageWithdraw 驳回申请
func (h *BrokerageWithdrawHandler) RejectBrokerageWithdraw(c *gin.Context) {
	var r trade.BrokerageWithdrawRejectReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteError(c, 400, "参数错误")
		return
	}

	// 20: AUDIT_FAIL
	if err := h.withdrawSvc.AuditBrokerageWithdraw(c, int64(r.ID), 20, r.AuditReason); err != nil {
		response.WriteBizError(c, err)
		return
	}

	response.WriteSuccess(c, true)
}

// UpdateBrokerageWithdrawTransferred 更新佣金提现的转账结果
func (h *BrokerageWithdrawHandler) UpdateBrokerageWithdrawTransferred(c *gin.Context) {
	// Placeholder: Req Struct
	var r struct {
		MerchantTransferID string `json:"merchantTransferId"`
		PayTransferID      int64  `json:"payTransferId"`
	}
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteError(c, 400, "参数错误")
		return
	}

	id := utils.ParseInt64(r.MerchantTransferID)
	if err := h.withdrawSvc.UpdateBrokerageWithdrawTransferred(c, id, r.PayTransferID); err != nil {
		response.WriteError(c, 500, err.Error())
		return
	}
	response.WriteSuccess(c, true)
}

// GetBrokerageWithdraw 获得佣金提现
func (h *BrokerageWithdrawHandler) GetBrokerageWithdraw(c *gin.Context) {
	id := utils.ParseInt64(c.Query("id"))
	if id == 0 {
		response.WriteError(c, 400, "参数错误")
		return
	}

	withdraw, err := h.withdrawSvc.GetBrokerageWithdraw(c, id)
	if err != nil {
		response.WriteError(c, 500, err.Error())
		return
	}
	if withdraw == nil {
		response.WriteError(c, 404, "提现记录不存在")
		return
	}

	res := h.convert(withdraw)
	response.WriteSuccess(c, res)
}

// GetBrokerageWithdrawPage 获得佣金提现分页
func (h *BrokerageWithdrawHandler) GetBrokerageWithdrawPage(c *gin.Context) {
	var r trade.BrokerageWithdrawPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteError(c, 400, "参数错误")
		return
	}

	pageResult, err := h.withdrawSvc.GetBrokerageWithdrawPage(c, &r)
	if err != nil {
		response.WriteError(c, 500, err.Error())
		return
	}

	// Aggregate User Info
	userIds := make([]int64, 0, len(pageResult.List))
	for _, item := range pageResult.List {
		userIds = append(userIds, item.UserID)
	}
	userMap, _ := h.memberUserSvc.GetUserMap(c, userIds)

	list := make([]trade.BrokerageWithdrawResp, len(pageResult.List))
	for i, item := range pageResult.List {
		res := h.convert(item)
		if u, ok := userMap[item.UserID]; ok {
			res.UserNickname = u.Nickname
		}
		list[i] = res
	}

	response.WriteSuccess(c, pagination.PageResult[trade.BrokerageWithdrawResp]{
		List:  list,
		Total: pageResult.Total,
	})
}

func (h *BrokerageWithdrawHandler) convert(do *brokerageModel.BrokerageWithdraw) trade.BrokerageWithdrawResp {
	payTransferId := int64(0)
	if do.PayTransferID > 0 {
		payTransferId = do.PayTransferID
	}
	return trade.BrokerageWithdrawResp{
		ID:               do.ID,
		UserID:           do.UserID,
		Price:            do.Price,
		FeePrice:         do.FeePrice,
		TotalPrice:       do.TotalPrice,
		Type:             do.Type,
		UserName:         do.UserName,
		UserAccount:      do.UserAccount,
		QRCodeUrl:        do.QrCodeURL,
		BankName:         do.BankName,
		BankAddress:      do.BankAddress,
		Status:           do.Status,
		AuditReason:      do.AuditReason,
		AuditTime:        do.AuditTime,
		Remark:           do.Remark,
		PayTransferID:    payTransferId,
		TransferErrorMsg: do.TransferErrorMsg,
		CreateTime:       do.CreateTime,
	}
}
