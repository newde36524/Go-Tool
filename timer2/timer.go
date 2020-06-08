package timer2

import (
	"container/list"
	"errors"
	"sync"
	"time"
)

//Entity .
type Entity struct {
	key             string
	getStartRunTime func(interface{}) time.Time
	start           time.Time
	delay           time.Duration //超过时间就执行命令
	task            func(remove func())
	value           interface{}
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
func (l *TimerTask) Add(key string, getStartRunTime func(interface{}) time.Time, value interface{}, delay time.Duration, task func(remove func())) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	e := &Entity{
		key:             key,
		getStartRunTime: getStartRunTime,
		delay:           delay,
		task:            task,
	}
	e.start = getStartRunTime(value)
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
	if _, ok := l.mp[key]; ok {
		return errors.New("key 已存在")
	}
	l.mp[key] = point
	l.cond.Signal()
	return nil
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
					l.Add(entity.key, entity.getStartRunTime, entity.value, entity.delay, entity.task)
				}
			} else {
				l.t.Reset(entity.delay - time.Now().Sub(entity.start))
				return
			}
		}
	})
}
