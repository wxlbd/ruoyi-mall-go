package promotion

import (
	"backend-go/internal/api/req"
	"backend-go/internal/api/resp"
	"backend-go/internal/pkg/core"
	"backend-go/internal/service/member"
	"backend-go/internal/service/promotion"

	"github.com/gin-gonic/gin"
)

type BargainHelpHandler struct {
	svc     *promotion.BargainHelpService
	userSvc *member.MemberUserService
}

func NewBargainHelpHandler(svc *promotion.BargainHelpService, userSvc *member.MemberUserService) *BargainHelpHandler {
	return &BargainHelpHandler{
		svc:     svc,
		userSvc: userSvc,
	}
}

// GetBargainHelpPage 获得砍价助力分页
func (h *BargainHelpHandler) GetBargainHelpPage(c *gin.Context) {
	var r req.BargainHelpPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		core.WriteError(c, 400, err.Error())
		return
	}

	// 1. Get Page
	pageResult, err := h.svc.GetBargainHelpPage(c, &r)
	if err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}
	if len(pageResult.List) == 0 {
		core.WriteSuccess(c, core.PageResult[resp.BargainHelpResp]{
			List:  []resp.BargainHelpResp{},
			Total: pageResult.Total,
		})
		return
	}

	// 2. Collect IDs
	userIds := make([]int64, 0, len(pageResult.List))
	for _, item := range pageResult.List {
		userIds = append(userIds, item.UserID)
	}

	// 3. Fetch Data
	userMap, _ := h.userSvc.GetUserMap(c, userIds)

	// 4. Assemble
	list := make([]resp.BargainHelpResp, len(pageResult.List))
	for i, item := range pageResult.List {
		vo := resp.BargainHelpResp{
			ID:          item.ID,
			UserID:      item.UserID,
			ActivityID:  item.ActivityID,
			RecordID:    item.RecordID,
			ReducePrice: item.ReducePrice,
			CreatedAt:   item.CreatedAt,
		}
		if u, ok := userMap[item.UserID]; ok {
			vo.UserNickname = u.Nickname
			vo.UserAvatar = u.Avatar
		}
		list[i] = vo
	}

	core.WriteSuccess(c, core.PageResult[resp.BargainHelpResp]{
		List:  list,
		Total: pageResult.Total,
	})
}
