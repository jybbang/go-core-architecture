package application

import (
	"sync"

	"github.com/google/uuid"
	"github.com/jybbang/go-core-architecture/application/contracts"
	"github.com/jybbang/go-core-architecture/domain"
	"github.com/sony/gobreaker"
)

type repositoryService struct {
	model             *domain.Entity
	queryRepository   contracts.QueryRepositoryAdapter
	commandRepository contracts.CommandRepositoryAdapter
	cb                *gobreaker.CircuitBreaker
	sync.RWMutex
}

func (r *repositoryService) SetQueryRepositoryAdapter(adapter contracts.QueryRepositoryAdapter) *repositoryService {
	r.queryRepository = adapter
	r.model.ID = uuid.Nil
	r.queryRepository.SetModel(r.model)
	return r
}

func (r *repositoryService) SetCommandRepositoryAdapter(adapter contracts.CommandRepositoryAdapter) *repositoryService {
	r.commandRepository = adapter
	r.model.ID = uuid.Nil
	r.commandRepository.SetModel(r.model)
	return r
}

func (r *repositoryService) Find(dto domain.Copyable, id uuid.UUID) error {
	_, err := r.cb.Execute(func() (interface{}, error) {
		err := r.queryRepository.Find(dto, id)
		return nil, err
	})
	return err
}

func (r *repositoryService) Any() (bool, error) {
	resp, err := r.cb.Execute(func() (interface{}, error) {
		return r.queryRepository.Any()
	})
	return resp.(bool), err
}

func (r *repositoryService) AnyWithFilter(query interface{}, args interface{}) (bool, error) {
	resp, err := r.cb.Execute(func() (interface{}, error) {
		return r.queryRepository.AnyWithFilter(query, args)
	})
	return resp.(bool), err
}

func (r *repositoryService) Count() (int64, error) {
	resp, err := r.cb.Execute(func() (interface{}, error) {
		return r.queryRepository.Count()
	})
	return resp.(int64), err
}

func (r *repositoryService) CountWithFilter(query interface{}, args interface{}) (int64, error) {
	resp, err := r.cb.Execute(func() (interface{}, error) {
		return r.queryRepository.CountWithFilter(query, args)
	})
	return resp.(int64), err
}

func (r *repositoryService) List(dtos []domain.Copyable) error {
	_, err := r.cb.Execute(func() (interface{}, error) {
		err := r.queryRepository.List(dtos)
		return nil, err
	})
	return err
}

func (r *repositoryService) ListWithFilter(dtos []domain.Copyable, query interface{}, args interface{}) error {
	_, err := r.cb.Execute(func() (interface{}, error) {
		err := r.queryRepository.ListWithFilter(dtos, query, args)
		return nil, err
	})
	return err
}

func (r *repositoryService) Remove(entity *domain.Entity) error {
	r.Lock()
	defer r.Unlock()

	_, err := r.cb.Execute(func() (interface{}, error) {
		err := r.commandRepository.Remove(entity)
		return nil, err
	})
	return err
}

func (r *repositoryService) RemoveRange(entities []*domain.Entity) error {
	r.Lock()
	defer r.Unlock()

	_, err := r.cb.Execute(func() (interface{}, error) {
		err := r.commandRepository.RemoveRange(entities)
		return nil, err
	})
	return err
}

func (r *repositoryService) Add(entity *domain.Entity) error {
	r.Lock()
	defer r.Unlock()

	_, err := r.cb.Execute(func() (interface{}, error) {
		err := r.commandRepository.Add(entity)
		return nil, err
	})
	return err
}

func (r *repositoryService) AddRange(entities []*domain.Entity) error {
	r.Lock()
	defer r.Unlock()

	_, err := r.cb.Execute(func() (interface{}, error) {
		err := r.commandRepository.AddRange(entities)
		return nil, err
	})
	return err
}

func (r *repositoryService) Update(entity *domain.Entity) error {
	r.Lock()
	defer r.Unlock()

	_, err := r.cb.Execute(func() (interface{}, error) {
		err := r.commandRepository.Update(entity)
		return nil, err
	})
	return err
}

func (r *repositoryService) UpdateRange(entities []*domain.Entity) error {
	r.Lock()
	defer r.Unlock()

	_, err := r.cb.Execute(func() (interface{}, error) {
		err := r.commandRepository.UpdateRange(entities)
		return nil, err
	})
	return err
}
