package main

import (
	"fmt"
	"time"

	"./BulkRunTool"
	"./arraytool"
	middle "./middleware"
)

func main() {
	TestMiddleware()
	TestBulkRunFuncs()
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
	fnArr := []func(){
		func() {
			fmt.Println("1")
			time.Sleep(1 * time.Second)
		},
		func() {
			fmt.Println("2")
			time.Sleep(1 * time.Second)
		},
		func() {
			fmt.Println("3")
			time.Sleep(1 * time.Second)
		},
		func() {
			fmt.Println("4")
			time.Sleep(1 * time.Second)
		},
		func() {
			fmt.Println("5")
			time.Sleep(1 * time.Second)
		},
		func() {
			fmt.Println("6")
			time.Sleep(1 * time.Second)
		},
		func() {
			fmt.Println("7")
			time.Sleep(1 * time.Second)
		},
		func() {
			fmt.Println("8")
			time.Sleep(1 * time.Second)
		},
		func() {
			fmt.Println("9")
			time.Sleep(1 * time.Second)
		},
	}
	BulkRunTool.RunTask(3, fnArr)
}

func TestMiddleware() {
	app := middle.NewApplication()
	app.Use(MiddlewareA)
	app.Use(MiddlewareB)
	app.Use(MiddlewareC)
	app.Build()(1)
}

func MiddlewareA(middleware middle.Middleware) middle.Middleware {
	return func(o interface{}) {
		fmt.Println("A")
		middleware(o)
	}
}
func MiddlewareB(middleware middle.Middleware) middle.Middleware {
	return func(o interface{}) {
		fmt.Println("B")
		middleware(o)
	}
}
func MiddlewareC(middleware middle.Middleware) middle.Middleware {
	return func(o interface{}) {
		fmt.Println("C")
		middleware(o)
	}
}
