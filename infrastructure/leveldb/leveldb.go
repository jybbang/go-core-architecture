package leveldb

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/jybbang/go-core-architecture/core"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

type adapter struct {
	client   *client
	settings LevelDbSettings
	mutex    sync.Mutex
}

type client struct {
	leveldb  *leveldb.DB
	isOpened bool
}

type clients struct {
	clients map[string]*client
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
					clients: make(map[string]*client),
				}
			})
	}
	return clientsInstance
}

func NewLevelDbAdapter(ctx context.Context, settings LevelDbSettings) *adapter {
	leveldbService := &adapter{
		settings: settings,
	}

	leveldbService.setClient(ctx)
	return leveldbService
}

func (a *adapter) setClient(ctx context.Context) {
	clientsInstance := getClients()

	clientsInstance.mutex.Lock()
	defer clientsInstance.mutex.Unlock()

	path := a.settings.Path

	if strings.TrimSpace(path) == "" {
		panic("path is required")
	}

	cli, ok := clientsInstance.clients[path]
	if !ok || !cli.isOpened {
		leveldbClient, err := leveldb.OpenFile(path, &opt.Options{
			ReadOnly: a.settings.ReadOnly,
		})
		if err != nil {
			panic(err)
		}
		// Check context cancellation
		if err := ctx.Err(); err != nil {
			panic(err)
		}

		clientsInstance.clients[path] = &client{
			leveldb:  leveldbClient,
			isOpened: true,
		}
	}

	client := clientsInstance.clients[path]

	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.client = client
}

func (a *adapter) OnCircuitOpen() {
	a.client.isOpened = false
}

func (a *adapter) Open() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	a.setClient(ctx)
}

func (a *adapter) Close() {
	a.client.leveldb.Close()
}

func (a *adapter) Has(ctx context.Context, key string) bool {
	if !a.client.isOpened {
		a.Open()
	}

	ok, err := a.client.leveldb.Has([]byte(key), nil)
	if err != nil {
		return false
	}
	return ok
}

func (a *adapter) Get(ctx context.Context, key string, dest interface{}) error {
	if !a.client.isOpened {
		a.Open()
	}

	value, err := a.client.leveldb.Get([]byte(key), nil)
	if err != nil {
		if errors.Is(err, leveldb.ErrNotFound) {
			return core.ErrNotFound
		}
		return err
	}
	return json.Unmarshal(value, dest)
}

func (a *adapter) Set(ctx context.Context, key string, value interface{}) error {
	if !a.client.isOpened {
		a.Open()
	}

	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return a.client.leveldb.Put([]byte(key), bytes, nil)
}

func (a *adapter) BatchSet(ctx context.Context, kvs []core.KV) error {
	if !a.client.isOpened {
		a.Open()
	}

	batch := new(leveldb.Batch)

	for _, v := range kvs {
		bytes, err := json.Marshal(v.V)
		if err != nil {
			return err
		}

		batch.Put([]byte(v.K), bytes)
	}

	err := a.client.leveldb.Write(batch, nil)
	return err
}

func (a *adapter) Delete(ctx context.Context, key string) error {
	if !a.client.isOpened {
		a.Open()
	}

	return a.client.leveldb.Delete([]byte(key), nil)
}
