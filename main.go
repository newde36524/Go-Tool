package main

import (
	"fmt"
	"runtime/debug"
	"time"

	"./arraytool"
	"./bulkruntool"
	middle "./middleware"
	"./redistool"
	"./task"
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
	TestBulkRunFuncs2()
	// TestRedis()
	// TestMiddleware2()
	// TestMiddleware3()
	// TestTask()

	for {
		<-time.After(time.Hour)
	}
}

//TestRedis .
func TestRedis() {
	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			debug.PrintStack()
		}
	}()
	client := redistool.NewRedisClient(&redistool.RedisClientOption{
		// Password: "123456",
	})
	err := client.Connect("127.0.0.1:6379")
	if err != nil {
		logs.Error(err)
	}
	fmt.Println("连接redis服务端")
	//=============== Set ======================
	res, err := client.Set("a", "hello")
	fmt.Println("Set", res, err)
	res, err = client.Get("a")
	fmt.Println("Get", res, err)
	//=============== PubSub ======================
	client.Publish("MyTopic", "hello world")
	go func() {
		for {
			client.Publish("MyTopic", "hello world")
			<-time.After(time.Second)
		}
	}()
	client.Subscript(func(msg string) {
		fmt.Println(msg)
	}, "MyTopic")

}

//TestRevertArray .
func TestRevertArray() {
	fmt.Println(arraytool.RevertArray([]interface{}{0x1, 0x2, 0x3}))
}

//TestBulkRunFuncs .
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

//TestBulkRunFuncs2 .
func TestBulkRunFuncs2() {
	ch := bulkruntool.CreateBulkRunFuncChannel(10, 10000)
	for index := 0; index < 10000; index++ {
		temp := index
		ch <- func() {
			fmt.Println(temp)
			time.Sleep(time.Second)
		}
	}
}

//TestMiddleware .
func TestMiddleware() {
	app := middle.NewApplication()
	app.Use(MiddlewareA)
	app.Use(MiddlewareB)
	app.Use(MiddlewareC)
	middleware := app.Build()
	middleware(1)
	// middleware(1)
}

//MiddlewareA .
func MiddlewareA(middleware middle.Middleware) middle.Middleware {
	return func(o interface{}) {
		logs.Info("A1")
		middleware(o)
		logs.Info("A2")
	}
}

//MiddlewareB .
func MiddlewareB(middleware middle.Middleware) middle.Middleware {
	return func(o interface{}) {
		logs.Info("B1")
		middleware(o)
		logs.Info("B2")
	}
}

//MiddlewareC .
func MiddlewareC(middleware middle.Middleware) middle.Middleware {
	return func(o interface{}) {
		logs.Info("C1")
		middleware(o)
		logs.Info("C2")
	}
}

//TestMiddleware2 .
func TestMiddleware2() {
	middleware := new(middle.Middleware2)
	middleware.Use(func(next func()) {
		fmt.Println("A1")
		next()
		fmt.Println("A2")
	})
	middleware.Use(func(next func()) {
		fmt.Println("B1")
		next()
		fmt.Println("B2")
	})
	middleware.Use(func(next func()) {
		fmt.Println("C1")
		next()
		fmt.Println("C2")
	})
	middleware.Invoke()
}

//TestMiddleware3 .
func TestMiddleware3() {
	middleware := new(middle.Middleware3)
	middleware.Use(func(o interface{}, next func()) {
		fmt.Println("A1")
		fmt.Println(o)
		next()
		fmt.Println("A2")
	})
	middleware.Use(func(o interface{}, next func()) {
		fmt.Println("B1")
		fmt.Println(o)
		next()
		fmt.Println("B2")
	})
	middleware.Use(func(o interface{}, next func()) {
		fmt.Println("C1")
		fmt.Println(o)
		next()
		fmt.Println("C2")
	})
	middleware.Invoke(111)
}

//TestTask .
func TestTask() {
	task.Run(func() {
		fmt.Println("1")
	}).Continue(func() {
		fmt.Println("2")
	}).Continue(func() {
		fmt.Println("3")
	}).Continue(func() {
		fmt.Println("4")
	})

	task := task.Run(func() {
		fmt.Println("666")
	})
	for index := 0; index < 100; index++ {
		temp := index
		task = task.Continue(func() {
			fmt.Println(temp)
		})
	}

}
