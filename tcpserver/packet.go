package tcpserver

// Packet 协议包内容
type Packet interface {
	SetBuffer(frame []byte)     // 设置客户端上传的数据帧
	Serialize() ([]byte, error) // 获取服务端解析后的数据帧
}
