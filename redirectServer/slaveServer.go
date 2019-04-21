package redirectServer

import (
	"net"

	"github.com/issue9/logs"
)

/**
1: 开启端口监听
2：连接远程主服务器
3：开始数据接收循环
4：发送数据给远程服务器
5：设置和获取分组唯一值
6：

*/

//SlaveServer 从转发服务器
type SlaveServer struct {
	MasterConn       *net.TCPConn //主服务器连接
	SlaveTCPListener *net.TCPListener
	ConnMap          map[string]*net.TCPConn
}

//NewSlaveServer 实例化一个SlaveServer
//@localAddr 本地监听地址
//@remoteAddr 远程服务器地址
func NewSlaveServer(localAddr *net.TCPAddr, remoteAddr *net.TCPAddr) (*SlaveServer, error) {
	server := new(SlaveServer)
	conn, err := net.DialTCP("tcp", nil, remoteAddr)
	if err != nil {
		return nil, err
	}
	server.MasterConn = conn
	tcpListener, err := net.ListenTCP("tcp", localAddr)
	if err != nil {
		return nil, err
	}
	server.SlaveTCPListener = tcpListener
	return server, err
}

//Start 开始运行服务端
func (srv *SlaveServer) Start() {
	go srv.onConnection()
}

//onConnection 连接处理程序
func (srv *SlaveServer) onConnection() {
	go srv.onReceiv(srv.MasterConn)
	for {
		conn, err := srv.SlaveTCPListener.AcceptTCP()
		if err != nil {
			logs.Error(err)
		} else {
			srv.ConnMap[conn.RemoteAddr().String()] = conn
			go srv.onSend(conn)
		}
	}
}

//onSend 业务服务器上发数据时处理
func (srv *SlaveServer) onSend(conn *net.TCPConn) { //业务服务器连接
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer) //读取业务服务器上发的数据并转发给主服务器
		if err != nil {
			logs.Errorf("%s:%s", conn.RemoteAddr().String(), err)
			delete(srv.ConnMap, conn.RemoteAddr().String())
			break
		}
		srv.MasterConn.Write(buffer[:n])
	}
	logs.Info("onSend exiting")
}

//onReceiv 远程服务器下发数据时处理
func (srv *SlaveServer) onReceiv(conn *net.TCPConn) { //远程服务器连接
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			logs.Errorf("%s:%s", conn.RemoteAddr().String(), err)
			break
		}
		for _, c := range srv.ConnMap { //读取远程服务器发送的数据并转发给所有业务服务器
			c.Write(buffer[:n])
		}
	}
	logs.Info("onReceiv exiting")
}
