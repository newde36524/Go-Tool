package task

import (
	"time"
)

//Task Task结构体
type Task struct {
	funcs   chan func()
	timeOut time.Duration
}

//Run 创建一个异步执行任务
func Run(fn func()) (result *Task) {
	result = NewTask()
	result.Start(fn)
	return
}

//NewTask 创建Task的实例
func NewTask() *Task {
	return &Task{
		funcs: make(chan func(), 1024),
	}
}

//Start 开始执行任务
//@fn 需要执行的任务，任务会在协程中执行
func (t *Task) Start(fn func()) {
	go func() {
		fn()
		ch := make(chan int, 1)
		for {
			if len(t.funcs) == 0 {
				return
			}
			select {
			case f, ok := <-t.funcs:
				if !ok {
					return
				}
				select {
				case ch <- 1:
					go func() {
						f()
						<-ch
					}()
				case <-time.After(t.timeOut):
					for len(t.funcs) > 0 {
						<-t.funcs
					}
					return
				}
			default:
				return
			}
		}
	}()
}

//Continue 执行延续任务，上一个任务完成时才会执行下一个
func (t *Task) Continue(fn func()) *Task {
	t.funcs <- fn
	return t
}
