package trade

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/mall/product"
	trade2 "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/mall/trade"
	memberContract "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/member"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/pay"
	"github.com/wxlbd/ruoyi-mall-go/internal/consts"
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	memberModel "github.com/wxlbd/ruoyi-mall-go/internal/model/member"
	tradeModel "github.com/wxlbd/ruoyi-mall-go/internal/model/trade"
	"github.com/wxlbd/ruoyi-mall-go/internal/pkg/area"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/member"
	pkgErrors "github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"go.uber.org/zap"
)

type memberAddressService interface {
	GetAddress(ctx context.Context, userId int64, id int64) (*memberContract.AppAddressResp, error)
	GetDefaultAddress(ctx context.Context, userId int64) (*memberContract.AppAddressResp, error)
}

// TradeOrderUpdateService 订单更新服务
// 使用责任链模式重构，内部使用订单处理器
type TradeOrderUpdateService struct {
	q            *query.Query
	manager      *OrderHandlerManager
	priceSvc     *TradePriceService
	cartSvc      *CartService
	addressSvc   memberAddressService
	paySvc       PayOrderServiceAPI
	payRefundSvc PayRefundServiceAPI
	payAppSvc    PayAppServiceAPI
	configSvc    TradeConfigServiceAPI
	skuSvc       ProductSkuServiceAPI
	commentSvc   ProductCommentServiceAPI
	couponSvc    CouponUserServiceAPI
	memberSvc    MemberUserServiceAPI
	logSvc       *TradeOrderLogService
	noDAO        TradeNoRedisDAOAPI
	logger       *zap.Logger
}

// NewTradeOrderUpdateService 创建订单更新服务
func NewTradeOrderUpdateService(
	q *query.Query,
	priceSvc *TradePriceService,
	cartSvc *CartService,
	addressSvc *member.MemberAddressService,
	paySvc PayOrderServiceAPI,
	payRefundSvc PayRefundServiceAPI,
	payAppSvc PayAppServiceAPI,
	configSvc TradeConfigServiceAPI,
	skuSvc ProductSkuServiceAPI,
	commentSvc ProductCommentServiceAPI,
	couponSvc CouponUserServiceAPI,
	memberSvc MemberUserServiceAPI,
	logSvc *TradeOrderLogService,
	noDAO TradeNoRedisDAOAPI,
	logger *zap.Logger,
) *TradeOrderUpdateService {
	service := &TradeOrderUpdateService{
		q:            q,
		manager:      NewOrderHandlerManager(logger),
		priceSvc:     priceSvc,
		cartSvc:      cartSvc,
		addressSvc:   addressSvc,
		paySvc:       paySvc,
		payRefundSvc: payRefundSvc,
		payAppSvc:    payAppSvc,
		configSvc:    configSvc,
		skuSvc:       skuSvc,
		commentSvc:   commentSvc,
		couponSvc:    couponSvc,
		memberSvc:    memberSvc,
		logSvc:       logSvc,
		noDAO:        noDAO,
		logger:       logger,
	}

	// 初始化订单处理器
	if err := service.initializeProcessors(); err != nil {
		logger.Error("初始化订单处理器失败", zap.Error(err))
	}

	return service
}

// initializeProcessors 初始化订单业务处理器
func (s *TradeOrderUpdateService) initializeProcessors() error {
	processors := []OrderHandler{
		NewCreateOrderProcessor(s.q, s.skuSvc, s.couponSvc, s.memberSvc, s.logger),
		NewPayOrderProcessor(s.q, s.paySvc, s.logger),
		NewDeliveryOrderProcessor(s.q, s.logger),
		NewReceiveOrderProcessor(s.q, s.logger),
		NewCancelOrderProcessor(s.q, s.skuSvc, s.couponSvc, s.memberSvc, s.logger),
		NewRefundOrderProcessor(s.q, s.logger),
		NewPickUpOrderProcessor(s.q, s.logger),
	}

	return s.manager.Initialize(processors)
}

// InitializeHandlers 初始化订单处理器（保留兼容性）
func (s *TradeOrderUpdateService) InitializeHandlers(handlers []OrderHandler) error {
	return s.manager.Initialize(handlers)
}

// UpdatePaidOrderRefunded 更新已支付订单为已退款状态
// 对应 Java: TradeOrderUpdateServiceImpl#updatePaidOrderRefunded
// 注意：这个方法只做校验，不更新订单状态（订单取消已在 cancelPaidOrder 中完成）
func (s *TradeOrderUpdateService) UpdatePaidOrderRefunded(ctx context.Context, orderId int64, payRefundId int64) error {
	s.logger.Info("开始校验售后退款状态",
		zap.Int64("orderId", orderId),
		zap.Int64("payRefundId", payRefundId),
	)

	// 1. 获取退款单信息
	if payRefundId <= 0 {
		return NewTradeErrorWithMsg(ErrorCodeOrderUpdateError, "退款单编号不能为空")
	}

	// 2. 校验退款单是否存在
	// 对应 Java: PayRefundRespDTO payRefund = payRefundApi.getRefund(payRefundId)
	payRefund, err := s.payRefundSvc.GetRefund(ctx, payRefundId)
	if err != nil {
		s.logger.Error("获取退款单失败",
			zap.Error(err),
			zap.Int64("payRefundId", payRefundId),
		)
		return NewTradeErrorWithMsg(ErrorCodeOrderUpdateError, "退款单不存在")
	}

	// 3. 校验退款单状态必须是成功
	// 对应 Java: 特殊：因为在 cancelPaidOrder 已经进行订单的取消，所以这里必须退款成功！！！
	if payRefund.Status != 2 { // 2 = PayRefundStatusSuccess
		s.logger.Error("退款单状态不是成功",
			zap.Int64("payRefundId", payRefundId),
			zap.Int("status", payRefund.Status),
		)
		return NewTradeErrorWithMsg(ErrorCodeOrderUpdateError, "退款单状态不是成功")
	}

	s.logger.Info("售后退款状态校验成功",
		zap.Int64("orderId", orderId),
		zap.Int64("payRefundId", payRefundId),
		zap.Int("refundPrice", payRefund.RefundPrice),
	)
	return nil
}

// DeliveryOrder 订单发货
func (s *TradeOrderUpdateService) DeliveryOrder(ctx context.Context, reqVO *trade2.TradeOrderDeliveryReq) error {
	var logisticsId int64
	if reqVO.LogisticsID != nil {
		logisticsId = *reqVO.LogisticsID
	}

	s.logger.Info("开始处理订单发货请求",
		zap.Int64("orderId", reqVO.ID),
		zap.Int64("logisticsId", logisticsId),
		zap.String("logisticsNo", reqVO.LogisticsNo),
	)

	req := &OrderHandleRequest{
		Operation:   "delivery",
		OrderID:     reqVO.ID,
		LogisticsID: logisticsId,
		TrackingNo:  reqVO.LogisticsNo,
	}

	_, err := s.manager.HandleOrder(ctx, req)
	if err != nil {
		s.logger.Error("订单发货处理失败",
			zap.Error(err),
			zap.Int64("orderId", reqVO.ID),
			zap.Int64("logisticsId", logisticsId),
			zap.String("logisticsNo", reqVO.LogisticsNo),
		)
		return err
	}

	s.logger.Info("订单发货处理成功",
		zap.Int64("orderId", reqVO.ID),
		zap.Int64("logisticsId", logisticsId),
		zap.String("logisticsNo", reqVO.LogisticsNo),
	)

	// 后置流程：执行发货后置钩子
	order, _ := s.q.TradeOrder.WithContext(ctx).Where(s.q.TradeOrder.ID.Eq(reqVO.ID)).First()
	if order != nil {
		_ = s.executeAfterDeliveryOrder(ctx, order)
	}

	return nil
}

// UpdateOrderRemark 更新订单备注
// 对应 Java: TradeOrderUpdateServiceImpl#updateOrderRemark
func (s *TradeOrderUpdateService) UpdateOrderRemark(ctx context.Context, reqVO *trade2.TradeOrderRemarkReq) error {
	s.logger.Info("开始处理订单备注更新",
		zap.Int64("orderId", reqVO.ID),
		zap.String("remark", reqVO.Remark),
	)

	err := s.q.Transaction(func(tx *query.Query) error {
		// 1. 校验并获得交易订单
		// 对应 Java: validateOrderExists(reqVO.getId())
		order, err := tx.TradeOrder.WithContext(ctx).
			Where(tx.TradeOrder.ID.Eq(reqVO.ID)).
			First()
		if err != nil {
			return ErrOrderNotExists()
		}

		// 2. 更新订单备注
		_, err = tx.TradeOrder.WithContext(ctx).
			Where(tx.TradeOrder.ID.Eq(reqVO.ID)).
			Update(tx.TradeOrder.Remark, reqVO.Remark)
		if err != nil {
			return err
		}

		// 记录旧备注用于日志
		s.logger.Debug("订单备注变更",
			zap.Int64("orderId", reqVO.ID),
			zap.String("oldRemark", order.Remark),
			zap.String("newRemark", reqVO.Remark),
		)

		return nil
	})

	if err != nil {
		s.logger.Error("订单备注更新失败",
			zap.Error(err),
			zap.Int64("orderId", reqVO.ID),
			zap.String("remark", reqVO.Remark),
		)
		return err
	}

	s.logger.Info("订单备注更新成功",
		zap.Int64("orderId", reqVO.ID),
		zap.String("remark", reqVO.Remark),
	)
	return nil
}

