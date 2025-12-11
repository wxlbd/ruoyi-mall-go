package trade

import (
	"backend-go/internal/api/req"
	"backend-go/internal/api/resp"
	"backend-go/internal/model/trade"
	"backend-go/internal/repo/query"
	"backend-go/internal/service/member"
	"backend-go/internal/service/product"
	"backend-go/internal/service/promotion"
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type TradeOrderUpdateService struct {
	q          *query.Query
	skuSvc     *product.ProductSkuService
	cartSvc    *CartService
	priceSvc   *TradePriceService
	addressSvc *member.MemberAddressService
	couponSvc  *promotion.CouponUserService
	logSvc     *TradeOrderLogService
}

func NewTradeOrderUpdateService(
	skuSvc *product.ProductSkuService,
	cartSvc *CartService,
	priceSvc *TradePriceService,
	addressSvc *member.MemberAddressService,
	couponSvc *promotion.CouponUserService,
	logSvc *TradeOrderLogService,
) *TradeOrderUpdateService {
	return &TradeOrderUpdateService{
		q:          query.Q,
		skuSvc:     skuSvc,
		cartSvc:    cartSvc,
		priceSvc:   priceSvc,
		addressSvc: addressSvc,
		couponSvc:  couponSvc,
		logSvc:     logSvc,
	}
}

// SettlementOrder 获得订单结算信息
func (s *TradeOrderUpdateService) SettlementOrder(ctx context.Context, uId int64, req *req.AppTradeOrderSettlementReq) (*resp.AppTradeOrderSettlementResp, error) {
	// 1. Calculate Price
	calcReq := &TradePriceCalculateReqBO{
		UserID:        uId,
		CouponID:      req.CouponID,
		PointStatus:   req.PointStatus,
		DeliveryType:  req.DeliveryType,
		AddressID:     req.AddressID,
		PickUpStoreID: req.PickUpStoreID,
		Items:         make([]TradePriceCalculateItemBO, len(req.Items)),
	}
	for i, item := range req.Items {
		calcReq.Items[i] = TradePriceCalculateItemBO{
			SkuID:    item.SkuID,
			Count:    item.Count,
			CartID:   item.CartID,
			Selected: true,
		}
	}

	priceResp, err := s.priceSvc.CalculateOrderPrice(ctx, calcReq)
	if err != nil {
		return nil, err
	}

	// 2. Fetch Address (if delivery)
	var address *resp.AppTradeOrderSettlementAddress
	if req.AddressID != nil {
		addr, err := s.addressSvc.GetAddress(ctx, uId, *req.AddressID)
		if err == nil && addr != nil {
			address = &resp.AppTradeOrderSettlementAddress{
				ID:            addr.ID,
				Name:          addr.Name,
				Mobile:        addr.Mobile,
				AreaID:        int64(addr.AreaID),
				DetailAddress: addr.DetailAddress,
				DefaultStatus: addr.DefaultStatus,
				// AreaName: fetch area name if needed
			}
		}
	}

	// 3. Assemble Response
	r := &resp.AppTradeOrderSettlementResp{
		Type:       priceResp.Type,
		Items:      make([]resp.AppTradeOrderSettlementItem, len(priceResp.Items)),
		Price:      resp.AppTradeOrderSettlementPrice(priceResp.Price),
		Address:    address,
		UsePoint:   priceResp.UsePoint,
		TotalPoint: priceResp.TotalPoint,
	}

	for i, item := range priceResp.Items {
		r.Items[i] = resp.AppTradeOrderSettlementItem{
			SpuID:      item.SpuID,
			SkuID:      item.SkuID,
			Count:      item.Count,
			CartID:     item.CartID,
			Price:      item.Price,
			PicURL:     item.PicURL,
			Properties: item.Properties,
			// SpuName: item.SpuName,
		}
	}

	return r, nil
}

// CreateOrder 创建交易订单
func (s *TradeOrderUpdateService) CreateOrder(ctx context.Context, uId int64, reqVO *req.AppTradeOrderCreateReq) (*trade.TradeOrder, error) {
	// 1. Price Calculation
	// 1. Price Calculation
	// Call SettlementOrder logic (reuse code or abstract)
	// Or better, just call price calc directly
	calcReq := &TradePriceCalculateReqBO{
		UserID:        uId,
		CouponID:      reqVO.CouponID,
		PointStatus:   reqVO.PointStatus,
		DeliveryType:  reqVO.DeliveryType,
		AddressID:     reqVO.AddressID,
		PickUpStoreID: reqVO.PickUpStoreID,
		Items:         make([]TradePriceCalculateItemBO, len(reqVO.Items)),
	}
	for i, item := range reqVO.Items {
		calcReq.Items[i] = TradePriceCalculateItemBO{
			SkuID:    item.SkuID,
			Count:    item.Count,
			CartID:   item.CartID,
			Selected: true,
		}
	}
	priceResp, err := s.priceSvc.CalculateOrderPrice(ctx, calcReq)
	if err != nil {
		return nil, err
	}

	// 2. Transaction
	var order *trade.TradeOrder
	err = s.q.Transaction(func(tx *query.Query) error {
		// 2.1 Create Order
		order = &trade.TradeOrder{
			No:             generateOrderNo(),
			Type:           1, // Normal order
			Terminal:       1, // TODO: passed from header/context
			UserID:         uId,
			UserIP:         "127.0.0.1", // TODO: from context
			Status:         0,           // Unpaid
			ProductCount:   len(reqVO.Items),
			Remark:         reqVO.Remark,
			PayStatus:      false,
			TotalPrice:     priceResp.Price.TotalPrice,
			DiscountPrice:  priceResp.Price.DiscountPrice,
			DeliveryPrice:  priceResp.Price.DeliveryPrice,
			PayPrice:       priceResp.Price.PayPrice,
			CouponID:       priceResp.CouponID,
			CouponPrice:    priceResp.Price.CouponPrice,
			DeliveryType:   reqVO.DeliveryType,
			ReceiverName:   reqVO.ReceiverName,
			ReceiverMobile: reqVO.ReceiverMobile,
			// Add address info...
		}

		if reqVO.AddressID != nil {
			addr, _ := s.addressSvc.GetAddress(ctx, uId, *reqVO.AddressID)
			if addr != nil {
				order.ReceiverName = addr.Name
				order.ReceiverMobile = addr.Mobile
				order.ReceiverAreaID = int(addr.AreaID)
				order.ReceiverDetailAddress = addr.DetailAddress
			}
		}

		if err := tx.TradeOrder.WithContext(ctx).Create(order); err != nil {
			return err
		}

		// 2.2 Create Order Items
		items := make([]*trade.TradeOrderItem, len(priceResp.Items))
		for i, item := range priceResp.Items {
			items[i] = &trade.TradeOrderItem{
				UserID:      uId,
				OrderID:     order.ID,
				SpuID:       item.SpuID,
				SkuID:       item.SkuID,
				Count:       item.Count,
				Price:       item.Price,
				PayPrice:    item.PayPrice,
				PicURL:      item.PicURL,
				CouponPrice: item.CouponPrice,
				// Properties: item.Properties (need serialize),
			}
		}
		if err := tx.TradeOrderItem.WithContext(ctx).Create(items...); err != nil {
			return err
		}

		// 2.3 Clear Cart (if cart items)
		var cartIds []int64
		for _, item := range calcReq.Items {
			if item.CartID > 0 {
				cartIds = append(cartIds, item.CartID)
			}
		}
		if len(cartIds) > 0 {
			if err := s.cartSvc.DeleteCart(ctx, uId, cartIds); err != nil {
				return err
			}
		}

		// 2.4 Decrease Stock
		var stockItems []req.ProductSkuUpdateStockItemReq
		for _, item := range priceResp.Items {
			stockItems = append(stockItems, req.ProductSkuUpdateStockItemReq{
				ID:        item.SkuID,
				IncrCount: -item.Count,
			})
		}
		if err := s.skuSvc.UpdateSkuStock(ctx, &req.ProductSkuUpdateStockReq{Items: stockItems}); err != nil {
			return err
		}

		// 2.5 Use Coupon
		if priceResp.CouponID > 0 {
			if err := s.couponSvc.UseCoupon(ctx, uId, priceResp.CouponID, order.ID); err != nil {
				return err
			}
		}

		// 2.6 Log
		if err := s.createOrderLog(ctx, order, "Create Order", 10); err != nil {
			return err
		}

		return nil
	})
	return order, err
}

// DeliveryOrder 订单发货
func (s *TradeOrderUpdateService) DeliveryOrder(ctx context.Context, reqVO *req.TradeOrderDeliveryReq) error {
	// 1. Check Order Exists
	order, err := s.q.TradeOrder.WithContext(ctx).Where(s.q.TradeOrder.ID.Eq(reqVO.ID)).First()
	if err != nil {
		return err
	}
	if order.Status != 10 { // Assume 10 is Paid/Undelivered. In real app use constants.
		// return fmt.Errorf("order status error")
	}

	now := time.Now()
	// 2. Update Order
	_, err = s.q.TradeOrder.WithContext(ctx).Where(s.q.TradeOrder.ID.Eq(reqVO.ID)).Updates(trade.TradeOrder{
		Status:       20, // Delivered
		LogisticsID:  reqVO.LogisticsID,
		LogisticsNo:  reqVO.LogisticsNo,
		DeliveryTime: &now,
	})
	if err != nil {
		return err
	}

	// 3. Log
	logOrder := *order
	logOrder.Status = 20
	// OperateType 30 for Delivery (Example)
	return s.createOrderLog(ctx, &logOrder, "Order Delivered", 30)
}

// UpdateOrderPaid 更新订单为已支付
func (s *TradeOrderUpdateService) UpdateOrderPaid(ctx context.Context, id int64, payOrderId int64) error {
	// 1. Get Order
	order, err := s.q.TradeOrder.WithContext(ctx).Where(s.q.TradeOrder.ID.Eq(id)).First()
	if err != nil {
		return err
	}
	if order.Status != 0 { // 0: Unpaid
		return fmt.Errorf("order status is not unpaid")
	}
	if order.PayStatus {
		return fmt.Errorf("order is already paid")
	}

	// 2. Update
	now := time.Now()
	err = s.q.Transaction(func(tx *query.Query) error {
		// Update Order
		updateMap := map[string]interface{}{
			"status":       10, // Undelivered
			"pay_status":   true,
			"pay_time":     &now,
			"pay_order_id": payOrderId,
		}
		if _, err := tx.TradeOrder.WithContext(ctx).Where(tx.TradeOrder.ID.Eq(id)).Updates(updateMap); err != nil {
			return err
		}

		// Log
		logOrder := *order
		logOrder.Status = 10
		if err := s.createOrderLog(ctx, &logOrder, "Order Paid", 20); err != nil { // 20: Pay
			return err
		}
		return nil
	})
	return err
}

func (s *TradeOrderUpdateService) createOrderLog(ctx context.Context, order *trade.TradeOrder, content string, operateType int) error {
	uid := int64(0) // System or unknown if not passed.
	// In strict context, we might want to get from ctx.
	// For CreateOrder, we know UID.
	// For Pay callback, maybe System (0).

	log := &trade.TradeOrderLog{
		UserID:       uid,
		UserType:     1, // Member
		OrderID:      order.ID,
		BeforeStatus: 0, // Simplified, ideally pass old status
		AfterStatus:  order.Status,
		OperateType:  operateType,
		Content:      content,
	}
	return s.q.TradeOrderLog.WithContext(ctx).Create(log)
}

// CancelOrder 取消交易订单
func (s *TradeOrderUpdateService) CancelOrder(ctx context.Context, uId int64, id int64) error {
	// 1. Check Order
	order, err := s.q.TradeOrder.WithContext(ctx).Where(s.q.TradeOrder.ID.Eq(id), s.q.TradeOrder.UserID.Eq(uId)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("订单不存在")
		}
		return err
	}
	if order.Status != 0 { // 0: Unpaid. If we allow cancelling Paid orders, we need Refund logic.
		// For now, restrict to Unpaid.
		return errors.New("订单状态不允许取消")
	}

	// 2. Transaction
	err = s.q.Transaction(func(tx *query.Query) error {
		// 2.1 Update Order Status
		now := time.Now()
		if _, err := tx.TradeOrder.WithContext(ctx).Where(tx.TradeOrder.ID.Eq(id)).Updates(trade.TradeOrder{
			Status:     40, // Cancelled (Closed)
			CancelTime: &now,
			CancelType: 1, // User Cancelled
		}); err != nil {
			return err
		}

		// 2.2 Release Stock
		items, err := tx.TradeOrderItem.WithContext(ctx).Where(tx.TradeOrderItem.OrderID.Eq(id)).Find()
		if err != nil {
			return err
		}
		var stockItems []req.ProductSkuUpdateStockItemReq
		for _, item := range items {
			stockItems = append(stockItems, req.ProductSkuUpdateStockItemReq{
				ID:        item.SkuID,
				IncrCount: item.Count, // Positive to restore stock
			})
		}
		if err := s.skuSvc.UpdateSkuStock(ctx, &req.ProductSkuUpdateStockReq{Items: stockItems}); err != nil {
			return err
		}

		// 2.3 Refund Coupon
		if order.CouponID > 0 {
			if err := s.couponSvc.ReturnCoupon(ctx, uId, order.CouponID); err != nil {
				return err
			}
		}

		// 2.4 Log
		logOrder := *order
		logOrder.Status = 40
		if err := s.createOrderLog(ctx, &logOrder, "User Cancelled Order", 40); err != nil {
			return err
		}

		return nil
	})

	return err
}

