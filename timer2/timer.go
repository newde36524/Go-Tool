package timer2

import (
	"container/list"
	"fmt"

	"sync"

	"time"

	"github.com/google/uuid"
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

//Modify .
func (l *TimerTask) Modify(key string, mod func(*Entity) error) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if v, ok := l.mp[key]; ok {
		entity := v.Value.(*Entity)
		err := mod(entity)
		if err != nil {
			return err
		}
		l.delete(key)
		return l.add(entity)
	}
	return fmt.Errorf("key is not exist")
}

//Add .
func (l *TimerTask) Add(delay time.Duration, getStartRunTime func(key string, v interface{}) (time.Time, error), v interface{}, task func(key string, remove func())) (string, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
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
	l.mu.Lock()
	defer l.mu.Unlock()
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
	t, err := e.GetStartRunTime(e.key, e.Value)
	if err != nil {
		return err
	}
	//todo 需要用跳表替换遍历，一定量级以下使用遍历，以上使用跳表
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
	if point == l.tasks.Front() {
		e := point.Value.(*Entity)
		l.t.Reset(time.Until(e.start.Add(e.Delay)))
	}
	return nil
}

//Delete .
func (l *TimerTask) Delete(key string) {
	l.atomicInvoce(func() { l.delete(key) })
}

//delete .
func (l *TimerTask) delete(key string) {
	element, ok := l.mp[key]
	if ok {
		delete(l.mp, key)
		l.tasks.Remove(element)
	}
}

//Len .
func (l *TimerTask) Len() int {
	return l.tasks.Len()
}

//atomicInvoce .
func (l *TimerTask) atomicInvoce(fn func()) {
	l.mu.Lock()
	defer l.mu.Unlock()
	fn()
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
			if time.Since(entity.start) >= entity.Delay {
				l.atomicInvoce(func() {
					entity.Task(entity.key, remove)
					l.delete(entity.key)
					if !isRemove {
						_ = l.add(entity)
					}
				})
			} else {
				l.t.Reset(entity.Delay - time.Since(entity.start))
				return
			}
		}
	})
}
