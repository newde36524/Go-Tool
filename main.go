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
	{
		// Pub(client.C, "MyTopic", "hello world")
		// c2, _ := client.Clone()
		// Sub(c2, "MyTopic")
	}
	{
		go func() {
			for {
				client.Publish("MyTopic", "hello world")
				<-time.After(time.Second)
			}
		}()
		c2, _ := client.Clone()
		Sub(c2, "MyTopic")
		// client.Subscript(func(msg string) {
		// 	fmt.Println(msg)
		// }, "MyTopic")
	}
	<-time.After(time.Hour)
}

func Pub(c redis.Conn, topic string, msg string) {
	go func(conn redis.Conn) {
		for {
			conn.Do("PUBLISH", topic, msg)
			conn.Flush()
			<-time.After(time.Second)
		}
	}(c)
}

func Sub(c redis.Conn, topic string) {
	psc := redis.PubSubConn{Conn: c}
	psc.Subscribe(topic)

	go func(psc *redis.PubSubConn) {
		for {
			switch v := psc.Receive().(type) {
			case redis.Message:
				fmt.Println(string(v.Data))
			case redis.Subscription:
				fmt.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
			case error:
				fmt.Println(v)
				psc.Close()
			}
		}
	}(&psc)
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
