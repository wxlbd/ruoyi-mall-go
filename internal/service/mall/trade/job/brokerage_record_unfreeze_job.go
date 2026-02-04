package job

import (
	"context"

	brokerageService "github.com/wxlbd/ruoyi-mall-go/internal/service/mall/trade/brokerage"
)

type brokerageRecordUnfreezeService interface {
	UnfreezeRecord(ctx context.Context, brokerageUserSvc *brokerageService.BrokerageUserService) (int, error)
}

// BrokerageRecordUnfreezeJob 佣金解冻 Job
// HandlerName 必须与数据库中的 handler_name 保持一致：brokerageRecordUnfreezeJob
type BrokerageRecordUnfreezeJob struct {
	recordService brokerageRecordUnfreezeService
	userService   *brokerageService.BrokerageUserService
}

func NewBrokerageRecordUnfreezeJob(
	recordService *brokerageService.BrokerageRecordService,
	userService *brokerageService.BrokerageUserService,
) *BrokerageRecordUnfreezeJob {
	return &BrokerageRecordUnfreezeJob{
		recordService: recordService,
		userService:   userService,
	}
}

func (j *BrokerageRecordUnfreezeJob) Execute(ctx context.Context, param string) error {
	_, err := j.recordService.UnfreezeRecord(ctx, j.userService)
	return err
}

func (j *BrokerageRecordUnfreezeJob) GetHandlerName() string {
	return "brokerageRecordUnfreezeJob"
}
