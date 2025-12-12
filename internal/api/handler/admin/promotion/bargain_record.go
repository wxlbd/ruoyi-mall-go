package promotion

import (
	"backend-go/internal/api/req"
	"backend-go/internal/api/resp"
	"backend-go/internal/pkg/core"
	"backend-go/internal/service/member"
	"backend-go/internal/service/promotion"

	"github.com/gin-gonic/gin"
)

type BargainRecordHandler struct {
	svc         *promotion.BargainRecordService
	activitySvc *promotion.BargainActivityService
	userSvc     *member.MemberUserService
}

func NewBargainRecordHandler(
	svc *promotion.BargainRecordService,
	activitySvc *promotion.BargainActivityService,
	userSvc *member.MemberUserService,
) *BargainRecordHandler {
	return &BargainRecordHandler{
		svc:         svc,
		activitySvc: activitySvc,
		userSvc:     userSvc,
	}
}

// GetBargainRecordPage 获得砍价记录分页
func (h *BargainRecordHandler) GetBargainRecordPage(c *gin.Context) {
	var r req.BargainRecordPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		core.WriteError(c, 400, err.Error())
		return
	}

	// 1. Get Page of DOs
	pageResult, err := h.svc.GetBargainRecordPageAdmin(c, &r)
	if err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}
	if len(pageResult.List) == 0 {
		core.WriteSuccess(c, core.PageResult[resp.BargainRecordResp]{
			List:  []resp.BargainRecordResp{},
			Total: pageResult.Total,
		})
		return
	}

	// 2. Collect IDs
	userIds := make([]int64, 0, len(pageResult.List))
	activityIds := make([]int64, 0, len(pageResult.List))
	for _, item := range pageResult.List {
		userIds = append(userIds, item.UserID)
		activityIds = append(activityIds, item.ActivityID)
	}

	// 3. Fetch Enriched Data
	userMap, _ := h.userSvc.GetUserMap(c, userIds)
	activityMap, _ := h.activitySvc.GetBargainActivityMap(c, activityIds)

	// 4. Assemble VOs
	list := make([]resp.BargainRecordResp, len(pageResult.List))
	for i, item := range pageResult.List {
		nickname := ""
		avatar := ""
		if u, ok := userMap[item.UserID]; ok {
			nickname = u.Nickname
			avatar = u.Avatar
		}
		activityName := ""
		if act, ok := activityMap[item.ActivityID]; ok {
			activityName = act.Name
		}

		list[i] = resp.BargainRecordResp{
			ID:                item.ID,
			UserID:            item.UserID,
			UserNickname:      nickname,
			UserAvatar:        avatar,
			ActivityID:        item.ActivityID,
			ActivityName:      activityName,
			SpuID:             item.SpuID,
			SkuID:             item.SkuID,
			BargainFirstPrice: item.BargainFirstPrice,
			BargainPrice:      item.BargainPrice,
			Status:            item.Status,
			EndTime:           item.EndTime,
			OrderID:           item.OrderID,
			CreatedAt:         item.CreatedAt,
		}
	}

	core.WriteSuccess(c, core.PageResult[resp.BargainRecordResp]{
		List:  list,
		Total: pageResult.Total,
	})
}