// UpdateOrderPrice 更新订单价格
// 对应 Java: TradeOrderUpdateServiceImpl#updateOrderPrice
func (s *TradeOrderUpdateService) UpdateOrderPrice(ctx context.Context, reqVO *trade2.TradeOrderUpdatePriceReq) error {
	s.logger.Info("开始处理订单价格调整",
		zap.Int64("orderId", reqVO.ID),
		zap.Int("adjustPrice", reqVO.AdjustPrice),
	)

	var oldPayPrice int
	var newPayPrice int

	err := s.q.Transaction(func(tx *query.Query) error {
		// 1. 校验交易订单
		order, err := tx.TradeOrder.WithContext(ctx).
			Where(tx.TradeOrder.ID.Eq(reqVO.ID)).
			First()
		if err != nil {
			return ErrOrderNotExists()
		}

		// 2. 校验订单未支付
		if order.PayStatus {
			return ErrOrderAlreadyPaid()
		}

		// 3. 校验是否已经调价过（对应 Java: if (order.getAdjustPrice() > 0)）
		if order.AdjustPrice > 0 {
			return NewTradeErrorWithMsg(ErrorCodeOrderUpdateError, "订单已经调价过，不能重复调价")
		}

		// 4. 计算新的支付金额，支付价格不能 <= 0
		newPayPrice = order.PayPrice + reqVO.AdjustPrice
		if newPayPrice <= 0 {
			return NewTradeErrorWithMsg(ErrorCodeOrderUpdateError, "调价后金额必须大于 0")
		}

		oldPayPrice = order.PayPrice

		// 5. 更新订单主表
		_, err = tx.TradeOrder.WithContext(ctx).
			Where(tx.TradeOrder.ID.Eq(reqVO.ID)).
			Updates(map[string]interface{}{
				"adjust_price": order.AdjustPrice + reqVO.AdjustPrice,
				"pay_price":    newPayPrice,
			})
		if err != nil {
			return err
		}

		// 6. 更新订单项表 - 需要做 adjustPrice 的分摊
		// 对应 Java: TradePriceCalculatorHelper.dividePrice2
		orderItems, err := tx.TradeOrderItem.WithContext(ctx).
			Where(tx.TradeOrderItem.OrderID.Eq(reqVO.ID)).
			Find()
		if err != nil {
			return err
		}

		if len(orderItems) > 0 {
			// 计算价格分摊
			dividePrices := s.dividePriceToOrderItems(orderItems, reqVO.AdjustPrice)

			// 批量更新订单项
			for i, item := range orderItems {
				_, err = tx.TradeOrderItem.WithContext(ctx).
					Where(tx.TradeOrderItem.ID.Eq(item.ID)).
					Updates(map[string]interface{}{
						"adjust_price": item.AdjustPrice + dividePrices[i],
						"pay_price":    item.PayPrice + dividePrices[i],
					})
				if err != nil {
					return err
				}
			}
		}

		// 7. 更新支付订单价格
		// 对应 Java: payOrderApi.updatePayOrderPrice(order.getPayOrderId(), newPayPrice)
		if order.PayOrderID != nil && *order.PayOrderID > 0 && s.paySvc != nil {
			// 调用支付服务更新支付订单价格
			err = s.paySvc.UpdatePayOrderPrice(ctx, *order.PayOrderID, newPayPrice)
			if err != nil {
				s.logger.Error("更新支付订单价格失败", zap.Error(err))
				return err
			}
			s.logger.Info("支付订单价格已更新",
				zap.Int64("payOrderId", *order.PayOrderID),
				zap.Int("newPayPrice", newPayPrice),
			)
		}

		return nil
	})

	if err != nil {
		s.logger.Error("订单价格调整失败",
			zap.Error(err),
			zap.Int64("orderId", reqVO.ID),
			zap.Int("adjustPrice", reqVO.AdjustPrice),
		)
		return err
	}

	// 8. 记录订单日志
	// 对应 Java: TradeOrderLogUtils.setOrderInfo
	logContent := fmt.Sprintf("订单调价：原价 %d 分，调价 %d 分，新价 %d 分",
		oldPayPrice, reqVO.AdjustPrice, newPayPrice)
	if err := s.createOrderLog(ctx, reqVO.ID, 10, logContent); err != nil {
		s.logger.Error("创建订单日志失败", zap.Error(err))
		// 日志创建失败不影响主流程
	}

	s.logger.Info("订单价格调整成功",
		zap.Int64("orderId", reqVO.ID),
		zap.Int("adjustPrice", reqVO.AdjustPrice),
		zap.Int("oldPayPrice", oldPayPrice),
		zap.Int("newPayPrice", newPayPrice),
	)
	return nil
}

// dividePriceToOrderItems 将价格分摊到订单项
// 对应 Java: TradePriceCalculatorHelper.dividePrice2
func (s *TradeOrderUpdateService) dividePriceToOrderItems(items []*tradeModel.TradeOrderItem, price int) []int {
	if len(items) == 0 {
		return []int{}
	}

	// 计算订单项总支付金额
	total := 0
	for _, item := range items {
		total += item.PayPrice
	}

	if total == 0 {
		// 如果总金额为0，平均分摊
		avgPrice := price / len(items)
		remainPrice := price - avgPrice*len(items)
		prices := make([]int, len(items))
		for i := range prices {
			prices[i] = avgPrice
			if i == len(items)-1 {
				prices[i] += remainPrice
			}
		}
		return prices
	}

	// 按比例分摊，最后一个用反减法避免精度问题
	prices := make([]int, len(items))
	remainPrice := price
	for i, item := range items {
		if i < len(items)-1 {
			// 按比例计算
			// partPrice := int(float64(price) * (float64(item.PayPrice) / float64(total)))
			// 改用整数运算避免精度问题
			partPrice := (price * item.PayPrice) / total
			prices[i] = partPrice
			remainPrice -= partPrice
		} else {
			// 最后一个用反减
			prices[i] = remainPrice
		}
	}

	return prices
}

// createOrderLog 创建订单日志（修改为支持订单ID）
func (s *TradeOrderUpdateService) createOrderLog(ctx context.Context, orderId int64, operateType int, content string) error {
	// 查询订单获取用户ID
	order, err := s.q.TradeOrder.WithContext(ctx).
		Where(s.q.TradeOrder.ID.Eq(orderId)).
		First()
	if err != nil {
		return err
	}

	log := &tradeModel.TradeOrderLog{
		OrderID:     orderId,
		UserID:      order.UserID,
		UserType:    2, // 2-管理员（调价是管理员操作）
		OperateType: operateType,
		Content:     content,
	}

	return s.q.TradeOrderLog.WithContext(ctx).Create(log)
}

// UpdateOrderAddress 更新订单收货地址
// 对应 Java: TradeOrderUpdateServiceImpl#updateOrderAddress
func (s *TradeOrderUpdateService) UpdateOrderAddress(ctx context.Context, reqVO *trade2.TradeOrderUpdateAddressReq) error {
	s.logger.Info("开始更新订单收货地址",
		zap.Int64("orderId", reqVO.ID),
	)

	err := s.q.Transaction(func(tx *query.Query) error {
		// 1. 校验交易订单
		order, err := tx.TradeOrder.WithContext(ctx).
			Where(tx.TradeOrder.ID.Eq(reqVO.ID)).
			First()
		if err != nil {
			return ErrOrderNotExists()
		}

		// 2. 只有待发货状态，才可以修改订单收货地址
		// 对应 Java: if (!TradeOrderStatusEnum.isUndelivered(order.getStatus()))
		if order.Status != consts.TradeOrderStatusUndelivered {
			return NewTradeErrorWithMsg(ErrorCodeOrderUpdateError, "只有待发货状态才能修改收货地址")
		}

		// 3. 更新收货地址
		_, err = tx.TradeOrder.WithContext(ctx).
			Where(tx.TradeOrder.ID.Eq(reqVO.ID)).
			Updates(map[string]interface{}{
				"receiver_name":           reqVO.ReceiverName,
				"receiver_mobile":         reqVO.ReceiverMobile,
				"receiver_area_id":        reqVO.ReceiverAreaID,
				"receiver_detail_address": reqVO.ReceiverDetailAddress,
			})
		return err
	})

	if err != nil {
		s.logger.Error("更新订单收货地址失败",
			zap.Error(err),
			zap.Int64("orderId", reqVO.ID),
		)
		return err
	}

	// 4. 记录订单日志
	// 对应 Java: @TradeOrderLog 注解
	logContent := fmt.Sprintf("修改收货地址：%s %s %s",
		reqVO.ReceiverName, reqVO.ReceiverMobile, reqVO.ReceiverDetailAddress)
	if err := s.createOrderLog(ctx, reqVO.ID, 11, logContent); err != nil {
		s.logger.Error("创建订单日志失败", zap.Error(err))
		// 日志创建失败不影响主流程
	}

	s.logger.Info("更新订单收货地址成功",
		zap.Int64("orderId", reqVO.ID),
	)
	return nil
}

// PickUpOrderByAdmin 管理员核销订单
func (s *TradeOrderUpdateService) PickUpOrderByAdmin(ctx context.Context, adminId int64, orderId int64) error {
	req := &OrderHandleRequest{
		Operation: "pickup",
		AdminID:   adminId,
		OrderID:   orderId,
	}

	_, err := s.manager.HandleOrder(ctx, req)
	return err
}

// PickUpOrderByVerifyCode 通过核销码核销订单
func (s *TradeOrderUpdateService) PickUpOrderByVerifyCode(ctx context.Context, adminId int64, verifyCode string) error {
	// 先根据核销码查找订单
	order, err := s.q.TradeOrder.WithContext(ctx).
		Where(s.q.TradeOrder.PickUpVerifyCode.Eq(verifyCode)).
		First()
	if err != nil {
		return err
	}

	// 调用管理员核销订单
	return s.PickUpOrderByAdmin(ctx, adminId, order.ID)
}

// GetByPickUpVerifyCode 根据核销码获取订单
func (s *TradeOrderUpdateService) GetByPickUpVerifyCode(ctx context.Context, verifyCode string) (*tradeModel.TradeOrder, error) {
	return s.q.TradeOrder.WithContext(ctx).
		Where(s.q.TradeOrder.PickUpVerifyCode.Eq(verifyCode)).
		First()
}

// SettlementOrder 获得订单结算信息
// 对应 Java: TradeOrderUpdateServiceImpl#settlementOrder
func (s *TradeOrderUpdateService) SettlementOrder(ctx context.Context, userId int64, settlementReq *trade2.AppTradeOrderSettlementReq) (*trade2.AppTradeOrderSettlementResp, error) {
	s.logger.Info("开始计算订单结算信息",
		zap.Int64("userId", userId),
		zap.Int("itemCount", len(settlementReq.Items)),
	)

	// 1. 获得收货地址
	var address *memberModel.MemberAddress
	if settlementReq.AddressID != nil && *settlementReq.AddressID > 0 {
		addressResp, _ := s.addressSvc.GetAddress(ctx, userId, *settlementReq.AddressID)
		if addressResp != nil {
			address = &memberModel.MemberAddress{
				ID:            addressResp.ID,
				Name:          addressResp.Name,
				Mobile:        addressResp.Mobile,
				AreaID:        addressResp.AreaID,
				DetailAddress: addressResp.DetailAddress,
				DefaultStatus: model.BitBool(addressResp.DefaultStatus),
			}
		}
	}
	if address == nil {
		addressResp, _ := s.addressSvc.GetDefaultAddress(ctx, userId)
		if addressResp != nil {
			address = &memberModel.MemberAddress{
				ID:            addressResp.ID,
				Name:          addressResp.Name,
				Mobile:        addressResp.Mobile,
				AreaID:        addressResp.AreaID,
				DetailAddress: addressResp.DetailAddress,
				DefaultStatus: model.BitBool(addressResp.DefaultStatus),
			}
		}
	}
	if address != nil {
		settlementReq.AddressID = &address.ID
	}

	// 2. 计算价格
	priceResp, err := s.calculatePrice(ctx, userId, settlementReq)
	if err != nil {
		s.logger.Error("计算订单价格失败", zap.Error(err))
		return nil, err
	}

	// 3. 拼接返回
	result := s.convertToSettlementResp(priceResp, address)

	s.logger.Info("订单结算信息计算完成",
		zap.Int64("userId", userId),
		zap.Int("payPrice", result.Price.PayPrice),
	)

	return result, nil
}

