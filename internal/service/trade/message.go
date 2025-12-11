package trade

import (
	"backend-go/internal/model/trade"
	"context"
)

type TradeMessageService struct {
}

func NewTradeMessageService() *TradeMessageService {
	return &TradeMessageService{}
}

// SendOrderCreateMessage 发送订单创建消息
func (s *TradeMessageService) SendOrderCreateMessage(ctx context.Context, order *trade.TradeOrder) error {
	// TODO: Implement message sending logic (SMS, Station Letter)
	return nil
}

// SendOrderPaySuccessMessage 发送订单支付成功消息
func (s *TradeMessageService) SendOrderPaySuccessMessage(ctx context.Context, order *trade.TradeOrder) error {
	// TODO: Implement
	return nil
}

// SendOrderDeliveryMessage 发送订单发货消息
func (s *TradeMessageService) SendOrderDeliveryMessage(ctx context.Context, order *trade.TradeOrder) error {
	// TODO: Implement
	return nil
}
