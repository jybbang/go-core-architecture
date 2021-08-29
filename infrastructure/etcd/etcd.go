package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/jybbang/go-core-architecture/core"
	etcd "go.etcd.io/etcd/client/v3"
	"gopkg.in/jeevatkm/go-model.v1"
)

type adapter struct {
	etcd     *etcd.Client
	settings EtcdSettings
	mutex    sync.Mutex
}

type clients struct {
	clients map[string]*etcd.Client
	mutex   sync.Mutex
}

type EtcdSettings struct {
	Endpoints   []string
	DialTimeout time.Duration `model:",omitempty"`
}

var clientsSync sync.Once

var clientsInstance *clients

func getClients() *clients {
	if clientsInstance == nil {
		clientsSync.Do(
			func() {
				clientsInstance = &clients{
					clients: make(map[string]*etcd.Client),
				}
			})
	}
	return clientsInstance
}

func NewEtcdAdapter(ctx context.Context, settings EtcdSettings) *adapter {
	s := &EtcdSettings{
		DialTimeout: time.Duration(5 * time.Second),
	}

	err := model.Copy(s, settings)
	if err != nil {
		panic(fmt.Errorf("mapping errors occurred: %v", err))
	}

	etcdService := &adapter{
		settings: *s,
	}
	etcdService.open(ctx)

	return etcdService
}

func (a *adapter) open(ctx context.Context) {
	clientsInstance := getClients()

	clientsInstance.mutex.Lock()
	defer clientsInstance.mutex.Unlock()

	if len(a.settings.Endpoints) == 0 {
		panic("at least 1 endpoint required")
	}

	endpoint := a.settings.Endpoints[0]

	if strings.TrimSpace(endpoint) == "" {
		panic("endpoint is required")
	}

	_, ok := clientsInstance.clients[endpoint]
	if !ok {
		etcdClient, err := etcd.New(etcd.Config{
			Endpoints:   a.settings.Endpoints,
			DialTimeout: a.settings.DialTimeout,
		})
		if err != nil {
			panic(err)
		}
		// Check context cancellation
		if err := ctx.Err(); err != nil {
			panic(err)
		}
		clientsInstance.clients[endpoint] = etcdClient
	}

	client := clientsInstance.clients[endpoint]

	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.etcd = client
}

func (a *adapter) OnCircuitOpen() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	a.open(ctx)
}

func (a *adapter) Close() {
	a.etcd.Close()
}

func (a *adapter) Has(ctx context.Context, key string) bool {
	value, err := a.etcd.Get(ctx, key)
	if err != nil {
		return false
	}
	return value.Count > 0
}

func (a *adapter) Get(ctx context.Context, key string, dest interface{}) error {
	value, err := a.etcd.Get(ctx, key)
	if value == nil || value.Count < 1 {
		return core.ErrNotFound
	}
	if err != nil {
		return err
	}
	return json.Unmarshal(value.Kvs[0].Value, dest)
}

func (a *adapter) Set(ctx context.Context, key string, value interface{}) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	_, err = a.etcd.Put(ctx, key, string(bytes))
	return err
}

func (a *adapter) BatchSet(ctx context.Context, kvs []core.KV) error {
	for _, v := range kvs {
		err := a.Set(ctx, v.K, v.V)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *adapter) Delete(ctx context.Context, key string) error {
	_, err := a.etcd.Delete(ctx, key)
	return err
}
