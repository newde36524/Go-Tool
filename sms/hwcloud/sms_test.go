package hwcloud

import (
	"fmt"
	"testing"
)

const (
	appkey     = "xxx"
	appsecret  = "xxxxxxx"
	apiAddress = "https://rtcsms.cn-north-1.myhuaweicloud.com:10743/sms/batchSendDiffSms/v1"
	from       = "xzxzxzx"
	templateId = "qwqwwwww"
)

func TestSms(t *testing.T) {
	sms := New(SmsOption{
		Appkey:         appkey,
		Appsecret:      appsecret,
		ApiAddress:     apiAddress,
		StatusCallback: "",
	})
	b, s, err := sms.SendSms([]string{
		"+8610086",
	}, from, templateId, "123459")
	fmt.Println(b)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(s)
}
