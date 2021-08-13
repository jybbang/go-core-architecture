package gorms

import (
	"sync"

	"github.com/google/uuid"
	"github.com/jybbang/go-core-architecture/domain"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type adapter struct {
	conn  *gorm.DB
	model domain.Entitier
}

type clients struct {
	clients map[string]*gorm.DB
	sync.Mutex
}

var log *zap.SugaredLogger

var clientsSync sync.Once

var clientsInstance *clients

func init() {
	logger, _ := zap.NewProduction()
	log = logger.Sugar()
}

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

func (a *adapter) SetModel(model domain.Entitier) {
	a.model = model
}

func (a *adapter) Find(dto domain.Entitier, id uuid.UUID) error {
	a.conn.Take(dto, id)
	if dto == nil {
		return domain.ErrNotFound
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

func (a *adapter) List(dtos []domain.Entitier) error {
	a.conn.Find(dtos)
	if dtos == nil {
		return domain.ErrNotFound
	}

	return nil
}

func (a *adapter) ListWithFilter(dtos []domain.Entitier, query interface{}, args interface{}) error {
	a.conn.Find(dtos).Where(query, args)
	if dtos == nil {
		return domain.ErrNotFound
	}

	return nil
}

func (a *adapter) Remove(entity domain.Entitier) error {
	a.conn.Delete(entity, entity.GetID())
	return nil
}

func (a *adapter) RemoveRange(entities []domain.Entitier) error {
	a.conn.Delete(entities)
	return nil
}

func (a *adapter) Add(entity domain.Entitier) error {
	a.conn.Create(entity)
	return nil
}

func (a *adapter) AddRange(entities []domain.Entitier) error {
	a.conn.Create(entities)
	return nil
}

func (a *adapter) Update(entity domain.Entitier) error {
	a.conn.Model(entity).Updates(entity)
	return nil
}

func (a *adapter) UpdateRange(entities []domain.Entitier) error {
	a.conn.Model(a.model).Updates(entities)
	return nil
}
