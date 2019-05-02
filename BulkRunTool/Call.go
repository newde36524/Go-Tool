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

//CreateBulkRunFuncChannel 创建一个指定并发数量处理,只允许发送的方法通道
func CreateBulkRunFuncChannel(maxTaskCount, maxFuncCount int) chan<- func() {
	return createBulkRunFuncChannel(maxTaskCount, maxFuncCount)
}

//createBulkRunFuncChannel 创建一个指定并发数量处理的方法通道
func createBulkRunFuncChannel(maxTaskCount, maxFuncCount int) (funcs chan func()) {
	funcs = make(chan func(), maxFuncCount)
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
			}
		}
	}(funcs, maxTaskCount)
	return funcs
}
func createBulkRunFuncChannel2(maxTaskCount, maxFuncCount int, done <-chan struct{}) (funcs chan func()) {
	funcs = make(chan func(), maxFuncCount)
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
