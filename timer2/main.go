package timer2

import (
	"fmt"
	"runtime"
	"strconv"
	"time"
)

func main() {
	task := NewTimerTask()
	for i := 1; i < 10000; i++ {
		// rd.ReadLine()
		task.Add(strconv.Itoa(i), func(v interface{}) time.Time {
			return time.Now()
		}, nil, 5*time.Second, func(remove func()) {
			num := runtime.NumGoroutine()
			fmt.Println("当前协程数:", num)
			remove()
		})
	}
}
