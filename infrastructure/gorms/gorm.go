package gorms

import (
	"sync"

	"github.com/google/uuid"
	"github.com/jybbang/go-core-architecture/core"
	"gorm.io/gorm"
)

type adapter struct {
	conn  *gorm.DB
	model core.Entitier
}

type clients struct {
	clients map[string]*gorm.DB
	sync.Mutex
}

var clientsSync sync.Once

var clientsInstance *clients

func getClients() *clients {
	if clientsInstance == nil {
		clientsSync.Do(
			func() {
				clientsInstance = &clients{
					clients: make(map[string]*gorm.DB),
				}
			})
	}
	return clientsInstance
}

func (a *adapter) SetModel(model core.Entitier) {
	a.model = model
}

func (a *adapter) Find(dto core.Entitier, id uuid.UUID) error {
	a.conn.Take(dto, id)
	if dto == nil {
		return core.ErrNotFound
	}

	return nil
}

func (a *adapter) Any() (bool, error) {
	count, err := a.Count()
	return count > 0, err
}

func (a *adapter) AnyWithFilter(query interface{}, args interface{}) (bool, error) {
	count, err := a.CountWithFilter(query, args)
	return count > 0, err
}

func (a *adapter) Count() (int64, error) {
	count := new(int64)
	a.conn.Model(a.model).Count(count)

	return *count, nil
}

func (a *adapter) CountWithFilter(query interface{}, args interface{}) (int64, error) {
	count := new(int64)
	a.conn.Model(a.model).Count(count).Where(query, args)

	return *count, nil
}

func (a *adapter) List(dtos []core.Entitier) error {
	a.conn.Find(dtos)
	if dtos == nil {
		return core.ErrNotFound
	}

	return nil
}

func (a *adapter) ListWithFilter(dtos []core.Entitier, query interface{}, args interface{}) error {
	a.conn.Find(dtos).Where(query, args)
	if dtos == nil {
		return core.ErrNotFound
	}

	return nil
}

func (a *adapter) Remove(entity core.Entitier) error {
	a.conn.Delete(entity, entity.GetID())
	return nil
}

func (a *adapter) RemoveRange(entities []core.Entitier) error {
	a.conn.Delete(entities)
	return nil
}

func (a *adapter) Add(entity core.Entitier) error {
	a.conn.Create(entity)
	return nil
}

func (a *adapter) AddRange(entities []core.Entitier) error {
	a.conn.Create(entities)
	return nil
}

func (a *adapter) Update(entity core.Entitier) error {
	a.conn.Model(entity).Updates(entity)
	return nil
}

func (a *adapter) UpdateRange(entities []core.Entitier) error {
	a.conn.Model(a.model).Updates(entities)
	return nil
}
