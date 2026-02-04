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

type BrokerageRecordHandler struct {
	logger        *zap.Logger
	recordSvc     *brokerage.BrokerageRecordService
	memberUserSvc *member.MemberUserService
}

func NewBrokerageRecordHandler(logger *zap.Logger, recordSvc *brokerage.BrokerageRecordService, memberUserSvc *member.MemberUserService) *BrokerageRecordHandler {
	return &BrokerageRecordHandler{
		logger:        logger,
		recordSvc:     recordSvc,
		memberUserSvc: memberUserSvc,
	}
}

// GetBrokerageRecord 获得分销记录
func (h *BrokerageRecordHandler) GetBrokerageRecord(c *gin.Context) {
	id := utils.ParseInt64(c.Query("id"))
	if id == 0 {
		response.WriteError(c, 400, "参数错误") // TODO: Error code
		return
	}

	record, err := h.recordSvc.GetBrokerageRecord(c, id)
	if err != nil {
		response.WriteError(c, 500, err.Error())
		return
	}
	if record == nil {
		response.WriteError(c, 404, "记录不存在")
		return
	}

	res := trade.BrokerageRecordResp{
		ID:              record.ID,
		UserID:          record.UserID,
		BizType:         record.BizType,
		BizID:           record.BizID,
		SourceUserID:    record.SourceUserID,
		SourceUserLevel: record.SourceUserLevel,
		Price:           record.Price,
		Status:          record.Status,
		FrozenDays:      record.FrozenDays,
		UnfreezeTime:    record.UnfreezeTime,
		Title:           record.Title,
		// ... copy fields
		CreateTime: record.CreateTime,
	}

	response.WriteSuccess(c, res)
}

// GetBrokerageRecordPage 获得分销记录分页
func (h *BrokerageRecordHandler) GetBrokerageRecordPage(c *gin.Context) {
	var r trade.BrokerageRecordPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteError(c, 400, "参数错误")
		return
	}
	r.CreateTime = resolveCreateTimeQuery(c, r.CreateTime)

	pageResult, err := h.recordSvc.GetBrokerageRecordPage(c, &r)
	if err != nil {
		response.WriteError(c, 500, err.Error())
		return
	}

	// Aggregate User Info
	userIds := make([]int64, 0, len(pageResult.List)*2)
	for _, item := range pageResult.List {
		userIds = append(userIds, item.UserID)
		if item.SourceUserID > 0 {
			userIds = append(userIds, item.SourceUserID)
		}
	}
	userMap, _ := h.memberUserSvc.GetUserMap(c, userIds)

	list := make([]trade.BrokerageRecordResp, len(pageResult.List))
	for i, item := range pageResult.List {
		res := trade.BrokerageRecordResp{
			ID:              item.ID,
			UserID:          item.UserID,
			BizType:         item.BizType,
			BizID:           item.BizID,
			SourceUserID:    item.SourceUserID,
			SourceUserLevel: item.SourceUserLevel,
			Price:           item.Price,
			Status:          item.Status,
			FrozenDays:      item.FrozenDays,
			UnfreezeTime:    item.UnfreezeTime,
			Title:           item.Title,
			CreateTime:      item.CreateTime,
		}
		if u, ok := userMap[item.UserID]; ok {
			res.UserNickname = u.Nickname
			res.UserAvatar = u.Avatar
		}
		if item.SourceUserID > 0 {
			if u, ok := userMap[item.SourceUserID]; ok {
				res.SourceUserNickname = u.Nickname
				res.SourceUserAvatar = u.Avatar
			}
		}

		list[i] = res
	}

	response.WriteSuccess(c, pagination.PageResult[trade.BrokerageRecordResp]{
		List:  list,
		Total: pageResult.Total,
	})
}

func resolveCreateTimeQuery(c *gin.Context, createTime []string) []string {
	if len(createTime) == 2 {
		return createTime
	}
	createTime = c.QueryArray("createTime")
	if len(createTime) == 2 {
		return createTime
	}
	return c.QueryArray("createTime[]")
}
