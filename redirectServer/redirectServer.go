package redirectServer

import (
	"net"

	"github.com/issue9/logs"
)

/*
	备注：发布订阅服务器本质上是一个广播服务器，同时也是个端口数据转发服务器，需要确保服务器稳定
	1: 每个客户端上传数据时必须指定 具备哪种标识的数据帧是给自己的
	2：客户端发送时 必须带上标识，以便于服务端转发
	3：服务端必须支持内网穿透，支持集群的方式相互通讯(GRPC)
	4：master slave 模式  1：无法连接外网服务器时， 内网集群，定时连接外网服务器，如果内网中断开外网服务器的服务器崩溃重启，
	那么内部集群中其他slave会连接外网


*/
//PubSuber 广播发布接收服务器
// 一般情况下 主从映射表都会找key对应的value ，外网服务器无法连接时，
type PubSuber struct {
	MasterConnMap map[string]string //集群标志和Master服务器IP地址的映射 key:IP
	SlaveConnMap  map[string]string //集群标志和Slave服务器IP地址的映射 key:IP
	SendFlag      string            //广播接收数据时带上的自身唯一发送标识
	RecvFlag      string            //广播反馈数据时带上的自身唯一反馈标识
	Key           string            //集群标志
	Conn          *net.Conn         //目标集群成员连接
	tcpListener   *net.TCPListener  //服务器监听对象
	isDebug       bool              //是否开启debug日志
	//需要一个数据帧 发送码区  反馈码区  转发服务器节点区
}

//Init x
//key 集群标志
//ip 自身地址或者目标集群成员
func (p *PubSuber) Init(key string, ip string, isMaster bool) error {
	p.Key = key
	addr, err := net.ResolveTCPAddr("tcp", ip)
	if err != nil {
		return err
	}
	if isMaster {
		p.MasterConnMap[key] = ip
	} else {
		p.SlaveConnMap[key] = ip
	}
	//todo 判断传进来的ip是否是本机地址  如果是就开启监听否则开启服务器并和ip建立tcp连接
	if true {
		p.tcpListener, err = net.ListenTCP("tcp", addr)
		//todo 获取本机的真实地址
	}
	return err
}

//Run 开启服务端监听
func (p *PubSuber) Run() {
	defer p.tcpListener.Close()
	//todo  建立连接后需要同步连接映射列表
	for {
		conn, err := p.tcpListener.AcceptTCP()
		if err != nil {
			logs.Error(err)
		}
		go p.OnMessage(conn)
		go p.OnReceiv(conn)
		go p.OnSend(conn)
	}
}

//RunAsync 异步执行Run函数
func (p *PubSuber) RunAsync() {
	go p.Run()
}

func (p *PubSuber) OnMessage(conn *net.TCPConn) {

}

func (p *PubSuber) OnReceiv(conn *net.TCPConn) {
	for {
		// conn.Read()
	}
}

func (p *PubSuber) OnSend(conn *net.TCPConn) {

}

//广播 可协程并发发送 谁先
func (*PubSuber) Publish() {

}

//广播接收  收到数据时需要等待 被告知谁收到了反馈。
func (*PubSuber) Sublish() {

}

//单点发送
func (*PubSuber) Send() {

}

//单点接收
func (*PubSuber) Recv() {

}

func (*PubSuber) AddOrUpdateLocalData() {

}
