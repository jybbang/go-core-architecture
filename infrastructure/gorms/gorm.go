package gorms

import (
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/jybbang/go-core-architecture/core"
	"gorm.io/gorm"
)

type adapter struct {
	tableName string
	model     core.Entitier
	conn      *gorm.DB
	settings  GormSettings
}

type GormSettings struct {
	ConnectionString string
	CanCreateTable   bool
}

type clients struct {
	clients map[string]*gorm.DB
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
				}
			})
	}
	return clientsInstance
}

func (a *adapter) SetModel(model core.Entitier, tableName string) {
	a.model = model
	a.tableName = tableName

	if !a.conn.Migrator().HasTable(a.tableName) && a.settings.CanCreateTable {
		if !a.conn.Migrator().HasTable(model) {
			a.conn.Migrator().CreateTable(model)
		}
		a.conn.Migrator().RenameTable(model, a.tableName)
	}
}

func (a *adapter) Find(ctx context.Context, id uuid.UUID, dest core.Entitier) (err error) {
	result := a.conn.WithContext(ctx).Table(a.tableName).Take(dest, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return core.ErrNotFound
	} else if result.Error != nil {
		return result.Error
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
	resp := new(int64)
	result := a.conn.WithContext(ctx).Table(a.tableName).Count(resp)
	if result.Error != nil {
		return 0, result.Error
	}
	return *resp, nil
}

func (a *adapter) CountWithFilter(ctx context.Context, query interface{}, args interface{}) (count int64, err error) {
	resp := new(int64)
	result := a.conn.WithContext(ctx).Table(a.tableName).Where(query, args).Count(resp)
	if result.Error != nil {
		return 0, result.Error
	}
	return *resp, nil
}

func (a *adapter) List(ctx context.Context, dest interface{}) (err error) {
	result := a.conn.WithContext(ctx).Table(a.tableName).Find(dest)
	return result.Error
}

func (a *adapter) ListWithFilter(ctx context.Context, query interface{}, args interface{}, dest interface{}) (err error) {
	result := a.conn.WithContext(ctx).Table(a.tableName).Where(query, args).Find(dest)
	return result.Error
}

func (a *adapter) Remove(ctx context.Context, id uuid.UUID) error {
	result := a.conn.WithContext(ctx).Table(a.tableName).Delete(a.model, id)
	return result.Error
}

func (a *adapter) RemoveRange(ctx context.Context, ids []uuid.UUID) error {
	err := a.conn.WithContext(ctx).Table(a.tableName).Transaction(func(tx *gorm.DB) error {
		for _, id := range ids {
			err := a.conn.WithContext(ctx).Table(a.tableName).Delete(a.model, id).Error
			if err != nil {
				return err
			}
		}
		return nil
	})

	return err
}

func (a *adapter) Add(ctx context.Context, entity core.Entitier) error {
	result := a.conn.WithContext(ctx).Table(a.tableName).Create(entity)
	return result.Error
}

func (a *adapter) AddRange(ctx context.Context, entities []core.Entitier) error {
	// result := a.conn.WithContext(ctx).Table(a.tableName).CreateInBatches(entities, 1000)
	err := a.conn.WithContext(ctx).Table(a.tableName).Transaction(func(tx *gorm.DB) error {
		for _, entity := range entities {
			err := a.conn.WithContext(ctx).Table(a.tableName).Create(entity).Error
			if err != nil {
				return err
			}
		}
		return nil
	})

	return err
}

func (a *adapter) Update(ctx context.Context, entity core.Entitier) error {
	result := a.conn.WithContext(ctx).Table(a.tableName).Updates(entity)
	return result.Error
}

func (a *adapter) UpdateRange(ctx context.Context, entities []core.Entitier) error {
	err := a.conn.WithContext(ctx).Table(a.tableName).Transaction(func(tx *gorm.DB) error {
		for _, entity := range entities {
			err := a.conn.WithContext(ctx).Table(a.tableName).Updates(entity).Error
			if err != nil {
				return err
			}
		}
		return nil
	})

	return err
}
