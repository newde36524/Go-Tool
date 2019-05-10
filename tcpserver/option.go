package tcpserver

import "time"

//ConnOption 连接配置项
type ConnOption struct {
	MaxSendChanCount int           //最大发包数
	MaxRecvChanCount int           //最大接包数
	SendTimeOut      time.Duration //发送消息超时时间
	RecvTimeOut      time.Duration //接收消息超时时间
	HandTimeOut      time.Duration //处理消息超时时间
	Logger           Logger        //日志打印对象
	Handle           TCPHandle     //包处理对象
}
