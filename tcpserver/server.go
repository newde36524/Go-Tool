package tcpserver

import (
	"net"
	"runtime/debug"
	"time"
)

//Server tcp服务器
type Server struct {
	tcpListener *net.TCPListener //TCP监听对象
	connOption  ConnOption
}

//NewServer 新服务
//@addr 服务器监听地址
//@connOption 客户端连接配置项
func NewServer(addr string, connOption ConnOption) (*Server, error) {
	// 根据服务器开启多CPU功能
	// runtime.GOMAXPROCS(runtime.NumCPU())
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return nil, err
	}
	return &Server{listener, connOption}, nil
}

//Binding 启动tcp服务器
func (s *Server) Binding() {
	defer s.tcpListener.Close()
	s.connOption.Logger.Info("服务器开始监听...")
	go func() {
		defer func() {
			if err := recover(); err != nil {
				s.connOption.Logger.Error(err)
				s.connOption.Logger.Error(debug.Stack())
			}
		}()
		for {
			conn, err := s.tcpListener.AcceptTCP()
			if err != nil {
				s.connOption.Logger.Error(err)
				<-time.After(time.Second)
				continue
			}
			c := NewConn(conn, s.connOption)
			c.run()
		}
	}()
}
