package tcpserver2

//TCPHandle 处理类
type TCPHandle interface {
	ReadPacket(conn *Conn) Packet            //读取包
	OnConnection(conn *Conn)                 //连接建立时处理
	OnMessage(conn *Conn, p Packet)          //每次获取到消息时处理
	OnClose(state ConnState)                 //连接关闭时处理
	OnTimeOut(conn *Conn, code TimeOutState) //超时处理
	OnPanic(conn *Conn, err error)           //Panic时处理
}