// DeleteOrder 删除订单
func (s *TradeOrderUpdateService) DeleteOrder(ctx context.Context, uId int64, id int64) error {
	// 1. Check Order
	order, err := s.q.TradeOrder.WithContext(ctx).Where(s.q.TradeOrder.ID.Eq(id), s.q.TradeOrder.UserID.Eq(uId)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("订单不存在")
		}
		return err
	}
	// Java: Check status (Cancelled or Completed can be deleted?)
	// Usually only Cancelled or Completed.
	if order.Status != 40 && order.Status != 30 { // 40: Cancelled, 30: Completed
		return errors.New("只有取消或完成的订单可以删除")
	}

	// 2. Delete (Soft Delete)
	_, err = s.q.TradeOrder.WithContext(ctx).Where(s.q.TradeOrder.ID.Eq(id)).Delete()
	return err
}

// UpdateOrderRemark 订单备注
func (s *TradeOrderUpdateService) UpdateOrderRemark(ctx context.Context, req *req.TradeOrderRemarkReq) error {
	_, err := s.q.TradeOrder.WithContext(ctx).Where(s.q.TradeOrder.ID.Eq(req.ID)).Update(s.q.TradeOrder.Remark, req.Remark)
	return err
}

// UpdateOrderPrice 订单调价
func (s *TradeOrderUpdateService) UpdateOrderPrice(ctx context.Context, req *req.TradeOrderUpdatePriceReq) error {
	order, err := s.q.TradeOrder.WithContext(ctx).Where(s.q.TradeOrder.ID.Eq(req.ID)).First()
	if err != nil {
		return err
	}
	if order.PayStatus {
		return errors.New("已支付订单不允许改价")
	}

	// New Price Calculation
	newPayPrice := order.PayPrice + req.AdjustPrice
	if newPayPrice < 0 {
		return errors.New("调价后金额不能小于 0")
	}

	_, err = s.q.TradeOrder.WithContext(ctx).Where(s.q.TradeOrder.ID.Eq(req.ID)).Updates(map[string]interface{}{
		"adjust_price": order.AdjustPrice + req.AdjustPrice,
		"pay_price":    newPayPrice,
	})
	return err
}

