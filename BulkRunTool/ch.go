package bulkruntool

import (
	"context"
	"time"
)

//OrChannel 演示如何组合多个通道
func OrChannel() {
	var or func(channels ...<-chan interface{}) <-chan interface{}
	or = func(channels ...<-chan interface{}) <-chan interface{} { //1
		switch len(channels) {
		case 0: //2
			return nil
		case 1: //3
			return channels[0]
		}
		orDone := make(chan interface{})
		go func() { //4
			defer close(orDone)
			switch len(channels) {
			case 2: //5
				select {
				case <-channels[0]:
				case <-channels[1]:
				}
			default: //6
				select {
				case <-channels[0]:
				case <-channels[1]:
				case <-channels[2]:
				case <-or(append(channels[3:], orDone)...): //6
				}
			}
		}()
		return orDone
	}
}

//MergeChan 组合多个通道
func MergeChan(ctx context.Context, do func(interface{}), chs ...<-chan interface{}) {
	var (
		index  = 0
		delay  = time.Second
		timer  = time.NewTimer(delay)
		result func()
	)
	result = func() {
		index = index%len(chs) + 1
		timer.Reset(delay)
		ch := chs[index-1]
		select {
		case <-ctx.Done():
			return
		default:
		}
		if len(ch) != 0 {
			select {
			case v := <-ch:
				do(v)
			default:
			}
		}
		go result()
	}
	result()
}
