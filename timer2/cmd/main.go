package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/newde36524/Go-Tool/timer2"
)

func main() {
	task := timer2.New()
	// for i := 1; i < 10000; i++ {
	// 	// rd.ReadLine()
	// 	task.Add(5*time.Second, func(key string, v interface{}) (time.Time, error) {
	// 		return time.Now(), nil
	// 	}, nil, func(key string, remove func()) {
	// 		num := runtime.NumGoroutine()
	// 		fmt.Println("当前协程数:", num)
	// 		remove()
	// 	})
	// }
	rd := bufio.NewReader(os.Stdin)
	for {
		rd.ReadLine()
		task.Add(2*time.Second, func(key string, v interface{}) (time.Time, error) {
			return time.Now(), nil
		}, nil, func(key string, remove func()) {
			fmt.Println(key, time.Now())
		})
	}
}
