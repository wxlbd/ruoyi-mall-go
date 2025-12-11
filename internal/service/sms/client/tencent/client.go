package tencent

import (
	"context"

	"backend-go/internal/model"
	"backend-go/internal/service/sms/client"
)

type SmsClient struct {
	id        int64
	sdkAppId  string
	apiKey    string
	apiSecret string
	signature string
}

func NewSmsClient(channel *model.SystemSmsChannel) *SmsClient {
	return &SmsClient{
		id:        channel.ID,
		apiKey:    channel.ApiKey,    // SecretId
		apiSecret: channel.ApiSecret, // SecretKey
		signature: channel.Signature,
		// sdkAppId usually also needed, might be in Extra or Code?
		// Java uses properties. For now, assume it's part of config or parsed from somewhere.
	}
}

func (c *SmsClient) GetId() int64 {
	return c.id
}

func (c *SmsClient) SendSms(ctx context.Context, sendLogId int64, mobile string, apiTemplateId string, templateParams map[string]interface{}) (*client.SmsSendResp, error) {
	// TODO: Integrate Tencent SDK
	return &client.SmsSendResp{
		ApiSendCode:  "SUCCESS",
		ApiSendMsg:   "Tencent Mock Success",
		ApiRequestId: "TENCENT_MOCK_ID",
		ApiSerialNo:  "TENCENT_MOCK_SERIAL",
	}, nil
}
