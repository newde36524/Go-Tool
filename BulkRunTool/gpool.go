package bulkruntool

import (
	"context"
	"runtime"
	"sync"
	"time"
)

//gPool .
//针对key值进行并行调用的协程池,同一个key下的任务串行,不同key下的任务并行
//实现协程池区分执行不同任务类型的调度
type gPool struct {
	ctx     context.Context
	taskNum int
	exp     time.Duration
	m       sync.Map
	sign    chan struct{}
}

func newgPoll(ctx context.Context, perItemTaskNum int, exp time.Duration, parallelSize int) *gPool {
	g := &gPool{
		ctx:     ctx,
		taskNum: perItemTaskNum,
		exp:     exp,
		sign:    make(chan struct{}, parallelSize), //创建的协程池数量
	}
	return g
}

//SchduleByKey 为不同key值下的任务并行调用,相同key值下的任务串行调用,并行任务量和串行任务量由配置参数决定
func (g *gPool) SchduleByKey(key interface{}, task func()) bool {
	if v, ok := g.m.Load(key); ok { //希望在同一协程下顺序执行
		gItem := v.(*gItem)
		return gItem.DoOrInChan(task)
	}
	select {
	case <-g.ctx.Done():
		return false
	case g.sign <- struct{}{}:
		gItem := newgItem(g.ctx, g.taskNum, g.exp, func() { g.m.Delete(key) })
		g.m.Store(key, gItem)
		return gItem.DoOrInChan(task)
	}
}

type gItem struct {
	tasks  chan func()     //任务通道
	sign   chan struct{}   //是否加入任务通道信号
	ctx    context.Context //退出协程信号
	exp    time.Duration
	onExit func()
}

func newgItem(ctx context.Context, taskNum int, exp time.Duration, onExit func()) *gItem {
	return &gItem{
		tasks:  make(chan func(), taskNum),
		sign:   make(chan struct{}, 1), //最多只创建一个协程
		ctx:    ctx,
		exp:    exp,
		onExit: onExit,
	}
}

func (g *gItem) DoOrInChan(task func()) bool {
	select {
	case <-g.ctx.Done():
		return false
	default:
	}
	select {
	case g.sign <- struct{}{}:
		go g.worker()
		runtime.Gosched()
		return g.DoOrInChan(task)
	default:
	}
	select {
	case g.tasks <- task:
		return true
	default:
		runtime.Gosched()
		return false
	}
}

func (g *gItem) worker() {
	timer := time.NewTimer(g.exp)
	defer func() {
		select {
		case <-g.sign:
		default:
		}
		timer.Stop()
		if g.onExit != nil {
			g.onExit()
		}
	}()
	for {
		select {
		case <-g.ctx.Done():
			return
		case task, ok := <-g.tasks: //执行任务优先
			if !ok {
				return
			}
			//timer.Reset(g.exp)
			/*
				1) 如果重置时间,那么会在任务全部处理完成后继续等待过期,虽然空闲等待是一种资源浪费,但这主要用于复用当前协程对任务队列的执行
				2) 如果不重置时间,那么会在任务队列为空时并且过期后退出协程
				3) 个人认为,不重置时间可均衡各个任务队列之间的任务调度
				4) 应根据实际应用场景设置过期时间,并且时间一般不宜过长,在3~5秒左右
			*/
			if task != nil {
				task()
			}
		case <-timer.C:
			return
		}
	}
}
