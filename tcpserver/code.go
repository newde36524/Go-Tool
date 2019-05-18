package tcpserver

//TimeOutState 当前状态: 1.发送超时 2.接收超时 3.处理超时
type TimeOutState byte

const (
	//SendTimeOut 发送超时
	SendTimeOut TimeOutState = 1
	//RecvTimeOut 接收超时
	RecvTimeOut TimeOutState = 2
	//HandTimeOut 处理超时
	HandTimeOut TimeOutState = 3
)
