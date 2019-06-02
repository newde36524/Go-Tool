package tcpserver3

import (
	"net"
	"runtime/debug"
	"time"
)

//Server tcp服务器
type Server struct {
	listener   net.Listener //TCP监听对象
	connOption ConnOption   //连接配置项
	pipe       *CoreHandle  //连接处理管道
	isDebug    bool         //是否开始debug日志
}

//New new server
//@network network 类型，具体参照ListenUDP ListenTCP等
//@addr local address
//@connOption connection options
func New(network, addr string, connOption ConnOption) (srv *Server, err error) {
	// 根据服务器开启多CPU功能
	// runtime.GOMAXPROCS(runtime.NumCPU())
	listener, err := net.Listen(network, addr)
	if err != nil {
		return
	}
	srv = &Server{
		listener:   listener,
		connOption: connOption,
	}
	// srv.Use(&DefaultTCPHandle{}) //默认占用第一个管道并调用下一个管道
	return
}

//Use middleware
func (s *Server) Use(h Handle) {
	tree := NewCoreHandle(h)
	if s.pipe != nil {
		s.pipe = s.pipe.Link(tree)
	} else {
		s.pipe = tree
	}
}

//UseDebug 开启debug日志
func (s *Server) UseDebug() {
	s.isDebug = true
}

//Binding start server
func (s *Server) Binding() {
	go func() {
		defer s.listener.Close()
		defer func() {
			defer recover()
			if err := recover(); err != nil {
				s.connOption.Logger.Error(err)
				s.connOption.Logger.Error(debug.Stack())
			}
		}()
		for {
			conn, err := s.listener.Accept()
			if err != nil {
				s.connOption.Logger.Error(err)
				<-time.After(time.Second)
				continue
			}
			c := NewConn(conn, s.connOption, First(s.pipe))
			if s.isDebug {
				c.UseDebug()
			}
			c.run()
		}
	}()
}
