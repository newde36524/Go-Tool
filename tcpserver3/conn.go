package tcpserver3

import (
	"context"
	"net"
	"runtime/debug"
	"time"
)

//Conn 连接代理对象
type Conn struct {
	rwc      net.Conn        //tcp原始连接对象
	option   ConnOption      //连接配置项
	handle   *CoreTCPHandle  //连接处理程序
	state    ConnState       //连接状态
	context  context.Context //全局上下文
	recvChan <-chan Packet   //接收消息队列
	sendChan chan<- Packet   //发送消息队列
	handChan chan<- Packet   //处理消息队列
	cancel   func()          //全局上下文取消函数
	isDebug  bool            //是否打印框架内部debug信息
}

//NewConn returns a wrapper of raw conn
func NewConn(rwc net.Conn, option ConnOption, h *CoreTCPHandle) (result *Conn) {
	result = &Conn{
		rwc:    rwc,
		option: option,
		handle: h,
		state: ConnState{
			ActiveTime: time.Now(),
			RemoteAddr: rwc.RemoteAddr().String(),
		},
		isDebug: false,
	}
	result.context, result.cancel = context.WithCancel(context.Background())
	return
}

//方法解读:
//框架希望实现类似管道处理，类似AOP的执行效果，
//虽然TCPHandle接口的实现类可以在应用层通过装饰器等其他手段实现AOP效果，但如果封装进框架中会有难度。
//这里采用链表的形式包装每个传进来的TCPHandle接口，并实现链式查找。
//框架遇到的问题是，在接口实现的方法内部，并不知道当前是哪个方法，所以希望这个接口的当前方法调用下一个接口的当前方法，有困难。
//这里采用闭包的形式，在框架的每一处调用接口方法的地方都创建一个闭包，并传入回调，把返回的方法再次传递给接口方法，
//那么每个接口方法的实现通过调用传递进去的方法，能各自访问各自创建的闭包，从而实现管道调用之间不会互相影响，
//至此完成管道的处理流程，关键是闭包的应用
//next 获取下一个连接处理程序,用于作为流式数据处理，实现类似AOP的效果
func (c *Conn) next(h func(*CoreTCPHandle)) func() {
	node := c.handle
	fn := func() {
		node = node.next
		h(node)
	}
	return fn
}

//fnProxy 代理执行方法,用于检测执行超时
func (c *Conn) fnProxy(fn func()) <-chan struct{} {
	result := make(chan struct{}, 1)
	go func() {
		defer func() {
			close(result)
			if err := recover(); err != nil {
				defer recover()
				var next func()
				next = c.next(func(h *CoreTCPHandle) { h.OnPanic(c, err.(error), next) })
				c.handle.OnPanic(c, err.(error), next)
				c.option.Logger.Error(string(debug.Stack()))
			}
		}()
		fn()
		result <- struct{}{}
	}()
	return result
}

//safeFn 代理方法，用于安全调用方法，恢复panic
func (c *Conn) safeFn(fn func()) {
	defer func() {
		if err := recover(); err != nil {
			defer recover()
			var next func()
			next = c.next(func(h *CoreTCPHandle) { h.OnPanic(c, err.(error), next) })
			c.handle.OnPanic(c, err.(error), next)
			c.option.Logger.Error(string(debug.Stack()))
		}
	}()
	fn()
}

//UseDebug 打开框架内部Debug信息
func (c *Conn) UseDebug() {
	c.isDebug = true
}

//Read 从tcp连接中读取数据帧
func (c *Conn) Read(b []byte) (n int, err error) {
	c.rwc.SetReadDeadline(time.Now().Add(c.option.ReadDataTimeOut))
	n, err = c.rwc.Read(b)
	if err != nil {
		var next func()
		next = c.next(func(h *CoreTCPHandle) { h.OnRecvError(c, err, next) })
		c.handle.OnRecvError(c, err, next)
	}
	return
}

//RemoteAddr 客户端IP地址
func (c *Conn) RemoteAddr() string {
	return c.rwc.RemoteAddr().String()
}

//LocalAddr 服务器IP地址
func (c *Conn) LocalAddr() string {
	return c.rwc.LocalAddr().String()
}

//Raw 获取原始连接
func (c *Conn) Raw() net.Conn {
	return c.rwc
}

