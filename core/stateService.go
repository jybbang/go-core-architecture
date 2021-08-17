package core

import (
	"context"
	"fmt"

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
	if key == "" {
		return false, fmt.Errorf("%w key is required", ErrInternalServerError)
	}

	resp, err := s.cb.Execute(func() (interface{}, error) {
		return s.state.Has(ctx, key)
	})
	return resp.(bool), err
}

func (s *stateService) Get(ctx context.Context, key string, dest interface{}) (err error) {
	if key == "" {
		return fmt.Errorf("%w key is required", ErrInternalServerError)
	}
	if dest == nil {
		return fmt.Errorf("%w dest is required", ErrInternalServerError)
	}

	_, err = s.cb.Execute(func() (interface{}, error) {
		return nil, s.state.Get(ctx, key, dest)
	})

	return err
}

func (s *stateService) Set(ctx context.Context, key string, value interface{}) error {
	if key == "" {
		return fmt.Errorf("%w key is required", ErrInternalServerError)
	}
	if value == nil {
		return fmt.Errorf("%w value is required", ErrInternalServerError)
	}

	_, err := s.cb.Execute(func() (interface{}, error) {
		err := s.state.Set(ctx, key, value)
		return nil, err
	})
	return err
}

func (s *stateService) Delete(ctx context.Context, key string) error {
	if key == "" {
		return fmt.Errorf("%w key is required", ErrInternalServerError)
	}

	_, err := s.cb.Execute(func() (interface{}, error) {
		err := s.state.Delete(ctx, key)
		return nil, err
	})
	return err
}
