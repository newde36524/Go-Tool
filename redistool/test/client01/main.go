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
	client.Subscript(func(msg string) {
		fmt.Println(msg)
	}, "client01")
	for {
		<-time.After(time.Hour)
	}
}
