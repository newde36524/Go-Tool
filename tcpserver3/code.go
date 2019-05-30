package tcpserver3

//TimeOutState 当前状态: 1.发送超时 2.接收超时 3.处理超时
type TimeOutState byte

const (
	//SendTimeOut 发送超时
	SendTimeOutCode TimeOutState = 1
	//RecvTimeOut 接收超时
	RecvTimeOutCode TimeOutState = 2
	//HandTimeOut 处理超时
	HandTimeOutCode TimeOutState = 3
)
