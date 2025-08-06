package main

import (
	"fmt"
	"time"

	"github.com/newde36524/Go-Tool/commontool"
)

func main() {
	locker := commontool.NewInteractionLocker(func(a, b int) bool {
		fmt.Println("嚶嚶嚶")
		return a == b
	})
	go func() {
		for i := 0; i < 10; i++ {
			locker.Left(i)
			fmt.Printf("Left1 %d \n", i)
			<-time.After(1 * time.Second)
		}
		time.After(time.Second)
		fmt.Println("========== left ==============")
		for i := 0; i < 10; i++ {
			locker.Right(i)
			fmt.Printf("Left2 %d \n", i)
			<-time.After(2 * time.Second)
		}
	}()
	go func() {
		for i := 0; i < 10; i++ {
			locker.Right(i)
			fmt.Printf("Right1 %d \n", i)
			<-time.After(2 * time.Second)
		}
		fmt.Println("========== right ==============")
		for i := 0; i < 10; i++ {
			locker.Left(i)
			fmt.Printf("Right2 %d \n", i)
			<-time.After(1 * time.Second)
		}
	}()
	<-time.After(time.Hour)
}
