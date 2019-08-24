package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/issue9/logs"
	middle "github.com/newde36524/Go-Tool/Middleware"
	"github.com/newde36524/Go-Tool/arraytool"
	"github.com/newde36524/Go-Tool/bulkruntool"
	"github.com/newde36524/Go-Tool/cryptotool"
	"github.com/newde36524/Go-Tool/filetool"
	"github.com/newde36524/Go-Tool/redistool"
	"github.com/newde36524/Go-Tool/task"
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
	// TestBulkRunFuncs2()
	// TestRedis()
	// TestMiddleware2()
	// TestMiddleware3()
	// TestMiddleware4()
	// TestMiddleware5()
	// TestTask()
	// err := fmt.Errorf("测试异常信息")
	// var err2 error
	// fmt.Printf("%s  %s  %#v", err2, err.Error(), err)
	// fmt.Scanln()
	// txtData, _ := ioutil.ReadFile("test.txt")
	// fmt.Println(strings.Split(string(txtData), "\r\n"))
	// TestReadLines()
	// TestCer()
	for index := 0; index < 20; index++ {
		time.Sleep(1 * time.Second)
		fmt.Println(index)
	}
	fmt.Println("============================")
	TestRunTaskAndAscCallBack()
	fmt.Println("============================")
	// TestRunTaskAndAscCallBack2()
	// TestCreateBulkRunFuncChannelAscCallBack()
	// TestReadPagingBuffer()
	<-time.After(time.Hour)
}

//TestRedis .
func TestRedis() {
	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			debug.PrintStack()
		}
	}()
	client := redistool.NewRedisClient("127.0.0.1:6379", &redistool.RedisClientOption{
		// Password: "123456",
	})
	err := client.Connect()
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
	ch := bulkruntool.CreateBulkRunFuncChannel(10, 10000, nil)
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

//TestMiddleware4 .
func TestMiddleware4() {
	middle.Do([]interface{}{1, 2, 3, 4}, func(o interface{}, next func()) {
		fmt.Println("q")
		fmt.Println(o)
		next()
		fmt.Println("w")
	})
}

//TestMiddleware5 .
func TestMiddleware5() {
	middleware := new(middle.Middleware5)
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
	t := task.Run(func() {
		fmt.Println("666")
	})
	for index := 0; index < 100; index++ {
		temp := index
		t = t.Continue(func() {
			fmt.Println(temp)
		})
	}
}

func TestReadLines() {
	lines := filetool.ReadLines(context.Background(), "test.txt")
	for line := range lines {
		fmt.Println(line)
	}
}

func TestReadPagingBuffer() {
	data := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	for i := 0; i < len(data); i++ {
		reader := bytes.NewReader(data)
		bs, n, e := filetool.ReadPagingBuffer(i, 3, 0, reader)
		if e != nil {
			break
		}
		fmt.Println(bs[:n], n, e)
	}
}

func TestCer() {

	// cryptotool.GenRsaKey(128)

	privateKey, err := rsa.GenerateKey(rand.Reader, 128)
	logs.Error(err)
	filePath := "./TestCer.cer"
	cryptotool.CreateX509Cer(filePath, privateKey, time.Now(), 999999999*time.Second, "测试证书", []string{"zsk"}, []string{"localhost"}, []byte{1, 2, 3, 4})

}

func TestRunTaskAndAscCallBack() {
	funcs := make([]func() interface{}, 0)
	for index := 0; index < 20; index++ {
		temp := index
		funcs = append(funcs, func() interface{} {
			time.Sleep(1 * time.Second)
			return temp
		})
	}
	bulkruntool.RunTaskAndAscCallBack(10, funcs, func(i interface{}) {
		fmt.Println(i)
	})
}

func TestRunTaskAndAscCallBack2() {
	funcs := make(chan func() interface{}, 10)
	go func() {
		for index := 0; ; index++ {
			temp := index
			funcs <- func() interface{} {
				return temp
			}
		}
	}()
	time.Sleep(1 * time.Second)
	bulkruntool.RunTaskAndAscCallBack2(10, funcs, func(i interface{}) {
		fmt.Println(i)
	})
}

func TestCreateBulkRunFuncChannelAscCallBack() {
	ch := bulkruntool.CreateBulkRunFuncChannelAscCallBack(10, 10, nil, func(i interface{}) {
		fmt.Println(i)
	})
	for index := 0; index < 10000; index++ {
		temp := index
		ch <- func() interface{} {
			time.Sleep(1 * time.Second)
			return temp
		}
	}
}