// calculatePrice 计算订单价格
// 对应 Java: TradeOrderUpdateServiceImpl#calculatePrice
func (s *TradeOrderUpdateService) calculatePrice(ctx context.Context, userId int64, settlementReq *trade2.AppTradeOrderSettlementReq) (*TradePriceCalculateRespBO, error) {
	// 1. 如果来自购物车，则获得购物车的商品
	var cartIDs []int64
	for _, item := range settlementReq.Items {
		if item.CartID > 0 {
			cartIDs = append(cartIDs, item.CartID)
		}
	}

	// 购物车信息在价格计算中不需要，仅用于验证
	// 如果需要购物车信息，可以在这里获取

	// 2. 构建价格计算请求
	calculateReq := s.convertToCalculateReq(userId, settlementReq)

	// 3. 验证所有商品都是选中的
	for _, item := range calculateReq.Items {
		if !item.Selected {
			return nil, pkgErrors.NewBizError(1004003001, fmt.Sprintf("商品(%d)未设置为选中", item.SkuID))
		}
	}

	// 4. 计算价格
	return s.priceSvc.CalculateOrderPrice(ctx, calculateReq)
}

// CreateOrder 创建交易订单
// 对应 Java: TradeOrderUpdateServiceImpl#createOrder
// 重要变更：对齐 Java 行为，将支付订单创建纳入事务流程
// 如果支付订单创建失败，整个订单创建失败并回滚
func (s *TradeOrderUpdateService) CreateOrder(ctx context.Context, userId int64, userIP string, terminal int, createReq *trade2.AppTradeOrderCreateReq) (*tradeModel.TradeOrder, error) {
	s.logger.Info("开始创建订单",
		zap.Int64("userId", userId),
		zap.Int("itemCount", len(createReq.Items)),
	)

	// 1.1 价格计算
	priceResp, err := s.calculatePrice(ctx, userId, &trade2.AppTradeOrderSettlementReq{
		CouponID:      createReq.CouponID,
		PointStatus:   createReq.PointStatus,
		DeliveryType:  createReq.DeliveryType,
		AddressID:     createReq.AddressID,
		PickUpStoreID: createReq.PickUpStoreID,
		Items:         createReq.Items,
	})
	if err != nil {
		s.logger.Error("订单价格计算失败", zap.Error(err))
		return nil, err
	}

	// 1.2 构建订单
	order := s.buildTradeOrder(ctx, userId, userIP, terminal, createReq, priceResp)
	orderItems := s.buildTradeOrderItems(order, priceResp)

	// 2. 订单创建前的逻辑（调用处理器）
	// 对应 Java: tradeOrderHandlers.forEach(handler -> handler.beforeOrderCreate(order, orderItems))
	if err := s.executeBeforeOrderCreate(ctx, order, orderItems); err != nil {
		s.logger.Error("订单创建前置处理失败", zap.Error(err))
		return nil, err
	}

	var createdOrder *tradeModel.TradeOrder

	// 3. 保存订单（事务）
	// 对齐 Java：整个订单创建流程（包括支付订单）在同一事务语义下
	err = s.q.Transaction(func(tx *query.Query) error {
		// 3.1 插入订单
		if err := tx.TradeOrder.WithContext(ctx).Create(order); err != nil {
			return err
		}

		// 3.2 插入订单项
		for i := range orderItems {
			orderItems[i].OrderID = order.ID
		}
		if err := tx.TradeOrderItem.WithContext(ctx).CreateInBatches(orderItems, len(orderItems)); err != nil {
			return err
		}

		// 3.3 创建支付订单（对齐 Java: afterCreateTradeOrder 中的 createPayOrder）
		// 重要：将支付订单创建移入事务，如果失败则回滚订单
		if order.PayPrice > 0 {
			if err := s.createPayOrderInTx(ctx, tx, order, orderItems); err != nil {
				s.logger.Error("创建支付订单失败，回滚订单", zap.Error(err))
				return err
			}
		}

		createdOrder = order
		return nil
	})

	if err != nil {
		s.logger.Error("订单保存失败", zap.Error(err))
		return nil, err
	}

	// 4. 订单创建后的非关键逻辑（不影响主流程）
	s.afterCreateTradeOrderNonCritical(ctx, createdOrder, orderItems, createReq)

	s.logger.Info("订单创建成功",
		zap.Int64("userId", userId),
		zap.Int64("orderId", createdOrder.ID),
		zap.String("orderNo", createdOrder.No),
	)

	return createdOrder, nil
}

// buildTradeOrder 构建订单
// 对应 Java: TradeOrderUpdateServiceImpl#buildTradeOrder
func (s *TradeOrderUpdateService) buildTradeOrder(ctx context.Context, userId int64, userIP string, terminal int, createReq *trade2.AppTradeOrderCreateReq, priceResp *TradePriceCalculateRespBO) *tradeModel.TradeOrder {
	order := &tradeModel.TradeOrder{
		UserID:       userId,
		UserIP:       userIP,
		Terminal:     terminal,
		Type:         priceResp.Type,
		No:           s.generateOrderNo(),
		Status:       consts.TradeOrderStatusUnpaid,
		RefundStatus: consts.OrderRefundStatusNone,
		Remark:       createReq.Remark,
		PayStatus:    false,
		AdjustPrice:  0,
		RefundPrice:  0,
		DeliveryType: createReq.DeliveryType,
	}

	// 计算商品总数量
	productCount := 0
	for _, item := range priceResp.Items {
		productCount += item.Count
	}
	order.ProductCount = productCount

	// 设置价格信息
	order.TotalPrice = priceResp.Price.TotalPrice
	order.DiscountPrice = priceResp.Price.DiscountPrice
	order.DeliveryPrice = priceResp.Price.DeliveryPrice
	order.CouponPrice = priceResp.Price.CouponPrice
	order.PointPrice = priceResp.Price.PointPrice
	order.VipPrice = priceResp.Price.VipPrice
	order.PayPrice = priceResp.Price.PayPrice
	order.UsePoint = priceResp.UsePoint
	order.GivePoint = priceResp.GivePoint

	// 设置优惠券ID
	if priceResp.CouponID > 0 {
		order.CouponID = priceResp.CouponID
	}

	// 设置配送信息
	switch createReq.DeliveryType {
	case consts.DeliveryTypeExpress:
		// 快递配送
		var address *memberContract.AppAddressResp
		if createReq.AddressID != nil && *createReq.AddressID > 0 {
			address, _ = s.addressSvc.GetAddress(ctx, userId, *createReq.AddressID)
		}
		if address == nil {
			address, _ = s.addressSvc.GetDefaultAddress(ctx, userId)
		}
		if address != nil {
			order.ReceiverName = address.Name
			order.ReceiverMobile = address.Mobile
			order.ReceiverAreaID = int(address.AreaID)
			order.ReceiverDetailAddress = address.DetailAddress
		}
	case consts.DeliveryTypePickUp:
		// 到店自提
		order.ReceiverName = createReq.ReceiverName
		order.ReceiverMobile = createReq.ReceiverMobile
		if createReq.PickUpStoreID != nil {
			order.PickUpStoreID = *createReq.PickUpStoreID
		}
		order.PickUpVerifyCode = s.generatePickUpVerifyCode()
	}

	// 设置活动信息
	if createReq.SeckillActivityID != nil {
		order.SeckillActivityID = *createReq.SeckillActivityID
	}
	if createReq.CombinationActivityID != nil {
		order.CombinationActivityID = *createReq.CombinationActivityID
	}
	if createReq.CombinationHeadID != nil {
		order.CombinationHeadID = *createReq.CombinationHeadID
	}
	if createReq.BargainRecordID != nil {
		order.BargainRecordID = *createReq.BargainRecordID
	}
	if createReq.PointActivityID != nil {
		order.PointActivityID = *createReq.PointActivityID
	}

	return order
}

// buildTradeOrderItems 构建订单项
// 对应 Java: TradeOrderUpdateServiceImpl#buildTradeOrderItems
func (s *TradeOrderUpdateService) buildTradeOrderItems(order *tradeModel.TradeOrder, priceResp *TradePriceCalculateRespBO) []*tradeModel.TradeOrderItem {
	orderItems := make([]*tradeModel.TradeOrderItem, 0, len(priceResp.Items))

	for _, item := range priceResp.Items {
		orderItem := &tradeModel.TradeOrderItem{
			OrderID:       order.ID,
			UserID:        order.UserID,
			SpuID:         item.SpuID,
			SkuID:         item.SkuID,
			SpuName:       item.SpuName,
			PicURL:        item.PicURL,
			Count:         item.Count,
			Price:         item.Price,
			DiscountPrice: item.DiscountPrice,
			DeliveryPrice: item.DeliveryPrice,
			CouponPrice:   item.CouponPrice,
			PointPrice:    item.PointPrice,
			VipPrice:      item.VipPrice,
			PayPrice:      item.PayPrice,
			UsePoint:      item.UsePoint,
			GivePoint:     item.GivePoint,
		}

		orderItems = append(orderItems, orderItem)
	}

	return orderItems
}

// afterCreateTradeOrder 订单创建后的后置逻辑
// 对应 Java: TradeOrderUpdateServiceImpl#afterCreateTradeOrder
func (s *TradeOrderUpdateService) afterCreateTradeOrder(ctx context.Context, order *tradeModel.TradeOrder, orderItems []*tradeModel.TradeOrderItem, createReq *trade2.AppTradeOrderCreateReq) error {
	// 1. 执行订单创建后置处理器
	// 对应 Java: tradeOrderHandlers.forEach(handler -> handler.afterOrderCreate(order, orderItems))
	if err := s.executeAfterOrderCreate(ctx, order, orderItems); err != nil {
		s.logger.Error("订单创建后置处理器执行失败", zap.Error(err))
		// 后置处理失败不影响主流程，仅记录日志
	}

	// 2. 删除购物车商品
	var cartIDs []int64
	for _, item := range createReq.Items {
		if item.CartID > 0 {
			cartIDs = append(cartIDs, item.CartID)
		}
	}
	if len(cartIDs) > 0 {
		if err := s.cartSvc.DeleteCart(ctx, order.UserID, cartIDs); err != nil {
			s.logger.Error("删除购物车失败", zap.Error(err))
			// 不影响主流程
		}
	}

	// 3. 生成预支付订单
	if order.PayPrice > 0 {
		if err := s.createPayOrder(ctx, order, orderItems); err != nil {
			s.logger.Error("创建支付订单失败", zap.Error(err))
			return err
		}
	}

	// 4. 插入订单日志
	// 订单操作类型：1-创建订单
	if err := s.createOrderLogWithOrder(ctx, order, 1, "用户下单"); err != nil {
		s.logger.Error("创建订单日志失败", zap.Error(err))
		// 日志创建失败不影响主流程
	}

	return nil
}

