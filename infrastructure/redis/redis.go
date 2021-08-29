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
	redis    *redis.Client
	pubsubs  cmap.ConcurrentMap
	handlers cmap.ConcurrentMap
	settings RedisSettings
	isOpened bool
	mutex    sync.Mutex
}

type clients struct {
	clients  map[string]*redis.Client
	pubsubs  map[string]cmap.ConcurrentMap
	handlers map[string]cmap.ConcurrentMap
	mutex    sync.Mutex
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
					clients:  make(map[string]*redis.Client),
					pubsubs:  make(map[string]cmap.ConcurrentMap),
					handlers: make(map[string]cmap.ConcurrentMap),
				}
			})
	}
	return clientsInstance
}

func NewRedisAdapter(ctx context.Context, settings RedisSettings) *adapter {
	redisService := &adapter{
		settings: settings,
	}

	return redisService
}

func (a *adapter) open(ctx context.Context) {
	clientsInstance := getClients()

	clientsInstance.mutex.Lock()
	defer clientsInstance.mutex.Unlock()

	host := a.settings.Host

	if strings.TrimSpace(host) == "" {
		panic("host is required")
	}

	password := a.settings.Password
	_, ok := clientsInstance.clients[host]
	if !ok || !a.isOpened {
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

		if _, ok := clientsInstance.handlers[host]; !ok {
			clientsInstance.handlers[host] = cmap.New()
		}

		clientsInstance.pubsubs[host] = cmap.New()
		clientsInstance.clients[host] = redisClient
		a.isOpened = true
	}

	client := clientsInstance.clients[host]
	pubsubs := clientsInstance.pubsubs[host]
	handlers := clientsInstance.handlers[host]

	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.redis = client
	a.pubsubs = pubsubs
	a.handlers = handlers

	for _, k := range handlers.Keys() {
		if v, ok := handlers.Get(k); ok {
			go a.Subscribe(context.Background(), k, v.(core.ReplyHandler))
		}
	}
}

func (a *adapter) OnCircuitOpen() {
	a.isOpened = false
}

func (a *adapter) Open() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	a.open(ctx)
}

func (a *adapter) Close() {
	a.redis.Close()
}

func (a *adapter) Has(ctx context.Context, key string) bool {
	if !a.isOpened {
		a.Open()
	}

	value, err := a.redis.Exists(ctx, key).Result()
	if err != nil {
		return false
	}
	return value > 0
}

func (a *adapter) Get(ctx context.Context, key string, dest interface{}) error {
	if !a.isOpened {
		a.Open()
	}

	value, err := a.redis.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return core.ErrNotFound
		}
		return err
	}
	return json.Unmarshal(value, dest)
}

func (a *adapter) Set(ctx context.Context, key string, value interface{}) error {
	if !a.isOpened {
		a.Open()
	}

	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	result := a.redis.Set(ctx, key, bytes, 0)
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
	if !a.isOpened {
		a.Open()
	}

	result := a.redis.Del(ctx, key)
	return result.Err()
}

func (a *adapter) Publish(ctx context.Context, coreEvent core.DomainEventer) error {
	if !a.isOpened {
		a.Open()
	}

	bytes, err := json.Marshal(coreEvent)
	if err != nil {
		return err
	}
	result := a.redis.Publish(ctx, coreEvent.GetTopic(), bytes)
	return result.Err()
}

func (a *adapter) Subscribe(ctx context.Context, topic string, handler core.ReplyHandler) error {
	if !a.isOpened {
		a.Open()
	}

	pubsub := a.redis.Subscribe(ctx, topic)
	a.pubsubs.Set(topic, pubsub)
	a.handlers.Set(topic, handler)

	ch := pubsub.Channel()

	for msg := range ch {
		handler(msg.Payload)
	}

	return nil
}

func (a *adapter) Unsubscribe(ctx context.Context, topic string) error {
	if !a.isOpened {
		a.Open()
	}

	if pubsub, ok := a.pubsubs.Get(topic); ok {
		if pubsub, ok := pubsub.(*redis.PubSub); ok {
			pubsub.Unsubscribe(ctx)
		}
	}

	a.pubsubs.Remove(topic)
	a.handlers.Remove(topic)
	return nil
}
