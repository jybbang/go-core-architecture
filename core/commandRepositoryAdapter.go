package core

import "context"

type CommandRepositoryAdapter interface {
	SetModel(model Entitier)
	Remove(ctx context.Context, entity Entitier) error
	RemoveRange(ctx context.Context, entities []Entitier) error
	Add(ctx context.Context, entity Entitier) error
	AddRange(ctx context.Context, entities []Entitier) error
	Update(ctx context.Context, entity Entitier) error
	UpdateRange(ctx context.Context, entities []Entitier) error
}
