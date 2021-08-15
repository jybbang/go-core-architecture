package core

import (
	"context"
	"sync"

	"github.com/sony/gobreaker"
)

type stateService struct {
	state stateAdapter
	cb    *gobreaker.CircuitBreaker
	sync.RWMutex
}

func (s *stateService) initialize() *stateService {
	return s
}

func (s *stateService) Has(ctx context.Context, key string) (ok bool, err error) {
	s.RLocker()
	defer s.RUnlock()

	resp, err := s.cb.Execute(func() (interface{}, error) {
		return s.state.Has(ctx, key)
	})
	return resp.(bool), err
}

func (s *stateService) Get(ctx context.Context, key string, dest interface{}) (ok bool, err error) {
	s.RLocker()
	defer s.RUnlock()

	resp, err := s.cb.Execute(func() (interface{}, error) {
		return s.state.Get(ctx, key, dest)
	})
	return resp.(bool), err
}

func (s *stateService) Set(ctx context.Context, key string, value interface{}) error {
	s.Lock()
	defer s.Unlock()

	_, err := s.cb.Execute(func() (interface{}, error) {
		err := s.state.Set(ctx, key, value)
		return nil, err
	})
	return err
}
