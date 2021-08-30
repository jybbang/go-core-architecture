package core

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/sony/gobreaker"
)

type repositoryService struct {
	tableName         string
	userIdKey         string
	queryRepository   queryRepositoryAdapter
	commandRepository commandRepositoryAdapter
	cb                *gobreaker.CircuitBreaker
	settings          RepositoryServiceSettings
}

func (r *repositoryService) initialize() *repositoryService {
	if err := r.queryRepositoryConnect(); err != nil {
		panic(err)
	}

	if err := r.commandRepositoryConnect(); err != nil {
		panic(err)
	}

	return r
}

func (r *repositoryService) queryRepositoryConnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), r.settings.ConnectionTimeout)
	defer cancel()

	return r.queryRepository.Connect(ctx)
}

func (r *repositoryService) commandRepositoryConnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), r.settings.ConnectionTimeout)
	defer cancel()

	return r.commandRepository.Connect(ctx)
}

func (r *repositoryService) onCircuitOpen() {
	r.queryRepository.Disconnect()

	r.commandRepository.Disconnect()

	if !r.queryRepository.IsConnected() {
		r.queryRepositoryConnect()
	}

	if !r.commandRepository.IsConnected() {
		r.commandRepositoryConnect()
	}
}

func (r *repositoryService) Find(ctx context.Context, id uuid.UUID, dest Entitier) Result {
	if dest == nil {
		return Result{E: fmt.Errorf("%w dest is required", ErrInternalServerError)}
	}

	span, ctx := opentracing.StartSpanFromContext(ctx, "Repository:Find")

	if span != nil {
		defer span.Finish()
	}

	_, err := r.cb.Execute(func() (interface{}, error) {
		return nil, r.queryRepository.Find(ctx, id, dest)
	})

	return Result{V: dest, E: err}
}

func (r *repositoryService) Any(ctx context.Context) Result {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Repository:Any")

	if span != nil {
		defer span.Finish()
	}

	resp, err := r.cb.Execute(func() (interface{}, error) {
		return r.queryRepository.Any(ctx)
	})

	if err != nil {
		return Result{V: false, E: err}
	}

	return Result{V: resp.(bool), E: err}
}

func (r *repositoryService) AnyWithFilter(ctx context.Context, query interface{}, args interface{}) Result {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Repository:AnyWithFilter")

	if span != nil {
		defer span.Finish()
	}

	resp, err := r.cb.Execute(func() (interface{}, error) {
		return r.queryRepository.AnyWithFilter(ctx, query, args)
	})

	if err != nil {
		return Result{V: false, E: err}
	}

	return Result{V: resp.(bool), E: err}
}

func (r *repositoryService) Count(ctx context.Context) Result {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Repository:Count")

	if span != nil {
		defer span.Finish()
	}

	resp, err := r.cb.Execute(func() (interface{}, error) {
		return r.queryRepository.Count(ctx)
	})

	if err != nil {
		return Result{V: int64(0), E: err}
	}

	return Result{V: resp.(int64), E: err}
}

func (r *repositoryService) CountWithFilter(ctx context.Context, query interface{}, args interface{}) Result {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Repository:CountWithFilter")

	if span != nil {
		defer span.Finish()
	}

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

	resultsVal := reflect.ValueOf(dest)

	if resultsVal.Kind() != reflect.Ptr {
		panic("dest must be a pointer to a slice")
	}

	sliceVal := resultsVal.Elem()

	if sliceVal.Kind() == reflect.Interface {
		sliceVal = sliceVal.Elem()
	}

	if sliceVal.Kind() != reflect.Slice {
		panic("dest must be a pointer to a slice")
	}

	span, ctx := opentracing.StartSpanFromContext(ctx, "Repository:List")

	if span != nil {
		defer span.Finish()
	}

	_, err := r.cb.Execute(func() (interface{}, error) {
		return nil, r.queryRepository.List(ctx, dest)
	})

	return Result{V: dest, E: err}
}

