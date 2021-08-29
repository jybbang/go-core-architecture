package redis

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	cmap "github.com/orcaman/concurrent-map"

	"github.com/jybbang/go-core-architecture/core"
)

type adapter struct {
	client   *client
	settings RedisSettings
	mutex    sync.Mutex
}

type client struct {
	redis    *redis.Client
	pubsubs  cmap.ConcurrentMap
	handlers cmap.ConcurrentMap
	isOpened bool
}

type clients struct {
	clients map[string]*client
	mutex   sync.Mutex
}

type RedisSettings struct {
	Host     string
	Password string
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

func NewRedisAdapter(ctx context.Context, settings RedisSettings) *adapter {
	redisService := &adapter{
		settings: settings,
	}

	redisService.setClient(ctx)
	return redisService
}

func (a *adapter) setClient(ctx context.Context) {
	clientsInstance := getClients()

	clientsInstance.mutex.Lock()
	defer clientsInstance.mutex.Unlock()

	host := a.settings.Host

	if strings.TrimSpace(host) == "" {
		panic("host is required")
	}

	password := a.settings.Password

	cli, ok := clientsInstance.clients[host]
	if !ok || !cli.isOpened {
		redisClient := redis.NewClient(&redis.Options{
			Addr:     host,
			Password: password,
			DB:       0,
		})
		redisClient.Conn(ctx)
		// Check context cancellation
		if err := ctx.Err(); err != nil {
			panic(err)
		}

		handlers := cli.handlers
		if handlers == nil {
			handlers = cmap.New()
		}

		clientsInstance.clients[host] = &client{
			redis:    redisClient,
			pubsubs:  cmap.New(),
			handlers: handlers,
			isOpened: true,
		}
	}

	client := clientsInstance.clients[host]

	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.client = client
	for _, k := range a.client.handlers.Keys() {
		v, _ := a.client.handlers.Get(k)
		a.Subscribe(context.Background(), k, v.(core.ReplyHandler))
	}
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
	a.client.redis.Close()
}

func (a *adapter) Has(ctx context.Context, key string) bool {
	if !a.client.isOpened {
		a.Open()
	}

	value, err := a.client.redis.Exists(ctx, key).Result()
	if err != nil {
		return false
	}
	return value > 0
}

func (a *adapter) Get(ctx context.Context, key string, dest interface{}) error {
	if !a.client.isOpened {
		a.Open()
	}

	value, err := a.client.redis.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
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

	result := a.client.redis.Set(ctx, key, bytes, 0)
	return result.Err()
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
	if !a.client.isOpened {
		a.Open()
	}

	result := a.client.redis.Del(ctx, key)
	return result.Err()
}

func (a *adapter) Publish(ctx context.Context, coreEvent core.DomainEventer) error {
	if !a.client.isOpened {
		a.Open()
	}

	bytes, err := json.Marshal(coreEvent)
	if err != nil {
		return err
	}
	result := a.client.redis.Publish(ctx, coreEvent.GetTopic(), bytes)
	return result.Err()
}

func (a *adapter) Subscribe(ctx context.Context, topic string, handler core.ReplyHandler) error {
	if !a.client.isOpened {
		a.Open()
	}

	pubsub := a.client.redis.Subscribe(ctx, topic)
	a.client.pubsubs.Set(topic, pubsub)
	a.client.handlers.Set(topic, handler)

	go func() {
		ch := pubsub.Channel()

		for msg := range ch {
			handler(msg.Payload)
		}
	}()

	return nil
}

func (a *adapter) Unsubscribe(ctx context.Context, topic string) error {
	if !a.client.isOpened {
		a.Open()
	}

	if pubsub, ok := a.client.pubsubs.Get(topic); ok {
		if pubsub, ok := pubsub.(*redis.PubSub); ok {
			pubsub.Unsubscribe(ctx)
		}
	}

	a.client.pubsubs.Remove(topic)
	a.client.handlers.Remove(topic)
	return nil
}
