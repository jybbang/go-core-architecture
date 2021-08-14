package core

import (
	"sync"

	"github.com/sony/gobreaker"
)

type StateService struct {
	state StateAdapter
	cb    *gobreaker.CircuitBreaker
	sync.RWMutex
}

func (s *StateService) SetStateAdapter(adapter StateAdapter) *StateService {
	s.state = adapter
	return s
}

func (s *StateService) Has(key string) (bool, error) {
	resp, err := s.cb.Execute(func() (interface{}, error) {
		resp, err := s.state.Has(key)
		if err != nil {
			return false, err
		}

		return resp, nil
	})
	return resp.(bool), err
}

func (s *StateService) Get(key string, dest Entitier) (bool, error) {
	resp, err := s.cb.Execute(func() (interface{}, error) {
		resp, err := s.state.Get(key, dest)
		if err != nil {
			return false, err
		}

		return resp, nil
	})
	return resp.(bool), err
}

func (s *StateService) Set(key string, item interface{}) error {
	s.Lock()
	defer s.Unlock()

	_, err := s.cb.Execute(func() (interface{}, error) {
		err := s.state.Set(key, item)
		return nil, err
	})
	return err
}
