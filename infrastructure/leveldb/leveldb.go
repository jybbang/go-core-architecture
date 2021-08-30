package leveldb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/jybbang/go-core-architecture/core"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

type adapter struct {
	client   *clientProxy
	settings LevelDbSettings
}

type clientProxy struct {
	leveldb     *leveldb.DB
	isConnected bool
}

type clients struct {
	clients map[string]*clientProxy
	sync.Mutex
}

type LevelDbSettings struct {
	Path     string
	ReadOnly bool
}

var clientsInstance *clients

func init() {
	clientsInstance = &clients{
		clients: make(map[string]*clientProxy),
	}
}

func NewLevelDbAdapter(settings LevelDbSettings) *adapter {
	return &adapter{
		settings: settings,
	}
}

func (a *adapter) IsConnected() bool {
	return a.client.isConnected
}

func (a *adapter) Connect(ctx context.Context) error {
	clientsInstance.Lock()
	defer clientsInstance.Unlock()

	path := a.settings.Path

	if strings.TrimSpace(path) == "" {
		return fmt.Errorf("path is required")
	}

	cli, ok := clientsInstance.clients[path]

	if !ok || !cli.isConnected {
		leveldbClient, err := leveldb.OpenFile(path, &opt.Options{
			ReadOnly: a.settings.ReadOnly,
		})

		if err != nil {
			return err
		}

		// Check context cancellation
		if err := ctx.Err(); err != nil {
			return err
		}

		clientsInstance.clients[path] = &clientProxy{
			leveldb:     leveldbClient,
			isConnected: true,
		}
	}

	a.client = clientsInstance.clients[path]

	return nil
}

func (a *adapter) Disconnect() {
	clientsInstance.Lock()
	defer clientsInstance.Unlock()

	a.client.leveldb.Close()

	a.client.isConnected = false
}

func (a *adapter) Has(ctx context.Context, key string) bool {
	ok, err := a.client.leveldb.Has([]byte(key), nil)

	if err != nil {
		return false
	}

	return ok
}

func (a *adapter) Get(ctx context.Context, key string, dest interface{}) error {
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
	bytes, err := json.Marshal(value)

	if err != nil {
		return err
	}

	return a.client.leveldb.Put([]byte(key), bytes, nil)
}

func (a *adapter) BatchSet(ctx context.Context, kvs []core.KV) error {
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
	return a.client.leveldb.Delete([]byte(key), nil)
}
