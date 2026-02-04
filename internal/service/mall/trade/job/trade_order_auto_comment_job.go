package job

import (
	"context"

	tradeService "github.com/wxlbd/ruoyi-mall-go/internal/service/mall/trade"
)

type tradeOrderAutoCommentService interface {
	CreateOrderItemCommentBySystem(ctx context.Context) (int64, error)
}

// TradeOrderAutoCommentJob 交易订单自动评论 Job
// HandlerName 必须与数据库中的 handler_name 保持一致：tradeOrderAutoCommentJob
type TradeOrderAutoCommentJob struct {
	orderService tradeOrderAutoCommentService
}

func NewTradeOrderAutoCommentJob(orderService *tradeService.TradeOrderUpdateService) *TradeOrderAutoCommentJob {
	return &TradeOrderAutoCommentJob{orderService: orderService}
}

func (j *TradeOrderAutoCommentJob) Execute(ctx context.Context, param string) error {
	_, err := j.orderService.CreateOrderItemCommentBySystem(ctx)
	return err
}

func (j *TradeOrderAutoCommentJob) GetHandlerName() string {
	return "tradeOrderAutoCommentJob"
}
