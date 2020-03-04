package bulkruntool

import (
	"sync"
	"time"
)

//RunTaskAndAscCallBack 启动指定数量的协程执行多个方法,并按顺序回调
func RunTaskAndAscCallBack(maxTaskCount int, funcs []func() interface{}, callback func(interface{})) {
	ch := make(chan struct{}, maxTaskCount)
	defer close(ch)
	ch2 := make(chan chan struct{}, maxTaskCount)
	once := sync.Once{}
	for _, fn := range funcs {
		ch <- struct{}{}
		sign := make(chan struct{})
		ch2 <- sign
		once.Do(func() {
			close(<-ch2)
		})
		go func(fn func() interface{}, sign chan struct{}) {
			result := fn()
			<-sign
			callback(result)
			close(<-ch2)
			_, ok := <-ch
			if !ok && len(ch2) == 0 {
				close(ch2)
			}
		}(fn, sign)
	}
}

//RunTaskAndAscCallBack2 启动指定数量的协程执行多个方法,并按顺序回调
func RunTaskAndAscCallBack2(maxTaskCount int, funcs <-chan func() interface{}, callback func(interface{})) {
	ch := make(chan struct{}, maxTaskCount)
	defer close(ch)
	ch2 := make(chan chan struct{}, maxTaskCount)
	once := sync.Once{}
	for len(funcs) > 0 {
		fn := <-funcs
		ch <- struct{}{}
		sign := make(chan struct{})
		ch2 <- sign
		once.Do(func() {
			close(<-ch2)
		})
		go func(fn func() interface{}, sign chan struct{}) {
			result := fn()
			<-sign
			callback(result)
			close(<-ch2)
			_, ok := <-ch
			if !ok && len(ch2) == 0 {
				close(ch2)
			}
		}(fn, sign)
	}
}

//RunTask 运行指定数量的协程执行多个方法
func RunTask(maxTaskCount int, funcs []func()) {
	ch := make(chan struct{}, maxTaskCount)
	defer close(ch)
	for _, fn := range funcs {
		ch <- struct{}{}
		go func(f func()) {
			f()
			<-ch
		}(fn)
	}
}

//RunTask2 运行指定数量的协程执行多个方法
func RunTask2(maxTaskCount int, funcs <-chan func()) {
	ch := make(chan struct{}, maxTaskCount)
	defer close(ch)
	for len(funcs) > 0 {
		fn := <-funcs
		ch <- struct{}{}
		go func(f func()) {
			f()
			<-ch
		}(fn)
	}
}

//CreateBulkRunFuncChannel 创建一个指定并发数量处理,只允许发送的方法通道
func CreateBulkRunFuncChannel(maxTaskCount, maxFuncCount int, done <-chan struct{}) chan<- func() {
	funcs := make(chan func(), maxFuncCount)
	go func(funcs chan func(), maxTaskCount int) {
		ch := make(chan struct{}, maxTaskCount)
		defer close(funcs)
		defer close(ch)
		for {
			select {
			case fn, ok := <-funcs:
				if !ok {
					return
				}
				ch <- struct{}{}
				go func(f func()) {
					f()
					<-ch
				}(fn)
			case <-done:
				return
			}
		}
	}(funcs, maxTaskCount)
	return funcs
}

func CreateBulkRunFuncChannelAscCallBack(maxTaskCount, maxFuncCount int, done <-chan struct{}, callBack func(interface{})) chan<- func() interface{} {
	funcs := make(chan func() interface{}, maxFuncCount)
	go func(funcs chan func() interface{}, maxTaskCount int) {
		ch := make(chan chan struct{}, maxTaskCount)
		once := sync.Once{}
		// defer close(funcs)
		for {
			select {
			case fn, ok := <-funcs:
				if !ok {
					return
				}
				sign := make(chan struct{})
				ch <- sign
				once.Do(func() { close(<-ch) })
				go func(fn func() interface{}, callBack func(interface{})) {
					result := fn()
					<-sign
					defer close(<-ch)
					callBack(result)
				}(fn, callBack)
			case <-done:
				return
			}
		}
	}(funcs, maxTaskCount)
	return funcs
}

//OrChannel 演示如何组合多个通道
func OrChannel() {
	var or func(channels ...<-chan interface{}) <-chan interface{}
	or = func(channels ...<-chan interface{}) <-chan interface{} { //1
		switch len(channels) {
		case 0: //2
			return nil
		case 1: //3
			return channels[0]
		}
		orDone := make(chan interface{})
		go func() { //4
			defer close(orDone)
			switch len(channels) {
			case 2: //5
				select {
				case <-channels[0]:
				case <-channels[1]:
				}
			default: //6
				select {
				case <-channels[0]:
				case <-channels[1]:
				case <-channels[2]:
				case <-or(append(channels[3:], orDone)...): //6
				}
			}
		}()
		return orDone
	}
}

type GoPoll struct {
	work    chan func()
	sem     chan struct{}
	timeout time.Duration
}

func NewGoPoll(size int, forExit time.Duration) *GoPoll {
	return &GoPoll{
		work:    make(chan func()),
		sem:     make(chan struct{}, size),
		timeout: forExit,
	}
}

//Grow .
func (p *GoPoll) Grow(num int) error {
	newSem := make(chan struct{}, num)
loop:
	for {
		select {
		case sign := <-p.sem:
			select {
			case newSem <- sign:
			default:
			}
		default:
			break loop
		}
	}
	p.sem = newSem
	return nil
}

//Schedule 把方法加入协程池并被执行
func (p *GoPoll) Schedule(task func()) error {
	select {
	case p.work <- task:
	case p.sem <- struct{}{}:
		go p.worker(p.timeout, task)
	}
	return nil
}

func (p *GoPoll) worker(delay time.Duration, task func()) {
	defer func() { <-p.sem }()
	timer := time.NewTimer(delay)
	for {
		task()
		timer.Reset(delay)
		select {
		case task = <-p.work:
		case <-timer.C:
			return
		}
	}
}

func Poll(size int, forExit time.Duration) func(func()) error {
	var (
		work    chan func()   = make(chan func())
		sem     chan struct{} = make(chan struct{}, size)
		timeout time.Duration = forExit
		worker                = func(delay time.Duration, task func()) {
			defer func() { <-sem }()
			timer := time.NewTimer(delay)
			for {
				task()
				timer.Reset(delay)
				select {
				case task = <-work:
				case <-timer.C:
					return
				}
			}
		}
	)
	return func(task func()) error {
		select {
		case work <- task:
		case sem <- struct{}{}:
			go worker(timeout, task)
		}
		return nil
	}
}
