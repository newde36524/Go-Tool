package hwcallback

import "strings"

//华为云短信回调
type HWYNote struct {
	SmsMsgID string `json:"smsMsgId"`
	Status   string `json:"status"`
}

//信息状态转为中文
//https://support.huaweicloud.com/devg-msgsms/sms_04_0010.html
func (this *HWYNote) String() *HWYNote {
	if this.Status == "DELIVRD" {
		this.Status = "用户已成功收到短信"
	} else if this.Status == "EXPIRED" {
		this.Status = "短信已超时"
	} else if this.Status == "DELETED" {
		this.Status = "短信已删除"
	} else if this.Status == "UNDELIV" {
		this.Status = "短信递送失败"
	} else if this.Status == "ACCEPTD" {
		this.Status = "短信已接受"
	} else if this.Status == "UNKNOWN" {
		this.Status = "短信状态未知"
	} else if this.Status == "REJECTD" {
		this.Status = "短信被拒绝"
	} else if this.Status == "RTE_ERR" {
		this.Status = "平台内部路由错误"
	} else if this.Status == "MILIMIT" {
		this.Status = "号码达到分钟下发限制"
	} else if this.Status == "LIMIT" || this.Status == "BEYONDN" {
		this.Status = "号码达到下发限制"
	} else if this.Status == "KEYWORD" {
		this.Status = "短信关键字拦截"
	} else if this.Status == "BLACK" || this.Status == "MBBLACK" {
		this.Status = "号码黑名单"
	} else if this.Status == "DJ:0255" || this.Status == "24" || this.Status == "LT:0001" {
		this.Status = "运营商拦截,一般因为短信内容不允许发送"
	} else if n := strings.Index(this.Status, "MK:"); n != -1 {
		this.Status = "运营商拦截,一般因为短信内容不允许发送"
	} else if n := strings.Index(this.Status, "MA:"); n != -1 {
		this.Status = "SNSC未返回响应消息"
	} else if n := strings.Index(this.Status, "MB:"); n != -1 {
		this.Status = "SNSC返回错误响应消息"
	} else if n := strings.Index(this.Status, "MC:"); n != -1 {
		this.Status = "未从SNSC接收到状态报告"
	} else if n := strings.Index(this.Status, "CA:"); n != -1 {
		this.Status = "SCP未返回响应消息"
	} else if n := strings.Index(this.Status, "CB:"); n != -1 {
		this.Status = "SCP返回错误响应消息"
	} else if n := strings.Index(this.Status, "DA:"); n != -1 {
		this.Status = "DSMP未返回响应消息"
	} else if n := strings.Index(this.Status, "DB:"); n != -1 {
		this.Status = "DSMP返回错误响应消息"
	} else if n := strings.Index(this.Status, "SA:"); n != -1 {
		this.Status = "SP未返回响应消息"
	} else if n := strings.Index(this.Status, "SB:"); n != -1 {
		this.Status = "SP返回错误响应消息"
	} else if n := strings.Index(this.Status, "IA:"); n != -1 {
		this.Status = "下一级ISMG未返回响应消息"
	} else if n := strings.Index(this.Status, "IB:"); n != -1 {
		this.Status = "下一级ISMG返回错误响应消息"
	} else if n := strings.Index(this.Status, "IC:"); n != -1 {
		this.Status = "没有下一级ISMG处接收到状态报告"
	} else {
		this.Status = "未知"
	}

	return this
}

//*************************************************************************//

//呼叫状态事件
type CallStatusInfo struct {
	TimesTamp string `json:"timestamp"` //呼叫时间
	SessionID string `json:"sessionId"` //唯一指定一条通话链路的标识ID
	StateCode int    `json:"stateCode"` //呼叫失败原因
	StateText string `json:"stateText"` //呼叫文本状态
}

//通知事件
type FeeInfo struct {
	Direction int    `json:"direction"` //0表示呼出,1表示来电
	SpID      string `json:"spId"`      //该参数标识开发者账号
	AppKey    string `json:"appKey"`    //业务应用的标识
	BindNum   string `json:"bindNum"`   //发起此次呼叫的业务号码
	SessionID string `json:"sessionId"` //该参数是唯一标识Enabler服务器的会话标识
	CallerNum string `json:"callerNum"` //呼叫发起的号码
	CalleeNum string `json:"calleeNum"` //呼叫的发起的被叫号码
}

//华为云语音回调
type HWYVoice struct {
	EventType  string         `json:"eventType"`  //标识api事件通知的类型
	StatusInfo CallStatusInfo `json:"statusInfo"` //呼叫状态事件
	FeeLst     FeeInfo        `json:"feeLst"`     //通知事件
}

