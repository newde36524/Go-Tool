package tcpserver

import (
	"fmt"
	"time"
)

//ConnState 连接状态
type ConnState struct {
	ActiveTime   time.Time //开始活动时间
	ComplateTime time.Time //结束活动时间
	InnerErr     error     //异常信息
	Message      string    //通知信息
}

//String 格式化输出结构体信息
func (state *ConnState) String() string {
	return fmt.Sprintf("开始活动时间:%s 结束活动时间:%s 异常信息:%s 通知信息:%s", state.ActiveTime, state.ComplateTime, state.InnerErr, state.Message)
}