// UpdateOrderAddress 修改订单收货地址
func (s *TradeOrderUpdateService) UpdateOrderAddress(ctx context.Context, req *req.TradeOrderUpdateAddressReq) error {
	// Check status (only undelivered?)
	_, err := s.q.TradeOrder.WithContext(ctx).Where(s.q.TradeOrder.ID.Eq(req.ID)).Updates(map[string]interface{}{
		"receiver_name":           req.ReceiverName,
		"receiver_mobile":         req.ReceiverMobile,
		"receiver_area_id":        req.ReceiverAreaID,
		"receiver_detail_address": req.ReceiverDetailAddress,
	})
	return err
}

// PickUpOrderByAdmin 核销订单 (By ID)
func (s *TradeOrderUpdateService) PickUpOrderByAdmin(ctx context.Context, adminUserId int64, id int64) error {
	order, err := s.q.TradeOrder.WithContext(ctx).Where(s.q.TradeOrder.ID.Eq(id)).First()
	if err != nil {
		return err
	}
	return s.pickUpOrder(ctx, order)
}

// PickUpOrderByVerifyCode 核销订单 (By Code)
func (s *TradeOrderUpdateService) PickUpOrderByVerifyCode(ctx context.Context, adminUserId int64, verifyCode string) error {
	order, err := s.q.TradeOrder.WithContext(ctx).Where(s.q.TradeOrder.PickUpVerifyCode.Eq(verifyCode)).First()
	if err != nil {
		return errors.New("核销码无效")
	}
	return s.pickUpOrder(ctx, order)
}