// createPayOrderInTx 在事务中创建支付订单
// 对齐 Java: 将支付订单创建作为订单创建事务的一部分
// 如果失败，整个事务回滚，订单不会被创建
func (s *TradeOrderUpdateService) createPayOrderInTx(ctx context.Context, tx *query.Query, order *tradeModel.TradeOrder, orderItems []*tradeModel.TradeOrderItem) error {
	s.logger.Info("创建支付订单（事务内）",
		zap.Int64("orderId", order.ID),
		zap.Int("payPrice", order.PayPrice),
	)

	// 1. 获取交易配置（用于获取超时时间等）
	tradeConfig, err := s.configSvc.GetTradeConfig(ctx)
	if err != nil {
		s.logger.Error("获取交易配置失败", zap.Error(err))
		return err
	}

	// 2. 对齐 Java: 使用 appKey 获取支付应用
	// Java 使用 TradeOrderProperties.getPayAppKey()，默认值为 "mall"
	// 参见: TradeOrderProperties.java 第 22-30 行
	const defaultPayAppKey = "mall"
	payApp, err := s.payAppSvc.GetAppByAppKey(ctx, defaultPayAppKey)
	if err != nil {
		s.logger.Error("获取支付应用失败，请确认支付应用 AppKey 已配置",
			zap.String("appKey", defaultPayAppKey),
			zap.Error(err),
		)
		return pkgErrors.NewBizError(1006000000, "支付应用不存在，请联系管理员配置")
	}

	// 3. 构建支付订单创建请求
	// 对齐 Java: Subject 使用商品名称，而非订单号
	subject := "未知商品"
	if len(orderItems) > 0 {
		subject = orderItems[0].SpuName
	}
	// 截取长度（简单处理，Java 使用了 StrUtils.maxLength）
	if len(subject) > 32 {
		r := []rune(subject)
		if len(r) > 32 {
			subject = string(r[:32])
		}
	}

	createReq := &pay.PayOrderCreateReq{
		AppKey:          payApp.AppKey,
		MerchantOrderId: order.No,
		Subject:         subject,
		Body:            s.buildPayBody(orderItems),
		Price:           order.PayPrice,
		ExpireTime:      time.Now().Add(time.Duration(tradeConfig.PayTimeoutMinutes) * time.Minute),
		UserIP:          order.UserIP,
	}

	// 4. 调用支付服务创建支付订单
	// 注意：外部服务调用无法回滚，但如果失败，事务会回滚订单数据
	payOrderID, err := s.paySvc.CreateOrder(ctx, createReq)
	if err != nil {
		s.logger.Error("调用支付系统创建订单失败", zap.Error(err))
		return err
	}

	// 5. 在事务中更新交易订单的支付单编号
	_, err = tx.TradeOrder.WithContext(ctx).
		Where(tx.TradeOrder.ID.Eq(order.ID)).
		Update(tx.TradeOrder.PayOrderID, payOrderID)
	if err != nil {
		s.logger.Error("更新交易订单支付单编号失败", zap.Error(err))
		return err
	}

	// 6. 更新内存中的订单对象，确保返回值包含 payOrderId
	order.PayOrderID = &payOrderID
	s.logger.Info("支付订单创建成功",
		zap.Int64("orderId", order.ID),
		zap.Int64("payOrderId", payOrderID),
	)

	return nil
}

// afterCreateTradeOrderNonCritical 订单创建后的非关键后置逻辑
// 这些操作失败不会影响订单创建的成功状态
// 对应 Java: afterCreateTradeOrder 中除 createPayOrder 外的其他逻辑
func (s *TradeOrderUpdateService) afterCreateTradeOrderNonCritical(ctx context.Context, order *tradeModel.TradeOrder, orderItems []*tradeModel.TradeOrderItem, createReq *trade2.AppTradeOrderCreateReq) {
	// 1. 执行订单创建后置处理器
	// 对应 Java: tradeOrderHandlers.forEach(handler -> handler.afterOrderCreate(order, orderItems))
	if err := s.executeAfterOrderCreate(ctx, order, orderItems); err != nil {
		s.logger.Error("订单创建后置处理器执行失败", zap.Error(err))
		// 后置处理失败不影响主流程，仅记录日志
	}

	// 2. 删除购物车商品
	var cartIDs []int64
	for _, item := range createReq.Items {
		if item.CartID > 0 {
			cartIDs = append(cartIDs, item.CartID)
		}
	}
	if len(cartIDs) > 0 {
		if err := s.cartSvc.DeleteCart(ctx, order.UserID, cartIDs); err != nil {
			s.logger.Error("删除购物车失败", zap.Error(err))
			// 不影响主流程
		}
	}

	// 3. 插入订单日志
	// 订单操作类型：1-创建订单
	if err := s.createOrderLogWithOrder(ctx, order, 1, "用户下单"); err != nil {
		s.logger.Error("创建订单日志失败", zap.Error(err))
		// 日志创建失败不影响主流程
	}
}

// createPayOrder 创建支付订单
// 对应 Java: payOrderApi.createOrder
func (s *TradeOrderUpdateService) createPayOrder(ctx context.Context, order *tradeModel.TradeOrder, orderItems []*tradeModel.TradeOrderItem) error {
	s.logger.Info("创建支付订单",
		zap.Int64("orderId", order.ID),
		zap.Int("payPrice", order.PayPrice),
	)

	// 1. 获取交易配置，得到 AppID
	tradeConfig, err := s.configSvc.GetTradeConfig(ctx)
	if err != nil {
		s.logger.Error("获取交易配置失败", zap.Error(err))
		return err
	}

	// 2. 获取支付应用信息，得到 AppKey
	payApp, err := s.payAppSvc.GetApp(ctx, tradeConfig.AppID)
	if err != nil {
		s.logger.Error("获取支付应用失败",
			zap.Int64("appId", tradeConfig.AppID),
			zap.Error(err),
		)
		return err
	}

	// 3. 构建支付订单创建请求
	createReq := &pay.PayOrderCreateReq{
		AppKey:          payApp.AppKey,
		MerchantOrderId: order.No,
		Subject:         fmt.Sprintf("订单编号：%s", order.No),
		Body:            s.buildPayBody(orderItems),
		Price:           order.PayPrice,
		ExpireTime:      time.Now().Add(time.Duration(tradeConfig.PayTimeoutMinutes) * time.Minute),
		UserIP:          order.UserIP,
	}

	// 4. 调用支付服务创建支付订单
	payOrderID, err := s.paySvc.CreateOrder(ctx, createReq)
	if err != nil {
		s.logger.Error("调用支付系统创建订单失败", zap.Error(err))
		return err
	}

	// 5. 更新交易订单的支付单编号
	_, err = s.q.TradeOrder.WithContext(ctx).
		Where(s.q.TradeOrder.ID.Eq(order.ID)).
		Update(s.q.TradeOrder.PayOrderID, payOrderID)
	if err != nil {
		s.logger.Error("更新交易订单支付单编号失败", zap.Error(err))
		return err
	}

	order.PayOrderID = &payOrderID
	s.logger.Info("支付订单创建成功",
		zap.Int64("orderId", order.ID),
		zap.Int64("payOrderId", payOrderID),
	)

	return nil
}

// buildPayBody 构建支付订单 Body
func (s *TradeOrderUpdateService) buildPayBody(items []*tradeModel.TradeOrderItem) string {
	if len(items) == 0 {
		return ""
	}
	if len(items) == 1 {
		return items[0].SpuName
	}
	return fmt.Sprintf("%s 等 %d 件商品", items[0].SpuName, len(items))
}

// UpdateOrderPaid 更新交易订单已支付
// 对应 Java: TradeOrderUpdateServiceImpl#updateOrderPaid
func (s *TradeOrderUpdateService) UpdateOrderPaid(ctx context.Context, orderId int64, payOrderId int64) error {
	s.logger.Info("开始更新订单为已支付",
		zap.Int64("orderId", orderId),
		zap.Int64("payOrderId", payOrderId),
	)

	req := &OrderHandleRequest{
		Operation:  "pay",
		OrderID:    orderId,
		PayOrderID: payOrderId,
	}

	_, err := s.manager.HandleOrder(ctx, req)
	if err != nil {
		s.logger.Error("更新订单为已支付失败",
			zap.Error(err),
			zap.Int64("orderId", orderId),
			zap.Int64("payOrderId", payOrderId),
		)
		return err
	}

	s.logger.Info("订单已支付更新成功",
		zap.Int64("orderId", orderId),
		zap.Int64("payOrderId", payOrderId),
	)

	// 后置流程：执行支付后置钩子
	order, _ := s.q.TradeOrder.WithContext(ctx).Where(s.q.TradeOrder.ID.Eq(orderId)).First()
	if order != nil {
		_ = s.executeAfterPayOrder(ctx, order)
	}

	return nil
}

// generateOrderNo 生成订单编号
func (s *TradeOrderUpdateService) generateOrderNo() string {
	// 格式: 时间戳 + 随机数
	// 例如: 20231225143012345678
	timestamp := time.Now().Format("20060102150405")
	random := rand.Intn(1000000)
	return fmt.Sprintf("%s%06d", timestamp, random)
}

// generatePickUpVerifyCode 生成自提核销码
func (s *TradeOrderUpdateService) generatePickUpVerifyCode() string {
	// 生成8位随机数字
	return fmt.Sprintf("%08d", rand.Intn(100000000))
}

