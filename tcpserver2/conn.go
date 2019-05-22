package tcpserver2

import (
	"context"
	"net"
	"time"
)

//Conn 连接代理对象
type Conn struct {
	conn     net.Conn        //tcp连接对象
	option   ConnOption      //连接配置项
	state    ConnState       //连接状态
	context  context.Context //全局上下文
	recvChan <-chan Packet   //接收消息队列
	sendChan chan<- Packet   //发送消息队列
	handChan chan<- Packet   //处理消息队列
	cancel   func()          //全局上下文取消函数
	isDebug  bool            //是否打印框架内部debug信息
}

//NewConn returns a wrapper of raw conn
func NewConn(conn net.Conn, option ConnOption) (result *Conn) {
	result = &Conn{
		conn:   conn,
		option: option,
		state: ConnState{
			ActiveTime: time.Now(),
			RemoteAddr: conn.RemoteAddr().String(),
		},
		isDebug: false,
	}
	result.context, result.cancel = context.WithCancel(context.Background())
	return
}

//fnProxy 代理执行方法,用于检测执行超时
func (c *Conn) fnProxy(fn func()) <-chan struct{} {
	result := make(chan struct{}, 1)
	go func() {
		defer func() {
			close(result)
			if err := recover(); err != nil {
				defer recover()
				c.option.Handle.OnPanic(c, err.(error))
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
			c.option.Handle.OnPanic(c, err.(error))
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
	n, err = c.conn.Read(b)
	return
}

//RemoteAddr 客户端IP地址
func (c *Conn) RemoteAddr() string {
	return c.conn.RemoteAddr().String()
}

//run 固定处理流程
func (c *Conn) run() {
	c.recvChan = c.recv(c.option.MaxRecvChanCount)(c.heartBeat(c.option.RecvTimeOut, func() { c.option.Handle.OnTimeOut(c, RecvTimeOutCode) }))
	c.sendChan = c.send(c.option.MaxSendChanCount)(c.heartBeat(c.option.SendTimeOut, func() { c.option.Handle.OnTimeOut(c, SendTimeOutCode) }))
	c.handChan = c.message(1)(c.heartBeat(c.option.HandTimeOut, func() { c.option.Handle.OnTimeOut(c, HandTimeOutCode) }))
	go c.safeFn(func() {
		select {
		case <-c.fnProxy(func() { c.option.Handle.OnConnection(c) }):
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
				c.option.Logger.Debugf("%s: Conn.Close: sendChan is closed", c.RemoteAddr())
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
		c.option.Logger.Errorf("%s: Conn.Send: packet is nil,do nothing", c.RemoteAddr())
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
	defer c.conn.Close()
	c.conn.SetDeadline(time.Now())      //set read timeout
	c.conn.SetWriteDeadline(time.Now()) //set write timeout
	c.state.Message = "conn is closed"
	c.state.ComplateTime = time.Now()
	c.option.Handle.OnClose(c.state)
	c.cancel()
	// runtime.GC()         //强制GC      待定可能有问题
	// debug.FreeOSMemory() //强制释放内存 待定可能有问题
}

//readPacket 读取一个包
// func (c *Conn) readPacket() Packet {
// 	logs.Info("readPacket")
// 	p, _ := c.option.Handle.ReadPacket(nil, c)
// 	// if err != nil {
// 	// 	c.option.Logger.Error(err)
// 	// }
// 	return p
// }
func (c *Conn) readPacket() <-chan Packet {
	result := make(chan Packet)
	go c.safeFn(func() {
		defer func() {
			close(result)
		}()
		p := c.option.Handle.ReadPacket(c)
		select {
		case <-c.context.Done():
		case result <- p:
		}
	})
	return result
}

//recv 创建一个可接收 packet channel
func (c *Conn) recv(maxRecvChanCount int) func(chan<- struct{}) <-chan Packet {
	return func(heartBeat chan<- struct{}) <-chan Packet {
		result := make(chan Packet, maxRecvChanCount)
		go c.safeFn(func() {
			defer func() {
				close(result)
				if c.isDebug {
					c.option.Logger.Debugf("%s: Conn.recv: recvChan is closed", c.RemoteAddr())
					c.option.Logger.Debugf("%s: Conn.recv: recv goruntinue exit", c.RemoteAddr())
				}
			}()
			for c.conn != nil {
				ch := c.readPacket()
				select {
				case <-c.context.Done():
					return
				case result <- <-ch:
					if c.isDebug {
						c.option.Logger.Debugf("%s: Conn.recv: read a packet", c.RemoteAddr())
					}
					select {
					case heartBeat <- struct{}{}:
					default:
					}
				}
			}
		})
		return result
	}
}

//send 创建一个可发送 packet channel
func (c *Conn) send(maxSendChanCount int) func(chan<- struct{}) chan<- Packet {
	return func(heartBeat chan<- struct{}) chan<- Packet {
		result := make(chan Packet, maxSendChanCount)
		go c.safeFn(func() {
			defer func() {
				if c.isDebug {
					c.option.Logger.Debugf("%s: Conn.send: send goruntinue exit", c.RemoteAddr())
				}
			}()
			for c.conn != nil {
				select {
				case <-c.context.Done():
					return
				case packet, ok := <-result:
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
					_, err = c.conn.Write(sendData)
					select {
					case heartBeat <- struct{}{}:
					default:
					}
					if err != nil {
						c.option.Logger.Error(err)
					} else {
						if c.isDebug {
							c.option.Logger.Debugf("%s: Conn.send: send a packet", c.RemoteAddr())
						}
					}
				}
			}
		})
		return result
	}
}

//message 创建一个可发送 hand packet channel
func (c *Conn) message(maxHandNum int) func(chan<- struct{}) chan<- Packet {
	return func(heartBeat chan<- struct{}) chan<- Packet {
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
					c.option.Handle.OnMessage(c, p)
					select {
					case heartBeat <- struct{}{}:
					default:
					}
					if c.isDebug {
						c.option.Logger.Debugf("%s: Conn.message: hand a packet", c.RemoteAddr())
					}
				}
			}
		})
		return result
	}
}

//heartBeat 协程心跳检测
func (c *Conn) heartBeat(timeOut time.Duration, callback func()) chan<- struct{} {
	result := make(chan struct{})
	go func() {
		defer close(result)
		for {
			select {
			case _, ok := <-result:
				if !ok {
					c.option.Logger.Error("heart beat chan was closed, are you sure all goruntinue is exit?")
				}
			case <-c.context.Done():
				return
			case <-time.After(timeOut):
				callback()
			}
		}
	}()
	return result
}