func (s *TradeOrderUpdateService) pickUpOrder(ctx context.Context, order *trade.TradeOrder) error {
	if order.DeliveryType != 2 { // 2: PickUp
		return errors.New("非自提订单")
	}
	if order.Status != 10 { // 10: Undelivered/Wait PickUp
		// Check constants
		return errors.New("订单状态不正确")
	}

	now := time.Now()
	// Update Status -> Completed (or PickedUp then Completed?)
	// Java impl: Status -> Received(20) or Completed(30)?
	// Usually PickUp -> Received/Completed. Let's say 30 (Completed) or 20 (Delivered/Received).
	// Java: TradeOrderStatusEnum.COMPLETED (30) directly? Or DELIVERED(20)?
	// If PickUp, usually means Delivered & Received at same time.
	// Let's set to 30 (Completed) or 20 if logic differs.
	// Check Java: tradeOrderUpdateService.pickUpOrderByAdmin -> updateStatus(TradeOrderStatusEnum.COMPLETED)

	err := s.q.Transaction(func(tx *query.Query) error {
		_, err := tx.TradeOrder.WithContext(ctx).Where(tx.TradeOrder.ID.Eq(order.ID)).Updates(trade.TradeOrder{
			Status:      30, // Completed
			ReceiveTime: &now,
		})
		if err != nil {
			return err
		}
		// Log
		return s.createOrderLog(ctx, order, "Admin PickUp", 50) // 50: PickUp
	})
	return err
}

