package tencent

import (
	"strconv"

	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

//=============================================================
//https://cloud.tencent.com/document/product/382/13297

//https://console.cloud.tencent.com/sms/smsContent/1400079010/1/10

type TencentModel struct {
	client *cvm.Client
	opt    SmsOption
}

type SmsOption struct {
	AppID  string
	AppKey string
}

func New(opt SmsOption) *TencentModel {
	return &TencentModel{
		opt: opt,
	}
}

//SendSms 发送分批短信
func (c *TencentModel) SendSms(phone []string, signatrue, templateId string, msgParams ...string) (bool, string, error) {
	i, err := strconv.Atoi(templateId)
	if err != nil {
		return false, "", err
	}
	client := NewSmsClient(c.opt.AppID, c.opt.AppKey, signatrue).UseTpl(uint(i))
	for _, p := range phone {
		if err := client.Send(p, msgParams...); err != nil {
			return false, "", err
		}
	}
	return true, "", nil
}
