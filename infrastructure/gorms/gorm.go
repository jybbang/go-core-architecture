package gorms

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jybbang/go-core-architecture/core"
	"gorm.io/gorm"
)

type adapter struct {
	tableName string
	model     core.Entitier
	dialector gorm.Dialector
	client    *client
	settings  GormSettings
	mutex     sync.Mutex
}

type client struct {
	db       *gorm.DB
	isOpened bool
}

type clients struct {
	clients map[string]*client
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
					clients: make(map[string]*client),
				}
			})
	}
	return clientsInstance
}

func (a *adapter) open(ctx context.Context) {
	clientsInstance := getClients()

	clientsInstance.mutex.Lock()
	defer clientsInstance.mutex.Unlock()

	connectionString := a.settings.ConnectionString

	if strings.TrimSpace(connectionString) == "" {
		panic("connectionString is required")
	}

	cli, ok := clientsInstance.clients[connectionString]
	if !ok || !cli.isOpened {
		db, err := gorm.Open(a.dialector, &gorm.Config{})
		if err != nil {
			panic(err)
		}
		// Check context cancellation
		if err := ctx.Err(); err != nil {
			panic(err)
		}
		tx := db.Session(&gorm.Session{SkipDefaultTransaction: true})

		clientsInstance.clients[connectionString] = &client{
			db:       tx,
			isOpened: true,
		}
	}

	client := clientsInstance.clients[connectionString]

	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.client = client
	if a.model != nil {
		a.SetModel(a.model, a.tableName)
	}
}

func (a *adapter) OnCircuitOpen() {
	a.client.isOpened = false
}

func (a *adapter) Open() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	a.open(ctx)
}

func (a *adapter) Close() {}

func (a *adapter) SetModel(model core.Entitier, tableName string) {
	a.model = model
	a.tableName = tableName

	if !a.client.db.Migrator().HasTable(a.tableName) && a.settings.CanCreateTable {
		if !a.client.db.Migrator().HasTable(model) {
			a.client.db.Migrator().CreateTable(model)
		}
		a.client.db.Migrator().RenameTable(model, a.tableName)
	}
}

func (a *adapter) Find(ctx context.Context, id uuid.UUID, dest core.Entitier) (err error) {
	if !a.client.isOpened {
		a.Open()
	}

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
	if !a.client.isOpened {
		a.Open()
	}

	resp := new(int64)
	result := a.client.db.WithContext(ctx).Table(a.tableName).Count(resp)
	if result.Error != nil {
		return 0, result.Error
	}
	return *resp, nil
}

func (a *adapter) CountWithFilter(ctx context.Context, query interface{}, args interface{}) (count int64, err error) {
	if !a.client.isOpened {
		a.Open()
	}

	resp := new(int64)
	result := a.client.db.WithContext(ctx).Table(a.tableName).Where(query, args).Count(resp)
	if result.Error != nil {
		return 0, result.Error
	}
	return *resp, nil
}

func (a *adapter) List(ctx context.Context, dest interface{}) (err error) {
	if !a.client.isOpened {
		a.Open()
	}

	result := a.client.db.WithContext(ctx).Table(a.tableName).Find(dest)
	return result.Error
}

func (a *adapter) ListWithFilter(ctx context.Context, query interface{}, args interface{}, dest interface{}) (err error) {
	if !a.client.isOpened {
		a.Open()
	}

	result := a.client.db.WithContext(ctx).Table(a.tableName).Where(query, args).Find(dest)
	return result.Error
}

func (a *adapter) Remove(ctx context.Context, id uuid.UUID) error {
	if !a.client.isOpened {
		a.Open()
	}

	result := a.client.db.WithContext(ctx).Table(a.tableName).Delete(a.model, id)
	return result.Error
}

func (a *adapter) RemoveRange(ctx context.Context, ids []uuid.UUID) error {
	if !a.client.isOpened {
		a.Open()
	}

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
	if !a.client.isOpened {
		a.Open()
	}

	result := a.client.db.WithContext(ctx).Table(a.tableName).Create(entity)
	return result.Error
}

func (a *adapter) AddRange(ctx context.Context, entities []core.Entitier) error {
	if !a.client.isOpened {
		a.Open()
	}

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
	if !a.client.isOpened {
		a.Open()
	}

	result := a.client.db.WithContext(ctx).Table(a.tableName).Updates(entity)
	return result.Error
}

func (a *adapter) UpdateRange(ctx context.Context, entities []core.Entitier) error {
	if !a.client.isOpened {
		a.Open()
	}

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
