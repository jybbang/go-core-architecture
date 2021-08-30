package core

import (
	"context"
	"fmt"
	"reflect"

	"github.com/sony/gobreaker"
)

type stateService struct {
	state    stateAdapter
	cb       *gobreaker.CircuitBreaker
	settings StateServiceSettings
}

func (s *stateService) initialize() *stateService {
	if err := s.connect(); err != nil {
		panic(err)
	}

	return s
}

func (s *stateService) connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.settings.ConnectionTimeout)
	defer cancel()

	return s.state.Connect(ctx)
}

func (s *stateService) onCircuitOpen() {
	s.state.Disconnect()

	if !s.state.IsConnected() {
		s.connect()
	}
}

func (s *stateService) Has(ctx context.Context, key string) Result {
	if key == "" {
		return Result{V: false, E: fmt.Errorf("%w key is required", ErrInternalServerError)}
	}

	resp, err := s.cb.Execute(func() (interface{}, error) {
		ok := s.state.Has(ctx, key)

		return ok, nil
	})

	if err != nil {
		return Result{V: false, E: err}
	}

	return Result{V: resp.(bool), E: nil}
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
		return nil, s.state.Set(ctx, key, value)
	})

	return Result{V: nil, E: err}
}

func (s *stateService) BatchSet(ctx context.Context, kvs []KV) Result {
	if len(kvs) == 0 {
		return Result{V: nil, E: fmt.Errorf("%w kvs is required", ErrInternalServerError)}
	}

	_, err := s.cb.Execute(func() (interface{}, error) {
		return nil, s.state.BatchSet(ctx, kvs)
	})

	return Result{V: nil, E: err}
}

func (s *stateService) Delete(ctx context.Context, key string) Result {
	if key == "" {
		return Result{V: nil, E: fmt.Errorf("%w key is required", ErrInternalServerError)}
	}

	_, err := s.cb.Execute(func() (interface{}, error) {
		return nil, s.state.Delete(ctx, key)
	})

	return Result{V: nil, E: err}
}
