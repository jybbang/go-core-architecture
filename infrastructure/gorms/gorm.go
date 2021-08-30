package gorms

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/jybbang/go-core-architecture/core"
	"gorm.io/gorm"
)

type adapter struct {
	tableName string
	model     core.Entitier
	dialector gorm.Dialector
	client    *clientProxy
	settings  GormSettings
}

type clientProxy struct {
	db          *gorm.DB
	isConnected bool
}

type clients struct {
	clients map[string]*clientProxy
	mutex   sync.Mutex
}

type GormSettings struct {
	ConnectionString string
	CanCreateTable   bool
}

var clientsSync sync.Once

var clientsInstance *clients

func getClients() *clients {
	if clientsInstance == nil {
		clientsSync.Do(
			func() {
				clientsInstance = &clients{
					clients: make(map[string]*clientProxy),
				}
			})
	}
	return clientsInstance
}

func (a *adapter) migration() {
	if !a.client.db.Migrator().HasTable(a.tableName) && a.settings.CanCreateTable {
		if !a.client.db.Migrator().HasTable(a.model) {
			a.client.db.Migrator().CreateTable(a.model)
		}

		a.client.db.Migrator().RenameTable(a.model, a.tableName)
	}
}

func (a *adapter) IsConnected() bool {
	return a.client.isConnected
}

func (a *adapter) Connect(ctx context.Context) error {
	clientsInstance := getClients()

	clientsInstance.mutex.Lock()
	defer clientsInstance.mutex.Unlock()

	connectionString := a.settings.ConnectionString

	if strings.TrimSpace(connectionString) == "" {
		return fmt.Errorf("connectionString is required")
	}

	cli, ok := clientsInstance.clients[connectionString]

	if !ok || !cli.isConnected {
		db, err := gorm.Open(a.dialector, &gorm.Config{})

		if err != nil {
			return err
		}

		tx := db.Session(&gorm.Session{SkipDefaultTransaction: true})

		// Check context cancellation
		if err := ctx.Err(); err != nil {
			return err
		}

		clientsInstance.clients[connectionString] = &clientProxy{
			db:          tx,
			isConnected: true,
		}
	}

	a.client = clientsInstance.clients[connectionString]

	if a.tableName != "" {
		a.migration()
	}

	return nil
}

func (a *adapter) Disconnect() {
	a.client.isConnected = false
}

func (a *adapter) SetModel(model core.Entitier, tableName string) {
	a.model = model
	a.tableName = tableName
}

func (a *adapter) Find(ctx context.Context, id uuid.UUID, dest core.Entitier) error {
	result := a.client.db.WithContext(ctx).Table(a.tableName).Take(dest, id)

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

	result := a.client.db.WithContext(ctx).Table(a.tableName).Count(resp)

	if result.Error != nil {
		return 0, result.Error
	}

	return *resp, nil
}

func (a *adapter) CountWithFilter(ctx context.Context, query interface{}, args interface{}) (count int64, err error) {
	resp := new(int64)

	result := a.client.db.WithContext(ctx).Table(a.tableName).Where(query, args).Count(resp)

	if result.Error != nil {
		return 0, result.Error
	}

	return *resp, nil
}

func (a *adapter) List(ctx context.Context, dest interface{}) error {
	result := a.client.db.WithContext(ctx).Table(a.tableName).Find(dest)

	return result.Error
}

func (a *adapter) ListWithFilter(ctx context.Context, query interface{}, args interface{}, dest interface{}) error {
	result := a.client.db.WithContext(ctx).Table(a.tableName).Where(query, args).Find(dest)

	return result.Error
}

func (a *adapter) Remove(ctx context.Context, id uuid.UUID) error {
	result := a.client.db.WithContext(ctx).Table(a.tableName).Delete(a.model, id)

	return result.Error
}

func (a *adapter) RemoveRange(ctx context.Context, ids []uuid.UUID) error {
	err := a.client.db.WithContext(ctx).Table(a.tableName).Transaction(func(tx *gorm.DB) error {
		for _, id := range ids {
			err := a.client.db.WithContext(ctx).Table(a.tableName).Delete(a.model, id).Error

			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

func (a *adapter) Add(ctx context.Context, entity core.Entitier) error {
	result := a.client.db.WithContext(ctx).Table(a.tableName).Create(entity)

	return result.Error
}

func (a *adapter) AddRange(ctx context.Context, entities []core.Entitier) error {
	err := a.client.db.WithContext(ctx).Table(a.tableName).Transaction(func(tx *gorm.DB) error {
		for _, entity := range entities {
			err := a.client.db.WithContext(ctx).Table(a.tableName).Create(entity).Error

			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

func (a *adapter) Update(ctx context.Context, entity core.Entitier) error {
	result := a.client.db.WithContext(ctx).Table(a.tableName).Updates(entity)

	return result.Error
}

func (a *adapter) UpdateRange(ctx context.Context, entities []core.Entitier) error {
	err := a.client.db.WithContext(ctx).Table(a.tableName).Transaction(func(tx *gorm.DB) error {
		for _, entity := range entities {
			err := a.client.db.WithContext(ctx).Table(a.tableName).Updates(entity).Error

			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}
