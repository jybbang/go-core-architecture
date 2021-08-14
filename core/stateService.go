package core

import (
	"sync"

	"github.com/sony/gobreaker"
)

type stateService struct {
	state StateAdapter
	cb    *gobreaker.CircuitBreaker
	sync.RWMutex
}

func (s *stateService) SetStateAdapter(adapter StateAdapter) *stateService {
	s.state = adapter
	return s
}

func (s *stateService) Has(key string) (bool, error) {
	resp, err := s.cb.Execute(func() (interface{}, error) {
		resp, err := s.state.Has(key)
		if err != nil {
			return false, err
		}

		return resp, nil
	})
	return resp.(bool), err
}

func (s *stateService) Get(key string, dest Entitier) (bool, error) {
	resp, err := s.cb.Execute(func() (interface{}, error) {
		resp, err := s.state.Get(key, dest)
		if err != nil {
			return false, err
		}

		return resp, nil
	})
	return resp.(bool), err
}

func (s *stateService) Set(key string, item interface{}) error {
	s.Lock()
	defer s.Unlock()

	_, err := s.cb.Execute(func() (interface{}, error) {
		err := s.state.Set(key, item)
		return nil, err
	})
	return err
}
