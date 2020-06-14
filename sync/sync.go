package sync

import (
	"context"
	"errors"
	"sync"
	"time"
)

type SyncOp struct {
	mu    sync.Mutex
	m     map[string]chan *Entity
	signs map[string]chan struct{}
}

func NewSyncOp() *SyncOp {
	return &SyncOp{
		m:     make(map[string]chan *Entity),
		signs: make(map[string]chan struct{}),
	}
}

func (s *SyncOp) Do(timeout time.Duration, key string, fn func()) (interface{}, error) {
	s.mu.Lock()
	ch, ok := s.m[key]
	if !ok {
		ch = make(chan *Entity)
		s.m[key] = ch
	}
	if _, ok := s.signs[key]; !ok {
		s.signs[key] = make(chan struct{}, 1)
	}
	s.mu.Unlock()
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	select {
	case <-ctx.Done():
		return nil, errors.New("任务超时1")
	case s.signs[key] <- struct{}{}:
		go fn()
	}
	defer func() {
		<-s.signs[key]
	}()
	select {
	case <-ctx.Done():
		return nil, errors.New("任务超时2")
	case v := <-ch:
		return v.value, v.err
	}
}

type Entity struct {
	value interface{}
	err   error
}

func (s *SyncOp) Back(key string, fn func() (interface{}, error)) (bool, error) {
	s.mu.Lock()
	v, ok := s.m[key]
	s.mu.Unlock()
	if !ok {
		return false, nil
	}
	value, err := fn()
	v <- &Entity{
		value: value,
		err:   err,
	}
	return true, nil
}
