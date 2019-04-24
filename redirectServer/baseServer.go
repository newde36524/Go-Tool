package redirectServer

import "net"

//BaseServer 服务端接口
type BaseServer interface {
	Start() //启动监听端口
}

//NewRedirectServer 创建一个中转服务器
//@localAddr 本机地址
//@remoteAddr 远程地址
//备注: 本机地址不能为nil，远程地址如果为nil，则定义为主服务器，否则定义为从服务器
func NewRedirectServer(localAddr, remoteAddr *net.TCPAddr) (server BaseServer, err error) {
	if localAddr == nil {
		panic("本机地址不能为空")
	}
	if remoteAddr == nil {
		server, err = NewMasterServer(localAddr)
	} else {
		server, err = NewSlaveServer(localAddr, remoteAddr)
	}
	return
}
