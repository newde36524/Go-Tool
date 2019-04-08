package main

import (
	"fmt"

	"./arraytool"
	middle "./middleware"
)

func main() {
	TestMiddleware()
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
