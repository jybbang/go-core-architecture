package core

import (
	"context"

	"github.com/google/uuid"
)

type queryRepositoryAdapter interface {
	SetModel(model Entitier)
	Find(ctx context.Context, id uuid.UUID, dest Entitier) (err error)
	Any(ctx context.Context) (ok bool, err error)
	AnyWithFilter(ctx context.Context, query interface{}, args interface{}) (ok bool, err error)
	Count(ctx context.Context) (count int64, err error)
	CountWithFilter(ctx context.Context, query interface{}, args interface{}) (count int64, err error)
	List(ctx context.Context) (result []Entitier, err error)
	ListWithFilter(ctx context.Context, query interface{}, args interface{}) (result []Entitier, err error)
}
