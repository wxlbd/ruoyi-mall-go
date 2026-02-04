package job

import (
	"context"

	infraService "github.com/wxlbd/ruoyi-mall-go/internal/service/infra"
)

const (
	jobCleanRetainDay = 14
	jobCleanDeleteLim = 100
)

type jobLogCleanService interface {
	CleanJobLog(ctx context.Context, exceedDay int, deleteLimit int) (int, error)
}

// JobLogCleanJob 任务日志清理 Job
// HandlerName 必须与数据库中的 handler_name 保持一致：jobLogCleanJob
type JobLogCleanJob struct {
	jobLogService jobLogCleanService
}

func NewJobLogCleanJob(jobLogService *infraService.JobLogService) *JobLogCleanJob {
	return &JobLogCleanJob{
		jobLogService: jobLogService,
	}
}

func (j *JobLogCleanJob) Execute(ctx context.Context, param string) error {
	_, err := j.jobLogService.CleanJobLog(ctx, jobCleanRetainDay, jobCleanDeleteLim)
	return err
}

func (j *JobLogCleanJob) GetHandlerName() string {
	return "jobLogCleanJob"
}
