package bulkruntool

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

//CreateBulkRunFuncChannel 创建一个指定并发数量处理的方法通道
func CreateBulkRunFuncChannel(maxTaskCount, maxFuncCount int) (funcs chan func()) {
	funcs = make(chan func(), maxFuncCount)
	ch := make(chan struct{}, maxTaskCount)
	go func(funcs chan func(), ch chan struct{}) {
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
			}
		}
	}(funcs, ch)
	return
}
