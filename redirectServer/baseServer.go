package redirectServer

import "net"

//BaseServer 服务端接口
type BaseServer interface {
	Run()           //启动监听端口
	OnMessage()     //处理消息
	OnSendPackage() //发送消息
	OnReceiv()      //接收消息

}

type ServerManage struct {
}

func (*ServerManage) CreateServer(isMaster bool) BaseServer {
	if isMaster {

	} else {

	}
	return nil
}

type Key string
type MasterServer struct {
	SlaveGroup map[Key]*net.Conn
}
type SlaveServer struct {
	MasterGroup map[Key]net.Addr
}
