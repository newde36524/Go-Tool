package customer

import (
	"time"

	"github.com/issue9/logs"

	srv "github.com/newde36524/Go-Tool/tcpserver3"
)

//LogTCPHandle tcpserver使用示例,打印相关日志
type LogTCPHandle struct {
	srv.TCPHandle
	//可增加新的属性
	//可增加全局属性，比如多个客户端连接可选择转发数据给其他连接，而增加一个全局map
}

//ReadPacket .
func (LogTCPHandle) ReadPacket(conn *srv.Conn, next func()) srv.Packet {
	next()
	return nil
}

//OnConnection .
func (LogTCPHandle) OnConnection(conn *srv.Conn, next func()) {
	next()
}

//OnMessage .
func (LogTCPHandle) OnMessage(conn *srv.Conn, p srv.Packet, next func()) {
	startTime := time.Now()
	next()
	endTime := time.Now()
	sub := endTime.Sub(startTime).Seconds() * 1000
	logs.Infof("日志模块输出:  开始时间:%s  结束时间:%s 总耗时: %f ms", startTime.Format("2006-01-02 15:04:05"), endTime.Format("2006-01-02 15:04:05"), sub)
}

//OnClose .
func (LogTCPHandle) OnClose(state srv.ConnState, next func()) {
	next()
}

//OnTimeOut .
func (LogTCPHandle) OnTimeOut(conn *srv.Conn, code srv.TimeOutState, next func()) {
	next()
}

//OnPanic .
func (LogTCPHandle) OnPanic(conn *srv.Conn, err error, next func()) {
	next()
}

//OnSendError .
func (LogTCPHandle) OnSendError(conn *srv.Conn, packet srv.Packet, err error, next func()) {
	next()
}

//OnRecvError .
func (LogTCPHandle) OnRecvError(conn *srv.Conn, err error, next func()) {
	next()
}
