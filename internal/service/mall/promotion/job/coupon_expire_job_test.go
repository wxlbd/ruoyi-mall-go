package job

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type fakeCouponExpireService struct {
	called bool
	err    error
}

func (f *fakeCouponExpireService) ExpireCoupon(context.Context) (int64, error) {
	f.called = true
	return 3, f.err
}

func TestCouponExpireJob_GetHandlerName(t *testing.T) {
	svc := &fakeCouponExpireService{}
	job := &CouponExpireJob{couponService: svc}

	require.Equal(t, "couponExpireJob", job.GetHandlerName())
}

func TestCouponExpireJob_Execute(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &fakeCouponExpireService{}
		job := &CouponExpireJob{couponService: svc}

		err := job.Execute(context.Background(), "")
		require.NoError(t, err)
		require.True(t, svc.called)
	})

	t.Run("service error", func(t *testing.T) {
		svc := &fakeCouponExpireService{err: errors.New("boom")}
		job := &CouponExpireJob{couponService: svc}

		err := job.Execute(context.Background(), "")
		require.Error(t, err)
		require.True(t, svc.called)
	})
}
