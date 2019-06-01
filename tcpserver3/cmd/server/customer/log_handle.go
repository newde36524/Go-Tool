package customer

import (
	"time"

	"github.com/issue9/logs"

	srv "github.com/newde36524/Go-Tool/tcpserver3"
)

//LogHandle tcpserver使用示例,打印相关日志
type LogHandle struct {
	srv.Handle
	//可增加新的属性
	//可增加全局属性，比如多个客户端连接可选择转发数据给其他连接，而增加一个全局map
}

//ReadPacket .
func (LogHandle) ReadPacket(conn *srv.Conn, next func()) srv.Packet {
	next()
	return nil
}

//OnConnection .
func (LogHandle) OnConnection(conn *srv.Conn, next func()) {
	next()
}

//OnMessage .
func (LogHandle) OnMessage(conn *srv.Conn, p srv.Packet, next func()) {
	startTime := time.Now()
	next()
	endTime := time.Now()
	sub := endTime.Sub(startTime).Seconds() * 1000
	logs.Infof("日志模块输出:  开始时间:%s  结束时间:%s 总耗时: %f ms", startTime.Format("2006-01-02 15:04:05"), endTime.Format("2006-01-02 15:04:05"), sub)
}

//OnClose .
func (LogHandle) OnClose(state srv.ConnState, next func()) {
	next()
}

//OnTimeOut .
func (LogHandle) OnTimeOut(conn *srv.Conn, code srv.TimeOutState, next func()) {
	next()
}

//OnPanic .
func (LogHandle) OnPanic(conn *srv.Conn, err error, next func()) {
	next()
}

//OnSendError .
func (LogHandle) OnSendError(conn *srv.Conn, packet srv.Packet, err error, next func()) {
	next()
}

//OnRecvError .
func (LogHandle) OnRecvError(conn *srv.Conn, err error, next func()) {
	next()
}
