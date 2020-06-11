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

	// for i := 0; i < 10000; i++ {
	// 	_, _ = task.Add(20*time.Second, func(key string, v interface{}) (time.Time, error) {
	// 		return time.Now(), nil
	// 	}, nil, func(key string, remove func()) {
	// 		fmt.Println(key, time.Now())
	// 		remove()
	// 	})
	// }
	// fmt.Println("end")

	key := "hehe"
	task.Sync(key, 1*time.Second, func(key string, v interface{}) (time.Time, error) {
		return time.Now(), nil
	}, nil, func(key string, remove func()) {
		fmt.Println(key, time.Now())
	})

	for {
		rd.ReadLine()
		task.Modify(key, func(e *timer2.Entity) error {
			fmt.Println(e.Delay)
			e.Delay = e.Delay + 1*time.Second
			fmt.Println(e.Delay)
			return nil
		})
	}
}
