package repo

import (
	"backend-go/internal/model/trade"
	"backend-go/internal/repo/query"
	"context"
)

type TradeOrderLogRepository struct {
	q *query.Query
}

func NewTradeOrderLogRepository(q *query.Query) *TradeOrderLogRepository {
	return &TradeOrderLogRepository{
		q: q,
	}
}

func (r *TradeOrderLogRepository) Create(ctx context.Context, log *trade.TradeOrderLog) error {
	return r.q.TradeOrderLog.WithContext(ctx).Create(log)
}

func (r *TradeOrderLogRepository) GetListByOrderId(ctx context.Context, orderId int64) ([]*trade.TradeOrderLog, error) {
	return r.q.TradeOrderLog.WithContext(ctx).Where(r.q.TradeOrderLog.OrderID.Eq(orderId)).Find()
}
