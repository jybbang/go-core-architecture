package core

import (
	"github.com/google/uuid"
)

type QueryRepositoryAdapter interface {
	SetModel(Entitier)
	Find(Entitier, uuid.UUID) error
	Any() (bool, error)
	AnyWithFilter(interface{}, interface{}) (bool, error)
	Count() (int64, error)
	CountWithFilter(interface{}, interface{}) (int64, error)
	List([]Entitier) error
	ListWithFilter([]Entitier, interface{}, interface{}) error
}
