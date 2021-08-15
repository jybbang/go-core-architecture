package core

import (
	"context"
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

func (r *RepositoryService) Find(ctx context.Context, dest Entitier, id uuid.UUID) (ok bool, err error) {
	resp, err := r.cb.Execute(func() (interface{}, error) {
		return r.queryRepository.Find(ctx, dest, id)
	})
	return resp.(bool), err
}

func (r *RepositoryService) Any(ctx context.Context) (ok bool, err error) {
	resp, err := r.cb.Execute(func() (interface{}, error) {
		return r.queryRepository.Any(ctx)
	})
	return resp.(bool), err
}

func (r *RepositoryService) AnyWithFilter(ctx context.Context, query interface{}, args interface{}) (ok bool, err error) {
	resp, err := r.cb.Execute(func() (interface{}, error) {
		return r.queryRepository.AnyWithFilter(ctx, query, args)
	})
	return resp.(bool), err
}

func (r *RepositoryService) Count(ctx context.Context) (count int64, err error) {
	resp, err := r.cb.Execute(func() (interface{}, error) {
		return r.queryRepository.Count(ctx)
	})
	return resp.(int64), err
}

func (r *RepositoryService) CountWithFilter(ctx context.Context, query interface{}, args interface{}) (count int64, err error) {
	resp, err := r.cb.Execute(func() (interface{}, error) {
		return r.queryRepository.CountWithFilter(ctx, query, args)
	})
	return resp.(int64), err
}

func (r *RepositoryService) List(ctx context.Context, dest []Entitier) error {
	_, err := r.cb.Execute(func() (interface{}, error) {
		err := r.queryRepository.List(ctx, dest)
		return nil, err
	})
	return err
}

func (r *RepositoryService) ListWithFilter(ctx context.Context, dest []Entitier, query interface{}, args interface{}) error {
	_, err := r.cb.Execute(func() (interface{}, error) {
		err := r.queryRepository.ListWithFilter(ctx, dest, query, args)
		return nil, err
	})
	return err
}

func (r *RepositoryService) Remove(ctx context.Context, entity Entitier) error {
	r.Lock()
	defer r.Unlock()

	_, err := r.cb.Execute(func() (interface{}, error) {
		err := r.commandRepository.Remove(ctx, entity)
		return nil, err
	})
	return err
}

func (r *RepositoryService) RemoveRange(ctx context.Context, entities []Entitier) error {
	r.Lock()
	defer r.Unlock()

	_, err := r.cb.Execute(func() (interface{}, error) {
		err := r.commandRepository.RemoveRange(ctx, entities)
		return nil, err
	})
	return err
}

func (r *RepositoryService) Add(ctx context.Context, entity Entitier) error {
	r.Lock()
	defer r.Unlock()

	entity.SetCreatedAt("", time.Now())
	_, err := r.cb.Execute(func() (interface{}, error) {
		err := r.commandRepository.Add(ctx, entity)
		return nil, err
	})
	return err
}

func (r *RepositoryService) AddRange(ctx context.Context, entities []Entitier) error {
	r.Lock()
	defer r.Unlock()

	user := ""
	now := time.Now()
	for _, v := range entities {
		v.SetCreatedAt(user, now)
	}

	_, err := r.cb.Execute(func() (interface{}, error) {
		err := r.commandRepository.AddRange(ctx, entities)
		return nil, err
	})
	return err
}

func (r *RepositoryService) Update(ctx context.Context, entity Entitier) error {
	r.Lock()
	defer r.Unlock()

	entity.SetUpdatedAt("", time.Now())

	_, err := r.cb.Execute(func() (interface{}, error) {
		err := r.commandRepository.Update(ctx, entity)
		return nil, err
	})
	return err
}

func (r *RepositoryService) UpdateRange(ctx context.Context, entities []Entitier) error {
	r.Lock()
	defer r.Unlock()

	user := ""
	now := time.Now()
	for _, v := range entities {
		v.SetUpdatedAt(user, now)
	}

	_, err := r.cb.Execute(func() (interface{}, error) {
		err := r.commandRepository.UpdateRange(ctx, entities)
		return nil, err
	})
	return err
}
