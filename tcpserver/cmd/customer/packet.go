package customer

import (
	tcp "Go-Tool/tcpserver"
	"context"
)

//CustomerPacket .
type CustomerPacket struct {
	tcp.Packet
	data []byte
}

//SetBuffer .
func (p CustomerPacket) SetBuffer(frame []byte) {
	//todo 解析数据包，并可根据需要在结构中定义多个字段存储
	p.data = frame
}

//GetBuffer .
func (p CustomerPacket) GetBuffer() []byte {
	//todo 解析数据包，并可根据需要在结构中定义多个字段存储
	return p.data
}

//Serialize .
func (p CustomerPacket) Serialize(ctx context.Context) ([]byte, error) {
	//todo 数据帧的业务处理逻辑
	return p.data, nil
}
