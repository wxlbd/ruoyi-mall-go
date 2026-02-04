package job

import (
	"context"

	tradeService "github.com/wxlbd/ruoyi-mall-go/internal/service/mall/trade"
)

type tradeOrderAutoCancelService interface {
	CancelOrderBySystem(ctx context.Context) (int64, error)
}

// TradeOrderAutoCancelJob 交易订单自动过期 Job
// HandlerName 必须与数据库中的 handler_name 保持一致：tradeOrderAutoCancelJob
type TradeOrderAutoCancelJob struct {
	orderService tradeOrderAutoCancelService
}

func NewTradeOrderAutoCancelJob(orderService *tradeService.TradeOrderUpdateService) *TradeOrderAutoCancelJob {
	return &TradeOrderAutoCancelJob{orderService: orderService}
}

func (j *TradeOrderAutoCancelJob) Execute(ctx context.Context, param string) error {
	_, err := j.orderService.CancelOrderBySystem(ctx)
	return err
}

func (j *TradeOrderAutoCancelJob) GetHandlerName() string {
	return "tradeOrderAutoCancelJob"
}