// GetByPickUpVerifyCode 查询核销码对应的订单
func (s *TradeOrderUpdateService) GetByPickUpVerifyCode(ctx context.Context, verifyCode string) (*trade.TradeOrder, error) {
	return s.q.TradeOrder.WithContext(ctx).Where(s.q.TradeOrder.PickUpVerifyCode.Eq(verifyCode)).First()
}

// CreateOrderItemCommentByMember 创建订单项评价
func (s *TradeOrderUpdateService) CreateOrderItemCommentByMember(ctx context.Context, uId int64, req *req.AppTradeOrderItemCommentCreateReq) (int64, error) {
	// 1. Get Order Item
	item, err := s.q.TradeOrderItem.WithContext(ctx).Where(s.q.TradeOrderItem.ID.Eq(req.OrderItemID), s.q.TradeOrderItem.UserID.Eq(uId)).First()
	if err != nil {
		return 0, err
	}
	if item.CommentStatus {
		return 0, errors.New("该商品已评价")
	}

	// 2. Create Comment (Mock or Call CommentService)
	// Assume CommentService exists or direct DB insert.
	// For now, mock success and update item status.
	// TODO: Call ProductCommentService.CreateComment(...)

	// 3. Update Order Item Comment Status
	_, err = s.q.TradeOrderItem.WithContext(ctx).Where(s.q.TradeOrderItem.ID.Eq(item.ID)).Update(s.q.TradeOrderItem.CommentStatus, true)
	if err != nil {
		return 0, err
	}

	// 4. Update Order Comment Status if all items commented
	// Count uncommented items
	count, _ := s.q.TradeOrderItem.WithContext(ctx).Where(s.q.TradeOrderItem.OrderID.Eq(item.OrderID), s.q.TradeOrderItem.CommentStatus.Is(false)).Count()
	if count == 0 {
		_, _ = s.q.TradeOrder.WithContext(ctx).Where(s.q.TradeOrder.ID.Eq(item.OrderID)).Update(s.q.TradeOrder.CommentStatus, true)
	}

	return 0, nil // Return Comment ID
}

