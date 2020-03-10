package hwcloud

type HwSms struct {
	opt SmsOption
}

type SmsRet struct {
	Result      []*reslts `json:"result"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
}

type smsrequest struct {
	From           string      `json:"from"`
	StatusCallback string      `json:"statusCallback"`
	SmsContent     []*scontent `json:"smsContent"`
}

type scontent struct {
	To            []string `json:"to"`
	TemplateId    string   `json:"templateId"`
	TemplateParas []string `json:"templateParas"`
	Signature     string   `json:"signature"`
}

type reslts struct {
	OriginTo   string `json:"originTo"`
	CreateTime string `json:"createTime"`
	From       string `json:"from"`
	SmsMsgId   string `json:"smsMsgId"`
	Status     string `json:"status"`
}

// //语音通话
// type access struct {
// 	Access_token  string
// 	Refresh_token string
// 	Resultcode    string
// 	Expires_in    string
// 	Resultdesc    string
// }

// type reflash struct {
// 	app_key string
// }

// type ret struct {
// 	Resultcode string
// 	Resultdesc string
// 	SessionId  string
// }

// type requestjson struct {
// 	BindNbr      string `json:"bindNbr"`
// 	CalleeNbr    string `json:"calleeNbr"`
// 	DisplayNbr   string `json:"displayNbr"`
// 	PlayInfoList []*pil `json:"playInfoList"`
// 	StatusUrl    string `json:"statusUrl"`
// }

// type pil struct {
// 	TemplateId    string   `json:"templateId"`
// 	TemplateParas []string `json:"templateParas"`
// }

// type VoiceHelper struct {
// 	pils         *pil
// 	requestjsons *requestjson
// 	rets         *ret
// 	access       *access

// 	app_key        string
// 	username       string
// 	app_secret     string
// 	authorization  string
// 	refresh_token  string
// 	expires_in     string
// 	bindNbr        string
// 	displayNbr     string
// 	templateId     string
// 	sendurl        string
// 	authurl        string
// 	reflashurl     string
// 	accesstoken    string
// 	tokenValieTime time.Time
// 	mux            sync.Mutex
// 	statusUrl      string
// 	debug          string
// }