//run 固定处理流程
func (c *Conn) run() {
	c.recvChan = c.recv(c.option.MaxRecvChanCount)(c.heartBeat(c.option.RecvTimeOut, func() {
		var next func()
		next = c.next(func(h *CoreTCPHandle) { h.OnTimeOut(c, RecvTimeOutCode, next) })
		c.handle.OnTimeOut(c, RecvTimeOutCode, next)
	}))
	c.sendChan = c.send(c.option.MaxSendChanCount)(c.heartBeat(c.option.SendTimeOut, func() {
		var next func()
		next = c.next(func(h *CoreTCPHandle) { h.OnTimeOut(c, SendTimeOutCode, next) })
		c.handle.OnTimeOut(c, SendTimeOutCode, next)
	}))
	c.handChan = c.message(1)(c.heartBeat(c.option.HandTimeOut, func() {
		var next func()
		next = c.next(func(h *CoreTCPHandle) { h.OnTimeOut(c, HandTimeOutCode, next) })
		c.handle.OnTimeOut(c, HandTimeOutCode, next)
	}))
	go c.safeFn(func() {
		select {
		case <-c.fnProxy(func() {
			var next func()
			next = c.next(func(h *CoreTCPHandle) { h.OnConnection(c, next) })
			c.handle.OnConnection(c, next)
		}):
		case <-time.After(c.option.SendTimeOut):
			c.option.Logger.Debugf("%s: Conn.run: OnConnection funtion invoke used time was too long", c.RemoteAddr())
		}
		defer func() {
			close(c.handChan)
			if c.isDebug {
				c.option.Logger.Debugf("%s: Conn.run: handChan is closed", c.RemoteAddr())
			}
			close(c.sendChan)
			if c.isDebug {
				c.option.Logger.Debugf("%s: Conn.run: sendChan is closed", c.RemoteAddr())
				c.option.Logger.Debugf("%s: Conn.run: proxy goruntinue exit", c.RemoteAddr())
			}
		}()
		for {
			select {
			case <-c.context.Done():
				return
			case p, ok := <-c.recvChan:
				if !ok {
					c.option.Logger.Errorf("%s: Conn.run: recvChan is closed", c.RemoteAddr())
				}
				select {
				case <-c.context.Done():
					return
				case c.handChan <- p:
				}
			}
		}
	})
}

//Write 发送消息到客户端
func (c *Conn) Write(packet Packet) {
	if packet == nil {
		c.option.Logger.Errorf("%s: Conn.Write: packet is nil,do nothing", c.RemoteAddr())
		return
	}
	select {
	case <-c.context.Done():
		return
	case c.sendChan <- packet:
	}
}

//Close 关闭服务器和客户端的连接
func (c *Conn) Close() {
	defer c.rwc.Close()
	c.rwc.SetReadDeadline(time.Time{})  //set read timeout
	c.rwc.SetWriteDeadline(time.Time{}) //set write timeout
	c.state.Message = "conn is closed"
	c.state.ComplateTime = time.Now()
	var next func()
	next = c.next(func(h *CoreTCPHandle) { h.OnClose(c.state, next) })
	c.handle.OnClose(c.state, next)
	c.cancel()
	// runtime.GC()         //强制GC      待定可能有问题
	// debug.FreeOSMemory() //强制释放内存 待定可能有问题
}

//readPacket 读取一个包
func (c *Conn) readPacket() <-chan Packet {
	result := make(chan Packet)
	go c.safeFn(func() {
		defer func() {
			close(result)
		}()
		select {
		case <-c.context.Done():
			return
		default:
		}
		var next func()
		var p Packet
		next = c.next(func(h *CoreTCPHandle) {
			_p := h.ReadPacket(c, next)
			//防止内部调用next()方法重复覆盖p的值
			//当前机制保证在管道处理流程中,只要有一个handle的ReadPacket方法返回值不为nil时才有效,之后无效
			if _p != nil && p != nil {
				panic("禁止在管道链路中重复读取生成Packet,在管道中读取数据帧，只能有一个管道返回Packet，其余只能返回nil")
			}
			if _p != nil && p == nil {
				p = _p
			}
		})
		c.handle.ReadPacket(c, next)
		result <- p
	})
	return result
}

