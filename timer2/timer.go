package timer2

import (
	"container/list"
	"sync"
	"time"
)

//Entity .
type Entity struct {
	key   string
	start time.Time
	delay time.Duration //超过时间就执行命令
	task  func(remove func())
}

//TimerTask .
type TimerTask struct {
	tasks *list.List
	mp    map[string]*list.Element
	mu    sync.Mutex
	t     *time.Timer
	cond  *sync.Cond
}

//NewTimerTask .
func NewTimerTask() *TimerTask {
	mu := &sync.Mutex{}
	mu.Lock()
	result := &TimerTask{
		tasks: list.New(),
		mp:    make(map[string]*list.Element),
		cond:  sync.NewCond(mu),
	}
	result.start()
	return result
}

//Add .
func (l *TimerTask) Add(key string, delay time.Duration, task func(remove func())) {
	l.mu.Lock()
	defer l.mu.Unlock()
	e := &Entity{
		key:   key,
		start: time.Now(),
		delay: delay,
		task:  task,
	}
	point := l.tasks.Front()
	for {
		if point == nil {
			point = l.tasks.PushBack(e)
			break
		}
		entity := point.Value.(*Entity)
		if entity.start.Add(entity.delay).Sub(time.Now().Add(delay)) >= 0 {
			point = l.tasks.InsertBefore(e, point)
			break
		}
		point = point.Next()
	}
	l.mp[key] = point
	l.cond.Signal()
}

//Delete .
func (l *TimerTask) Delete(key string) {
	l.mu.Lock()
	element, ok := l.mp[key]
	if ok {
		l.tasks.Remove(element)
		delete(l.mp, key)
	}
	l.mu.Unlock()
}

//Len .
func (l *TimerTask) Len() int {
	return l.tasks.Len()
}

//start .
func (l *TimerTask) start() {
	l.t = time.AfterFunc(0, func() {
		for {
			point := l.tasks.Front()
			if point == nil {
				l.cond.Wait()
				continue
			}
			var (
				entity   = point.Value.(*Entity)
				isRemove = false
				remove   = func() { isRemove = true }
			)
			if time.Now().Sub(entity.start) >= entity.delay {
				entity.task(remove)
				l.Delete(entity.key)
				if !isRemove {
					l.Add(entity.key, entity.delay, entity.task)
				}
			} else {
				l.t.Reset(entity.delay - time.Now().Sub(entity.start))
				return
			}
		}
	})
}
