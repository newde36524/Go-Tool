package main

import (
	"fmt"
	"time"

	"github.com/issue9/logs"

	"./arraytool"
	"./bulkruntool"
	middle "./middleware"
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

	mp := make(map[string]string, 1024)
	mp[""] = "hehe"
	fmt.Println(mp[""])

	<-time.After(24 * time.Hour)
}

// func ReadFile(index, pagnum int, filePath string) {
// 	data, err := filetool.ReadPagingFile(index, pagnum, filePath)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Println(data)
// }

// func TestRedis() {
// 	client := new(redistool.RedisClient)
// 	fmt.Println("连接redis服务端")

// 	client.Login("ip:port", &redistool.RedisClientOption{
// 		Password: "password",
// 	})

// 	res, err := client.Set("a", "hello")
// 	fmt.Println("Set", res, err)

// 	res, err = client.Get("a")
// 	fmt.Println("Get", res, err)
// }
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
