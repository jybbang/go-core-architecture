package gorms

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/jybbang/go-core-architecture/core"
	"gorm.io/gorm"
)

type adapter struct {
	model    core.Entitier
	conn     *gorm.DB
	rw       *sync.RWMutex
	settings GormSettings
}

type GormSettings struct {
	ConnectionString string
	CanCreateTable   bool
}

type clients struct {
	clients map[string]*gorm.DB
	mutexes map[string]*sync.RWMutex
	mutex   sync.Mutex
}

var clientsSync sync.Once

var clientsInstance *clients

func getClients() *clients {
	if clientsInstance == nil {
		clientsSync.Do(
			func() {
				clientsInstance = &clients{
					clients: make(map[string]*gorm.DB),
					mutexes: make(map[string]*sync.RWMutex),
				}
			})
	}
	return clientsInstance
}

func (a *adapter) SetModel(model core.Entitier) {
	a.model = model

	if !a.conn.Migrator().HasTable(model) && a.settings.CanCreateTable {
		a.conn.Migrator().CreateTable(model)
	}
}

func (a *adapter) Find(ctx context.Context, id uuid.UUID, dest core.Entitier) (err error) {
	a.rw.RLock()
	defer a.rw.RUnlock()

	a.conn.WithContext(ctx).Take(dest, id)
	if dest == nil {
		return core.ErrNotFound
	}
	return nil
}

func (a *adapter) Any(ctx context.Context) (ok bool, err error) {
	count, err := a.Count(ctx)
	return count > 0, err
}

func (a *adapter) AnyWithFilter(ctx context.Context, query interface{}, args interface{}) (ok bool, err error) {
	count, err := a.CountWithFilter(ctx, query, args)
	return count > 0, err
}

func (a *adapter) Count(ctx context.Context) (count int64, err error) {
	a.rw.RLock()
	defer a.rw.RUnlock()

	resp := new(int64)
	a.conn.WithContext(ctx).Model(a.model).Count(resp)

	return *resp, nil
}

func (a *adapter) CountWithFilter(ctx context.Context, query interface{}, args interface{}) (count int64, err error) {
	a.rw.RLock()
	defer a.rw.RUnlock()

	resp := new(int64)
	a.conn.WithContext(ctx).Model(a.model).Count(resp).Where(query, args)

	return *resp, nil
}

func (a *adapter) List(ctx context.Context, dest interface{}) (err error) {
	a.rw.RLock()
	defer a.rw.RUnlock()

	a.conn.WithContext(ctx).Find(dest)

	return nil
}

func (a *adapter) ListWithFilter(ctx context.Context, query interface{}, args interface{}, dest interface{}) (err error) {
	a.rw.RLock()
	defer a.rw.RUnlock()

	a.conn.WithContext(ctx).Find(dest).Where(query, args)

	return nil
}

func (a *adapter) Remove(ctx context.Context, entity core.Entitier) error {
	a.rw.Lock()
	defer a.rw.Unlock()

	a.conn.WithContext(ctx).Delete(entity, entity.GetID())
	return nil
}

func (a *adapter) RemoveRange(ctx context.Context, entities []core.Entitier) error {
	a.rw.Lock()
	defer a.rw.Unlock()

	a.conn.WithContext(ctx).Delete(entities)
	return nil
}

func (a *adapter) Add(ctx context.Context, entity core.Entitier) error {
	a.rw.Lock()
	defer a.rw.Unlock()

	a.conn.WithContext(ctx).Create(entity)
	return nil
}

func (a *adapter) AddRange(ctx context.Context, entities []core.Entitier) error {
	a.rw.Lock()
	defer a.rw.Unlock()

	a.conn.WithContext(ctx).Create(entities)
	return nil
}

func (a *adapter) Update(ctx context.Context, entity core.Entitier) error {
	a.rw.Lock()
	defer a.rw.Unlock()

	a.conn.WithContext(ctx).Model(entity).Updates(entity)
	return nil
}

func (a *adapter) UpdateRange(ctx context.Context, entities []core.Entitier) error {
	a.rw.Lock()
	defer a.rw.Unlock()

	a.conn.WithContext(ctx).Model(a.model).Updates(entities)
	return nil
}
