package aliyun

import (
	"context"

	"backend-go/internal/model"
	"backend-go/internal/service/sms/client"
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
	// TODO: Integrate Aliyun SDK
	// For now, return mock success or error to avoid compilation failure if SDK missing
	return &client.SmsSendResp{
		ApiSendCode:  "SUCCESS",
		ApiSendMsg:   "Aliyun Mock Success",
		ApiRequestId: "ALIYUN_MOCK_ID",
		ApiSerialNo:  "ALIYUN_MOCK_SERIAL",
	}, nil
}
