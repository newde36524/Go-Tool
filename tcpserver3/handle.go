package tcpserver3

//Handle 处理类
type Handle interface {
	ReadPacket(conn *Conn, next func()) Packet                     //读取包
	OnConnection(conn *Conn, next func())                          //连接建立时处理
	OnMessage(conn *Conn, p Packet, next func())                   //每次获取到消息时处理
	OnClose(state ConnState, next func())                          //连接关闭时处理
	OnTimeOut(conn *Conn, code TimeOutState, next func())          //超时处理
	OnPanic(conn *Conn, err error, next func())                    //Panic时处理
	OnSendError(conn *Conn, packet Packet, err error, next func()) //连接数据发送异常 发送和接收的超时不应该超过其他packet的超时机制
	OnRecvError(conn *Conn, err error, next func())                //连接数据接收异常
}

//CoreHandle 包装接口实现类
type CoreHandle struct {
	handle Handle
	prev   *CoreHandle
	next   *CoreHandle
}

//NewCoreHandle 实例化
//@h 连接处理程序接口
//@return 返回实例
func NewCoreHandle(h Handle) *CoreHandle {
	return &CoreHandle{
		handle: h,
		prev:   nil,
		next:   nil,
	}
}

//NextHandle 链式调用
//方法解读:框架希望实现类似管道处理，类似AOP的执行效果，
//虽然TCPHandle接口的实现类可以在应用层通过装饰器等其他手段实现AOP效果，但如果封装进框架中会有难度。
//这里采用链表的形式包装每个传进来的TCPHandle接口，并实现链式查找。
//框架遇到的问题是，在接口实现的方法内部，并不知道当前是哪个方法，所以希望这个接口的当前方法调用下一个接口的当前方法，有困难。
//这里采用闭包的形式，在框架的每一处调用接口方法的地方都创建一个闭包，并传入回调，把返回的方法再次传递给接口方法，
//那么每个接口方法的实现通过调用传递进去的方法，能各自访问各自创建的闭包，从而实现管道调用之间不会互相影响，
//至此完成管道的处理流程，关键是闭包的应用
func (h *CoreHandle) NextHandle(callback func(*CoreHandle)) func() {
	node := h
	fn := func() {
		if node.next != nil {
			node = node.next
			callback(node.prev)
		} else {
			callback(node)
		}
	}
	return fn
}

//Link 为当前节点连接并返回下一个节点
func (h *CoreHandle) Link(next *CoreHandle) *CoreHandle {
	h.next = next
	next.prev = h
	return next
}

//First 获取传入节点链路中第一个节点
func First(curr *CoreHandle) *CoreHandle {
	for {
		if curr.prev != nil {
			curr = curr.prev
		} else {
			return curr
		}
	}
}

//Last 获取传入节点链路中最后一个节点
func Last(curr *CoreHandle) *CoreHandle {
	for {
		if curr.next != nil {
			curr = curr.next
		} else {
			return curr
		}
	}
}

//Next 获取当前节点的下一个节点
func (h *CoreHandle) Next() *CoreHandle { return h.next }

//Prev 获取当前节点的上一个节点
func (h *CoreHandle) Prev() *CoreHandle { return h.prev }

//ReadPacket .
func (h *CoreHandle) ReadPacket(conn *Conn, next func()) Packet {
	p := h.handle.ReadPacket(conn, next)
	return p
}

//OnConnection .
func (h *CoreHandle) OnConnection(conn *Conn, next func()) { h.handle.OnConnection(conn, next) }

//OnMessage .
func (h *CoreHandle) OnMessage(conn *Conn, p Packet, next func()) {
	h.handle.OnMessage(conn, p, next)
}

//OnClose .
func (h *CoreHandle) OnClose(state ConnState, next func()) { h.handle.OnClose(state, next) }

//OnTimeOut .
func (h *CoreHandle) OnTimeOut(conn *Conn, code TimeOutState, next func()) {
	h.handle.OnTimeOut(conn, code, next)
}

//OnPanic .
func (h *CoreHandle) OnPanic(conn *Conn, err error, next func()) { h.handle.OnPanic(conn, err, next) }

//OnRecvError .
func (h *CoreHandle) OnRecvError(conn *Conn, err error, next func()) {
	h.handle.OnRecvError(conn, err, next)
}

//OnSendError .
func (h *CoreHandle) OnSendError(conn *Conn, p Packet, err error, next func()) {
	h.handle.OnSendError(conn, p, err, next)
}
