package tencent

import (
	"errors"
	"log"
)

//使用方法：
/*
sms.NewSmsClient(APPID,APPSECRET,SIGN).UseTpl().Send("13800138000",a,b,c...)
*/
type SmsClient struct {
	tplId  uint
	client *QcloudSMS
}

func NewSmsClient(appId string, appSecret string, sign string) *SmsClient {
	return &SmsClient{
		client: NewClient(NewOptions(appId, appSecret, sign)),
	}
}
func (c *SmsClient) UseTpl(tplId uint) *SmsClient {
	c.tplId = tplId
	return c
}

func (c *SmsClient) SendCode(phone string, code string) error {
	if c.tplId == 0 {
		return errors.New("请调用 UseTpl 设置短信模版ID FIRST!!")
	}

	var vr = SMSSingleReq{
		TplID: c.tplId, //模版Id
	}
	vr.Params = []string{code} //, "10"} //验证码，有效时间
	vr.Tel.Nationcode = "86"
	vr.Tel.Mobile = phone
	err := c.client.SendSMSSingle(vr)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (c *SmsClient) Send(phone string, params ...string) error {
	if c.tplId == 0 {
		return errors.New("请调用 UseTpl 设置短信模版ID FIRST!!")
	}
	var vr = SMSSingleReq{
		TplID: c.tplId, //模版Id
	}
	if len(params) > 0 {
		vr.Params = params
	} else {
		vr.Params = []string{}
	}

	vr.Tel.Nationcode = "86"
	vr.Tel.Mobile = phone
	err := c.client.SendSMSSingle(vr)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
