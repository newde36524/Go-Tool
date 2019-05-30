package tcpserver3

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

type CoreTCPHandle struct {
	handle  []TCPHandle
	current TCPHandle
	index   int
}

func NewCoreTCPHandle(h []TCPHandle) *CoreTCPHandle {
	return &CoreTCPHandle{
		handle:  h,
		current: h[0],
		index:   0,
	}
}

func (h *CoreTCPHandle) ReadPacket(conn *Conn) Packet {
	p := h.handle[0].ReadPacket(conn)

	return p
}
func (h *CoreTCPHandle) OnConnection(conn *Conn)                 { h.handle[0].OnConnection(conn) }
func (h *CoreTCPHandle) OnMessage(conn *Conn, p Packet)          { h.handle[0].OnMessage(conn, p) }
func (h *CoreTCPHandle) OnClose(state ConnState)                 { h.handle[0].OnClose(state) }
func (h *CoreTCPHandle) OnTimeOut(conn *Conn, code TimeOutState) { h.handle[0].OnTimeOut(conn, code) }
func (h *CoreTCPHandle) OnPanic(conn *Conn, err error)           { h.handle[0].OnPanic(conn, err) }
func (h *CoreTCPHandle) OnRecvError(conn *Conn, err error)       { h.handle[0].OnRecvError(conn, err) }
func (h *CoreTCPHandle) OnSendError(conn *Conn, p Packet, err error) {
	h.handle[0].OnSendError(conn, p, err)
}
func (h *CoreTCPHandle) Next() {

}