//recv 创建一个可接收 packet channel
func (c *Conn) recv(maxRecvChanCount int) func(<-chan struct{}) <-chan Packet {
	return func(heartBeat <-chan struct{}) <-chan Packet {
		result := make(chan Packet, maxRecvChanCount)
		go c.safeFn(func() {
			defer func() {
				close(result)
				if c.isDebug {
					c.option.Logger.Debugf("%s: Conn.recv: recvChan is closed", c.RemoteAddr())
					c.option.Logger.Debugf("%s: Conn.recv: recv goruntinue exit", c.RemoteAddr())
				}
			}()
			for c.rwc != nil {
				ch := c.readPacket()
				select {
				case <-c.context.Done():
					return
				case result <- <-ch:
					c.state.RecvPacketCount++
					if c.isDebug {
						c.option.Logger.Debugf("%s: Conn.recv: read a packet", c.RemoteAddr())
					}
					select {
					case <-heartBeat:
					default:
					}
				}
			}
		})
		return result
	}
}

//send 创建一个可发送 packet channel
func (c *Conn) send(maxSendChanCount int) func(<-chan struct{}) chan<- Packet {
	return func(heartBeat <-chan struct{}) chan<- Packet {
		result := make(chan Packet, maxSendChanCount)
		go c.safeFn(func() {
			defer func() {
				if c.isDebug {
					c.option.Logger.Debugf("%s: Conn.send: send goruntinue exit", c.RemoteAddr())
				}
			}()
			for c.rwc != nil {
				select {
				case <-c.context.Done():
					return
				case packet, ok := <-result:
					c.state.SendPacketCount++
					if !ok {
						if c.isDebug {
							c.option.Logger.Errorf("%s: Conn.send: send packet chan was closed", c.RemoteAddr())
						}
						return
					}
					if packet == nil {
						c.option.Logger.Errorf("%s: Conn.send: the send packet is nil", c.RemoteAddr())
						break
					}
					sendData, err := packet.Serialize(nil)
					if err != nil {
						c.option.Logger.Error(err)
					}
					c.rwc.SetWriteDeadline(time.Now().Add(c.option.WriteDataTimeOut))
					_, err = c.rwc.Write(sendData)
					if err != nil {
						var next func()
						next = c.next(func(h *CoreTCPHandle) { h.OnSendError(c, packet, err, next) })
						c.handle.OnSendError(c, packet, err, next)
					} else {
						if c.isDebug {
							c.option.Logger.Debugf("%s: Conn.send: send a packet", c.RemoteAddr())
						}
					}
					select {
					case <-heartBeat:
					default:
					}
				}
			}
		})
		return result
	}
}

//message 创建一个可发送 hand packet channel
func (c *Conn) message(maxHandNum int) func(<-chan struct{}) chan<- Packet {
	return func(heartBeat <-chan struct{}) chan<- Packet {
		result := make(chan Packet, maxHandNum)
		go c.safeFn(func() {
			defer func() {
				if c.isDebug {
					c.option.Logger.Debugf("%s: Conn.message: hand goruntinue exit", c.RemoteAddr())
				}
			}()
			for {
				select {
				case <-c.context.Done():
					return
				case p, ok := <-result:
					if !ok {
						c.option.Logger.Errorf("%s: Conn.message: hand packet chan was closed", c.RemoteAddr())
						break
					}
					var next func()
					next = c.next(func(h *CoreTCPHandle) { h.OnMessage(c, p, next) })
					c.handle.OnMessage(c, p, next)
					if c.isDebug {
						c.option.Logger.Debugf("%s: Conn.message: hand a packet", c.RemoteAddr())
					}
					select {
					case <-heartBeat:
					default:
					}
				}
			}
		})
		return result
	}
}

//heartBeat 协程心跳检测
func (c *Conn) heartBeat(timeOut time.Duration, callback func()) <-chan struct{} {
	result := make(chan struct{})
	go func() {
		defer func() {
			close(result)
			if c.isDebug {
				c.option.Logger.Debugf("%s: Conn.heartBeat: heartBeat goruntinue exit", c.RemoteAddr())
			}
		}()
		for {
			select {
			case result <- struct{}{}:
			case <-c.context.Done():
				return
			case <-time.After(timeOut):
				callback()
			}
		}
	}()
	return result
}
