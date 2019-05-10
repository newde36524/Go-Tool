package tcpserver

import (
	"context"
	"net"
	"time"
)

//fnProxy 代理执行方法,用于检测执行超时
func fnProxy(fn func()) <-chan struct{} {
	result := make(chan struct{}, 1)
	go func() {
		defer close(result)
		fn()
		result <- struct{}{}
	}()
	return result
}

//Conn 连接代理对象
type Conn struct {
	conn     net.Conn        //tcp连接对象
	option   ConnOption      //连接配置项
	state    ConnState       //连接状态
	context  context.Context //全局上下文
	recvChan <-chan Packet   //接收消息队列
	sendChan chan<- Packet   //发送消息队列
	handChan chan<- Packet   //处理消息队列
	cancel   func()
}

//NewConn returns a wrapper of raw conn
func NewConn(conn net.Conn, option ConnOption) (result *Conn) {
	result = &Conn{
		conn:   conn,
		option: option,
		state: ConnState{
			ActiveTime: time.Now(),
		},
	}
	result.context, result.cancel = context.WithCancel(context.Background())
	return
}

//run 固定处理流程
func (c *Conn) run() {
	c.recvChan = c.recv(c.option.MaxRecvChanCount)
	c.sendChan = c.send(c.option.MaxSendChanCount)
	c.handChan = c.message()
	go func() {
		select {
		case p, ok := <-c.recvChan:
			if !ok {
				c.option.Logger.Debug("Conn.run: recvChan is closed")
			}
			select {
			case <-fnProxy(func() { c.option.Handle.OnFirst(c, p) }):
			case <-time.After(c.option.HandTimeOut):
				c.option.Handle.OnTimeOut(c, FirstHandTimeOut)
			}
		}
		for {
			select {
			case <-c.context.Done():
				return
			case p, ok := <-c.recvChan:
				if !ok {
					c.option.Logger.Debug("Conn.run: recvChan is closed")
				}
				select {
				case <-c.context.Done():
					close(c.handChan)
					return
				case c.handChan <- p:
				}
			}
		}
	}()
}

//Send 发送消息到设备
func (c *Conn) Send(packet Packet) {
	select {
	case <-c.context.Done():
		return
	case c.sendChan <- packet:
	}
}

// Close 关闭服务器和设备的连接
func (c *Conn) Close() {
	defer c.conn.Close()
	c.option.Handle.OnClose()
	c.state.Message = "conn is closed"
	c.state.ComplateTime = time.Now()
	c.cancel()
	close(c.sendChan)
	c.option.Logger.Info(c.state.String())
	// runtime.GC()         //强制GC      待定可能有问题
	// debug.FreeOSMemory() //强制释放内存 待定可能有问题
}

//ReadPacket 读取一个包
func (c *Conn) readPacket() <-chan Packet {
	result := make(chan Packet)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		defer close(result)
		p, err := c.option.Handle.ReadPacket(c, ctx)
		if err != nil {
			c.option.Logger.Error(err)
		}
		select {
		case <-ctx.Done():
		case result <- p:
			cancel()
		}
	}()
	return result
}

//recv 创建一个包接收channel
func (c *Conn) recv(maxRecvChanCount int) <-chan Packet {
	result := make(chan Packet, maxRecvChanCount)
	go func() {
		defer close(result)
		defer func() {
			c.option.Logger.Debugf("%s: recv goruntinue exit", c.conn.RemoteAddr().String())
		}()
		for {
			select {
			case <-c.context.Done():
				return
			case <-time.After(c.option.RecvTimeOut):
				c.option.Handle.OnTimeOut(c, RecvTimeOut)
			case p, ok := <-c.readPacket():
				if ok {
					result <- p
				}
			}
		}
	}()
	return result
}

//send 创建一个包发送channel
func (c *Conn) send(maxSendChanCount int) chan<- Packet {
	result := make(chan Packet, maxSendChanCount)
	go func() {
		defer func() {
			c.option.Logger.Debugf("%s:send goruntinue exit", c.conn.RemoteAddr().String())
		}()
		for {
			select {
			case <-c.context.Done():
				return
			case packet, ok := <-result:
				if !ok {
					c.option.Logger.Debugf("%s: Conn.Send:send chan is closed", c.conn.RemoteAddr().String())
					return
				}
				ctx, cancel := context.WithTimeout(context.Background(), c.option.SendTimeOut)
				select {
				case <-c.context.Done():
					cancel()
					return
				case <-time.After(c.option.SendTimeOut):
					c.option.Handle.OnTimeOut(c, SendTimeOut)
				case <-fnProxy(func() {
					sendData, err := packet.Serialize(ctx)
					if err != nil {
						c.option.Logger.Error(err)
					} else {
						select {
						case <-ctx.Done():
						default:
							_, err = c.conn.Write(sendData)
							if err != nil {
								c.option.Logger.Error(err)
							}
							cancel()
						}
					}
				}):
				}
			}
		}
	}()
	return result
}

//message 创建一个消息处理channel
func (c *Conn) message() chan<- Packet {
	result := make(chan Packet, 1)
	go func() {
		defer func() {
			c.option.Logger.Debugf("%s: message goruntinue exit", c.conn.RemoteAddr().String())
		}()
		for {
			select {
			case <-c.context.Done():
				return
			case p, ok := <-result:
				if !ok {
					c.option.Logger.Debugf("%s: Conn.Message: hand packet chan was closed", c.conn.RemoteAddr().String())
					return
				}
				select {
				case <-c.context.Done():
					return
				case <-time.After(c.option.HandTimeOut):
					c.option.Handle.OnTimeOut(c, HandTimeOut)
				case <-fnProxy(func() {
					c.option.Handle.OnMessage(c, p)
				}):
				}
			}
		}
	}()
	return result
}
