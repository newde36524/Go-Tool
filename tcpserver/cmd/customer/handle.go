package customer

import (
	tcp "Go-Tool/tcpserver"
	"context"

	"github.com/issue9/logs"
)

//CustomerTCPHandle tcpserver使用示例，回复相同的内容
type CustomerTCPHandle struct {
	tcp.TCPHandle
}

//ReadPacket .
func (CustomerTCPHandle) ReadPacket(context context.Context, conn *tcp.Conn) (tcp.Packet, error) {
	//todo 定义读取数据帧的规则
	b := make([]byte, 1024)
	n, err := conn.Read(b)
	if err != nil {
		logs.Error(err)
	}
	p := &CustomerPacket{}
	p.SetBuffer(b[:n])

	return p, err
}

//OnMessage .
func (CustomerTCPHandle) OnMessage(conn *tcp.Conn, p tcp.Packet) error {
	//todo 处理接收的包
	sendP := CustomerPacket{}
	data := p.GetBuffer()
	sendP.SetBuffer(data)
	conn.Send(p) //回复客户端发送的内容
	return nil
}

//OnClose .
func (CustomerTCPHandle) OnClose(state tcp.ConnState) {
	logs.Infof("客户端退出，当前连接状态:%s", state.String())
}

//OnTimeOut .
func (CustomerTCPHandle) OnTimeOut(conn *tcp.Conn, code tcp.TimeOutState) {
	logs.Infof("%s: 触发超時，超时类型:%d", conn.RemoteAddr(), code)
}
