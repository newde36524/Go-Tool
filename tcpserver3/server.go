package tcpserver3

import (
	"net"
	"runtime/debug"
	"time"
)

//Server tcp服务器
type Server struct {
	tcpListener *net.TCPListener //TCP监听对象
	connOption  ConnOption       //连接配置项
	pipe        *CoreTCPHandle   //连接处理管道
}

//New new server
//@addr local address
//@connOption connection options
func New(addr string, connOption ConnOption) (srv *Server, err error) {
	// 根据服务器开启多CPU功能
	// runtime.GOMAXPROCS(runtime.NumCPU())
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return
	}
	srv = &Server{
		tcpListener: listener,
		connOption:  connOption,
	}
	srv.Use(&DefaultTCPHandle{})
	return
}

//Use middleware
func (s *Server) Use(h TCPHandle) {
	tree := NewCoreTCPHandle(h)
	if s.pipe != nil {
		s.pipe = s.pipe.Link(tree)
	} else {
		s.pipe = tree
	}
}

//Binding start server
func (s *Server) Binding() {
	go func() {
		defer s.tcpListener.Close()
		defer func() {
			defer recover()
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
			c := NewConn(conn, s.connOption, First(s.pipe))
			c.UseDebug()
			c.run()
		}
	}()
}