// convertToCalculateReq 转换为价格计算请求
func (s *TradeOrderUpdateService) convertToCalculateReq(userId int64, settlementReq *trade2.AppTradeOrderSettlementReq) *TradePriceCalculateReqBO {
	pointStatus := false
	if settlementReq.PointStatus != nil {
		pointStatus = *settlementReq.PointStatus
	}

	// 辅助函数：安全解引用 int64 指针
	getValue := func(ptr *int64) int64 {
		if ptr == nil {
			return 0
		}
		return *ptr
	}

	calculateReq := &TradePriceCalculateReqBO{
		UserID:                userId,
		CouponID:              settlementReq.CouponID,
		PointStatus:           pointStatus,
		DeliveryType:          settlementReq.DeliveryType,
		AddressID:             settlementReq.AddressID,
		PickUpStoreID:         settlementReq.PickUpStoreID,
		SeckillActivityId:     getValue(settlementReq.SeckillActivityID),
		CombinationActivityId: getValue(settlementReq.CombinationActivityID),
		CombinationHeadId:     getValue(settlementReq.CombinationHeadID),
		BargainRecordId:       getValue(settlementReq.BargainRecordID),
		PointActivityId:       getValue(settlementReq.PointActivityID),
		Items:                 make([]TradePriceCalculateItemBO, 0),
	}

	// 转换商品项
	for _, item := range settlementReq.Items {
		calculateReq.Items = append(calculateReq.Items, TradePriceCalculateItemBO{
			SkuID:    item.SkuID,
			Count:    item.Count,
			CartID:   item.CartID,
			Selected: true, // 默认选中
		})
	}

	return calculateReq
}

// convertToSettlementResp 转换为结算响应
func (s *TradeOrderUpdateService) convertToSettlementResp(priceResp *TradePriceCalculateRespBO, address *memberModel.MemberAddress) *trade2.AppTradeOrderSettlementResp {
	result := &trade2.AppTradeOrderSettlementResp{
		Type: priceResp.Type,
		Price: trade2.AppTradeOrderSettlementPrice{
			TotalPrice:    priceResp.Price.TotalPrice,
			DiscountPrice: priceResp.Price.DiscountPrice,
			DeliveryPrice: priceResp.Price.DeliveryPrice,
			CouponPrice:   priceResp.Price.CouponPrice,
			PointPrice:    priceResp.Price.PointPrice,
			VipPrice:      priceResp.Price.VipPrice,
			PayPrice:      priceResp.Price.PayPrice,
		},
		Items:      make([]trade2.AppTradeOrderSettlementItemResp, 0),
		Coupons:    make([]trade2.AppTradeOrderSettlementCoupon, 0),
		Promotions: make([]trade2.AppTradeOrderSettlementPromotion, 0),
		UsePoint:   priceResp.UsePoint,
		TotalPoint: priceResp.TotalPoint,
	}

	// 转换商品项
	for _, item := range priceResp.Items {
		settlementItem := trade2.AppTradeOrderSettlementItemResp{
			CategoryID: item.CategoryID,
			SpuID:      item.SpuID,
			SpuName:    item.SpuName,
			SkuID:      item.SkuID,
			PicURL:     item.PicURL,
			Price:      item.Price,
			Count:      item.Count,
			Properties: item.Properties, // 填充 SKU 属性
		}
		// 处理 CartID：如果为 0 则设为 nil（Java 返回 null）
		if item.CartID > 0 {
			cID := item.CartID
			settlementItem.CartID = &cID
		}
		result.Items = append(result.Items, settlementItem)
	}

	// 转换优惠券列表
	for _, coupon := range priceResp.Coupons {
		var mismatchReason *string
		if !coupon.Match && coupon.MismatchReason != nil {
			mismatchReason = coupon.MismatchReason
		}
		result.Coupons = append(result.Coupons, trade2.AppTradeOrderSettlementCoupon{
			ID:                 coupon.ID,
			Name:               coupon.Name,
			UsePrice:           coupon.UsePrice,
			ValidStartTime:     coupon.ValidStartTime,
			ValidEndTime:       coupon.ValidEndTime,
			DiscountType:       coupon.DiscountType,
			DiscountPercent:    coupon.DiscountPercent,
			DiscountPrice:      coupon.DiscountPrice,
			DiscountLimitPrice: coupon.DiscountLimitPrice,
			Match:              coupon.Match,
			MismatchReason:     mismatchReason,
		})
	}

	// 转换促销活动列表（对齐 Java: TradePriceCalculateRespBO.Promotion）
	for _, promotion := range priceResp.Promotions {
		promotionItem := trade2.AppTradeOrderSettlementPromotion{
			ID:            promotion.ID,
			Name:          promotion.Name,
			Type:          promotion.Type,
			TotalPrice:    promotion.TotalPrice,
			DiscountPrice: promotion.DiscountPrice,
			Match:         promotion.Match,
			Description:   promotion.Description,
			Items:         make([]trade2.AppTradeOrderSettlementPromotionItem, 0, len(promotion.Items)),
		}
		// 转换促销活动商品项
		for _, item := range promotion.Items {
			promotionItem.Items = append(promotionItem.Items, trade2.AppTradeOrderSettlementPromotionItem{
				SkuID:         item.SkuID,
				TotalPrice:    item.TotalPrice,
				DiscountPrice: item.DiscountPrice,
				PayPrice:      item.PayPrice,
			})
		}
		result.Promotions = append(result.Promotions, promotionItem)
	}

	// 设置地址信息（对齐 Java: 使用 AreaUtils.format() 填充 areaName）
	if address != nil {
		result.Address = &trade2.AppTradeOrderSettlementAddress{
			ID:            address.ID,
			Name:          address.Name,
			Mobile:        address.Mobile,
			AreaID:        int32(address.AreaID),
			AreaName:      area.Format(int(address.AreaID)), // 填充地区格式化名称
			DetailAddress: address.DetailAddress,
			DefaultStatus: bool(address.DefaultStatus),
		}
	}

	return result
}

// executeBeforeOrderCreate 执行订单创建前置处理
func (s *TradeOrderUpdateService) executeBeforeOrderCreate(ctx context.Context, order *tradeModel.TradeOrder, orderItems []*tradeModel.TradeOrderItem) error {
	s.logger.Info("执行订单创建前置处理",
		zap.Int64("orderId", order.ID),
		zap.String("orderNo", order.No),
	)

	req := &OrderHandleRequest{
		Operation:  "create",
		UserID:     order.UserID,
		OrderID:    order.ID,
		OrderItems: orderItems,
	}

	for _, handler := range s.manager.GetFactory().GetHandlers() {
		if err := handler.BeforeOrderCreate(ctx, req); err != nil {
			s.logger.Error("处理器 BeforeOrderCreate 执行失败",
				zap.String("handler", handler.GetHandlerType()),
				zap.Error(err),
			)
			return err
		}
	}

	return nil
}

// executeAfterOrderCreate 执行订单创建后置处理
func (s *TradeOrderUpdateService) executeAfterOrderCreate(ctx context.Context, order *tradeModel.TradeOrder, orderItems []*tradeModel.TradeOrderItem) error {
	s.logger.Info("执行订单创建后置处理",
		zap.Int64("orderId", order.ID),
		zap.String("orderNo", order.No),
		zap.Int("itemCount", len(orderItems)),
	)

	req := &OrderHandleRequest{
		Operation:  "create",
		UserID:     order.UserID,
		OrderID:    order.ID,
		OrderItems: orderItems,
	}
	resp := &OrderHandleResponse{
		Order:   order,
		Success: true,
	}

	for _, handler := range s.manager.GetFactory().GetHandlers() {
		if err := handler.AfterOrderCreate(ctx, req, resp); err != nil {
			s.logger.Error("处理器 AfterOrderCreate 执行失败",
				zap.String("handler", handler.GetHandlerType()),
				zap.Error(err),
			)
		}
	}

	return nil
}

// executeAfterCancelOrder 执行订单取消后置处理
func (s *TradeOrderUpdateService) executeAfterCancelOrder(ctx context.Context, order *tradeModel.TradeOrder, orderItems []*tradeModel.TradeOrderItem) error {
	s.logger.Info("执行订单取消后置处理",
		zap.Int64("orderId", order.ID),
		zap.String("orderNo", order.No),
		zap.Int("itemCount", len(orderItems)),
	)

	req := &OrderHandleRequest{
		Operation:  "cancel",
		UserID:     order.UserID,
		OrderID:    order.ID,
		OrderItems: orderItems,
	}
	resp := &OrderHandleResponse{
		Order:   order,
		Success: true,
	}

	for _, handler := range s.manager.GetFactory().GetHandlers() {
		if err := handler.AfterCancelOrder(ctx, req, resp); err != nil {
			s.logger.Error("处理器 AfterCancelOrder 执行失败",
				zap.String("handler", handler.GetHandlerType()),
				zap.Error(err),
			)
		}
	}

	return nil
}

// executeAfterPayOrder 执行订单支付后置处理
func (s *TradeOrderUpdateService) executeAfterPayOrder(ctx context.Context, order *tradeModel.TradeOrder) error {
	s.logger.Info("执行订单支付后置处理",
		zap.Int64("orderId", order.ID),
		zap.String("orderNo", order.No),
	)

	req := &OrderHandleRequest{
		Operation: "pay",
		UserID:    order.UserID,
		OrderID:   order.ID,
	}
	resp := &OrderHandleResponse{
		Order:   order,
		Success: true,
	}

	for _, handler := range s.manager.GetFactory().GetHandlers() {
		if err := handler.AfterPayOrder(ctx, req, resp); err != nil {
			s.logger.Error("处理器 AfterPayOrder 执行失败",
				zap.String("handler", handler.GetHandlerType()),
				zap.Error(err),
			)
		}
	}

	return nil
}

// executeAfterReceiveOrder 执行订单收货后置处理
func (s *TradeOrderUpdateService) executeAfterReceiveOrder(ctx context.Context, order *tradeModel.TradeOrder) error {
	s.logger.Info("执行订单收货后置处理",
		zap.Int64("orderId", order.ID),
		zap.String("orderNo", order.No),
	)

	req := &OrderHandleRequest{
		Operation: "receive",
		UserID:    order.UserID,
		OrderID:   order.ID,
	}
	resp := &OrderHandleResponse{
		Order:   order,
		Success: true,
	}

	for _, handler := range s.manager.GetFactory().GetHandlers() {
		if err := handler.AfterReceiveOrder(ctx, req, resp); err != nil {
			s.logger.Error("处理器 AfterReceiveOrder 执行失败",
				zap.String("handler", handler.GetHandlerType()),
				zap.Error(err),
			)
		}
	}

	return nil
}

// executeAfterDeliveryOrder 执行订单发货后置处理
func (s *TradeOrderUpdateService) executeAfterDeliveryOrder(ctx context.Context, order *tradeModel.TradeOrder) error {
	s.logger.Info("执行订单发货后置处理",
		zap.Int64("orderId", order.ID),
		zap.String("orderNo", order.No),
	)

	req := &OrderHandleRequest{
		Operation: "delivery",
		UserID:    order.UserID,
		OrderID:   order.ID,
	}
	resp := &OrderHandleResponse{
		Order:   order,
		Success: true,
	}

	for _, handler := range s.manager.GetFactory().GetHandlers() {
		if err := handler.AfterDeliveryOrder(ctx, req, resp); err != nil {
			s.logger.Error("处理器 AfterDeliveryOrder 执行失败",
				zap.String("handler", handler.GetHandlerType()),
				zap.Error(err),
			)
		}
	}

	return nil
}

