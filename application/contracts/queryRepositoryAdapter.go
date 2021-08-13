package contracts

import (
	"github.com/google/uuid"
	"github.com/jybbang/go-core-architecture/domain"
)

type QueryRepositoryAdapter interface {
	SetModel(domain.Entitier)
	Find(domain.Entitier, uuid.UUID) error
	Any() (bool, error)
	AnyWithFilter(interface{}, interface{}) (bool, error)
	Count() (int64, error)
	CountWithFilter(interface{}, interface{}) (int64, error)
	List([]domain.Entitier) error
	ListWithFilter([]domain.Entitier, interface{}, interface{}) error
}
