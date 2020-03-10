package tencent

import (
	"fmt"
	"testing"
)

const (
	appID  = "xxx"
	appKey = "xxxxxxx"
	sign   = "xxxxx"
)

func TestSms(t *testing.T) {
	sms := New(SmsOption{
		AppID:  appID,
		AppKey: appKey,
	})
	b, s, err := sms.SendSms([]string{
		"10086",
	}, sign, "558496")
	fmt.Println(b)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(s)
}
func TestSms2(t *testing.T) {
	sms := New(SmsOption{
		AppID:  appID,
		AppKey: appKey,
	})
	b, s, err := sms.SendSms([]string{
		"10086",
	}, sign, "123654", "2225")
	fmt.Println(b)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(s)
}