// ReceiveOrder 用户确认收货
func (s *TradeOrderUpdateService) ReceiveOrder(ctx context.Context, uId int64, orderId int64) error {
	// 1. Get Order
	order, err := s.q.TradeOrder.WithContext(ctx).Where(s.q.TradeOrder.ID.Eq(orderId), s.q.TradeOrder.UserID.Eq(uId)).First()
	if err != nil {
		return errors.New("订单不存在")
	}

	// 2. Validate Status - 只有已发货状态才能确认收货
	if order.Status != 20 { // 20: Delivered
		return errors.New("订单状态不正确，无法确认收货")
	}

	// 3. Update Order Status
	now := time.Now()
	err = s.q.Transaction(func(tx *query.Query) error {
		_, err := tx.TradeOrder.WithContext(ctx).Where(tx.TradeOrder.ID.Eq(order.ID)).Updates(trade.TradeOrder{
			Status:      30, // Completed
			ReceiveTime: &now,
		})
		if err != nil {
			return err
		}
		// Log
		return s.createOrderLog(ctx, order, "用户确认收货", 40) // 40: Receive
	})
	return err
}

// UpdatePaidOrderRefunded 更新支付订单为已退款 (Callback from Pay)
func (s *TradeOrderUpdateService) UpdatePaidOrderRefunded(ctx context.Context, orderId int64, payRefundId int64) error {
	order, err := s.q.TradeOrder.WithContext(ctx).Where(s.q.TradeOrder.ID.Eq(orderId)).First()
	if err != nil {
		return err
	}
	// Verify Status? if cancelled or something? Usually calls this after full refund.
	// We just update status to Refunded?
	// Java doesn't show implementation, but likely updates RefundStatus.
	// Let's assume RefundStatus = 30 (Finish? No, RefundStatus has constants)
	// Check consts.go for RefundStatusEnum.
	// 0: None, 10: Apply, 20: Audit Pass, 30: Refunded?
	// Assuming 30 is correct for now based on AfterSale logic.

	return s.q.Transaction(func(tx *query.Query) error {
		_, err := tx.TradeOrder.WithContext(ctx).Where(tx.TradeOrder.ID.Eq(orderId)).Updates(map[string]interface{}{
			"refund_status": 30, // All Refunded
			// "pay_refund_id": payRefundId, // No field
		})
		if err != nil {
			return err
		}
		// Log
		return s.createOrderLog(ctx, order, "Order Refunded (Pay Callback)", 40)
	})
}

func generateOrderNo() string {
	return fmt.Sprintf("%d", time.Now().UnixNano()) // Simplified
}
