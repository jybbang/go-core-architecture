package core

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sony/gobreaker"
)

type repositoryService struct {
	queryRepository   queryRepositoryAdapter
	commandRepository commandRepositoryAdapter
	cb                *gobreaker.CircuitBreaker
}

func (r *repositoryService) initialize() *repositoryService {
	return r
}

func (r *repositoryService) Find(ctx context.Context, id uuid.UUID, dest Entitier) Result {
	if dest == nil {
		return Result{E: fmt.Errorf("%w dest is required", ErrInternalServerError)}
	}

	_, err := r.cb.Execute(func() (interface{}, error) {
		err := r.queryRepository.Find(ctx, id, dest)
		return nil, err
	})

	return Result{V: dest, E: err}
}

func (r *repositoryService) Any(ctx context.Context) Result {
	resp, err := r.cb.Execute(func() (interface{}, error) {
		return r.queryRepository.Any(ctx)
	})
	if err != nil {
		return Result{V: false, E: err}
	}

	return Result{V: resp.(bool), E: err}
}

func (r *repositoryService) AnyWithFilter(ctx context.Context, query interface{}, args interface{}) Result {
	resp, err := r.cb.Execute(func() (interface{}, error) {
		return r.queryRepository.AnyWithFilter(ctx, query, args)
	})
	if err != nil {
		return Result{V: false, E: err}
	}

	return Result{V: resp.(bool), E: err}
}

func (r *repositoryService) Count(ctx context.Context) Result {
	resp, err := r.cb.Execute(func() (interface{}, error) {
		return r.queryRepository.Count(ctx)
	})
	if err != nil {
		return Result{V: int64(0), E: err}
	}

	return Result{V: resp.(int64), E: err}
}

func (r *repositoryService) CountWithFilter(ctx context.Context, query interface{}, args interface{}) Result {
	resp, err := r.cb.Execute(func() (interface{}, error) {
		return r.queryRepository.CountWithFilter(ctx, query, args)
	})
	if err != nil {
		return Result{V: int64(0), E: err}
	}

	return Result{V: resp.(int64), E: err}
}

func (r *repositoryService) List(ctx context.Context, dest interface{}) Result {
	if dest == nil {
		return Result{E: fmt.Errorf("%w dest is required", ErrInternalServerError)}
	}

	_, err := r.cb.Execute(func() (interface{}, error) {
		err := r.queryRepository.List(ctx, dest)
		return nil, err
	})

	return Result{V: dest, E: err}
}

func (r *repositoryService) ListWithFilter(ctx context.Context, query interface{}, args interface{}, dest interface{}) Result {
	if dest == nil {
		return Result{E: fmt.Errorf("%w dest is required", ErrInternalServerError)}
	}

	_, err := r.cb.Execute(func() (interface{}, error) {
		err := r.queryRepository.ListWithFilter(ctx, query, args, dest)
		return nil, err
	})

	return Result{V: dest, E: err}
}

func (r *repositoryService) Remove(ctx context.Context, entity Entitier) Result {
	if entity == nil {
		return Result{E: fmt.Errorf("%w entity is required", ErrInternalServerError)}
	}

	_, err := r.cb.Execute(func() (interface{}, error) {
		err := r.commandRepository.Remove(ctx, entity)
		return nil, err
	})

	return Result{V: nil, E: err}
}

func (r *repositoryService) RemoveRange(ctx context.Context, entities []Entitier) Result {
	if len(entities) == 0 {
		return Result{E: fmt.Errorf("%w entities is required", ErrInternalServerError)}
	}

	_, err := r.cb.Execute(func() (interface{}, error) {
		err := r.commandRepository.RemoveRange(ctx, entities)
		return nil, err
	})

	return Result{V: nil, E: err}
}

func (r *repositoryService) Add(ctx context.Context, entity Entitier) Result {
	if entity == nil {
		return Result{E: fmt.Errorf("%w entity is required", ErrInternalServerError)}
	}

	user, _ := ctx.Value("userId").(string)
	entity.SetCreatedAt(user, time.Now())
	_, err := r.cb.Execute(func() (interface{}, error) {
		err := r.commandRepository.Add(ctx, entity)
		return nil, err
	})

	return Result{V: nil, E: err}
}

func (r *repositoryService) AddRange(ctx context.Context, entities []Entitier) Result {
	if len(entities) == 0 {
		return Result{E: fmt.Errorf("%w entities is required", ErrInternalServerError)}
	}

	user, _ := ctx.Value("userId").(string)
	now := time.Now()
	for _, v := range entities {
		v.SetCreatedAt(user, now)
	}

	_, err := r.cb.Execute(func() (interface{}, error) {
		err := r.commandRepository.AddRange(ctx, entities)
		return nil, err
	})

	return Result{V: nil, E: err}
}

func (r *repositoryService) Update(ctx context.Context, entity Entitier) Result {
	if entity == nil {
		return Result{E: fmt.Errorf("%w entity is required", ErrInternalServerError)}
	}

	user, _ := ctx.Value("userId").(string)
	entity.SetUpdatedAt(user, time.Now())

	_, err := r.cb.Execute(func() (interface{}, error) {
		err := r.commandRepository.Update(ctx, entity)
		return nil, err
	})

	return Result{V: nil, E: err}
}

func (r *repositoryService) UpdateRange(ctx context.Context, entities []Entitier) Result {
	if len(entities) == 0 {
		return Result{E: fmt.Errorf("%w entities is required", ErrInternalServerError)}
	}

	user, _ := ctx.Value("userId").(string)
	now := time.Now()
	for _, v := range entities {
		v.SetUpdatedAt(user, now)
	}

	_, err := r.cb.Execute(func() (interface{}, error) {
		err := r.commandRepository.UpdateRange(ctx, entities)
		return nil, err
	})

	return Result{V: nil, E: err}
}
