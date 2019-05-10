package tcpserver

//TimeOutState 当前状态: 1.发送超时 2.接收超时 3.第一次处理超时 4.处理超时
type TimeOutState byte

const (
	//SendTimeOut .
	SendTimeOut TimeOutState = 1
	//RecvTimeOut .
	RecvTimeOut TimeOutState = 2
	//FirstHandTimeOut .
	FirstHandTimeOut TimeOutState = 3
	//HandTimeOut .
	HandTimeOut TimeOutState = 4
)
