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
	GetStartRunTime func(key string, v interface{}) (time.Time, error)
	start           time.Time
	Delay           time.Duration
	Task            func(key string, remove func())
	Value           interface{}
}

//TimerTask .
type TimerTask struct {
	tasks *list.List
	mp    map[string]*list.Element
	mu    sync.Mutex
	t     *time.Timer
	cond  *sync.Cond
}

//New .
func New() *TimerTask {
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

func (l *TimerTask) Modify(key string, mod func(*Entity) error) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if v, ok := l.mp[key]; ok {
		entity := v.Value.(*Entity)
		err := mod(entity)
		if err != nil {
			return err
		}
		l.Delete(key)
		l.add(entity)
		return nil
	} else {
		return fmt.Errorf("key is not exist")
	}
}

//Add .
func (l *TimerTask) Add(delay time.Duration, getStartRunTime func(key string, v interface{}) (time.Time, error), v interface{}, task func(key string, remove func())) (string, error) {
	key := uuid.New().String()
	e := &Entity{
		key:             key,
		GetStartRunTime: getStartRunTime,
		Delay:           delay,
		Value:           v,
		Task:            task,
	}
	err := l.add(e)
	if err != nil {
		return "", err
	}
	return key, nil
}

//Sync .
func (l *TimerTask) Sync(key string, delay time.Duration, getStartRunTime func(key string, v interface{}) (time.Time, error), v interface{}, task func(key string, remove func())) error {
	e := &Entity{
		key:             key,
		GetStartRunTime: getStartRunTime,
		Delay:           delay,
		Value:           v,
		Task:            task,
	}
	return l.add(e)
}

func (l *TimerTask) add(e *Entity) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	t, err := e.GetStartRunTime(e.key, e.Value)
	if err != nil {
		return err
	}
	e.start = t
	point := l.tasks.Front()
	for {
		if point == nil {
			point = l.tasks.PushBack(e)
			break
		}
		entity := point.Value.(*Entity)
		if entity.start.Add(entity.Delay).Sub(time.Now().Add(e.Delay)) >= 0 {
			point = l.tasks.InsertBefore(e, point)
			break
		}
		point = point.Next()
	}
	l.mp[e.key] = point
	l.cond.Signal()
	return nil
}

//Delete .
func (l *TimerTask) Delete(key string) {
	l.mu.Lock()
	element, ok := l.mp[key]
	if ok {
		delete(l.mp, key)
		l.tasks.Remove(element)
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
			if time.Now().Sub(entity.start) >= entity.Delay {
				entity.Task(entity.key, remove)
				l.Delete(entity.key)
				if !isRemove {
					l.add(entity)
				}
			} else {
				l.t.Reset(entity.Delay - time.Now().Sub(entity.start))
				return
			}
		}
	})
}
