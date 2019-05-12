package main

import (
	"Go-Tool/redistool"
	"fmt"
	"time"

	"github.com/issue9/logs"
)

func main() {
	client := redistool.NewRedisClient(&redistool.RedisClientOption{
		// Password: "123456",
	})
	err := client.Connect("127.0.0.1:6379")
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
