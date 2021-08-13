package contracts

import (
	"github.com/jybbang/go-core-architecture/domain"
)

type CommandRepositoryAdapter interface {
	SetModel(domain.Entity)
	Remove(*domain.Entity) error
	RemoveRange([]*domain.Entity) error
	Add(*domain.Entity) error
	AddRange([]*domain.Entity) error
	Update(*domain.Entity) error
	UpdateRange([]*domain.Entity) error
}
