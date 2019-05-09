package tcpserver

import "time"

//ConnState 连接状态
type ConnState struct {
	IsExit       chan struct{} //是否退出
	ActiveTime   time.Time     //开始活动时间
	ComplateTime time.Time     //结束活动时间
	InnerErr     error         //异常信息
	Message      string        //通知信息
}
