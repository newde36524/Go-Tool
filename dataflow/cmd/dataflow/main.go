package main

import (
	"fmt"
	"time"

	"github.com/newde36524/Go-Tool/dataflow"
)

func main() {
	done := make(chan interface{})
	rand := func() interface{} {
		return 1
	}
	for v := range dataflow.Take(done, dataflow.RepeatFunc(done, rand), 10) {
		fmt.Println(v)
	}
	close(done)
	fmt.Println("end")
	<-time.After(time.Hour)
}
