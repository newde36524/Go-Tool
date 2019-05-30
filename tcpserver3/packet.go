package tcpserver3

import "context"

//Packet 协议包内容
type Packet interface {
	SetBuffer(frame []byte)                        // 设置客户端上传的数据帧
	GetBuffer() []byte                             // 获取客户端上传的数据帧
	Serialize(ctx context.Context) ([]byte, error) // 获取服务端解析后的数据帧
}

//BasePacket .
type BasePacket struct {
	Packet
	data []byte
}

//SetBuffer .
func (p *BasePacket) SetBuffer(frame []byte) {
	p.data = frame
}

//GetBuffer .
func (p *BasePacket) GetBuffer() []byte {
	return p.data
}

//Serialize .
func (p *BasePacket) Serialize(ctx context.Context) ([]byte, error) {
	return p.data, nil
}
