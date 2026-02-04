package job

import (
	"context"

	infraService "github.com/wxlbd/ruoyi-mall-go/internal/service/infra"
)

const (
	accessLogCleanRetainDay   = 14
	accessLogCleanDeleteLimit = 100
)

type accessLogCleanService interface {
	CleanAccessLog(ctx context.Context, exceedDay int, deleteLimit int) (int, error)
}

// AccessLogCleanJob 访问日志清理 Job
// HandlerName 必须与数据库中的 handler_name 保持一致：accessLogCleanJob
type AccessLogCleanJob struct {
	apiAccessLogService accessLogCleanService
}

func NewAccessLogCleanJob(apiAccessLogService *infraService.ApiAccessLogService) *AccessLogCleanJob {
	return &AccessLogCleanJob{
		apiAccessLogService: apiAccessLogService,
	}
}

func (j *AccessLogCleanJob) Execute(ctx context.Context, param string) error {
	_, err := j.apiAccessLogService.CleanAccessLog(ctx, accessLogCleanRetainDay, accessLogCleanDeleteLimit)
	return err
}

func (j *AccessLogCleanJob) GetHandlerName() string {
	return "accessLogCleanJob"
}
