package etcd

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	etcd "go.etcd.io/etcd/client/v3"
)

type adapter struct {
	etcd     *etcd.Client
	settings EtcdSettings
}

type clients struct {
	clients map[string]*etcd.Client
	mutex   sync.Mutex
}

type EtcdSettings struct {
	Endpoints   []string
	DialTimeout time.Duration
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

func getEtcdClient(ctx context.Context, settings EtcdSettings) *etcd.Client {
	clientsInstance := getClients()

	clientsInstance.mutex.Lock()
	defer clientsInstance.mutex.Unlock()

	if settings.Endpoints == nil || len(settings.Endpoints) == 0 {
		panic("at least 1 endpoint required")
	}
	endpoint := settings.Endpoints[0]
	_, ok := clientsInstance.clients[endpoint]
	if !ok {
		etcdClient, err := etcd.New(etcd.Config{
			Endpoints:   settings.Endpoints,
			DialTimeout: settings.DialTimeout,
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
	return client
}

func NewEtcdAdapter(ctx context.Context, settings EtcdSettings) *adapter {
	client := getEtcdClient(ctx, settings)
	etcdService := &adapter{
		etcd:     client,
		settings: settings,
	}

	return etcdService
}

func (a *adapter) Has(ctx context.Context, key string) (ok bool, err error) {
	value, err := a.etcd.Get(ctx, key)
	if err != nil {
		return false, err
	}
	return value.Count > 0, err
}

func (a *adapter) Get(ctx context.Context, key string, dest interface{}) (ok bool, err error) {
	value, err := a.etcd.Get(ctx, key)
	if err != nil {
		return false, err
	}

	return true, json.Unmarshal(value.XXX_unrecognized, dest)
}

func (a *adapter) Set(ctx context.Context, key string, value interface{}) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	_, err = a.etcd.Put(ctx, key, string(bytes))
	return err
}

func (a *adapter) Delete(ctx context.Context, key string) error {
	_, err := a.etcd.Delete(ctx, key)
	return err
}
