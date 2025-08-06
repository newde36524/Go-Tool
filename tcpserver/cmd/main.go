package main

import (
	"context"
	"fmt"
	"time"

	"github.com/issue9/logs"
	tcp "github.com/newde36524/Go-Tool/tcpserver"
	customer "github.com/newde36524/Go-Tool/tcpserver/cmd/customer"
)

func init() {
	err := logs.InitFromXMLFile("./logs.xml")
	if err != nil {
		fmt.Println(err)
		<-time.After(10 * time.Second)
		return
	}
}
func main() {
	address := "0.0.0.0:12336"
	logger, err := tcp.NewDefaultLogger()
	server, err := tcp.New(address, tcp.ConnOption{
		MaxSendChanCount: 100,
		MaxRecvChanCount: 100,                  //最大接包数
		SendTimeOut:      1 * time.Minute,      //发送消息超时时间
		RecvTimeOut:      1 * time.Minute,      //接收消息超时时间
		HandTimeOut:      1 * time.Minute,      //处理消息超时时间
		Logger:           logger,               //日志打印对象
		Handle:           customer.TCPHandle{}, //包处理对象
	})
	if err != nil {
		logs.Error(err)
	}
	server.Binding()
	logs.Infof("服务器开始监听...  监听地址:%s", address)
	fmt.Scanln()
	<-context.Background().Done()
}
