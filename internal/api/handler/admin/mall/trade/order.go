package trade

import (
	trade2 "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/mall/trade"
	memberContract "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/member"
	"github.com/wxlbd/ruoyi-mall-go/internal/pkg/area"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/mall/trade"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/member"
	"github.com/wxlbd/ruoyi-mall-go/pkg/context"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"
	"github.com/wxlbd/ruoyi-mall-go/pkg/utils"

	"github.com/gin-gonic/gin"
)

type TradeOrderHandler struct {
	svc                        *trade.TradeOrderUpdateService
	querySvc                   *trade.TradeOrderQueryService
	memberSvc                  *member.MemberUserService
	deliveryFreightTemplateSvc *trade.DeliveryExpressTemplateService
}

func NewTradeOrderHandler(svc *trade.TradeOrderUpdateService, querySvc *trade.TradeOrderQueryService, memberSvc *member.MemberUserService, deliveryFreightTemplateSvc *trade.DeliveryExpressTemplateService) *TradeOrderHandler {
	return &TradeOrderHandler{
		svc:                        svc,
		querySvc:                   querySvc,
		memberSvc:                  memberSvc,
		deliveryFreightTemplateSvc: deliveryFreightTemplateSvc,
	}
}

// GetOrderPage 获得订单分页
func (h *TradeOrderHandler) GetOrderPage(c *gin.Context) {
	var r trade2.TradeOrderPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	// Call Service
	pageResult, err := h.querySvc.GetOrderPageForAdmin(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	// Fetch Items
	var resultList []trade2.TradeOrderPageItemResp
	if len(pageResult.List) > 0 {
		orderIds := make([]int64, len(pageResult.List))
		userIds := make([]int64, 0, len(pageResult.List))
		brokerageUserIds := make([]int64, 0, len(pageResult.List))

		for i, o := range pageResult.List {
			orderIds[i] = o.ID
			userIds = append(userIds, o.UserID)
			if o.BrokerageUserID != nil {
				brokerageUserIds = append(brokerageUserIds, *o.BrokerageUserID)
			}
		}

		// Query Items
		items, err := h.querySvc.GetOrderItemListByOrderIds(c, orderIds)
		if err != nil {
			response.WriteBizError(c, err)
			return
		}
		itemMap := make(map[int64][]trade2.TradeOrderItemBase)
		for _, item := range items {
			itemResp := trade2.TradeOrderItemBase{
				ID:       item.ID,
				UserID:   item.UserID,
				OrderID:  item.OrderID,
				SpuID:    item.SpuID,
				SpuName:  item.SpuName,
				SkuID:    item.SkuID,
				PicURL:   item.PicURL,
				Count:    item.Count,
				Price:    item.Price,
				PayPrice: item.PayPrice,
			}
			itemMap[item.OrderID] = append(itemMap[item.OrderID], itemResp)
		}

		// Query Users
		userMap, err := h.memberSvc.GetUserRespMap(c, userIds)
		if err != nil {
			// Log error but continue? Or fail? Java fails if user query fails usually.
			response.WriteBizError(c, err)
			return
		}

		// Query Brokerage Users
		var brokerageUserMap map[int64]*memberContract.MemberUserResp
		if len(brokerageUserIds) > 0 {
			brokerageUserMap, err = h.memberSvc.GetUserRespMap(c, brokerageUserIds)
			if err != nil {
				response.WriteBizError(c, err)
				return
			}
		}

		resultList = make([]trade2.TradeOrderPageItemResp, len(pageResult.List))
		for i, o := range pageResult.List {
			var payOrderID int64
			if o.PayOrderID != nil {
				payOrderID = *o.PayOrderID
			}

			var brokerageUser *memberContract.MemberUserResp
			if o.BrokerageUserID != nil {
				brokerageUser = brokerageUserMap[*o.BrokerageUserID]
			}

			resultList[i] = trade2.TradeOrderPageItemResp{
				TradeOrderBase: trade2.TradeOrderBase{
					ID:                    o.ID,
					No:                    o.No,
					CreateTime:            o.CreateTime,
					Type:                  o.Type,
					Terminal:              o.Terminal,
					UserID:                o.UserID,
					UserIP:                o.UserIP,
					UserRemark:            o.UserRemark,
					Status:                o.Status,
					ProductCount:          o.ProductCount,
					FinishTime:            o.FinishTime,
					CancelTime:            o.CancelTime,
					CancelType:            o.CancelType,
					Remark:                o.Remark,
					PayOrderID:            payOrderID,
					PayStatus:             bool(o.PayStatus),
					PayTime:               o.PayTime,
					PayChannelCode:        o.PayChannelCode,
					TotalPrice:            o.TotalPrice,
					DiscountPrice:         o.DiscountPrice,
					DeliveryPrice:         o.DeliveryPrice,
					AdjustPrice:           o.AdjustPrice,
					PayPrice:              o.PayPrice,
					DeliveryType:          o.DeliveryType,
					LogisticsID:           o.LogisticsID,
					LogisticsNo:           o.LogisticsNo,
					DeliveryTime:          o.DeliveryTime,
					ReceiveTime:           o.ReceiveTime,
					ReceiverName:          o.ReceiverName,
					ReceiverMobile:        o.ReceiverMobile,
					ReceiverAreaID:        int32(o.ReceiverAreaID),
					ReceiverDetailAddress: o.ReceiverDetailAddress,
					RefundPrice:           o.RefundPrice,
					CouponID:              o.CouponID,
					CouponPrice:           o.CouponPrice,
					PointPrice:            o.PointPrice,
					VipPrice:              o.VipPrice,
				},
				Items:            itemMap[o.ID],
				User:             userMap[o.UserID],
				BrokerageUser:    brokerageUser,
				ReceiverAreaName: area.Format(int(o.ReceiverAreaID)), // 地区名称查询
			}
		}
	} else {
		resultList = []trade2.TradeOrderPageItemResp{}
	}

	response.WriteSuccess(c, pagination.PageResult[trade2.TradeOrderPageItemResp]{
		List:  resultList,
		Total: pageResult.Total,
	})
}

// GetOrderDetail 获得订单详情
func (h *TradeOrderHandler) GetOrderDetail(c *gin.Context) {
	id := utils.ParseInt64(c.Query("id"))
	if id == 0 {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	// 1. Get Order
	order, err := h.querySvc.GetOrder(c, id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	if order == nil {
		response.WriteBizError(c, errors.ErrNotFound)
		return
	}

	// 2. Get Items
	items, err := h.querySvc.GetOrderItemListByOrderId(c, order.ID)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	// 3. Get Logs
	logs, err := h.querySvc.GetOrderLogListByOrderId(c, order.ID)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	// 4. Query User and Brokerage User
	var user *memberContract.MemberUserResp
	if order.UserID > 0 {
		userMap, err := h.memberSvc.GetUserRespMap(c, []int64{order.UserID})
		if err == nil && userMap != nil {
			user = userMap[order.UserID]
		}
	}

	var brokerageUser *memberContract.MemberUserResp
	if order.BrokerageUserID != nil && *order.BrokerageUserID > 0 {
		brokerageUserMap, err := h.memberSvc.GetUserRespMap(c, []int64{*order.BrokerageUserID})
		if err == nil && brokerageUserMap != nil {
			brokerageUser = brokerageUserMap[*order.BrokerageUserID]
		}
	}

	// 5. Get Receiver Area Name
	receiverAreaName := area.Format(int(order.ReceiverAreaID))
	itemResps := make([]trade2.TradeOrderItemBase, len(items))
	for i, item := range items {
		itemResps[i] = trade2.TradeOrderItemBase{
			ID:       item.ID,
			UserID:   item.UserID,
			OrderID:  item.OrderID,
			SpuID:    item.SpuID,
			SpuName:  item.SpuName,
			SkuID:    item.SkuID,
			PicURL:   item.PicURL,
			Count:    item.Count,
			Price:    item.Price,
			PayPrice: item.PayPrice,
			// ... other fields
		}
	}

	logResps := make([]trade2.TradeOrderLogResp, len(logs))
	for i, l := range logs {
		logResps[i] = trade2.TradeOrderLogResp{
			Content:    l.Content,
			CreateTime: l.CreateTime,
			UserType:   l.UserType,
		}
	}

	var payOrderID int64
	if order.PayOrderID != nil {
		payOrderID = *order.PayOrderID
	}

	// Prepare User Info if needed, skip for now

	res := trade2.TradeOrderDetailResp{
		TradeOrderBase: trade2.TradeOrderBase{
			ID:                    order.ID,
			No:                    order.No,
			Type:                  order.Type,
			Terminal:              order.Terminal,
			UserID:                order.UserID,
			UserIP:                order.UserIP,
			UserRemark:            order.UserRemark,
			Status:                order.Status,
			ProductCount:          order.ProductCount,
			FinishTime:            order.FinishTime,
			CancelTime:            order.CancelTime,
			CancelType:            order.CancelType,
			Remark:                order.Remark,
			PayOrderID:            payOrderID,
			PayStatus:             bool(order.PayStatus),
			PayTime:               order.PayTime,
			PayChannelCode:        order.PayChannelCode,
			TotalPrice:            order.TotalPrice,
			DiscountPrice:         order.DiscountPrice,
			DeliveryPrice:         order.DeliveryPrice,
			AdjustPrice:           order.AdjustPrice,
			PayPrice:              order.PayPrice,
			DeliveryType:          order.DeliveryType,
			LogisticsID:           order.LogisticsID,
			LogisticsNo:           order.LogisticsNo,
			DeliveryTime:          order.DeliveryTime,
			ReceiveTime:           order.ReceiveTime,
			ReceiverName:          order.ReceiverName,
			ReceiverMobile:        order.ReceiverMobile,
			ReceiverAreaID:        int32(order.ReceiverAreaID),
			ReceiverDetailAddress: order.ReceiverDetailAddress,
			RefundPrice:           order.RefundPrice,
			CouponID:              order.CouponID,
			CouponPrice:           order.CouponPrice,
			PointPrice:            order.PointPrice,
			VipPrice:              order.VipPrice,
		},
		Items:            itemResps,
		Logs:             logResps,
		User:             user,
		BrokerageUser:    brokerageUser,
		ReceiverAreaName: receiverAreaName, // 地区名称
	}

	response.WriteSuccess(c, res)
}

// GetOrderExpressTrackList 获得交易订单的物流轨迹
func (h *TradeOrderHandler) GetOrderExpressTrackList(c *gin.Context) {
	id := utils.ParseInt64(c.Query("id"))
	if id == 0 {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	tracks, err := h.querySvc.GetExpressTrackListById(c, id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, tracks)
}

// DeliveryOrder 订单发货
func (h *TradeOrderHandler) DeliveryOrder(c *gin.Context) {
	var r trade2.TradeOrderDeliveryReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	if err := h.svc.DeliveryOrder(c, &r); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// UpdateOrderRemark 订单备注
func (h *TradeOrderHandler) UpdateOrderRemark(c *gin.Context) {
	var r trade2.TradeOrderRemarkReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	if err := h.svc.UpdateOrderRemark(c, &r); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// UpdateOrderPrice 订单调价
func (h *TradeOrderHandler) UpdateOrderPrice(c *gin.Context) {
	var r trade2.TradeOrderUpdatePriceReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	if err := h.svc.UpdateOrderPrice(c, &r); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// UpdateOrderAddress 修改订单收货地址
func (h *TradeOrderHandler) UpdateOrderAddress(c *gin.Context) {
	var r trade2.TradeOrderUpdateAddressReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	if err := h.svc.UpdateOrderAddress(c, &r); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// PickUpOrderById 订单核销 (By ID)
func (h *TradeOrderHandler) PickUpOrderById(c *gin.Context) {
	id := utils.ParseInt64(c.Query("id"))
	if id == 0 {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	if err := h.svc.PickUpOrderByAdmin(c, context.GetUserId(c), id); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// PickUpOrderByVerifyCode 订单核销 (By Code)
func (h *TradeOrderHandler) PickUpOrderByVerifyCode(c *gin.Context) {
	code := c.Query("pickUpVerifyCode")
	if code == "" {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	if err := h.svc.PickUpOrderByVerifyCode(c, context.GetUserId(c), code); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// GetByPickUpVerifyCode 查询核销码对应的订单
func (h *TradeOrderHandler) GetByPickUpVerifyCode(c *gin.Context) {
	code := c.Query("pickUpVerifyCode")
	if code == "" {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	res, err := h.svc.GetByPickUpVerifyCode(c, code)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, res)
}

// GetOrderSummary 获得交易订单统计
func (h *TradeOrderHandler) GetOrderSummary(c *gin.Context) {
	var r trade2.TradeOrderPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	res, err := h.querySvc.GetOrderSummary(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, res)
}