func (r *repositoryService) ListWithFilter(ctx context.Context, query interface{}, args interface{}, dest interface{}) Result {
	if dest == nil {
		return Result{E: fmt.Errorf("%w dest is required", ErrInternalServerError)}
	}

	resultsVal := reflect.ValueOf(dest)

	if resultsVal.Kind() != reflect.Ptr {
		panic("dest must be a pointer to a slice")
	}

	sliceVal := resultsVal.Elem()

	if sliceVal.Kind() == reflect.Interface {
		sliceVal = sliceVal.Elem()
	}

	if sliceVal.Kind() != reflect.Slice {
		panic("dest must be a pointer to a slice")
	}

	span, ctx := opentracing.StartSpanFromContext(ctx, "Repository:ListWithFilter")

	if span != nil {
		defer span.Finish()
	}

	_, err := r.cb.Execute(func() (interface{}, error) {
		return nil, r.queryRepository.ListWithFilter(ctx, query, args, dest)
	})

	return Result{V: dest, E: err}
}

func (r *repositoryService) Remove(ctx context.Context, id uuid.UUID) Result {
	if id == uuid.Nil {
		return Result{E: fmt.Errorf("%w id is required", ErrInternalServerError)}
	}

	span, ctx := opentracing.StartSpanFromContext(ctx, "Repository:Remove")

	if span != nil {
		defer span.Finish()
	}

	_, err := r.cb.Execute(func() (interface{}, error) {
		return nil, r.commandRepository.Remove(ctx, id)
	})

	return Result{V: nil, E: err}
}

func (r *repositoryService) RemoveRange(ctx context.Context, ids []uuid.UUID) Result {
	if len(ids) == 0 {
		return Result{E: fmt.Errorf("%w ids is required", ErrInternalServerError)}
	}

	span, ctx := opentracing.StartSpanFromContext(ctx, "Repository:RemoveRange")

	if span != nil {
		defer span.Finish()
	}

	_, err := r.cb.Execute(func() (interface{}, error) {
		return nil, r.commandRepository.RemoveRange(ctx, ids)
	})

	return Result{V: nil, E: err}
}

func (r *repositoryService) Add(ctx context.Context, entity Entitier) Result {
	if entity == nil {
		return Result{E: fmt.Errorf("%w entity is required", ErrInternalServerError)}
	}

	span, ctx := opentracing.StartSpanFromContext(ctx, "Repository:Add")

	if span != nil {
		defer span.Finish()
	}

	_, err := r.cb.Execute(func() (interface{}, error) {
		user, _ := ctx.Value(r.userIdKey).(string)

		entity.SetCreatedAt(user, time.Now())

		return nil, r.commandRepository.Add(ctx, entity)
	})

	return Result{V: nil, E: err}
}

func (r *repositoryService) AddRange(ctx context.Context, entities []Entitier) Result {
	if len(entities) == 0 {
		return Result{E: fmt.Errorf("%w entities is required", ErrInternalServerError)}
	}

	span, ctx := opentracing.StartSpanFromContext(ctx, "Repository:AddRange")

	if span != nil {
		defer span.Finish()
	}

	_, err := r.cb.Execute(func() (interface{}, error) {
		user, _ := ctx.Value(r.userIdKey).(string)

		now := time.Now()

		for _, v := range entities {
			v.SetCreatedAt(user, now)
		}

		return nil, r.commandRepository.AddRange(ctx, entities)
	})

	return Result{V: nil, E: err}
}

func (r *repositoryService) Update(ctx context.Context, entity Entitier) Result {
	if entity == nil {
		return Result{E: fmt.Errorf("%w entity is required", ErrInternalServerError)}
	}

	span, ctx := opentracing.StartSpanFromContext(ctx, "Repository:Update")

	if span != nil {
		defer span.Finish()
	}

	_, err := r.cb.Execute(func() (interface{}, error) {
		user, _ := ctx.Value(r.userIdKey).(string)

		entity.SetUpdatedAt(user, time.Now())

		return nil, r.commandRepository.Update(ctx, entity)
	})

	return Result{V: nil, E: err}
}

func (r *repositoryService) UpdateRange(ctx context.Context, entities []Entitier) Result {
	if len(entities) == 0 {
		return Result{E: fmt.Errorf("%w entities is required", ErrInternalServerError)}
	}

	span, ctx := opentracing.StartSpanFromContext(ctx, "Repository:UpdateRange")

	if span != nil {
		defer span.Finish()
	}

	_, err := r.cb.Execute(func() (interface{}, error) {
		user, _ := ctx.Value(r.userIdKey).(string)

		now := time.Now()

		for _, v := range entities {
			v.SetUpdatedAt(user, now)
		}

		return nil, r.commandRepository.UpdateRange(ctx, entities)
	})

	return Result{V: nil, E: err}
}
