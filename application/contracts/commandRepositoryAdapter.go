package contracts

import "github.com/jybbang/go-core-architecture/domain"

type CommandRepositoryAdapter interface {
	SetModel(domain.Entitier)
	Remove(domain.Entitier) error
	RemoveRange([]domain.Entitier) error
	Add(domain.Entitier) error
	AddRange([]domain.Entitier) error
	Update(domain.Entitier) error
	UpdateRange([]domain.Entitier) error
}
