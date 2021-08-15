package core

import (
	"context"
	"sync"

	"github.com/sony/gobreaker"
)

type stateService struct {
	state stateAdapter
	cb    *gobreaker.CircuitBreaker
	mutex sync.RWMutex
}

func (s *stateService) initialize() *stateService {
	return s
}

func (s *stateService) Has(ctx context.Context, key string) (ok bool, err error) {
	s.mutex.RLocker()
	defer s.mutex.RUnlock()

	resp, err := s.cb.Execute(func() (interface{}, error) {
		return s.state.Has(ctx, key)
	})
	return resp.(bool), err
}

func (s *stateService) Get(ctx context.Context, key string, dest interface{}) (ok bool, err error) {
	s.mutex.RLocker()
	defer s.mutex.RUnlock()

	resp, err := s.cb.Execute(func() (interface{}, error) {
		return s.state.Get(ctx, key, dest)
	})
	return resp.(bool), err
}

func (s *stateService) Set(ctx context.Context, key string, value interface{}) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, err := s.cb.Execute(func() (interface{}, error) {
		err := s.state.Set(ctx, key, value)
		return nil, err
	})
	return err
}
