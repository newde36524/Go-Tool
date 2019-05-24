package tcpserver2

//TCPHandle 处理类
type TCPHandle interface {
	ReadPacket(conn *Conn) Packet            //读取包
	OnConnection(conn *Conn)                 //连接建立时处理
	OnMessage(conn *Conn, p Packet)          //每次获取到消息时处理
	OnClose(state ConnState)                 //连接关闭时处理
	OnTimeOut(conn *Conn, code TimeOutState) //超时处理
	OnPanic(conn *Conn, err error)           //Panic时处理
	//发送和接收的超时不应该超过其他packet的超时机制
	// OnSendTimeOut(conn *Conn)                //连接发送超时
	// OnRecvTimeOut(conn *Conn)                //连接接收超时
}
