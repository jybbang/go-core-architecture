package redis

import (
	"context"
	"encoding/json"
	"errors"
	"sync"

	"github.com/go-redis/redis/v8"
	cmap "github.com/orcaman/concurrent-map"

	"github.com/jybbang/go-core-architecture/core"
)

type adapter struct {
	redis    *redis.Client
	pubsubs  cmap.ConcurrentMap
	settings RedisSettings
}

type clients struct {
	clients map[string]*redis.Client
	pubsubs map[string]cmap.ConcurrentMap
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
					clients: make(map[string]*redis.Client),
					pubsubs: make(map[string]cmap.ConcurrentMap),
				}
			})
	}
	return clientsInstance
}

func getRedisClient(ctx context.Context, settings RedisSettings) (*redis.Client, cmap.ConcurrentMap) {
	clientsInstance := getClients()

	clientsInstance.mutex.Lock()
	defer clientsInstance.mutex.Unlock()

	host := settings.Host
	password := settings.Password
	_, ok := clientsInstance.clients[host]
	if !ok {
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

		clientsInstance.clients[host] = redisClient
		clientsInstance.pubsubs[host] = cmap.New()
	}

	client := clientsInstance.clients[host]
	pubsub := clientsInstance.pubsubs[host]
	return client, pubsub
}

func NewRedisAdapter(ctx context.Context, settings RedisSettings) *adapter {
	client, pubsub := getRedisClient(ctx, settings)
	redisService := &adapter{
		redis:    client,
		pubsubs:  pubsub,
		settings: settings,
	}

	return redisService
}

func (a *adapter) Close() {
	a.redis.Close()
}

func (a *adapter) Has(ctx context.Context, key string) (ok bool, err error) {
	value, err := a.redis.Exists(ctx, key).Result()
	if err == redis.Nil {
		return false, core.ErrNotFound
	}
	return value > 0, err
}

func (a *adapter) Get(ctx context.Context, key string, dest interface{}) (err error) {
	value, err := a.redis.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return core.ErrNotFound
	} else if err != nil {
		return err
	}
	return json.Unmarshal(value, dest)
}

func (a *adapter) Set(ctx context.Context, key string, value interface{}) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	result := a.redis.Set(ctx, key, bytes, 0)
	return result.Err()
}

func (a *adapter) Delete(ctx context.Context, key string) error {
	result := a.redis.Del(ctx, key)
	return result.Err()
}

func (a *adapter) Publish(ctx context.Context, coreEvent core.DomainEventer) error {
	bytes, err := json.Marshal(coreEvent)
	if err != nil {
		return err
	}
	result := a.redis.Publish(ctx, coreEvent.GetTopic(), bytes)
	return result.Err()
}

func (a *adapter) Subscribe(ctx context.Context, topic string, handler core.ReplyHandler) error {
	pubsub := a.redis.Subscribe(ctx, topic)
	a.pubsubs.Set(topic, pubsub)

	ch := pubsub.Channel()

	for msg := range ch {
		handler(msg.Payload)
	}

	return nil
}

func (a *adapter) Unsubscribe(ctx context.Context, topic string) error {
	if pubsub, ok := a.pubsubs.Get(topic); ok {
		if pubsub, ok := pubsub.(*redis.PubSub); ok {
			pubsub.Unsubscribe(ctx)
		}
	}

	return core.ErrNotFound
}
