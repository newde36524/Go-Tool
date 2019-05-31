package tcpserver3

//TCPHandle 处理类
type TCPHandle interface {
	ReadPacket(conn *Conn, next func()) Packet                     //读取包
	OnConnection(conn *Conn, next func())                          //连接建立时处理
	OnMessage(conn *Conn, p Packet, next func())                   //每次获取到消息时处理
	OnClose(state ConnState, next func())                          //连接关闭时处理
	OnTimeOut(conn *Conn, code TimeOutState, next func())          //超时处理
	OnPanic(conn *Conn, err error, next func())                    //Panic时处理
	OnSendError(conn *Conn, packet Packet, err error, next func()) //连接数据发送异常 发送和接收的超时不应该超过其他packet的超时机制
	OnRecvError(conn *Conn, err error, next func())                //连接数据接收异常
}

/*
	issure:
	1:  无法得知在接口的接口实现方法中调用next知否能够调用下一个接口同样的接口方法
	2:  无法在方法内部得知调用了哪个方法
*/

//CoreTCPHandle 包装接口实现类
type CoreTCPHandle struct {
	handle TCPHandle
	prev   *CoreTCPHandle
	next   *CoreTCPHandle
}

//NewCoreTCPHandle 实例化
//@h 连接处理程序接口
//@return 返回实例
func NewCoreTCPHandle(h TCPHandle) *CoreTCPHandle {
	return &CoreTCPHandle{
		handle: h,
	}
}

//Link 为当前节点连接并返回下一个节点
func (h *CoreTCPHandle) Link(next *CoreTCPHandle) *CoreTCPHandle {
	h.next = next
	next.prev = h
	return next
}

//First 获取传入节点链路中第一个节点
func First(curr *CoreTCPHandle) *CoreTCPHandle {
	if curr.prev != nil {
		First(curr.prev)
	}
	return curr
}

//Last 获取传入节点链路中最后一个节点
func Last(curr *CoreTCPHandle) *CoreTCPHandle {
	if curr.next != nil {
		Last(curr.next)
	}
	return curr
}

//Next 获取当前节点的下一个节点
func (h *CoreTCPHandle) Next() *CoreTCPHandle { return h.next }

//Prev 获取当前节点的上一个节点
func (h *CoreTCPHandle) Prev() *CoreTCPHandle { return h.prev }

func (h *CoreTCPHandle) ReadPacket(conn *Conn, next func()) Packet {
	return h.handle.ReadPacket(conn, next)
}
func (h *CoreTCPHandle) OnConnection(conn *Conn, next func()) { h.handle.OnConnection(conn, next) }
func (h *CoreTCPHandle) OnMessage(conn *Conn, p Packet, next func()) {
	h.handle.OnMessage(conn, p, next)
}
func (h *CoreTCPHandle) OnClose(state ConnState, next func()) { h.handle.OnClose(state, next) }
func (h *CoreTCPHandle) OnTimeOut(conn *Conn, code TimeOutState, next func()) {
	h.handle.OnTimeOut(conn, code, next)
}
func (h *CoreTCPHandle) OnPanic(conn *Conn, err error, next func()) { h.handle.OnPanic(conn, err, next) }
func (h *CoreTCPHandle) OnRecvError(conn *Conn, err error, next func()) {
	h.handle.OnRecvError(conn, err, next)
}
func (h *CoreTCPHandle) OnSendError(conn *Conn, p Packet, err error, next func()) {
	h.handle.OnSendError(conn, p, err, next)
}