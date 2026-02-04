package job

import (
	"context"

	tradeService "github.com/wxlbd/ruoyi-mall-go/internal/service/mall/trade"
)

type tradeOrderAutoReceiveService interface {
	ReceiveOrderBySystem(ctx context.Context) (int64, error)
}

// TradeOrderAutoReceiveJob 交易订单自动收货 Job
// HandlerName 必须与数据库中的 handler_name 保持一致：tradeOrderAutoReceiveJob
type TradeOrderAutoReceiveJob struct {
	orderService tradeOrderAutoReceiveService
}

func NewTradeOrderAutoReceiveJob(orderService *tradeService.TradeOrderUpdateService) *TradeOrderAutoReceiveJob {
	return &TradeOrderAutoReceiveJob{orderService: orderService}
}

func (j *TradeOrderAutoReceiveJob) Execute(ctx context.Context, param string) error {
	_, err := j.orderService.ReceiveOrderBySystem(ctx)
	return err
}

func (j *TradeOrderAutoReceiveJob) GetHandlerName() string {
	return "tradeOrderAutoReceiveJob"
}
