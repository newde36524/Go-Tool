package redirectServer

import (
	"net"

	"github.com/issue9/logs"
)

// type MasterServer struct {
// 	SlaveGroup map[Key][]*net.Conn //从服务器列表
// }

type MasterServer struct {
	MasterServerTCPListener *net.TCPListener
	ConnList                []*net.TCPConn
}

//NewMasterServer 实例化一个MasterServer
//@localAddr 本地监听地址
//@remoteAddr 远程服务器地址
func NewMasterServer(localAddr *net.TCPAddr) (*MasterServer, error) {
	server := new(MasterServer)
	tcpListener, err := net.ListenTCP("tcp", localAddr)
	if err != nil {
		return nil, err
	}
	server.MasterServerTCPListener = tcpListener
	return server, err
}

//Start 开始运行服务端
func (srv *MasterServer) Start() {
	go srv.OnConnection()
}

//OnConnection 连接处理程序
func (srv *MasterServer) OnConnection() {
	for {
		conn, err := srv.MasterServerTCPListener.AcceptTCP()
		if err != nil {
			logs.Error(err)
		} else {
			srv.ConnList = append(srv.ConnList, conn)
			go srv.OnSend(conn)
		}
	}
}

//OnSend 业务服务器上发数据时处理
func (srv *MasterServer) OnSend(conn *net.TCPConn) { //业务服务器连接
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer) //读取业务服务器上发的数据并转发给其他业务服务器
		if err != nil {
			logs.Error(err)
		}
		for _, c := range srv.ConnList {
			if c.RemoteAddr().String() != conn.RemoteAddr().String() {
				c.Write(buffer[:n])
			}
		}
	}
}
