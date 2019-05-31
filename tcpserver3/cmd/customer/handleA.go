package customer

import (
	tcp "github.com/newde36524/Go-Tool/tcpserver3"

	"github.com/issue9/logs"
)

//TCPHandleA tcpserver使用示例,回复相同的内容
type TCPHandleA struct {
	tcp.TCPHandle
	//可增加新的属性
	//可增加全局属性，比如多个客户端连接可选择转发数据给其他连接，而增加一个全局map
}

//ReadPacket .
func (TCPHandleA) ReadPacket(conn *tcp.Conn, next func()) tcp.Packet {
	// //todo 定义读取数据帧的规则
	// b := make([]byte, 1024)
	// n, err := conn.Read(b)
	// if err != nil {
	// 	switch e := err.(type) {
	// 	case net.Error:
	// 		if !e.Timeout() {
	// 			logs.Error(err)
	// 			conn.Close()
	// 		}
	// 	}
	// }
	// p := &Packet{}
	// p.SetBuffer(b[:n])
	next()
	return nil
}

//OnConnection .
func (TCPHandleA) OnConnection(conn *tcp.Conn, next func()) {
	//todo 连接建立时处理,用于一些建立连接时,需要主动下发数据包的场景,可以在这里开启心跳协程,做登录验证等等
	logs.Infof("%s: 对方好像对你很感兴趣呦~~AAAAAAAAAAAAAA   start", conn.RemoteAddr())
	next()
	logs.Infof("%s: 对方好像对你很感兴趣呦~~AAAAAAAAAAAAAA   end", conn.RemoteAddr())
}

//OnMessage .
func (TCPHandleA) OnMessage(conn *tcp.Conn, p tcp.Packet, next func()) {
	//todo 处理接收的包
	logs.Infof("%s:我好像收到了不知名快递哦~~AAAAAAAAAA start", conn.RemoteAddr())

	sendP := &Packet{}
	if p != nil {
		data := p.GetBuffer()
		sendP.SetBuffer(data)
	}
	conn.Write(sendP) //回复客户端发送的内容
	next()
	logs.Infof("%s:我好像收到了不知名快递哦~~AAAAAAAAAA end", conn.RemoteAddr())
}

//OnClose .
func (TCPHandleA) OnClose(state tcp.ConnState, next func()) {
	logs.Infof("对方好像撤退了呦~~,连接状态:%s", state.String())
}

//OnTimeOut .
func (TCPHandleA) OnTimeOut(conn *tcp.Conn, code tcp.TimeOutState, next func()) {
	logs.Infof("%s: 对方好像在做一些灰暗的事情呢~~,超时类型:%d", conn.RemoteAddr(), code)
}

//OnPanic .
func (TCPHandleA) OnPanic(conn *tcp.Conn, err error, next func()) {
	logs.Errorf("%s: 对方好像发生了一些不得了的事情哦~~,错误信息:%s", conn.RemoteAddr(), err)
}

//OnSendError .
func (TCPHandleA) OnSendError(conn *tcp.Conn, packet tcp.Packet, err error, next func()) {
	logs.Errorf("%s: 发送数据的时间好像有点久诶~~,错误信息:%s", conn.RemoteAddr(), err)
}

//OnRecvError .
func (TCPHandleA) OnRecvError(conn *tcp.Conn, err error, next func()) {
	logs.Errorf("%s: 接收数据的时间好像有点久诶~~,错误信息:%s", conn.RemoteAddr(), err)
}
