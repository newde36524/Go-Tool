package bulkruntool

//RunTask 运行指定数量的协程执行多个方法
func RunTask(maxTask int, funcs []func()) {
	ch := make(chan int, maxTask)
	for _, fn := range funcs {
		ch <- 1
		go func(f func()) {
			f()
			<-ch
		}(fn)
	}
	close(ch)
}