//语音状态转中文
func (this *HWYVoice) String() *HWYVoice {
	if this.StatusInfo.StateCode == 7001 {
		this.StatusInfo.StateText = "开发者呼叫频次管控"
	} else if this.StatusInfo.StateCode == 7002 {
		this.StatusInfo.StateText = "应用呼叫频次管控"
	} else if this.StatusInfo.StateCode == 7003 {
		this.StatusInfo.StateText = "显示号码频次管控"
	} else if this.StatusInfo.StateCode == 7004 {
		this.StatusInfo.StateText = "被叫黑名单呼叫管控"
	} else if this.StatusInfo.StateCode == 7005 {
		this.StatusInfo.StateText = "主叫黑名单安全管控"
	} else if this.StatusInfo.StateCode == 7108 {
		this.StatusInfo.StateText = "用户状态已冻结"
	} else if this.StatusInfo.StateCode == 7109 {
		this.StatusInfo.StateText = "语音端口不足"
	} else if this.StatusInfo.StateCode == 8000 {
		this.StatusInfo.StateText = "内部错误"
	} else if this.StatusInfo.StateCode == 8001 {
		this.StatusInfo.StateText = "用户未接续成功"
	} else if this.StatusInfo.StateCode == 8002 {
		this.StatusInfo.StateText = "接续用户时听失败放音"
	} else if this.StatusInfo.StateCode == 8003 {
		this.StatusInfo.StateText = "用户振铃超时"
	} else if this.StatusInfo.StateCode == 8004 {
		this.StatusInfo.StateText = "用户振铃时挂机"
	} else if this.StatusInfo.StateCode == 8005 {
		this.StatusInfo.StateText = "TTS转换失败"
	} else if this.StatusInfo.StateCode == 8006 {
		this.StatusInfo.StateText = "放音文件不存在"
	} else if this.StatusInfo.StateCode == 8007 {
		this.StatusInfo.StateText = "给用户放音失败"
	} else if this.StatusInfo.StateCode == 8008 {
		this.StatusInfo.StateText = "给用户放音收号失败"
	} else if this.StatusInfo.StateCode == 8009 {
		this.StatusInfo.StateText = "主叫用户主动挂机"
	} else if this.StatusInfo.StateCode == 8010 {
		this.StatusInfo.StateText = "超过最大时长挂机"
	} else if this.StatusInfo.StateCode == 8012 {
		this.StatusInfo.StateText = "无效的app_key"
	} else if this.StatusInfo.StateCode == 8013 {
		this.StatusInfo.StateText = "无效的个人小号呼叫"
	} else if this.StatusInfo.StateCode == 8014 {
		this.StatusInfo.StateText = "无效的关系小号呼叫"
	} else if this.StatusInfo.StateCode == 8015 {
		this.StatusInfo.StateText = "给用户录音失败"
	} else if this.StatusInfo.StateCode == 8016 {
		this.StatusInfo.StateText = "关系小号呼叫方向不允许"
	} else if this.StatusInfo.StateCode == 8017 {
		this.StatusInfo.StateText = "SP指示挂机"
	} else if this.StatusInfo.StateCode == 8018 {
		this.StatusInfo.StateText = "业务无权限"
	} else if this.StatusInfo.StateCode == 8019 {
		this.StatusInfo.StateText = "绑定关系不存在"
	} else if this.StatusInfo.StateCode == 8020 {
		this.StatusInfo.StateText = "业务异常"
	} else if this.StatusInfo.StateCode == 8022 {
		this.StatusInfo.StateText = "无效的分机号码"
	} else if this.StatusInfo.StateCode == 8100 {
		this.StatusInfo.StateText = "被叫号码不存在"
	} else if this.StatusInfo.StateCode == 8101 {
		this.StatusInfo.StateText = "被叫无应答"
	} else if this.StatusInfo.StateCode == 8102 {
		this.StatusInfo.StateText = "被叫用户正忙"
	} else if this.StatusInfo.StateCode == 8103 {
		this.StatusInfo.StateText = "被叫用户暂时无法接通"
	} else if this.StatusInfo.StateCode == 8104 {
		this.StatusInfo.StateText = "被叫已关机"
	} else if this.StatusInfo.StateCode == 8105 {
		this.StatusInfo.StateText = "被叫已停机"
	} else if this.StatusInfo.StateCode == 8106 {
		this.StatusInfo.StateText = "主叫号码不存在"
	} else if this.StatusInfo.StateCode == 8107 {
		this.StatusInfo.StateText = "主叫无应答"
	} else if this.StatusInfo.StateCode == 8108 {
		this.StatusInfo.StateText = "主叫用户正忙"
	} else if this.StatusInfo.StateCode == 8109 {
		this.StatusInfo.StateText = "主叫用户暂时无法接通"
	} else if this.StatusInfo.StateCode == 8110 {
		this.StatusInfo.StateText = "主叫已关机"
	} else if this.StatusInfo.StateCode == 8111 {
		this.StatusInfo.StateText = "主叫已停机"
	} else {
		this.StatusInfo.StateText = "成功"
	}

	return this
}
