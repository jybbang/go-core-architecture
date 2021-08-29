package core

import (
	"context"

	"github.com/google/uuid"
)

type commandRepositoryAdapter interface {
	OnCircuitOpen()
	Open()
	Close()
	SetModel(model Entitier, tableName string)
	Remove(ctx context.Context, id uuid.UUID) error
	RemoveRange(ctx context.Context, ids []uuid.UUID) error
	Add(ctx context.Context, entity Entitier) error
	AddRange(ctx context.Context, entities []Entitier) error
	Update(ctx context.Context, entity Entitier) error
	UpdateRange(ctx context.Context, entities []Entitier) error
}
