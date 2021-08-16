package core

import (
	"context"

	"github.com/sony/gobreaker"
)

type stateService struct {
	state stateAdapter
	cb    *gobreaker.CircuitBreaker
}

func (s *stateService) initialize() *stateService {
	return s
}

func (s *stateService) Has(ctx context.Context, key string) (ok bool, err error) {
	resp, err := s.cb.Execute(func() (interface{}, error) {
		return s.state.Has(ctx, key)
	})
	return resp.(bool), err
}

func (s *stateService) Get(ctx context.Context, key string, dest interface{}) (err error) {
	_, err = s.cb.Execute(func() (interface{}, error) {
		return nil, s.state.Get(ctx, key, dest)
	})

	return err
}

func (s *stateService) Set(ctx context.Context, key string, value interface{}) error {
	_, err := s.cb.Execute(func() (interface{}, error) {
		err := s.state.Set(ctx, key, value)
		return nil, err
	})
	return err
}

func (s *stateService) Delete(ctx context.Context, key string) error {
	_, err := s.cb.Execute(func() (interface{}, error) {
		err := s.state.Delete(ctx, key)
		return nil, err
	})
	return err
}
