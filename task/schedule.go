package task

import (
	"sync"
	"sync/atomic"
)

// Schedule 任务调度器 使用指定的方法列表均衡处理数据源
func Schedule[T any](src []T, actions []func(v T)) {
	taskEntrys := make([]*taskEntry[T], len(actions))
	for i, v := range actions {
		taskEntrys[i] = newTaskEntry(v)
	}
	wg := new(sync.WaitGroup)
	for j := 0; j < len(src); j++ {
		item := src[j]
		isUse := false //数据项是否被使用
		for _, entry := range taskEntrys {
			if entry.CanUse() {
				//1.占用
				wg.Add(1)
				entry.Use()
				isUse = true
				//2.异步执行后解除占用
				go func(item T, taskEntry *taskEntry[T]) {
					defer wg.Done()
					defer taskEntry.UnUse()
					taskEntry.Invoke(item)
				}(item, entry)
				//3.数据项只被使用一次
				break
			}
		}
		if !isUse {
			j-- //如果任务都被占用则不跳过数据项
		}
	}
	wg.Wait()
}

type taskEntry[T any] struct {
	canUse int64
	action func(v T)
}

func newTaskEntry[T any](action func(v T)) *taskEntry[T] {
	return &taskEntry[T]{
		action: action,
	}
}

func (e *taskEntry[T]) CanUse() bool {
	return atomic.LoadInt64(&e.canUse) == 0
}

func (e *taskEntry[T]) Invoke(v T) {
	e.action(v)
}

func (e *taskEntry[T]) UnUse() {
	atomic.StoreInt64(&e.canUse, 0)
}

func (e *taskEntry[T]) Use() {
	atomic.StoreInt64(&e.canUse, 1)
}
