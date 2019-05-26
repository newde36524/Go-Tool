package websocket

import (
	"crypto/sha1"
	"encoding/base64"
	"io/ioutil"
	"net"

	"github.com/issue9/logs"
	tcp "github.com/newde36524/Go-Tool/tcpserver2"
)

var keyGUID = []byte("258EAFA5-E914-47DA-95CA-C5AB0DC85B11")

func computeAcceptKey(challengeKey string) string {
	h := sha1.New()
	h.Write([]byte(challengeKey))
	h.Write(keyGUID)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

//WebSocketHandle .
type WebSocketHandle struct {
	handle        tcp.TCPHandle
	isFirstHandle bool
}

//ReadPacket .
func (h WebSocketHandle) ReadPacket(conn *tcp.Conn) tcp.Packet {
	//todo 定义读取数据帧的规则
	b, err := ioutil.ReadAll(conn.Raw())
	if err != nil {
		switch e := err.(type) {
		case net.Error:
			if !e.Timeout() {
				logs.Error(err)
				conn.Close()
			}
		}
	}
	p := &tcp.BasePacket{}
	p.SetBuffer(b)

	return h.handle.ReadPacket(conn)
}

//OnConnection .
func (WebSocketHandle) OnConnection(conn *tcp.Conn) {
	//todo 连接建立时处理,用于一些建立连接时,需要主动下发数据包的场景,可以在这里开启心跳协程,做登录验证等等
	logs.Infof("%s: 对方好像对你很感兴趣呦~~", conn.RemoteAddr())
}

//OnMessage .
func (h WebSocketHandle) OnMessage(conn *tcp.Conn, pkt tcp.Packet) {
	sendP := &tcp.BasePacket{}
	if h.isFirstHandle {
		h.isFirstHandle = false
		data := pkt.GetBuffer()
		challengeKey := string(data) //todo 从第一次读取的数据帧中获得http头中"Sec-Websocket-Key"的属性值
		p := make([]byte, 0)
		p = append(p, "HTTP/1.1 101 Switching Protocols\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Accept: "...)
		p = append(p, computeAcceptKey(challengeKey)...)
		p = append(p, "\r\n"...)
		// if c.subprotocol != "" {
		// 	p = append(p, "Sec-WebSocket-Protocol: "...)
		// 	p = append(p, c.subprotocol...)
		// 	p = append(p, "\r\n"...)
		// }
		// if compress {
		// 	p = append(p, "Sec-WebSocket-Extensions: permessage-deflate; server_no_context_takeover; client_no_context_takeover\r\n"...)
		// }
		// for k, vs := range responseHeader {
		// 	if k == "Sec-Websocket-Protocol" {
		// 		continue
		// 	}
		// 	for _, v := range vs {
		// 		p = append(p, k...)
		// 		p = append(p, ": "...)
		// 		for i := 0; i < len(v); i++ {
		// 			b := v[i]
		// 			if b <= 31 {
		// 				// prevent response splitting.
		// 				b = ' '
		// 			}
		// 			p = append(p, b)
		// 		}
		// 		p = append(p, "\r\n"...)
		// 	}
		// }
		p = append(p, "\r\n"...)
		sendP.SetBuffer(p)
	}
	conn.Write(sendP) //回复客户端发送的内容
}

//OnClose .
func (WebSocketHandle) OnClose(state tcp.ConnState) {
	logs.Infof("对方好像撤退了呦~~,连接状态:%s", state.String())
}

//OnTimeOut .
func (WebSocketHandle) OnTimeOut(conn *tcp.Conn, code tcp.TimeOutState) {
	logs.Infof("%s: 对方好像在做一些灰暗的事情呢~~,超时类型:%d", conn.RemoteAddr(), code)
}

//OnPanic .
func (WebSocketHandle) OnPanic(conn *tcp.Conn, err error) {
	logs.Errorf("%s: 对方好像发生了一些不得了的事情哦~~,错误信息:%s", conn.RemoteAddr(), err)
}

//OnSendError .
func (WebSocketHandle) OnSendError(conn *tcp.Conn, packet tcp.Packet, err error) {
	logs.Errorf("%s: 发送数据的时间好像有点久诶~~,错误信息:%s", conn.RemoteAddr(), err)
}

//OnRecvError .
func (WebSocketHandle) OnRecvError(conn *tcp.Conn, err error) {
	logs.Errorf("%s: 接收数据的时间好像有点久诶~~,错误信息:%s", conn.RemoteAddr(), err)
}

//firstHandShaking 第一次"握手"
func (WebSocketHandle) firstHandShaking(conn *tcp.Conn) {

}
