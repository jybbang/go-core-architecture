package core

import (
	"context"
	"sync"

	"github.com/sony/gobreaker"
)

type StateService struct {
	state StateAdapter
	cb    *gobreaker.CircuitBreaker
	sync.RWMutex
}

func (s *StateService) Setup() *StateService {
	return s
}

func (s *StateService) SetStateAdapter(adapter StateAdapter) *StateService {
	s.state = adapter
	return s
}

func (s *StateService) Has(ctx context.Context, key string) (ok bool, err error) {
	resp, err := s.cb.Execute(func() (interface{}, error) {
		return s.state.Has(ctx, key)
	})
	return resp.(bool), err
}

func (s *StateService) Get(ctx context.Context, key string, dest interface{}) (ok bool, err error) {
	resp, err := s.cb.Execute(func() (interface{}, error) {
		return s.state.Get(ctx, key, dest)
	})
	return resp.(bool), err
}

func (s *StateService) Set(ctx context.Context, key string, value interface{}) error {
	s.Lock()
	defer s.Unlock()

	_, err := s.cb.Execute(func() (interface{}, error) {
		err := s.state.Set(ctx, key, value)
		return nil, err
	})
	return err
}
