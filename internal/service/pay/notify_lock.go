package pay

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// PayNotifyLock 支付通知分布式锁
type PayNotifyLock struct {
	rdb *redis.Client
}

func NewPayNotifyLock(rdb *redis.Client) *PayNotifyLock {
	return &PayNotifyLock{rdb: rdb}
}

const (
	payNotifyLockKeyPrefix = "pay:notify:lock:"
	payNotifyLockTimeout   = 120 * time.Second
)

// Lock 加锁并执行函数
// 对齐 Java: PayNotifyLockRedisDAO.lock
func (l *PayNotifyLock) Lock(ctx context.Context, taskID int64, fn func() error) error {
	lockKey := fmt.Sprintf("%s%d", payNotifyLockKeyPrefix, taskID)

	// 尝试获取锁
	acquired, err := l.rdb.SetNX(ctx, lockKey, "1", payNotifyLockTimeout).Result()
	if err != nil {
		return fmt.Errorf("failed to acquire lock: %w", err)
	}

	if !acquired {
		// 锁已被其他进程持有
		return fmt.Errorf("lock already held for task %d", taskID)
	}

	// 确保释放锁
	defer func() {
		l.rdb.Del(ctx, lockKey)
	}()

	// 执行业务逻辑
	return fn()
}
