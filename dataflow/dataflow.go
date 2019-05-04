package dataflow

import "fmt"

// Take 从通道中获取指定数量的数据
func Take(done <-chan interface{}, valueStream <-chan interface{}, num int) <-chan interface{} {
	takeStream := make(chan interface{})
	go func() {
		defer close(takeStream)
		for i := 0; i < num; i++ {
			select {
			case <-done:
				fmt.Println("Take done")
				return
			case takeStream <- <-valueStream:
			}
		}
	}()
	return takeStream
}

//Repeat 重复生成数据流
func Repeat(done <-chan interface{}, values ...interface{}) <-chan interface{} {
	repeatStream := make(chan interface{})
	go func() {
		defer close(repeatStream)
		for {
			for _, v := range values {
				select {
				case <-done:
					fmt.Println("Repeat done")
					return
				case repeatStream <- v:
				}
			}
		}
	}()
	return repeatStream
}

//RepeatFunc 通过方法重复生成数据流
func RepeatFunc(done <-chan interface{}, fn func() interface{}) <-chan interface{} {
	repeatStream := make(chan interface{})
	go func() {
		defer close(repeatStream)
		for {
			select {
			case <-done:
				fmt.Println("RepeatFunc done")
				return
			case repeatStream <- fn():
			}
		}
	}()
	return repeatStream
}
