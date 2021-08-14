package core

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sony/gobreaker"
)

type RepositoryService struct {
	model             Entitier
	queryRepository   QueryRepositoryAdapter
	commandRepository CommandRepositoryAdapter
	cb                *gobreaker.CircuitBreaker
	sync.RWMutex
}

func (r *RepositoryService) Setup() *RepositoryService {
	return r
}

func (r *RepositoryService) SetQueryRepositoryAdapter(adapter QueryRepositoryAdapter) *RepositoryService {
	r.queryRepository = adapter
	r.model.SetID(uuid.Nil)
	r.queryRepository.SetModel(r.model)
	return r
}

func (r *RepositoryService) SetCommandRepositoryAdapter(adapter CommandRepositoryAdapter) *RepositoryService {
	r.commandRepository = adapter
	r.model.SetID(uuid.Nil)
	r.commandRepository.SetModel(r.model)
	return r
}

func (r *RepositoryService) Find(dto Entitier, id uuid.UUID) error {
	_, err := r.cb.Execute(func() (interface{}, error) {
		err := r.queryRepository.Find(dto, id)
		return nil, err
	})
	return err
}

func (r *RepositoryService) Any() (bool, error) {
	resp, err := r.cb.Execute(func() (interface{}, error) {
		return r.queryRepository.Any()
	})
	return resp.(bool), err
}

func (r *RepositoryService) AnyWithFilter(query interface{}, args interface{}) (bool, error) {
	resp, err := r.cb.Execute(func() (interface{}, error) {
		return r.queryRepository.AnyWithFilter(query, args)
	})
	return resp.(bool), err
}

func (r *RepositoryService) Count() (int64, error) {
	resp, err := r.cb.Execute(func() (interface{}, error) {
		return r.queryRepository.Count()
	})
	return resp.(int64), err
}

func (r *RepositoryService) CountWithFilter(query interface{}, args interface{}) (int64, error) {
	resp, err := r.cb.Execute(func() (interface{}, error) {
		return r.queryRepository.CountWithFilter(query, args)
	})
	return resp.(int64), err
}

func (r *RepositoryService) List(dtos []Entitier) error {
	_, err := r.cb.Execute(func() (interface{}, error) {
		err := r.queryRepository.List(dtos)
		return nil, err
	})
	return err
}

func (r *RepositoryService) ListWithFilter(dtos []Entitier, query interface{}, args interface{}) error {
	_, err := r.cb.Execute(func() (interface{}, error) {
		err := r.queryRepository.ListWithFilter(dtos, query, args)
		return nil, err
	})
	return err
}

func (r *RepositoryService) Remove(entity Entitier) error {
	r.Lock()
	defer r.Unlock()

	_, err := r.cb.Execute(func() (interface{}, error) {
		err := r.commandRepository.Remove(entity)
		return nil, err
	})
	return err
}

func (r *RepositoryService) RemoveRange(entities []Entitier) error {
	r.Lock()
	defer r.Unlock()

	_, err := r.cb.Execute(func() (interface{}, error) {
		err := r.commandRepository.RemoveRange(entities)
		return nil, err
	})
	return err
}

func (r *RepositoryService) Add(entity Entitier) error {
	r.Lock()
	defer r.Unlock()

	entity.SetCreatedAt(time.Now())
	_, err := r.cb.Execute(func() (interface{}, error) {
		err := r.commandRepository.Add(entity)
		return nil, err
	})
	return err
}

func (r *RepositoryService) AddRange(entities []Entitier) error {
	r.Lock()
	defer r.Unlock()

	now := time.Now()
	for _, v := range entities {
		v.SetCreatedAt(now)
	}

	_, err := r.cb.Execute(func() (interface{}, error) {
		err := r.commandRepository.AddRange(entities)
		return nil, err
	})
	return err
}

func (r *RepositoryService) Update(entity Entitier) error {
	r.Lock()
	defer r.Unlock()

	entity.SetUpdatedAt(time.Now())

	_, err := r.cb.Execute(func() (interface{}, error) {
		err := r.commandRepository.Update(entity)
		return nil, err
	})
	return err
}

func (r *RepositoryService) UpdateRange(entities []Entitier) error {
	r.Lock()
	defer r.Unlock()

	now := time.Now()
	for _, v := range entities {
		v.SetUpdatedAt(now)
	}

	_, err := r.cb.Execute(func() (interface{}, error) {
		err := r.commandRepository.UpdateRange(entities)
		return nil, err
	})
	return err
}
