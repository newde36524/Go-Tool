package main

import (
	"fmt"
	"time"
)

func main() {
	funcs := make([]func(), 0)
	for index := 0; index < 100; index++ {
		temp := index
		funcs = append(funcs, func() {
			fmt.Println(temp)
		})
	}
	ch := make(chan chan struct{}, len(funcs))
	for _, fn := range funcs {
		sign := make(chan struct{}, 0)
		ch <- sign
		go func(fn func(), s chan struct{}) {
			fn()
			<-s
			close(<-ch)
		}(fn, sign)
	}
	close(<-ch)
	<-time.After(time.Hour)
}

type Task interface {
	Run() (int, error)
}

type DownloadTask struct {
	Task
	Name string
}

func (task *DownloadTask) Run() (int, error) {
	//todo
	return 1, nil
}

type UploadTask struct {
	Task
	Name string
}

func (task *UploadTask) Run() (int, error) {
	//todo
	return 2, nil
}
