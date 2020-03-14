package hwcloud

import (
	"bytes"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"strings"
	"time"

	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	uuid "github.com/satori/go.uuid"
)

const (
	// MULTISMSMAX 群发短信单批次最大手机号数量
	MULTISMSMAX int = 1000
)

var (
	//ErrMultiCount 群发号码数量错误
	ErrMultiCount = errors.New("单次提交不超过1000个手机号")

	//ErrRequest 请求失败
	ErrRequest = errors.New("请求失败")
)

type SmsOption struct {
	Appkey         string
	Appsecret      string
	ApiAddress     string
	StatusCallback string
}

//New 初始化
func New(opt SmsOption) *HwSms {
	result := &HwSms{
		opt: opt,
	}
	return result
}

//https://support.huaweicloud.com/devg-msgsms/sms_04_0005.html

//https://console.huaweicloud.com/message/?region=cn-north-4&subType=cn_sms#/msgSms/signatureManage

//SendSms 发送分批短信
func (c *HwSms) SendSms(phone []string, signatrue, templateId string, msgParams ...string) (bool, string, error) {
	if len(phone) > MULTISMSMAX {
		return false, "", ErrMultiCount
	}
	c.checkZN(phone) //检查手机号是否是国内，并且格式设置正确

	digest := c.buildWSSEHeader()
	reqBody, err := json.Marshal(&smsrequest{
		From:           signatrue,
		StatusCallback: c.opt.StatusCallback,
		SmsContent: []*scontent{
			&scontent{
				TemplateId:    templateId,
				TemplateParas: msgParams,
				To:            phone,
				Signature:     "",
			},
		},
	})
	if err != nil {
		return false, "", err
	}
	req, err := http.NewRequest("POST", c.opt.ApiAddress, bytes.NewBuffer(reqBody))
	if err != nil {
		return false, "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", `WSSE realm="SDP",profile="UsernameToken",type="Appkey"`)
	req.Header.Set("X-WSSE", digest)
	//让client略过对证书的校检
	httpClient := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}
	resp, err := httpClient.Do(req)
	if err != nil {
		return false, "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, "", err
	}
	respBody := new(SmsRet)
	err = json.Unmarshal(body, respBody)
	if err != nil {
		return false, "", err
	}
	if len(respBody.Result) > 0 && respBody.Result[0].Status == "000000" {
		return true, respBody.Result[0].SmsMsgId, nil
	}
	return false, "", errors.New(string(body))
}

//checkZN 检查手机号是否是国内，并且格式设置正确
func (c *HwSms) checkZN(phone []string) {
	for i := 0; i < len(phone); i++ {
		if !strings.HasPrefix(phone[i], "+") {
			phone[i] = "+86" + phone[i]
		}
	}
}

//加密生成PasswordDigest
//buildWSSEHeader 构造X-WSSE参数值
func (c *HwSms) buildWSSEHeader() string {
	var (
		now       = time.Now().Format("2006-01-02T15:04:05Z")           //created
		nonce     = strings.Replace(uuid.NewV4().String(), "-", "", -1) //nonce
		material  = []byte(nonce + now + c.opt.Appsecret)
		hashed    = sha256.Sum256(material)
		hexdigest = strings.ToUpper(hex.EncodeToString(hashed[:]))
		base64    = base64.StdEncoding.EncodeToString([]byte(hexdigest)) //PasswordDigest
	)
	return fmt.Sprintf(`UsernameToken Username="%s",PasswordDigest="%s",Nonce="%s",Created="%s"`, c.opt.Appkey, base64, nonce, now)
}