// createOrderLogWithOrder 创建订单日志（使用订单对象）
func (s *TradeOrderUpdateService) createOrderLogWithOrder(ctx context.Context, order *tradeModel.TradeOrder, operateType int, content string) error {
	log := &tradeModel.TradeOrderLog{
		OrderID:     order.ID,
		UserID:      order.UserID,
		UserType:    consts.UserTypeMember, // 1-会员 2-管理员
		OperateType: operateType,
		Content:     content,
	}

	return s.q.TradeOrderLog.WithContext(ctx).Create(log)
}

// parseTimeString 解析时间字符串
func (s *TradeOrderUpdateService) parseTimeString(timeStr string) *time.Time {
	if timeStr == "" {
		return nil
	}

	// 尝试多种时间格式
	formats := []string{
		"2006-01-02 15:04:05",
		time.RFC3339,
		"2006-01-02T15:04:05Z",
		"2006-01-02",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, timeStr); err == nil {
			return &t
		}
	}

	s.logger.Warn("时间字符串解析失败",
		zap.String("timeStr", timeStr),
	)
	return nil
}

// CancelOrder 取消订单
// 对应 Java: TradeOrderUpdateServiceImpl#cancelOrderByMember
func (s *TradeOrderUpdateService) CancelOrder(ctx context.Context, userId int64, orderId int64) error {
	s.logger.Info("开始取消订单",
		zap.Int64("userId", userId),
		zap.Int64("orderId", orderId),
	)

	// 1. 查询订单
	order, err := s.q.TradeOrder.WithContext(ctx).
		Where(s.q.TradeOrder.ID.Eq(orderId), s.q.TradeOrder.UserID.Eq(userId)).
		First()
	if err != nil {
		return pkgErrors.NewBizError(1004001001, "订单不存在")
	}

	// 2. 校验订单状态必须是待支付
	if order.Status != consts.TradeOrderStatusUnpaid {
		return pkgErrors.NewBizError(1004001002, "订单状态不是待支付，不允许取消")
	}

	// 3. 校验支付延迟（防止支付回调延迟导致已支付订单被取消）
	// 对应 Java: if (TradeOrderStatusEnum.isUnpaid(order.getStatus()))
	if order.PayOrderID != nil && *order.PayOrderID > 0 {
		payOrder, err := s.paySvc.GetOrder(ctx, *order.PayOrderID)
		if err == nil && payOrder != nil && payOrder.Status == consts.PayOrderStatusSuccess {
			s.logger.Warn("订单支付单已支付（支付回调延迟），不支持取消",
				zap.Int64("orderId", orderId),
				zap.Int64("payOrderId", *order.PayOrderID),
			)
			return pkgErrors.NewBizError(1004001002, "订单已支付，不允许取消")
		}
	}

	// 4. 取消订单（使用事务）
	err = s.q.Transaction(func(tx *query.Query) error {
		// 4.1 更新订单状态为已取消
		now := time.Now()
		_, err := tx.TradeOrder.WithContext(ctx).
			Where(tx.TradeOrder.ID.Eq(orderId), tx.TradeOrder.Status.Eq(order.Status)).
			Updates(map[string]interface{}{
				"status":      consts.TradeOrderStatusCanceled,
				"cancel_time": now,
				"cancel_type": consts.OrderCancelTypeMember,
			})
		if err != nil {
			return err
		}

		// 4.2 执行后置处理器（对应 Java: tradeOrderHandlers.forEach(handler -> handler.afterCancelOrder)）
		orderItems, err := tx.TradeOrderItem.WithContext(ctx).
			Where(tx.TradeOrderItem.OrderID.Eq(orderId)).
			Find()
		if err != nil {
			return err
		}

		// 更新订单对象状态用于后置处理
		order.Status = consts.TradeOrderStatusCanceled
		order.CancelTime = &now
		order.CancelType = consts.OrderCancelTypeMember

		// 执行后置处理（库存回滚、优惠券回滚等）
		if err := s.executeAfterCancelOrder(ctx, order, orderItems); err != nil {
			s.logger.Error("执行取消订单后置处理失败", zap.Error(err))
			return err
		}

		return nil
	})

	if err != nil {
		s.logger.Error("取消订单失败", zap.Error(err))
		return err
	}

	// 5. 记录订单日志
	if err := s.createOrderLogWithOrder(ctx, order, 5, "用户取消订单"); err != nil {
		s.logger.Error("创建订单日志失败", zap.Error(err))
	}

	s.logger.Info("取消订单成功",
		zap.Int64("userId", userId),
		zap.Int64("orderId", orderId),
	)
	return nil
}

// CancelPaidOrder 取消已支付订单（拼团失败等场景）
// 对应 Java: TradeOrderUpdateServiceImpl#cancelPaidOrder
func (s *TradeOrderUpdateService) CancelPaidOrder(ctx context.Context, userID int64, orderID int64, cancelType int) error {
	s.logger.Info("开始取消已支付订单",
		zap.Int64("userId", userID),
		zap.Int64("orderId", orderID),
		zap.Int("cancelType", cancelType),
	)

	// 1. 查询订单
	order, err := s.q.TradeOrder.WithContext(ctx).
		Where(s.q.TradeOrder.ID.Eq(orderID), s.q.TradeOrder.UserID.Eq(userID)).
		First()
	if err != nil {
		return pkgErrors.NewBizError(1004001001, "订单不存在")
	}

	// 2. 校验状态：已支付但未收货（或根据业务需要更细粒度校验）
	if !order.PayStatus {
		return pkgErrors.NewBizError(1004001002, "订单未支付，不能使用此接口取消")
	}
	if order.Status == consts.TradeOrderStatusCanceled {
		return nil // 已经取消了
	}
	if order.Status == consts.TradeOrderStatusCompleted {
		return pkgErrors.NewBizError(1004001002, "订单已完成，不允许取消")
	}

	// 3. 执行取消并退款
	err = s.q.Transaction(func(tx *query.Query) error {
		now := time.Now()
		// 3.1 更新订单状态
		_, err := tx.TradeOrder.WithContext(ctx).
			Where(tx.TradeOrder.ID.Eq(orderID), tx.TradeOrder.Status.Eq(order.Status)).
			Updates(map[string]interface{}{
				"status":      consts.TradeOrderStatusCanceled,
				"cancel_time": now,
				"cancel_type": cancelType,
			})
		if err != nil {
			return err
		}

		// 3.2 准备后置处理
		orderItems, err := tx.TradeOrderItem.WithContext(ctx).
			Where(tx.TradeOrderItem.OrderID.Eq(orderID)).
			Find()
		if err != nil {
			return err
		}

		// 更新内存对象状态
		order.Status = consts.TradeOrderStatusCanceled
		order.CancelTime = &now
		order.CancelType = cancelType

		// 3.3 执行后置动作（回滚库存、回滚优惠券、回滚积分等）
		if err := s.executeAfterCancelOrder(ctx, order, orderItems); err != nil {
			return err
		}

		// 3.4 发起支付退款
		if order.PayOrderID != nil && *order.PayOrderID > 0 && s.payRefundSvc != nil {
			// 1. 获取支付单信息
			payOrder, err := s.paySvc.GetOrder(ctx, *order.PayOrderID)
			if err != nil {
				s.logger.Error("取消订单失败：获取支付单失败", zap.Error(err), zap.Int64("orderId", orderID), zap.Int64("payOrderId", *order.PayOrderID))
				return err
			}
			if payOrder == nil {
				s.logger.Error("取消订单失败：支付单不存在", zap.Int64("orderId", orderID), zap.Int64("payOrderId", *order.PayOrderID))
				return fmt.Errorf("支付单不存在")
			}

			// 2. 获取支付应用信息 (为了拿到 AppKey)
			payApp, err := s.payAppSvc.GetApp(ctx, payOrder.AppID)
			if err != nil {
				s.logger.Error("取消订单失败：获取支付应用失败", zap.Error(err), zap.Int64("appId", payOrder.AppID))
				return err
			}
			if payApp == nil {
				s.logger.Error("取消订单失败：支付应用不存在", zap.Int64("appId", payOrder.AppID))
				return fmt.Errorf("支付应用不存在")
			}

			// 3. 发起退款
			// 生成唯一退款单号
			refundNo, err := s.noDAO.Generate(ctx, "R")
			if err != nil {
				s.logger.Error("取消订单失败：生成退款单号失败", zap.Error(err), zap.Int64("orderId", orderID))
				return err
			}

			_, err = s.payRefundSvc.CreateRefund(ctx, &pay.PayRefundCreateReq{
				AppKey:           payApp.AppKey,
				MerchantOrderId:  order.No,
				MerchantRefundId: refundNo,
				Price:            order.PayPrice,
				Reason:           "订单取消退款",
				UserIP:           order.UserIP,
			})
			if err != nil {
				s.logger.Error("取消订单发起退款失败", zap.Error(err), zap.Int64("orderId", orderID))
				return err
			}
		}

		return nil
	})

	if err != nil {
		s.logger.Error("取消已支付订单失败", zap.Error(err))
		return err
	}

	// 4. 记录日志
	logContent := fmt.Sprintf("取消已支付订单，类型：%d", cancelType)
	_ = s.createOrderLogWithOrder(ctx, order, 5, logContent)

	return nil
}

// DeleteOrder 删除订单
// 对应 Java: TradeOrderUpdateServiceImpl#deleteOrder
func (s *TradeOrderUpdateService) DeleteOrder(ctx context.Context, userId int64, orderId int64) error {
	s.logger.Info("开始删除订单",
		zap.Int64("userId", userId),
		zap.Int64("orderId", orderId),
	)

	// 1. 校验订单存在
	order, err := s.q.TradeOrder.WithContext(ctx).
		Where(s.q.TradeOrder.ID.Eq(orderId), s.q.TradeOrder.UserID.Eq(userId)).
		First()
	if err != nil {
		return pkgErrors.NewBizError(1004001001, "订单不存在")
	}

	// 2. 校验订单状态必须是已取消
	// 对应 Java: if (ObjectUtil.notEqual(order.getStatus(), TradeOrderStatusEnum.CANCELED.getStatus()))
	if order.Status != consts.TradeOrderStatusCanceled {
		return pkgErrors.NewBizError(1004001003, "订单状态不是已取消，不允许删除")
	}

	// 3. 删除订单
	err = s.q.Transaction(func(tx *query.Query) error {
		_, err := tx.TradeOrder.WithContext(ctx).
			Where(tx.TradeOrder.ID.Eq(orderId)).
			Delete()
		return err
	})

	if err != nil {
		s.logger.Error("删除订单失败", zap.Error(err))
		return err
	}

	// 4. 记录订单日志
	if err := s.createOrderLogWithOrder(ctx, order, 6, "用户删除订单"); err != nil {
		s.logger.Error("创建订单日志失败", zap.Error(err))
	}

	s.logger.Info("删除订单成功",
		zap.Int64("userId", userId),
		zap.Int64("orderId", orderId),
	)
	return nil
}

