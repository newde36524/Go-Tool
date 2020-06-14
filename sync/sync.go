package sync

import (
	"context"
	"errors"
	"sync"
)

type SyncOp struct {
	mu sync.Mutex
	m  map[string]chan interface{}
}

func NewSyncOp() *SyncOp {
	return &SyncOp{
		m: make(map[string]chan interface{}),
	}
}

func (s *SyncOp) Do(ctx context.Context, key string, fn func()) (interface{}, error) {
	s.mu.Lock()
	ch := make(chan interface{})
	s.m[key] = ch
	s.mu.Unlock()
	fn()
	select {
	case <-ctx.Done():
		return 0, errors.New("中断操作")
	case v := <-ch:
		return v, nil
	}
}

func (s *SyncOp) Back(key string, fn func() interface{}) (bool, error) {
	s.mu.Lock()
	v, ok := s.m[key]
	s.mu.Unlock()
	if !ok {
		return false, nil
	}
	value := fn()
	select {
	case v <- value:
		return true, nil
	default:
		return false, nil
	}
}
