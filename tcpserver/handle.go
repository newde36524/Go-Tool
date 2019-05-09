package tcpserver

// TCPHandle 处理类
type TCPHandle interface {
	ReadPacket(conn *Conn) (Packet, error) //读取包
	OnFirst(conn *Conn, p Packet) error    //处理第一个包,全局只执行一次
	OnMessage(conn *Conn, p Packet) error  //每次获取到消息时处理
	OnClose()                              //连接关闭时处理
}
