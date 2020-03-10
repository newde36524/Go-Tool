package reference

import (
	"sync"
	"time"
)

/*
防抖节流函数的区别:
1. 两个方法的共同点都是降低操作频率,但在不同的业务场景下会有细节区别
2. 防抖函数不会在一开始执行(防止无意义的重复触发操作),而节流函数会(旨在降低调用频率)
3. 防抖函数必须是间隔时间内不再调用时才会触发操作,否则永远不执行,而节流函数会至少在间隔超时后执行
4. 防抖和节流函数解决的是操作重复触发的问题,在调用频率上如果超过函数触发时间间隔,那么两个函数表现一致
*/

//Throttle 方法节流,降低方法调用频率
//一段时间内只会调用一次
/*
	example:
	fn := A(func(){fmt.Println(time.Now())},time.Second)
	for {
		fn()
	}
*/
func Throttle(fn func(), delay time.Duration) func() {
	canDo := true
	return func() {
		if canDo {
			canDo = false
			fn()
			<-time.After(delay)
			canDo = true
		}
	}
}

func Throttle2(fn func(), delay time.Duration) func() {
	tricker := time.NewTicker(delay)
	once := sync.Once{}
	return func() {
		select {
		case <-tricker.C:
			fn()
		default:
			once.Do(fn)
		}
	}
}

//Debounce 防抖方法
//超过时间间隔后才会被执行,防止短时间内重复调用
func Debounce(fn func(), timeout time.Duration) func() {
	timeOutCh := time.After(timeout)
	sign := make(chan struct{}, 1)
	sign <- struct{}{}
	return func() {
		select {
		case <-sign:
			go func() {
				for {
					select {
					case <-timeOutCh:
						fn()
						sign <- struct{}{}
						return
					default:
					}
				}
			}()
		default:
			timeOutCh = time.After(timeout)
		}
	}
}

//Debounce2 防抖方法
//超过时间间隔后才会被执行,防止短时间内重复调用
func Debounce2(fn func(), delay time.Duration) func() {
	prev := time.Unix(0, 0)
	return func() {
		curr := time.Now()
		sub := curr.Sub(prev).Seconds()
		if sub < float64(delay) {
			return
		}
		prev = curr // 可执行了之后，在刷新计时
		fn()
	}
}

//Debounce3 防抖方法
//超过时间间隔后才会被执行,防止短时间内重复调用
func Debounce3(fn func(), delay time.Duration) func() {
	var timer *time.Timer
	var once sync.Once
	return func() {
		once.Do(func() {
			timer = time.AfterFunc(delay, fn)
		})
		timer.Reset(delay)
	}
}

type throttleImpl struct {
	once    sync.Once
	fn      func()
	tricker *time.Ticker
}

func NewThrottleImpl(delay time.Duration, fn func()) *throttleImpl {
	return &throttleImpl{
		fn:      fn,
		tricker: time.NewTicker(delay),
	}
}

func (t *throttleImpl) Do() {
	select {
	case <-t.tricker.C:
		t.fn()
	default:
		t.once.Do(t.fn)
	}
}

type debounceImpl struct {
	timer *time.Timer
	once  sync.Once
	fn    func()
	delay time.Duration
}

func NewDebounceImpl(delay time.Duration, fn func()) *debounceImpl {
	return &debounceImpl{
		fn:    fn,
		delay: delay,
	}
}

func (d *debounceImpl) Do() {
	d.once.Do(func() {
		d.timer = time.AfterFunc(d.delay, d.fn)
	})
	d.timer.Reset(d.delay)
}
