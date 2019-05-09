package tcpserver

import (
	"net"
	"time"
)

//Conn 连接代理对象
type Conn struct {
	conn     net.Conn      //tcp连接对象
	option   ConnOption    //连接配置项
	state    ConnState     //连接状态
	recvChan <-chan Packet //接收消息队列
	sendChan chan<- Packet //发送消息队列
	handChan chan<- Packet //处理消息队列
}

// NewConn returns a wrapper of raw conn
func NewConn(conn net.Conn, option ConnOption) *Conn {
	return &Conn{
		conn:   conn,
		option: option,
		state: ConnState{
			IsExit:     make(chan struct{}),
			ActiveTime: time.Now(),
		},
	}
}

//Run 固定处理流程
func (c *Conn) Run() {
	c.recvChan = c.recv(c.option.MaxRecvChanCount)
	c.sendChan = c.send(c.option.MaxSendChanCount)
	c.handChan = c.message()
	go func() {
		p, ok := <-c.recvChan
		if !ok {
			c.option.Logger.Error("Conn.Run: recvChan is closed")
		}
		c.option.Handle.OnFirst(c, p)
		for {
			p, ok := <-c.recvChan
			if !ok {
				c.option.Logger.Error("Conn.Run: recvChan is closed")
			}
			c.handChan <- p
		}
	}()
}

//Send 发送消息到设备
func (c *Conn) Send(packet Packet) {
	select {
	case <-c.state.IsExit:
		close(c.sendChan)
	case c.sendChan <- packet: //消息入列
	}
}

// Close 关闭服务器和设备的连接
func (c *Conn) Close() {
	defer c.conn.Close()
	c.option.Handle.OnClose() //调用接口连接关闭时的处理函数
	c.state.Message = "连接退出"
	close(c.state.IsExit)
	// runtime.GC()         //强制GC      待定可能有问题
	// debug.FreeOSMemory() //强制释放内存 待定可能有问题
}

//ReadPacket 读取一个包
func (c *Conn) readPacket() Packet {
	p, err := c.option.Handle.ReadPacket(c)
	if err != nil {
		c.option.Logger.Error(err)
	}
	return p
}

//recv 创建一个包接收channel
func (c *Conn) recv(maxRecvChanCount int) <-chan Packet {
	result := make(chan Packet, maxRecvChanCount)
	go func() {
		for {
			select {
			case <-c.state.IsExit:
				close(result)
			case <-time.After(c.option.RecvTimeOut):
				c.option.Logger.Info("Conn.Recv:recvChan locked used time was too long ")
				close(result)
				return
			case result <- c.readPacket():
			}
		}
	}()
	return result
}

//send 创建一个包发送channel
func (c *Conn) send(maxSendChanCount int) chan<- Packet {
	result := make(chan Packet, maxSendChanCount)
	go func() {
		for {
			select {
			case <-time.After(c.option.SendTimeOut):
			case packet, ok := <-result:
				if !ok {
					c.option.Logger.Info("Conn.Send:send chan is closed")

					return
				}
				sendData, err := packet.Serialize()
				if err != nil {
					c.option.Logger.Error(err)
				}
				c.conn.Write(sendData)
			}
		}
	}()
	return result
}

//message 创建一个消息处理channel
func (c *Conn) message() chan<- Packet {
	result := make(chan Packet, 1)
	go func() {
		for {
			select {
			case <-c.state.IsExit:
				return
			case p, ok := <-result:
				if !ok {
					c.option.Logger.Info("Conn.Message: hand packet chan was closed")
					return
				}
				ch := make(chan struct{})
				go func() {
					defer close(ch)
					c.option.Handle.OnMessage(c, p)
					ch <- struct{}{}
				}()
				select {
				case <-ch:
				case <-time.After(c.option.HandTimeOut):
					c.option.Logger.Info("Conn.Message: hand packet used time was too long")
				}
			}
		}
	}()
	return result
}
