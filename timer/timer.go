package timer

import (
	"sync"
	"time"
)

//entity .
type entity struct {
	start time.Time
	delay time.Duration
	task  func(remove func())
}

type loopTask struct {
	tasks []entity
	delay time.Duration
	mu    sync.Mutex
}

func (l *loopTask) Add(delay time.Duration, task func(remove func())) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.tasks = append(l.tasks, entity{
		start: time.Now(),
		delay: delay,
		task:  task,
	})
}

func (l *loopTask) Len() int {
	return len(l.tasks)
}

func (l *loopTask) Start() {
	var t *time.Timer
	t = time.AfterFunc(l.delay, func() {
		for i := 0; i < len(l.tasks); i++ {
			var (
				once   sync.Once
				entity = l.tasks[i]
				remove = func() {
					once.Do(func() {
						front := l.tasks[:i]
						back := l.tasks[i+1:]
						l.tasks = append(front, back...)
						i--
					})
				}
			)
			if time.Now().Sub(entity.start) > entity.delay {
				entity.start = time.Now()
				entity.task(remove)
			}
		}
		t.Reset(l.delay)
	})
}

//loopTaskPool 循环任务池
type loopTaskPool struct {
	pool  sync.Pool
	once  sync.Once
	mu    sync.Mutex
	loops []*loopTask
	idx   int
}

//Schdule .
func (l *loopTaskPool) Schdule(delay time.Duration, task func(remove func())) {
	l.once.Do(func() {
		l.pool.New = func() interface{} {
			v := &loopTask{
				delay: delay,
			}
			v.Start()
			l.loops = append(l.loops, v)
			return v
		}
	})
	v := l.pool.Get()
	l.pool.Put(&v)
	v = l.loops[l.idx%len(l.loops)] //012012012012012
	loopTask := v.(*loopTask)
	loopTask.Add(delay, task)
	l.mu.Lock()
	l.idx++
	l.mu.Unlock()
}
