package redis

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/go-redis/redis/v8"
	cmap "github.com/orcaman/concurrent-map"

	"github.com/jybbang/go-core-architecture/core"
)

type adapter struct {
	redis   *redis.Client
	pubsubs cmap.ConcurrentMap
}

type clients struct {
	clients map[string]*redis.Client
	pubsubs map[string]cmap.ConcurrentMap
	mutex   sync.Mutex
}

var clientsSync sync.Once

var clientsInstance *clients

func getClients() *clients {
	if clientsInstance == nil {
		clientsSync.Do(
			func() {
				clientsInstance = &clients{
					clients: make(map[string]*redis.Client),
				}
			})
	}
	return clientsInstance
}

func getRedisClient(ctx context.Context, host string, password string) (*redis.Client, cmap.ConcurrentMap) {
	clientsInstance := getClients()

	clientsInstance.mutex.Lock()
	defer clientsInstance.mutex.Unlock()

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

		core.Log.Infow("redisClient created")
		clientsInstance.clients[host] = redisClient
		clientsInstance.pubsubs[host] = cmap.New()
	}

	client := clientsInstance.clients[host]
	pubsub := clientsInstance.pubsubs[host]
	return client, pubsub
}

func NewRedisAdapter(ctx context.Context, host string, password string) *adapter {
	client, pubsub := getRedisClient(ctx, host, password)
	redisService := &adapter{
		redis:   client,
		pubsubs: pubsub,
	}

	return redisService
}

func (a *adapter) Has(ctx context.Context, key string) (ok bool, err error) {
	value, err := a.redis.Exists(ctx, key).Result()
	if err == redis.Nil {
		return false, core.ErrNotFound
	}
	return value > 0, err
}

func (a *adapter) Get(ctx context.Context, key string, dest interface{}) (ok bool, err error) {
	value, err := a.redis.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return false, core.ErrNotFound
	} else if err != nil {
		return false, err
	}

	return true, json.Unmarshal(value, dest)
}

func (a *adapter) Set(ctx context.Context, key string, value interface{}) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return a.redis.Set(ctx, key, bytes, 0).Err()
}

func (a *adapter) Publish(ctx context.Context, coreEvent core.DomainEventer) error {
	err := a.redis.Publish(ctx, coreEvent.GetTopic(), coreEvent).Err()
	return err
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
