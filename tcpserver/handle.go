package tcpserver

import "context"

//TCPHandle 处理类
type TCPHandle interface {
	ReadPacket(conn *Conn, context context.Context) (Packet, error) //读取包
	OnMessage(conn *Conn, p Packet) error                           //每次获取到消息时处理
	OnClose()                                                       //连接关闭时处理
	OnTimeOut(conn *Conn, code TimeOutState)                        //超时处理
}
