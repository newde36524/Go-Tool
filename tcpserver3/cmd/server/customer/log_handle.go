package customer

import (
	"time"

	"github.com/issue9/logs"

	tcp "github.com/newde36524/Go-Tool/tcpserver3"
)

//LogTCPHandle tcpserver使用示例,打印相关日志
type LogTCPHandle struct {
	tcp.TCPHandle
	//可增加新的属性
	//可增加全局属性，比如多个客户端连接可选择转发数据给其他连接，而增加一个全局map
}

//ReadPacket .
func (LogTCPHandle) ReadPacket(conn *tcp.Conn, next func()) tcp.Packet {
	next()
	return nil
}

//OnConnection .
func (LogTCPHandle) OnConnection(conn *tcp.Conn, next func()) {
	next()
}

//OnMessage .
func (LogTCPHandle) OnMessage(conn *tcp.Conn, p tcp.Packet, next func()) {
	startTime := time.Now()
	next()
	endTime := time.Now()
	sub := endTime.Sub(startTime).Seconds() * 1000
	logs.Infof("日志模块输出:  开始时间:%s  结束时间:%s 总耗时: %f ms", startTime.Format("2006-01-02 15:04:05"), endTime.Format("2006-01-02 15:04:05"), sub)
}

//OnClose .
func (LogTCPHandle) OnClose(state tcp.ConnState, next func()) {
	next()
}

//OnTimeOut .
func (LogTCPHandle) OnTimeOut(conn *tcp.Conn, code tcp.TimeOutState, next func()) {
	next()
}

//OnPanic .
func (LogTCPHandle) OnPanic(conn *tcp.Conn, err error, next func()) {
	next()
}

//OnSendError .
func (LogTCPHandle) OnSendError(conn *tcp.Conn, packet tcp.Packet, err error, next func()) {
	next()
}

//OnRecvError .
func (LogTCPHandle) OnRecvError(conn *tcp.Conn, err error, next func()) {
	next()
}
