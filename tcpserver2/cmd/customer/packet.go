package customer

import (
	tcp "Go-Tool/tcpserver2"
	"context"
)

//Packet .
type Packet struct {
	tcp.Packet
	data []byte
}

//SetBuffer .
func (p *Packet) SetBuffer(frame []byte) {
	//todo 解析数据包，并可根据需要在结构中定义多个字段存储
	p.data = frame
}

//GetBuffer .
func (p *Packet) GetBuffer() []byte {
	//todo 解析数据包，并可根据需要在结构中定义多个字段存储
	return p.data
}

//Serialize .
func (p *Packet) Serialize(ctx context.Context) ([]byte, error) {
	//todo 数据帧的业务处理逻辑
	return p.data, nil
}
