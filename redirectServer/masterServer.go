package redirectServer

import "net"

type MasterServer struct {
	SlaveGroup map[Key][]*net.Conn //从服务器列表
}
