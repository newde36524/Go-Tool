package tcpserver

import "context"

//TCPHandle 处理类
type TCPHandle interface {
	ReadPacket(context context.Context, conn *Conn) (Packet, error) //读取包
	OnConnection(conn *Conn)                                        //连接建立时处理
	OnMessage(conn *Conn, p Packet) error                           //每次获取到消息时处理
	OnClose(state ConnState)                                        //连接关闭时处理
	OnTimeOut(conn *Conn, code TimeOutState)                        //超时处理
	OnPanic(conn *Conn, err error)                                  //Panic时处理
}
