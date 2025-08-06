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
	client.Subscript(func(msg string) {
		fmt.Println(msg)
	}, "client03")
	for {
		<-time.After(time.Hour)
	}
}
