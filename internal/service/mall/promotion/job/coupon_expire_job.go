package job

import (
	"context"

	"github.com/wxlbd/ruoyi-mall-go/internal/service/mall/promotion"
)

type couponExpireService interface {
	ExpireCoupon(ctx context.Context) (int64, error)
}

// CouponExpireJob 优惠券过期任务
// HandlerName 必须与数据库中的 handler_name 保持一致：couponExpireJob
type CouponExpireJob struct {
	couponService couponExpireService
}

func NewCouponExpireJob(couponService *promotion.CouponService) *CouponExpireJob {
	return &CouponExpireJob{
		couponService: couponService,
	}
}

func (j *CouponExpireJob) Execute(ctx context.Context, param string) error {
	if _, err := j.couponService.ExpireCoupon(ctx); err != nil {
		return err
	}
	return nil
}

func (j *CouponExpireJob) GetHandlerName() string {
	return "couponExpireJob"
}
