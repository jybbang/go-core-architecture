package leveldb

import (
	"context"
	"encoding/json"
	"errors"
	"sync"

	"github.com/jybbang/go-core-architecture/core"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

type adapter struct {
	leveldb  *leveldb.DB
	settings LevelDbSettings
}

type clients struct {
	clients map[string]*leveldb.DB
	mutex   sync.Mutex
}

type LevelDbSettings struct {
	Path     string
	ReadOnly bool
}

var clientsSync sync.Once

var clientsInstance *clients

func getClients() *clients {
	if clientsInstance == nil {
		clientsSync.Do(
			func() {
				clientsInstance = &clients{
					clients: make(map[string]*leveldb.DB),
				}
			})
	}
	return clientsInstance
}

func getLevelDbClient(ctx context.Context, settings LevelDbSettings) *leveldb.DB {
	clientsInstance := getClients()

	clientsInstance.mutex.Lock()
	defer clientsInstance.mutex.Unlock()

	path := settings.Path
	_, ok := clientsInstance.clients[path]
	if !ok {
		leveldbClient, err := leveldb.OpenFile(path, &opt.Options{
			ReadOnly: settings.ReadOnly,
		})
		if err != nil {
			panic(err)
		}
		// Check context cancellation
		if err := ctx.Err(); err != nil {
			panic(err)
		}
		clientsInstance.clients[path] = leveldbClient
	}

	client := clientsInstance.clients[path]
	return client
}

func NewLevelDbAdapter(ctx context.Context, settings LevelDbSettings) *adapter {
	client := getLevelDbClient(ctx, settings)
	leveldbService := &adapter{
		leveldb:  client,
		settings: settings,
	}

	return leveldbService
}

func (a *adapter) Close() {
	a.leveldb.Close()
}

func (a *adapter) Has(ctx context.Context, key string) bool {
	ok, err := a.leveldb.Has([]byte(key), nil)
	if err != nil {
		return false
	}
	return ok
}

func (a *adapter) Get(ctx context.Context, key string, dest interface{}) error {
	value, err := a.leveldb.Get([]byte(key), nil)
	if err != nil {
		if errors.Is(err, leveldb.ErrNotFound) {
			return core.ErrNotFound
		}
	}
	return json.Unmarshal(value, dest)
}

func (a *adapter) Set(ctx context.Context, key string, value interface{}) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return a.leveldb.Put([]byte(key), bytes, nil)
}

func (a *adapter) Delete(ctx context.Context, key string) error {
	return a.leveldb.Delete([]byte(key), nil)
}
