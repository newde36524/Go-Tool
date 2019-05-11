package main

import (
	tcp "Go-Tool/tcpserver"
	customer "Go-Tool/tcpserver/cmd/customer"
	"context"
	"fmt"
	"time"

	"github.com/issue9/logs"
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
	address := "127.0.0.1:12336"
	logger, err := tcp.NewDefaultLogger()
	server, err := tcp.NewServer(address, tcp.ConnOption{
		MaxSendChanCount: 100,
		MaxRecvChanCount: 100,                  //最大接包数
		SendTimeOut:      5 * time.Minute,      //发送消息超时时间
		RecvTimeOut:      5 * time.Minute,      //接收消息超时时间
		HandTimeOut:      5 * time.Minute,      //处理消息超时时间
		Logger:           logger,               //日志打印对象
		Handle:           customer.TCPHandle{}, //包处理对象
	})
	if err != nil {
		logs.Error(err)
	}
	server.Binding()
	logs.Infof("服务器开始监听...  监听地址:%s", address)
	<-context.Background().Done()
}
