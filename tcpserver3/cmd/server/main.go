package main

import (
	"context"
	"fmt"
	"time"

	srv "github.com/newde36524/Go-Tool/tcpserver3"
	customer "github.com/newde36524/Go-Tool/tcpserver3/cmd/Server/customer"

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
	address := "0.0.0.0:12336"
	logger, err := srv.NewDefaultLogger()
	opt := srv.ConnOption{
		MaxSendChanCount: 100,             //最大发包数
		MaxRecvChanCount: 100,             //最大接包数
		SendTimeOut:      1 * time.Minute, //发送消息包超时时间
		RecvTimeOut:      1 * time.Minute, //接收消息包超时时间
		HandTimeOut:      1 * time.Minute, //处理消息包超时时间
		WriteDataTimeOut: 1 * time.Minute, //发送数据超时时间
		ReadDataTimeOut:  1 * time.Minute, //接收数据超时时间
		Logger:           logger,          //日志打印对象
	}
	server, err := srv.New("tcp", address, opt)
	if err != nil {
		logs.Error(err)
	}
	server.Use(customer.LogHandle{})
	server.Use(customer.RootHandle{})
	server.UseDebug()
	server.Binding()
	logs.Infof("服务器开始监听...  监听地址:%s", address)
	fmt.Scanln()
	<-context.Background().Done()
}
