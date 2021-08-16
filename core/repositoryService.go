package core

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/sony/gobreaker"
)

type repositoryService struct {
	model             Entitier
	queryRepository   queryRepositoryAdapter
	commandRepository commandRepositoryAdapter
	cb                *gobreaker.CircuitBreaker
}

func (r *repositoryService) initialize() *repositoryService {
	return r
}

func (r *repositoryService) Find(ctx context.Context, id uuid.UUID, dest Entitier) (err error) {
	_, err = r.cb.Execute(func() (interface{}, error) {
		err = r.queryRepository.Find(ctx, id, dest)
		return nil, err
	})
	return err
}

func (r *repositoryService) Any(ctx context.Context) (ok bool, err error) {
	resp, err := r.cb.Execute(func() (interface{}, error) {
		return r.queryRepository.Any(ctx)
	})
	return resp.(bool), err
}

func (r *repositoryService) AnyWithFilter(ctx context.Context, query interface{}, args interface{}) (ok bool, err error) {
	resp, err := r.cb.Execute(func() (interface{}, error) {
		return r.queryRepository.AnyWithFilter(ctx, query, args)
	})
	return resp.(bool), err
}

func (r *repositoryService) Count(ctx context.Context) (count int64, err error) {
	resp, err := r.cb.Execute(func() (interface{}, error) {
		return r.queryRepository.Count(ctx)
	})
	return resp.(int64), err
}

func (r *repositoryService) CountWithFilter(ctx context.Context, query interface{}, args interface{}) (count int64, err error) {
	resp, err := r.cb.Execute(func() (interface{}, error) {
		return r.queryRepository.CountWithFilter(ctx, query, args)
	})
	return resp.(int64), err
}

func (r *repositoryService) List(ctx context.Context, dest interface{}) (err error) {
	_, err = r.cb.Execute(func() (interface{}, error) {
		err = r.queryRepository.List(ctx, dest)
		return nil, err
	})
	return err
}

func (r *repositoryService) ListWithFilter(ctx context.Context, query interface{}, args interface{}, dest interface{}) (err error) {
	_, err = r.cb.Execute(func() (interface{}, error) {
		err = r.queryRepository.ListWithFilter(ctx, query, args, dest)
		return nil, err
	})
	return err
}

func (r *repositoryService) Remove(ctx context.Context, entity Entitier) error {
	_, err := r.cb.Execute(func() (interface{}, error) {
		err := r.commandRepository.Remove(ctx, entity)
		return nil, err
	})
	return err
}

func (r *repositoryService) RemoveRange(ctx context.Context, entities []Entitier) error {
	_, err := r.cb.Execute(func() (interface{}, error) {
		err := r.commandRepository.RemoveRange(ctx, entities)
		return nil, err
	})
	return err
}

func (r *repositoryService) Add(ctx context.Context, entity Entitier) error {
	entity.SetCreatedAt("", time.Now())
	_, err := r.cb.Execute(func() (interface{}, error) {
		err := r.commandRepository.Add(ctx, entity)
		return nil, err
	})
	return err
}

func (r *repositoryService) AddRange(ctx context.Context, entities []Entitier) error {
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

func (r *repositoryService) Update(ctx context.Context, entity Entitier) error {
	entity.SetUpdatedAt("", time.Now())

	_, err := r.cb.Execute(func() (interface{}, error) {
		err := r.commandRepository.Update(ctx, entity)
		return nil, err
	})
	return err
}

func (r *repositoryService) UpdateRange(ctx context.Context, entities []Entitier) error {
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
