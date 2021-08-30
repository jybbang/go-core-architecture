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
	client   *clientProxy
	settings EtcdSettings
}

type clientProxy struct {
	etcd        *etcd.Client
	isConnected bool
}

type clients struct {
	clients map[string]*clientProxy
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
					clients: make(map[string]*clientProxy),
				}
			})
	}
	return clientsInstance
}

func NewEtcdAdapter(settings EtcdSettings) *adapter {
	s := &EtcdSettings{
		DialTimeout: time.Duration(5 * time.Second),
	}

	errs := model.Copy(s, settings)
	if errs != nil {
		panic(fmt.Errorf("mapping errors occurred: %v", errs))
	}

	return &adapter{
		settings: *s,
	}
}

func (a *adapter) IsConnected() bool {
	return a.client.isConnected
}

func (a *adapter) Connect(ctx context.Context) error {
	clientsInstance := getClients()

	clientsInstance.mutex.Lock()
	defer clientsInstance.mutex.Unlock()

	if len(a.settings.Endpoints) == 0 {
		return fmt.Errorf("at least 1 endpoint required")
	}

	endpoint := a.settings.Endpoints[0]

	if strings.TrimSpace(endpoint) == "" {
		return fmt.Errorf("endpoint is required")
	}

	cli, ok := clientsInstance.clients[endpoint]

	if !ok || !cli.isConnected {
		etcdClient, err := etcd.New(etcd.Config{
			Endpoints:   a.settings.Endpoints,
			DialTimeout: a.settings.DialTimeout,
		})

		if err != nil {
			return err
		}

		// Check context cancellation
		if err := ctx.Err(); err != nil {
			return err
		}

		clientsInstance.clients[endpoint] = &clientProxy{
			etcd:        etcdClient,
			isConnected: true,
		}
	}

	a.client = clientsInstance.clients[endpoint]

	return nil
}

func (a *adapter) Disconnect() {
	clientsInstance.mutex.Lock()
	defer clientsInstance.mutex.Unlock()

	a.client.etcd.Close()

	a.client.isConnected = false
}

func (a *adapter) Has(ctx context.Context, key string) bool {
	value, err := a.client.etcd.Get(ctx, key)

	if err != nil {
		return false
	}

	return value.Count > 0
}

func (a *adapter) Get(ctx context.Context, key string, dest interface{}) error {
	value, err := a.client.etcd.Get(ctx, key)

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

	_, err = a.client.etcd.Put(ctx, key, string(bytes))

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
	_, err := a.client.etcd.Delete(ctx, key)

	return err
}
