package job

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type fakeTradeOrderAutoCancelService struct {
	called bool
	err    error
}

func (f *fakeTradeOrderAutoCancelService) CancelOrderBySystem(context.Context) (int64, error) {
	f.called = true
	return 3, f.err
}

type fakeTradeOrderAutoReceiveService struct {
	called bool
	err    error
}

func (f *fakeTradeOrderAutoReceiveService) ReceiveOrderBySystem(context.Context) (int64, error) {
	f.called = true
	return 2, f.err
}

type fakeTradeOrderAutoCommentService struct {
	called bool
	err    error
}

func (f *fakeTradeOrderAutoCommentService) CreateOrderItemCommentBySystem(context.Context) (int64, error) {
	f.called = true
	return 1, f.err
}

func TestTradeOrderAutoCancelJob(t *testing.T) {
	t.Run("handler name", func(t *testing.T) {
		job := &TradeOrderAutoCancelJob{}
		require.Equal(t, "tradeOrderAutoCancelJob", job.GetHandlerName())
	})
	t.Run("execute", func(t *testing.T) {
		svc := &fakeTradeOrderAutoCancelService{}
		job := &TradeOrderAutoCancelJob{orderService: svc}
		err := job.Execute(context.Background(), "")
		require.NoError(t, err)
		require.True(t, svc.called)
	})
	t.Run("execute error", func(t *testing.T) {
		svc := &fakeTradeOrderAutoCancelService{err: errors.New("boom")}
		job := &TradeOrderAutoCancelJob{orderService: svc}
		err := job.Execute(context.Background(), "")
		require.Error(t, err)
		require.True(t, svc.called)
	})
}

func TestTradeOrderAutoReceiveJob(t *testing.T) {
	t.Run("handler name", func(t *testing.T) {
		job := &TradeOrderAutoReceiveJob{}
		require.Equal(t, "tradeOrderAutoReceiveJob", job.GetHandlerName())
	})
	t.Run("execute", func(t *testing.T) {
		svc := &fakeTradeOrderAutoReceiveService{}
		job := &TradeOrderAutoReceiveJob{orderService: svc}
		err := job.Execute(context.Background(), "")
		require.NoError(t, err)
		require.True(t, svc.called)
	})
	t.Run("execute error", func(t *testing.T) {
		svc := &fakeTradeOrderAutoReceiveService{err: errors.New("boom")}
		job := &TradeOrderAutoReceiveJob{orderService: svc}
		err := job.Execute(context.Background(), "")
		require.Error(t, err)
		require.True(t, svc.called)
	})
}

func TestTradeOrderAutoCommentJob(t *testing.T) {
	t.Run("handler name", func(t *testing.T) {
		job := &TradeOrderAutoCommentJob{}
		require.Equal(t, "tradeOrderAutoCommentJob", job.GetHandlerName())
	})
	t.Run("execute", func(t *testing.T) {
		svc := &fakeTradeOrderAutoCommentService{}
		job := &TradeOrderAutoCommentJob{orderService: svc}
		err := job.Execute(context.Background(), "")
		require.NoError(t, err)
		require.True(t, svc.called)
	})
	t.Run("execute error", func(t *testing.T) {
		svc := &fakeTradeOrderAutoCommentService{err: errors.New("boom")}
		job := &TradeOrderAutoCommentJob{orderService: svc}
		err := job.Execute(context.Background(), "")
		require.Error(t, err)
		require.True(t, svc.called)
	})
}
