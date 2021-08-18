package core

import (
	"context"
	"fmt"
	"reflect"

	"github.com/sony/gobreaker"
)

type stateService struct {
	state stateAdapter
	cb    *gobreaker.CircuitBreaker
}

func (s *stateService) initialize() *stateService {
	return s
}

func (s *stateService) close() {
	s.state.Close()
}

func (s *stateService) Has(ctx context.Context, key string) Result {
	if key == "" {
		return Result{V: false, E: fmt.Errorf("%w key is required", ErrInternalServerError)}
	}

	resp, err := s.cb.Execute(func() (interface{}, error) {
		return s.state.Has(ctx, key)
	})
	if err != nil {
		return Result{V: false, E: err}
	}
	return Result{V: resp.(bool), E: err}
}

func (s *stateService) Get(ctx context.Context, key string, dest interface{}) Result {
	if key == "" {
		return Result{V: nil, E: fmt.Errorf("%w key is required", ErrInternalServerError)}
	}
	if dest == nil {
		return Result{V: nil, E: fmt.Errorf("%w dest is required", ErrInternalServerError)}
	}

	resultsVal := reflect.ValueOf(dest)
	if resultsVal.Kind() == reflect.Interface {
		resultsVal = resultsVal.Elem()
	}
	if resultsVal.Kind() != reflect.Ptr {
		panic("dest must be a pointer")
	}

	_, err := s.cb.Execute(func() (interface{}, error) {
		return nil, s.state.Get(ctx, key, dest)
	})

	return Result{V: dest, E: err}
}

func (s *stateService) Set(ctx context.Context, key string, value interface{}) Result {
	if key == "" {
		return Result{V: nil, E: fmt.Errorf("%w key is required", ErrInternalServerError)}
	}
	if value == nil {
		return Result{V: nil, E: fmt.Errorf("%w value is required", ErrInternalServerError)}
	}

	_, err := s.cb.Execute(func() (interface{}, error) {
		err := s.state.Set(ctx, key, value)
		return nil, err
	})

	return Result{V: nil, E: err}
}

func (s *stateService) Delete(ctx context.Context, key string) Result {
	if key == "" {
		return Result{V: nil, E: fmt.Errorf("%w key is required", ErrInternalServerError)}
	}

	_, err := s.cb.Execute(func() (interface{}, error) {
		err := s.state.Delete(ctx, key)
		return nil, err
	})

	return Result{V: nil, E: err}
}
