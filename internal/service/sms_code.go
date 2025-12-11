package service

import (
	"backend-go/internal/pkg/core"
	"backend-go/internal/repo/query"
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const (
	SmsCodeCacheKeyPrefix = "sms:code:"
	SmsCodeExpire         = 5 * time.Minute
)

type SmsCodeService struct {
	q       *query.Query
	rdb     *redis.Client
	factory *SmsClientFactory
}

func NewSmsCodeService(q *query.Query, rdb *redis.Client, factory *SmsClientFactory) *SmsCodeService {
	return &SmsCodeService{
		q:       q,
		rdb:     rdb,
		factory: factory,
	}
}

// SendSmsCode 发送短信验证码
func (s *SmsCodeService) SendSmsCode(ctx context.Context, mobile string, scene int) error {
	// 1. 校验频率 (1分钟内只能发一次)
	rateLimitKey := fmt.Sprintf("sms:rate:%s:%d", mobile, scene)
	exists, err := s.rdb.Exists(ctx, rateLimitKey).Result()
	if err != nil {
		return err
	}
	if exists > 0 {
		return core.NewBizError(1004003001, "发送过于频繁，请稍后再试")
	}
	// Set rate limit key with 60s expiry
	if err := s.rdb.Set(ctx, rateLimitKey, "1", 60*time.Second).Err(); err != nil {
		return err
	}

	// 2. 生成验证码
	code := fmt.Sprintf("%06d", rand.Intn(1000000))

	// 3. 保存到 Redis
	key := s.getCacheKey(mobile, scene)
	if err := s.rdb.Set(ctx, key, code, SmsCodeExpire).Err(); err != nil {
		return err
	}

	// 4. 发送短信
	// 查询启用的渠道
	channel, err := s.q.SystemSmsChannel.WithContext(ctx).Where(s.q.SystemSmsChannel.Status.Eq(0)).First()
	if err != nil {
		zap.L().Error("No enabled SMS channel found", zap.Error(err))
		// For development, allow fallback or just return error
		// return err
	}

	// Get client from factory
	client := s.factory.GetClient(channel.ID)
	if client == nil {
		zap.L().Info("SMS Client not found in factory, initializing...", zap.Int64("channelId", channel.ID))
		s.factory.CreateOrUpdateClient(channel)
		client = s.factory.GetClient(channel.ID)
	}

	if client != nil {
		// Prepare template params (mock)
		params := map[string]interface{}{
			"code": code,
		}
		// TODO: Retrieve valid apiTemplateId from SystemSmsTemplate based on scene

		_, err := client.SendSms(ctx, 0, mobile, "TEMPLATE_ID", params)
		if err != nil {
			zap.L().Error("Send SMS failed", zap.Error(err))
			return err
		}
	} else {
		zap.L().Warn("No SMS client available for channel", zap.String("code", channel.Code))
	}

	return nil
}

// ValidateSmsCode 校验验证码
func (s *SmsCodeService) ValidateSmsCode(ctx context.Context, mobile string, scene int, code string) error {
	key := s.getCacheKey(mobile, scene)
	val, err := s.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return core.NewBizError(1004003003, "验证码已过期或不存在")
	}
	if err != nil {
		return err
	}
	if val != code {
		return core.NewBizError(1004003004, "验证码错误")
	}

	// 验证成功后删除，避免重复使用
	s.rdb.Del(ctx, key)
	return nil
}

func (s *SmsCodeService) getCacheKey(mobile string, scene int) string {
	return fmt.Sprintf("%s%s:%d", SmsCodeCacheKeyPrefix, mobile, scene)
}
