package pay

import (
	"backend-go/internal/api/req"
	"backend-go/internal/api/resp"
	"backend-go/internal/model/pay"
	"backend-go/internal/pkg/core"
	"backend-go/internal/repo/query"
	"backend-go/internal/service/pay/client"
	"backend-go/pkg/config"
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type PayOrderService struct {
	q          *query.Query
	appSvc     *PayAppService
	channelSvc *PayChannelService
	clientFac  *client.PayClientFactory
}

func NewPayOrderService(q *query.Query, appSvc *PayAppService, channelSvc *PayChannelService, clientFac *client.PayClientFactory) *PayOrderService {
	return &PayOrderService{
		q:          q,
		appSvc:     appSvc,
		channelSvc: channelSvc,
		clientFac:  clientFac,
	}
}

// GetOrder 获得支付订单
func (s *PayOrderService) GetOrder(ctx context.Context, id int64) (*pay.PayOrder, error) {
	return s.q.PayOrder.WithContext(ctx).Where(s.q.PayOrder.ID.Eq(id)).First()
}

// GetOrderPage 获得支付订单分页
func (s *PayOrderService) GetOrderPage(ctx context.Context, req *req.PayOrderPageReq) (*core.PageResult[*pay.PayOrder], error) {
	q := s.q.PayOrder.WithContext(ctx)
	if req.AppID > 0 {
		q = q.Where(s.q.PayOrder.AppID.Eq(req.AppID))
	}
	if req.ChannelCode != "" {
		q = q.Where(s.q.PayOrder.ChannelCode.Eq(req.ChannelCode))
	}
	if req.MerchantOrderId != "" {
		q = q.Where(s.q.PayOrder.MerchantOrderId.Eq(req.MerchantOrderId))
	}
	if req.Subject != "" {
		q = q.Where(s.q.PayOrder.Subject.Like("%" + req.Subject + "%"))
	}
	if req.No != "" {
		q = q.Where(s.q.PayOrder.No.Eq(req.No))
	}
	if req.Status != nil {
		q = q.Where(s.q.PayOrder.Status.Eq(*req.Status))
	}

	total, err := q.Count()
	if err != nil {
		return nil, err
	}
	list, err := q.Limit(req.GetLimit()).Offset(req.GetOffset()).Order(s.q.PayOrder.ID.Desc()).Find()
	if err != nil {
		return nil, err
	}
	return &core.PageResult[*pay.PayOrder]{
		List:  list,
		Total: total,
	}, nil
}

// CreateOrder 创建支付单
func (s *PayOrderService) CreateOrder(ctx context.Context, reqDTO *req.PayOrderCreateReq) (int64, error) {
	app, err := s.appSvc.GetApp(ctx, reqDTO.AppID)
	if err != nil {
		return 0, err
	}
	if app == nil || app.Status != 0 {
		return 0, errors.New("App disabled or not found")
	}

	existOrder, _ := s.q.PayOrder.WithContext(ctx).
		Where(s.q.PayOrder.AppID.Eq(app.ID), s.q.PayOrder.MerchantOrderId.Eq(reqDTO.MerchantOrderId)).
		First()
	if existOrder != nil {
		return existOrder.ID, nil
	}

	// 创建支付交易单 (对齐 Java: 使用 app.OrderNotifyURL)
	order := &pay.PayOrder{
		AppID:           app.ID,
		MerchantOrderId: reqDTO.MerchantOrderId,
		Subject:         reqDTO.Subject,
		Body:            reqDTO.Body,
		NotifyURL:       app.OrderNotifyURL, // 对齐 Java: 使用 app 的回调地址
		Price:           reqDTO.Price,
		ExpireTime:      time.Now().Add(2 * time.Hour),
		Status:          PayOrderStatusWaiting,
		RefundPrice:     0,
		UserIP:          reqDTO.UserIP,
	}

	if err := s.q.PayOrder.WithContext(ctx).Create(order); err != nil {
		return 0, err
	}
	return order.ID, nil
}

// ... GetOrderCountByAppId ...

// SubmitOrder 提交支付订单
func (s *PayOrderService) SubmitOrder(ctx context.Context, reqVO *req.PayOrderSubmitReq, userIP string) (*resp.PayOrderSubmitResp, error) {
	order, err := s.validateOrderCanSubmit(ctx, reqVO.ID)
	if err != nil {
		return nil, err
	}

	channel, err := s.validateChannelCanSubmit(ctx, order.AppID, reqVO.ChannelCode)
	if err != nil {
		return nil, err
	}

	// Generate No
	no := s.generateNo()

	// Create Extension
	ext := &pay.PayOrderExtension{
		OrderID:     order.ID,
		No:          no,
		ChannelID:   channel.ID,
		ChannelCode: channel.Code,
		UserIP:      userIP,
		Status:      PayOrderStatusWaiting,
	}
	if err := s.q.PayOrderExtension.WithContext(ctx).Create(ext); err != nil {
		return nil, err
	}

	// Get Pay Client
	payClient := s.clientFac.GetPayClient(channel.ID)
	if payClient == nil {
		// Lazy create if not exists
		var err error
		payClient, err = s.clientFac.CreateOrUpdatePayClient(channel.ID, channel.Code, channel.Config)
		if err != nil {
			return nil, err
		}
	}

	// Call UnifiedOrder (对齐 Java: 使用渠道特定的回调 URL)
	unifiedReq := &client.UnifiedOrderReq{
		UserIP:     userIP,
		OutTradeNo: no,
		Subject:    order.Subject,
		Body:       order.Body,
		NotifyURL:  s.genChannelOrderNotifyUrl(channel), // 对齐 Java: 渠道回调 URL
		// ReturnURL:   reqVO.ReturnUrl,
		Price:       order.Price,
		ExpireTime:  order.ExpireTime,
		DisplayMode: reqVO.DisplayMode,
	}
	unifiedResp, err := payClient.UnifiedOrder(ctx, unifiedReq)
	if err != nil {
		return nil, err
	}

	// Return response
	return &resp.PayOrderSubmitResp{
		Status:         unifiedResp.Status,
		DisplayMode:    unifiedResp.DisplayMode,
		DisplayContent: unifiedResp.DisplayContent,
	}, nil
}

func (s *PayOrderService) validateOrderCanSubmit(ctx context.Context, id int64) (*pay.PayOrder, error) {
	order, err := s.q.PayOrder.WithContext(ctx).Where(s.q.PayOrder.ID.Eq(id)).First()
	if err != nil {
		return nil, gorm.ErrRecordNotFound
	}
	if order.Status == PayOrderStatusSuccess {
		return nil, errors.New("Order already paid")
	}
	if order.Status != PayOrderStatusWaiting {
		return nil, errors.New("Order status not waiting")
	}
	if order.ExpireTime.Before(time.Now()) {
		return nil, errors.New("Order expired")
	}
	return order, nil
}

func (s *PayOrderService) validateChannelCanSubmit(ctx context.Context, appId int64, channelCode string) (*pay.PayChannel, error) {
	// app validation is implicit or done separately
	return s.channelSvc.GetChannelByAppIdAndCode(ctx, appId, channelCode)
}

// genChannelOrderNotifyUrl 根据支付渠道生成回调地址
// 对齐 Java: payProperties.getOrderNotifyUrl() + "/" + channel.getId()
func (s *PayOrderService) genChannelOrderNotifyUrl(channel *pay.PayChannel) string {
	return fmt.Sprintf("%s/%d", config.C.Pay.OrderNotifyURL, channel.ID)
}

func (s *PayOrderService) generateNo() string {
	// Simple timestamp + random for now.
	// Java uses Redis. We can use core.RDB.Incr if we want strictly strict.
	// For MVP: P + yyyyMMddHHmmss + 6 digit random
	return "P" + time.Now().Format("20060102150405") + core.GenerateRandomString(6) // Need helper?
	// Let's use simplified version
	return "P" + time.Now().Format("20060102150405") + "000000"
}

// GetOrderExtension 获得支付订单拓展
func (s *PayOrderService) GetOrderExtension(ctx context.Context, id int64) (*pay.PayOrderExtension, error) {
	return s.q.PayOrderExtension.WithContext(ctx).Where(s.q.PayOrderExtension.ID.Eq(id)).First()
}

// SyncOrderQuietly 同步订单的支付状态 (Quietly)
func (s *PayOrderService) SyncOrderQuietly(ctx context.Context, id int64) {
	// TODO: Implement Sync Logic (Requires Pay Client)
	// For now, just a placeholder or minimal logic if possible.
	// In Java, this calls payOrderService.syncOrder(id).
}

// GetOrderList 获得支付订单列表 (Export)
func (s *PayOrderService) GetOrderList(ctx context.Context, req *req.PayOrderExportReq) ([]*pay.PayOrder, error) {
	q := s.q.PayOrder.WithContext(ctx)
	if req.AppID > 0 {
		q = q.Where(s.q.PayOrder.AppID.Eq(req.AppID))
	}
	if req.ChannelCode != "" {
		q = q.Where(s.q.PayOrder.ChannelCode.Eq(req.ChannelCode))
	}
	if req.MerchantOrderId != "" {
		q = q.Where(s.q.PayOrder.MerchantOrderId.Eq(req.MerchantOrderId))
	}
	if req.Subject != "" {
		q = q.Where(s.q.PayOrder.Subject.Like("%" + req.Subject + "%"))
	}
	if req.No != "" {
		q = q.Where(s.q.PayOrder.No.Eq(req.No))
	}
	if req.Status != nil {
		q = q.Where(s.q.PayOrder.Status.Eq(*req.Status))
	}
	return q.Order(s.q.PayOrder.ID.Desc()).Find()
}