// ReceiveOrder 确认收货
// 对应 Java: TradeOrderUpdateServiceImpl#receiveOrderByMember
func (s *TradeOrderUpdateService) ReceiveOrder(ctx context.Context, userId int64, orderId int64) error {
	s.logger.Info("开始确认收货",
		zap.Int64("userId", userId),
		zap.Int64("orderId", orderId),
	)

	// 1. 查询订单
	order, err := s.q.TradeOrder.WithContext(ctx).
		Where(s.q.TradeOrder.ID.Eq(orderId), s.q.TradeOrder.UserID.Eq(userId)).
		First()
	if err != nil {
		return pkgErrors.NewBizError(1004001001, "订单不存在")
	}

	// 2. 校验订单状态必须是已发货
	// 对应 Java: if (!TradeOrderStatusEnum.isDelivered(order.getStatus()))
	if order.Status != consts.TradeOrderStatusDelivered {
		return pkgErrors.NewBizError(1004001004, "订单状态不是已发货，不允许确认收货")
	}

	// 3. 更新订单状态为已完成（使用事务）
	err = s.q.Transaction(func(tx *query.Query) error {
		now := time.Now()
		// 3.1 更新订单状态为已完成（使用乐观锁）
		_, err := tx.TradeOrder.WithContext(ctx).
			Where(tx.TradeOrder.ID.Eq(orderId), tx.TradeOrder.Status.Eq(order.Status)).
			Updates(map[string]interface{}{
				"status":       consts.TradeOrderStatusCompleted,
				"receive_time": now,
				"finish_time":  now,
			})
		if err != nil {
			return err
		}

		// 3.2 更新订单对象状态用于后置处理
		order.Status = consts.TradeOrderStatusCompleted
		order.ReceiveTime = &now
		order.FinishTime = &now

		// 3.3 执行后置处理器（对应 Java: tradeOrderHandlers.forEach(handler -> handler.afterReceiveOrder(order))）
		if err := s.executeAfterReceiveOrder(ctx, order); err != nil {
			s.logger.Error("执行确认收货后置处理失败", zap.Error(err))
			return err
		}

		return nil
	})

	if err != nil {
		s.logger.Error("确认收货失败", zap.Error(err))
		return err
	}

	// 4. 记录订单日志
	if err := s.createOrderLogWithOrder(ctx, order, 4, "用户确认收货"); err != nil {
		s.logger.Error("创建订单日志失败", zap.Error(err))
	}

	s.logger.Info("确认收货成功",
		zap.Int64("userId", userId),
		zap.Int64("orderId", orderId),
	)
	return nil
}

// CreateOrderItemCommentByMember 创建订单项评价
func (s *TradeOrderUpdateService) CreateOrderItemCommentByMember(ctx context.Context, userId int64, createReq *trade2.AppTradeOrderItemCommentCreateReq) (int64, error) {
	s.logger.Info("开始创建订单项评价",
		zap.Int64("userId", userId),
		zap.Int64("orderItemId", createReq.OrderItemID),
	)

	// 1. 查询订单项
	orderItem, err := s.q.TradeOrderItem.WithContext(ctx).
		Where(s.q.TradeOrderItem.ID.Eq(createReq.OrderItemID)).
		First()
	if err != nil {
		return 0, pkgErrors.NewBizError(1004002001, "订单项不存在")
	}

	// 2. 查询订单并验证权限
	order, err := s.q.TradeOrder.WithContext(ctx).
		Where(s.q.TradeOrder.ID.Eq(orderItem.OrderID), s.q.TradeOrder.UserID.Eq(userId)).
		First()
	if err != nil {
		return 0, pkgErrors.NewBizError(1004001001, "订单不存在或无权限")
	}

	// 3. 验证订单状态（必须是已完成）
	if order.Status != consts.TradeOrderStatusCompleted {
		return 0, pkgErrors.NewBizError(1004002002, "订单未完成，不能评价")
	}

	// 4. 创建评价记录
	// 委托给 ProductCommentService 创建评价
	commentReq := &product.AppProductCommentCreateReq{
		OrderItemID:       createReq.OrderItemID,
		Anonymous:         createReq.Anonymous,
		Content:           createReq.Content,
		PicURLs:           createReq.PicUrls,
		Scores:            (createReq.DescriptionScores + createReq.BenefitScores) / 2,
		DescriptionScores: createReq.DescriptionScores,
		BenefitScores:     createReq.BenefitScores,
	}

	comment, err := s.commentSvc.CreateAppComment(ctx, userId, commentReq)
	if err != nil {
		s.logger.Error("创建评价失败", zap.Error(err))
		return 0, err
	}

	s.logger.Info("订单项评价创建成功",
		zap.Int64("orderItemId", createReq.OrderItemID),
		zap.Int64("commentId", comment.ID),
	)

	return comment.ID, nil
}

// UpdateOrderItemWhenAfterSaleCreate 更新订单项在售后创建时状态
// 对应 Java: TradeOrderUpdateServiceImpl#updateOrderItemWhenAfterSaleCreate
func (s *TradeOrderUpdateService) UpdateOrderItemWhenAfterSaleCreate(ctx context.Context, orderId int64, orderItemId int64, afterSaleId int64) error {
	s.logger.Info("开始更新订单项售后状态(创建)",
		zap.Int64("orderId", orderId),
		zap.Int64("orderItemId", orderItemId),
		zap.Int64("afterSaleId", afterSaleId),
	)

	return s.q.Transaction(func(tx *query.Query) error {
		// 1. 更新订单项售后状态
		if _, err := tx.TradeOrderItem.WithContext(ctx).
			Where(tx.TradeOrderItem.ID.Eq(orderItemId)).
			Updates(map[string]interface{}{
				"after_sale_status": tradeModel.TradeOrderItemAfterSaleStatusApply,
				"after_sale_id":     afterSaleId,
			}); err != nil {
			return err
		}

		// 2. 更新订单售后状态为【申请退款】
		if _, err := tx.TradeOrder.WithContext(ctx).
			Where(tx.TradeOrder.ID.Eq(orderId)).
			Update(tx.TradeOrder.RefundStatus, consts.OrderRefundStatusApply); err != nil {
			return err
		}

		return nil
	})
}

// UpdateOrderItemWhenAfterSaleSuccess 更新订单项在售后成功时状态
// 对应 Java: TradeOrderUpdateServiceImpl#updateOrderItemWhenAfterSaleSuccess
func (s *TradeOrderUpdateService) UpdateOrderItemWhenAfterSaleSuccess(ctx context.Context, orderId int64, orderItemId int64, refundPrice int) error {
	s.logger.Info("开始更新订单项售后状态(成功)",
		zap.Int64("orderId", orderId),
		zap.Int64("orderItemId", orderItemId),
		zap.Int("refundPrice", refundPrice),
	)

	return s.q.Transaction(func(tx *query.Query) error {
		// 1. 更新订单项售后状态
		_, err := tx.TradeOrderItem.WithContext(ctx).
			Where(tx.TradeOrderItem.ID.Eq(orderItemId)).
			Updates(map[string]interface{}{
				"after_sale_status": tradeModel.TradeOrderItemAfterSaleStatusSuccess,
			})
		if err != nil {
			return err
		}

		// 2. 更新订单的退款金额和积分
		order, err := tx.TradeOrder.WithContext(ctx).Where(tx.TradeOrder.ID.Eq(orderId)).First()
		if err != nil {
			return err
		}

		updates := map[string]interface{}{
			"refund_price":  order.RefundPrice + refundPrice,
			"refund_status": consts.OrderRefundStatusRefunded,
		}

		// 3. 检查是否所有订单项都售后成功，如果是则取消订单
		items, err := tx.TradeOrderItem.WithContext(ctx).Where(tx.TradeOrderItem.OrderID.Eq(orderId)).Find()
		if err != nil {
			return err
		}

		allSuccess := true
		for _, item := range items {
			if item.ID == orderItemId {
				// 当前项还没提交更新到DB，但在内存中应视为成功
				continue
			}
			if item.AfterSaleStatus != tradeModel.TradeOrderItemAfterSaleStatusSuccess {
				allSuccess = false
				break
			}
		}

		if allSuccess {
			updates["status"] = consts.TradeOrderStatusCanceled
			updates["cancel_type"] = consts.OrderCancelTypeAfterSaleClose
			now := time.Now()
			updates["cancel_time"] = &now
		}

		_, err = tx.TradeOrder.WithContext(ctx).Where(tx.TradeOrder.ID.Eq(orderId)).Updates(updates)
		return err
	})
}

// UpdateOrderItemWhenAfterSaleCancel 更新订单项在售后取消时状态
// 对应 Java: TradeOrderUpdateServiceImpl#updateOrderItemWhenAfterSaleCancel
func (s *TradeOrderUpdateService) UpdateOrderItemWhenAfterSaleCancel(ctx context.Context, orderId int64, orderItemId int64) error {
	s.logger.Info("开始更新订单项售后状态(取消)",
		zap.Int64("orderId", orderId),
		zap.Int64("orderItemId", orderItemId),
	)

	return s.q.Transaction(func(tx *query.Query) error {
		// 1. 更新订单项售后状态为无
		if _, err := tx.TradeOrderItem.WithContext(ctx).
			Where(tx.TradeOrderItem.ID.Eq(orderItemId)).
			Updates(map[string]interface{}{
				"after_sale_status": tradeModel.TradeOrderItemAfterSaleStatusNone,
				"after_sale_id":     0,
			}); err != nil {
			return err
		}

		// 2. 检查是否还有其他订单项处于售后中
		count, err := tx.TradeOrderItem.WithContext(ctx).
			Where(tx.TradeOrderItem.OrderID.Eq(orderId), tx.TradeOrderItem.AfterSaleStatus.Neq(tradeModel.TradeOrderItemAfterSaleStatusNone)).
			Count()
		if err != nil {
			return err
		}

		// 3. 如果没有其他售后项，更新订单售后状态为【无】
		if count == 0 {
			if _, err := tx.TradeOrder.WithContext(ctx).
				Where(tx.TradeOrder.ID.Eq(orderId)).
				Update(tx.TradeOrder.RefundStatus, consts.OrderRefundStatusNone); err != nil {
				return err
			}
		}

		return nil
	})
}

