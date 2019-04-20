package BulkRunTool

//RunTask 运行指定数量的协程执行多个方法
func RunTask(maxTask int32, funcs []func()) {
	ch := make(chan int, maxTask)
	for _, fn := range funcs {
		temp := fn
		ch <- 1
		go func() {
			temp()
			<-ch
		}()
	}
}
