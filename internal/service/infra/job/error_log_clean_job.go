package job

import (
	"context"

	infraService "github.com/wxlbd/ruoyi-mall-go/internal/service/infra"
)

const (
	errorLogCleanRetainDay   = 14
	errorLogCleanDeleteLimit = 100
)

type errorLogCleanService interface {
	CleanErrorLog(ctx context.Context, exceedDay int, deleteLimit int) (int, error)
}

// ErrorLogCleanJob 错误日志清理 Job
// HandlerName 必须与数据库中的 handler_name 保持一致：errorLogCleanJob
type ErrorLogCleanJob struct {
	apiErrorLogService errorLogCleanService
}

func NewErrorLogCleanJob(apiErrorLogService *infraService.ApiErrorLogService) *ErrorLogCleanJob {
	return &ErrorLogCleanJob{
		apiErrorLogService: apiErrorLogService,
	}
}

func (j *ErrorLogCleanJob) Execute(ctx context.Context, param string) error {
	_, err := j.apiErrorLogService.CleanErrorLog(ctx, errorLogCleanRetainDay, errorLogCleanDeleteLimit)
	return err
}

func (j *ErrorLogCleanJob) GetHandlerName() string {
	return "errorLogCleanJob"
}
