package main

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/gomodule/redigo/redis"

	"./arraytool"
	"./bulkruntool"
	middle "./middleware"
	"./redistool"
	"github.com/issue9/logs"
)

func init() {
	err := logs.InitFromXMLFile("config/logs.xml")
	if err != nil {
		fmt.Println(err)
	}
}
func main() {
	// TestMiddleware()
	// TestBulkRunFuncs()
	TestRedis()
	for {
		<-time.After(24 * time.Hour)
	}
}

// func ReadFile(index, pagnum int, filePath string) {
// 	data, err := filetool.ReadPagingFile(index, pagnum, filePath)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Println(data)
// }

func TestRedis() {
	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			debug.PrintStack()
		}
	}()
	client := redistool.NewRedisClient(&redistool.RedisClientOption{
		// Password: "pCy1@nr#86z12%v",
	})
	err := client.Connect("127.0.0.1:6379")
	// err := client.Connect("10.66.178.38:6379")
	if err != nil {
		logs.Error(err)
	}
	fmt.Println("连接redis服务端")
	//=============== Set ======================
	// res, err := client.Set("a", "hello")
	// fmt.Println("Set", res, err)
	// res, err = client.Get("a")
	// fmt.Println("Get", res, err)
	//=============== PubSub ======================
	// go func() {
	// 	for {

	// 	}
	// }()
	// client.Subscript(func(msg string) {
	// 	fmt.Println(msg)
	// }, "MyTopic")
	// time.Sleep(1 * time.Second)
	topic := "MyTopic"
	c, _ := client.Clone()
	go func(a *redistool.RedisClient, b redis.Conn) {
		for {
			a.Publish(topic, "hello world")
			fmt.Println("1")
			time.Sleep(time.Second)
		}
	}(client, c)
	// fmt.Println(client.Publish(topic, "hello world"))
	client.Subscript(func(msg interface{}) {
		fmt.Println("666", msg)
	}, topic)
	<-time.After(time.Hour)

}
func TestRevertArray() {
	fmt.Println(arraytool.RevertArray([]interface{}{0x1, 0x2, 0x3}))
}

func TestBulkRunFuncs() {
	fnArr := []func(){}
	for index := 0; index < 10000; index++ {
		temp := index
		fnArr = append(fnArr, func() {
			fmt.Println(temp)
			time.Sleep(time.Second)
		})
	}
	bulkruntool.RunTask(2, fnArr)
}

func TestMiddleware() {
	app := middle.NewApplication()
	app.Use(MiddlewareA)
	app.Use(MiddlewareB)
	app.Use(MiddlewareC)
	middleware := app.Build()
	middleware(1)
	// middleware(1)
}

func MiddlewareA(middleware middle.Middleware) middle.Middleware {
	return func(o interface{}) {
		logs.Info("A1")
		middleware(o)
		logs.Info("A2")
	}
}
func MiddlewareB(middleware middle.Middleware) middle.Middleware {
	return func(o interface{}) {
		logs.Info("B1")
		middleware(o)
		logs.Info("B2")
	}
}
func MiddlewareC(middleware middle.Middleware) middle.Middleware {
	return func(o interface{}) {
		logs.Info("C1")
		middleware(o)
		logs.Info("C2")
	}
}
