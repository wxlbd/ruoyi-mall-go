package client

import "context"

// SmsSendResp 短信发送结果
type SmsSendResp struct {
	ApiSendCode  string
	ApiSendMsg   string
	ApiRequestId string
	ApiSerialNo  string
}

// SmsClient 短信客户端接口
type SmsClient interface {
	// GetId 获得渠道编号
	GetId() int64
	// SendSms 发送消息
	SendSms(ctx context.Context, sendLogId int64, mobile string, apiTemplateId string, templateParams map[string]interface{}) (*SmsSendResp, error)
}
