package main

import (
	"fmt"
	"time"

	"github.com/issue9/logs"
	"github.com/newde36524/Go-Tool/redistool"
)

func main() {
	client := redistool.NewRedisClient("127.0.0.1:6379", &redistool.RedisClientOption{
		// Password: "123456",
	})
	err := client.Connect()
	if err != nil {
		logs.Error(err)
	}
	fmt.Println("连接redis服务端")
	go func() {
		for {
			client.Publish("client01", "hello world 001")
			client.Publish("client02", "hello world 002")
			client.Publish("client03", "hello world 003")
			<-time.After(time.Second)
		}
	}()
	for {
		<-time.After(time.Hour)
	}
}
