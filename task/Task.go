package task

//Task Task结构体
type Task struct {
	ch    chan int
	funcs []func()
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
		ch: make(chan int, 1),
	}
}

//Start 开始执行任务
//@fn 需要执行的任务，任务会在协程中执行
func (t *Task) Start(fn func()) {
	go func(ch chan int) {
		fn()
		ch <- 1
	}(t.ch)
	go func(_task *Task) {
		for _, f := range _task.funcs {
			<-_task.ch
			go func(ch chan int, _f func()) {
				_f()
				ch <- 1
			}(_task.ch, f)
		}
		close(_task.ch)
	}(t)
}

//Continue 执行延续任务，上一个任务完成时才会执行下一个
func (t *Task) Continue(fn func()) *Task {
	t.funcs = append(t.funcs, fn)
	return t
}