// CancelOrderBySystem 系统自动取消订单
// 对应 Java: TradeOrderUpdateServiceImpl#cancelOrderBySystem
func (s *TradeOrderUpdateService) CancelOrderBySystem(ctx context.Context) (int64, error) {
	s.logger.Info("开始执行系统自动取消订单任务")

	// 1. 获取过期未支付订单
	// 配置的超时时间，这里暂时硬编码或者后续从配置服务获取，假设为 30 分钟
	// Java版是从 TradeConfig 获取 tradeOrderExpireTime
	tradeConfig, err := s.configSvc.GetTradeConfig(ctx)
	if err != nil {
		s.logger.Error("获取交易配置失败", zap.Error(err))
		return 0, err
	}
	expireTime := time.Now().Add(-time.Duration(tradeConfig.PayTimeoutMinutes) * time.Minute)

	orders, err := s.q.TradeOrder.WithContext(ctx).
		Where(s.q.TradeOrder.Status.Eq(consts.TradeOrderStatusUnpaid), s.q.TradeOrder.CreateTime.Lt(expireTime)).
		Find()
	if err != nil {
		return 0, err
	}

	if len(orders) == 0 {
		return 0, nil
	}

	// 2. 遍历取消
	count := int64(0)
	for _, order := range orders {
		// 避免影响主流程，单个失败不中断
		if err := s.cancelOrderBySystemSingle(ctx, order); err != nil {
			s.logger.Error("系统取消订单失败", zap.Int64("orderId", order.ID), zap.Error(err))
		} else {
			count++
		}
	}

	s.logger.Info("系统自动取消订单任务完成", zap.Int64("count", count))
	return count, nil
}

// cancelOrderBySystemSingle 单个订单系统取消逻辑
func (s *TradeOrderUpdateService) cancelOrderBySystemSingle(ctx context.Context, order *tradeModel.TradeOrder) error {
	// 1. 获取订单项
	items, err := s.q.TradeOrderItem.WithContext(ctx).
		Where(s.q.TradeOrderItem.OrderID.Eq(order.ID)).
		Find()
	if err != nil {
		s.logger.Error("获取订单项失败", zap.Int64("orderId", order.ID), zap.Error(err))
		return err
	}
	// 2. 将 items 设置到 order 对象中，用于后续流程（虽然 HandleOrder 内可能不用，但 executeAfter 需要）
	// 注意：OrderHandleRequest 需要包含 CancelType

	handleReq := &OrderHandleRequest{
		Operation:    "cancel",
		OrderID:      order.ID,
		UserID:       order.UserID,
		CancelType:   consts.OrderCancelTypeTimeout,
		CancelReason: "支付超时取消",
		OrderItems:   items, // 放入 Request 以便 Processor 可能使用
	}

	_, err = s.manager.HandleOrder(ctx, handleReq)
	if err != nil {
		return err
	}

	// 3. 执行后置处理（库存回滚、优惠券回滚等）
	// 注意：订单对象的状态在 HandleOrder 中已被更新，但为了确保 executeAfterCancelOrder 获取到最新状态
	order.Status = consts.TradeOrderStatusCanceled
	now := time.Now()
	order.CancelTime = &now
	order.CancelType = consts.OrderCancelTypeTimeout

	if err := s.executeAfterCancelOrder(ctx, order, items); err != nil {
		s.logger.Error("执行取消订单后置处理失败", zap.Error(err))
		return err
	}

	return nil
}

// ReceiveOrderBySystem 系统自动确认收货
// 对应 Java: TradeOrderUpdateServiceImpl#receiveOrderBySystem
func (s *TradeOrderUpdateService) ReceiveOrderBySystem(ctx context.Context) (int64, error) {
	s.logger.Info("开始执行系统自动确认收货任务")

	// 1. 获取过期未收货订单
	tradeConfig, err := s.configSvc.GetTradeConfig(ctx)
	if err != nil {
		return 0, err
	}
	expireTime := time.Now().Add(-time.Duration(tradeConfig.AutoReceiveDays) * time.Hour * 24) // Day to Day

	orders, err := s.q.TradeOrder.WithContext(ctx).
		Where(s.q.TradeOrder.Status.Eq(consts.TradeOrderStatusDelivered), s.q.TradeOrder.DeliveryTime.Lt(expireTime)).
		Find()
	if err != nil {
		return 0, err
	}

	if len(orders) == 0 {
		return 0, nil
	}

	// 2. 遍历确认收货
	count := int64(0)
	for _, order := range orders {
		if err := s.ReceiveOrder(ctx, order.UserID, order.ID); err != nil {
			s.logger.Error("系统确认收货失败", zap.Int64("orderId", order.ID), zap.Error(err))
		} else {
			count++
		}
	}
	s.logger.Info("系统自动确认收货任务完成", zap.Int64("count", count))
	return count, nil
}

// CreateOrderItemCommentBySystem 系统自动创建评价
// 对应 Java: TradeOrderUpdateServiceImpl#createOrderItemCommentBySystem
func (s *TradeOrderUpdateService) CreateOrderItemCommentBySystem(ctx context.Context) (int64, error) {
	s.logger.Info("开始执行系统自动评价任务")

	// 1. 获取过期未评价订单
	tradeConfig, err := s.configSvc.GetTradeConfig(ctx)
	if err != nil {
		return 0, err
	}
	expireTime := time.Now().Add(-time.Duration(tradeConfig.AutoCommentDays) * time.Hour * 24)

	// 查询 Status=Completed AND CommentStatus=false AND FinishTime < expireTime, Limit ?
	orders, err := s.q.TradeOrder.WithContext(ctx).
		Where(s.q.TradeOrder.Status.Eq(consts.TradeOrderStatusCompleted),
			s.q.TradeOrder.CommentStatus.Eq(model.NewBitBool(false)),
			s.q.TradeOrder.FinishTime.Lt(expireTime)).
		Find()
	if err != nil {
		return 0, err
	}

	if len(orders) == 0 {
		return 0, nil
	}

	// 2. 遍历创建评价
	count := int64(0)
	for _, order := range orders {
		// 查询订单项
		items, _ := s.q.TradeOrderItem.WithContext(ctx).Where(s.q.TradeOrderItem.OrderID.Eq(order.ID)).Find()
		for _, item := range items {
			if item.CommentStatus {
				continue
			}
			// 创建好评
			req := &trade2.AppTradeOrderItemCommentCreateReq{
				OrderItemID:       item.ID,
				Content:           "好评！系统默认好评。",
				BenefitScores:     5,
				DescriptionScores: 5,
				Anonymous:         true,
			}
			if _, err := s.CreateOrderItemCommentByMember(ctx, order.UserID, req); err != nil {
				s.logger.Error("系统创建评价失败", zap.Int64("orderItemId", item.ID), zap.Error(err))
			} else {
				count++
			}
		}
	}
	s.logger.Info("系统自动评价任务完成", zap.Int64("count", count))
	return count, nil
}

// UpdateOrderCombinationInfo 更新订单拼团信息
// 对应 Java: TradeOrderUpdateServiceImpl#updateOrderCombinationInfo
func (s *TradeOrderUpdateService) UpdateOrderCombinationInfo(ctx context.Context, orderId int64, activityId int64, combinationRecordId int64, headId int64) error {
	_, err := s.q.TradeOrder.WithContext(ctx).
		Where(s.q.TradeOrder.ID.Eq(orderId)).
		Updates(map[string]interface{}{
			"combination_activity_id": activityId,
			"combination_record_id":   combinationRecordId,
			"combination_head_id":     headId,
		})
	return err
}

// UpdateOrderGiveCouponIds 更新订单赠送优惠券
// 对应 Java: TradeOrderUpdateServiceImpl#updateOrderGiveCouponIds
func (s *TradeOrderUpdateService) UpdateOrderGiveCouponIds(ctx context.Context, userId int64, orderId int64, couponIds []int64) error {
	// 校验订单
	_, err := s.q.TradeOrder.WithContext(ctx).
		Where(s.q.TradeOrder.ID.Eq(orderId), s.q.TradeOrder.UserID.Eq(userId)).
		First()
	if err != nil {
		return pkgErrors.NewBizError(1004001001, "订单不存在")
	}

	if len(couponIds) == 0 {
		return nil
	}

	// 转换 []int64 -> model.IntListFromCSV (assuming alias to []int)
	intCouponIds := make([]int, len(couponIds))
	for i, v := range couponIds {
		intCouponIds[i] = int(v)
	}

	_, err = s.q.TradeOrder.WithContext(ctx).
		Where(s.q.TradeOrder.ID.Eq(orderId)).
		Updates(map[string]interface{}{
			"give_coupon_ids": model.IntListFromCSV(intCouponIds),
		})
	return err
}

// SyncOrderPayStatusQuietly 静默同步订单支付状态
func (s *TradeOrderUpdateService) SyncOrderPayStatusQuietly(ctx context.Context, orderId int64) {
	// 调用 PayOrderProcessor 的逻辑或者 PayOrderService 的 sync logic
	// 这里主要是为了防止支付回调丢失，主动去查支付单状态并同步
	// 对应 Java: tradeOrderUpdateService.syncOrderPayStatus(id) inside try-catch
	// 1. 查询订单
	order, err := s.q.TradeOrder.WithContext(ctx).Where(s.q.TradeOrder.ID.Eq(orderId)).First()
	if err != nil {
		s.logger.Error("静默同步支付状态失败：订单不存在", zap.Int64("orderId", orderId), zap.Error(err))
		return
	}

	if order.PayStatus { // Already paid
		return
	}
	if order.PayOrderID == nil {
		return
	}

	// 2. 查询支付单
	// 需要 PayOrderService
	payOrder, err := s.paySvc.GetOrder(ctx, *order.PayOrderID)
	if err != nil {
		s.logger.Error("静默同步支付状态失败：支付单查询失败", zap.Int64("payOrderId", *order.PayOrderID), zap.Error(err))
		return
	}

	// 3. 如果支付成功，则更新我们系统的订单状态
	if payOrder != nil && payOrder.Status == consts.PayOrderStatusSuccess {
		// 调用 UpdateOrderPaid
		if err := s.UpdateOrderPaid(ctx, orderId, payOrder.ID); err != nil {
			s.logger.Error("静默同步支付状态失败：更新订单支付状态失败", zap.Int64("orderId", orderId), zap.Error(err))
		}
	}
}
