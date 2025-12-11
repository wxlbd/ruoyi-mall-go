package debug

import (
	"context"

	"backend-go/internal/model"
	"backend-go/internal/service/sms/client"

	"go.uber.org/zap"
)

type SmsClient struct {
	id        int64
	apiKey    string
	apiSecret string
	signature string
}

func NewSmsClient(channel *model.SystemSmsChannel) *SmsClient {
	return &SmsClient{
		id:        channel.ID,
		apiKey:    channel.ApiKey,
		apiSecret: channel.ApiSecret,
		signature: channel.Signature,
	}
}

func (c *SmsClient) GetId() int64 {
	return c.id
}

func (c *SmsClient) SendSms(ctx context.Context, sendLogId int64, mobile string, apiTemplateId string, templateParams map[string]interface{}) (*client.SmsSendResp, error) {
	zap.L().Info("Debug Sms Client Send Sms",
		zap.Int64("sendLogId", sendLogId),
		zap.String("mobile", mobile),
		zap.String("apiTemplateId", apiTemplateId),
		zap.Any("params", templateParams),
	)
	return &client.SmsSendResp{
		ApiSendCode:  "SUCCESS",
		ApiSendMsg:   "Debug Send Success",
		ApiRequestId: "DEBUG_REQUEST_ID",
		ApiSerialNo:  "DEBUG_SERIAL_NO",
	}, nil
}
