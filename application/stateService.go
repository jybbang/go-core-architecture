package application

import (
	"sync"

	"github.com/jybbang/go-core-architecture/application/contracts"
	"github.com/jybbang/go-core-architecture/domain"
	"github.com/sony/gobreaker"
)

type stateService struct {
	state contracts.StateAdapter
	cb    *gobreaker.CircuitBreaker
	sync.RWMutex
}

func (s *stateService) SetStateAdapter(adapter contracts.StateAdapter) *stateService {
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

func (s *stateService) Get(key string, dest domain.Entitier) (bool, error) {
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
