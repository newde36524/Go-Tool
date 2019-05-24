package tcpserver2

//TCPHandle 处理类
type TCPHandle interface {
	ReadPacket(conn *Conn) Packet                     //读取包
	OnConnection(conn *Conn)                          //连接建立时处理
	OnMessage(conn *Conn, p Packet)                   //每次获取到消息时处理
	OnClose(state ConnState)                          //连接关闭时处理
	OnTimeOut(conn *Conn, code TimeOutState)          //超时处理
	OnPanic(conn *Conn, err error)                    //Panic时处理
	OnSendError(conn *Conn, packet Packet, err error) //连接数据发送异常 发送和接收的超时不应该超过其他packet的超时机制
	OnRecvError(conn *Conn, err error)                //连接数据接收异常
}
