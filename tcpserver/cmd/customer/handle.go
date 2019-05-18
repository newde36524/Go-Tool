package customer

import (
	tcp "Go-Tool/tcpserver"
	"context"

	"github.com/issue9/logs"
)

//TCPHandle tcpserver使用示例，回复相同的内容
type TCPHandle struct {
	tcp.TCPHandle
}

//ReadPacket .
func (TCPHandle) ReadPacket(context context.Context, conn *tcp.Conn) (tcp.Packet, error) {
	//todo 定义读取数据帧的规则
	b := make([]byte, 1024)
	n, err := conn.Read(b)
	if err != nil {
		logs.Error(err)
		conn.Close()
	}
	p := &Packet{}
	p.SetBuffer(b[:n])

	return p, err
}

//OnConnection .
func (TCPHandle) OnConnection(conn *tcp.Conn) {
	//todo 连接建立时处理，用于一些建立连接时，需要主动下发数据包的场景
	logs.Infof("客户端:%s 连接上来了呦~~~", conn.RemoteAddr())
}

//OnMessage .
func (TCPHandle) OnMessage(conn *tcp.Conn, p tcp.Packet) error {
	//todo 处理接收的包
	sendP := Packet{}
	data := p.GetBuffer()
	sendP.SetBuffer(data)
	conn.Send(p) //回复客户端发送的内容
	return nil
}

//OnClose .
func (TCPHandle) OnClose(state tcp.ConnState) {
	logs.Infof("客户端退出，当前连接状态:%s", state.String())
}

//OnTimeOut .
func (TCPHandle) OnTimeOut(conn *tcp.Conn, code tcp.TimeOutState) {
	logs.Infof("%s: 触发超時，超时类型:%d", conn.RemoteAddr(), code)
}

//OnPanic .
func (TCPHandle) OnPanic(conn *tcp.Conn, err error) {
	logs.Error(err)
	conn.Close()
}
