package redirectserver

// import (
// 	"net"
// 	"strconv"

// 	"github.com/issue9/logs"
// )

// /*
// 	备注：发布订阅服务器本质上是一个广播服务器，同时也是个端口数据转发服务器，需要确保服务器稳定
// 	1: 每个客户端上传数据时必须指定 具备哪种标识的数据帧是给自己的
// 	2：客户端发送时 必须带上标识，以便于服务端转发
// 	3：服务端必须支持内网穿透，支持集群的方式相互通讯(GRPC)
// 	4：master slave 模式  1：无法连接外网服务器时， 内网集群，定时连接外网服务器，如果内网中断开外网服务器的服务器崩溃重启，
// 	那么内部集群中其他slave会连接外网

// */
// //RedirectServer 广播发布接收服务器
// // 一般情况下 主从映射表都会找key对应的value ，外网服务器无法连接时，
// //master服务器管理所有客户端连接并转发数据，slave服务器负责发送数据，主从服务器不接触任何数据帧转换处理，
// //业务服务器需要和主从服务器建立连接，从主从服务器中接收和发送数据
// //在主从服务器之下 需要提供数据帧转换收发规则，便于业务服务器之间通过主从服务器交互
// type RedirectServer struct {
// 	MasterConnMap map[Key][]*SlaveConn //集群标志和Master服务器IP地址的映射 key:IP
// 	RemoteConn    *MasterConn          //远程主服务器链接
// 	tcpListener   *net.TCPListener     //主从服务器监听对象
// 	Code          int                  // 0 表示主服务器 slaveIP失效  1  表示从服务器  masterIP失效  2  表示主从同时生效
// 	isDebug       bool                 //是否开启debug日志
// 	//需要一个数据帧 发送码区  反馈码区  转发服务器节点区
// }
// type MasterConn net.TCPConn //从服务器和主服务器建立的连接
// type SlaveConn net.TCPConn  //主服务器和从服务器建立的连接
// type Key string

// func NewRedirectServer() *RedirectServer {
// 	return &RedirectServer{}
// }

// //Init x
// //key 集群标志
// //ip 自身地址或者目标集群成员
// //isMaster true 用ip自己做监听 false 自己建立监听并连接ip
// //无论主从 都会需要自己监听自己服务器的端口
// //服务器自身可作为主服务器也可做其他主服务器的从服务器
// //code  0 表示主服务器 slaveIP失效  1  表示从服务器  masterIP失效  2  表示主从同时生效
// func (p *RedirectServer) Init(masterIP string, slaveIP string, code int) (err error) {
// 	masterAddr, err := net.ResolveTCPAddr("tcp", masterIP)
// 	if err != nil {
// 		logs.Error(err)
// 	}
// 	p.InitMasterServer(masterAddr)
// 	slaveAddr, err := net.ResolveTCPAddr("tcp", slaveIP)
// 	if err != nil {
// 		logs.Error(err)
// 	}
// 	p.InitSlaveServer(masterAddr)
// 	if err != nil {
// 		return err
// 	}
// 	//todo 判断传进来的ip是否是本机地址  如果是就开启监听否则开启服务器并和ip建立tcp连接
// 	if isMaster {
// 		p.tcpListener, err = net.ListenTCP("tcp", addr)
// 		//todo 获取本机的真实地址
// 	} else {
// 		addr, err := net.ResolveTCPAddr("tcp", ":"+strconv.Itoa(port))
// 		if err != nil {
// 			return err
// 		}
// 		p.tcpListener, err = net.ListenTCP("tcp", addr)
// 		//todo  如果是slave服务器就和master服务器建立长连接

// 	}
// 	return err
// }

// func (p *RedirectServer) InitMasterServer(masterAddr *net.TCPAddr) (conns chan *SlaveConn, err error) {
// 	//todo master服务器监听自身的端口
// 	p.tcpListener, err = net.ListenTCP("tcp", masterAddr)
// 	if err != nil {
// 		return
// 	}
// 	conns = make(chan *SlaveConn, 1024)
// 	go func(connCh chan *SlaveConn, tcpListen *net.TCPListener) {
// 		for {
// 			conn, err := tcpListen.AcceptTCP()
// 			if err != nil {
// 				logs.Error(err)
// 			}
// 			connCh <- conn //建立的链接全部扔进队列
// 		}
// 	}(conns, p.tcpListener)

// 	return
// }

// func (p *RedirectServer) ConnHandle(conns chan *net.TCPConn) {

// 	go func(conns chan *net.TCPConn) {
// 		for {
// 			select {
// 			case conn, ok := <-conns:
// 				if ok {
// 					conn.Read()
// 				}
// 			}
// 		}
// 	}(conns)
// }

// func (p *RedirectServer) InitSlaveServer(slaveAddr *net.TCPAddr, masterAddr *net.TCPAddr) (err error) {
// 	//todo  slave服务器需要自己建立服务端链接 并监听自身的端口
// 	p.tcpListener, err = net.ListenTCP("tcp", slaveAddr)
// 	if err != nil {
// 		return
// 	}
// 	return
// }

// //Run 开启服务端监听
// func (p *RedirectServer) Run() {
// 	defer p.tcpListener.Close()
// 	//todo  建立连接后需要同步连接映射列表
// 	for {
// 		conn, err := p.tcpListener.AcceptTCP()
// 		if err != nil {
// 			logs.Error(err)
// 		} //建立连接以后需要获取一个key 之后才能做群组映射
// 		go p.OnReceiv(conn)
// 	}
// }

// //RunAsync 异步执行Run函数
// func (p *RedirectServer) RunAsync() {
// 	go p.Run()
// }

// func (p *RedirectServer) OnReceiv(conn *net.TCPConn) {
// 	for { //如果主服务器 要给所有的客户端发送数据
// 		buffer := make([]byte, 1024)
// 		conn.Read(buffer)
// 		// for _, clientConn := range p. {

// 		// }

// 	}
// }
