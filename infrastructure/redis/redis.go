package redis

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/go-redis/redis/v8"
	cmap "github.com/orcaman/concurrent-map"
	"go.uber.org/zap"

	"github.com/jybbang/core-architecture/application/contracts"
	"github.com/jybbang/core-architecture/domain"
)

type adapter struct {
	redis   *redis.Client
	pubsubs cmap.ConcurrentMap
}

type clients struct {
	clients map[string]*redis.Client
	sync.Mutex
}

var log *zap.SugaredLogger

var clientsSync sync.Once

var clientsInstance *clients

var ctx context.Context

func init() {
	logger, _ := zap.NewProduction()
	log = logger.Sugar()
}

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

func getRedisClient(host string, password string) *redis.Client {
	clientsInstance := getClients()

	clientsInstance.Lock()
	defer clientsInstance.Unlock()

	_, ok := clientsInstance.clients[host]
	if !ok {
		ctx = context.Background()
		redisClient := redis.NewClient(&redis.Options{
			Addr:     host,
			Password: password,
			DB:       0,
		})
		log.Info("redisClient created")
		clientsInstance.clients[host] = redisClient
	}

	client := clientsInstance.clients[host]
	return client
}

func NewRedisAdapter(host string, password string) *adapter {
	redisService := &adapter{
		redis:   getRedisClient(host, password),
		pubsubs: cmap.New(),
	}

	return redisService
}

func (a *adapter) Has(key string) (bool, error) {
	value, err := a.redis.Exists(ctx, key).Result()
	if err == redis.Nil {
		return false, domain.ErrNotFound
	}
	return value > 0, err
}

func (a *adapter) Get(key string, dest domain.Copyable) (bool, error) {
	value, err := a.redis.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return false, domain.ErrNotFound
	} else if err != nil {
		return false, err
	}

	return true, json.Unmarshal(value, dest)
}

func (a *adapter) Set(key string, value interface{}) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return a.redis.Set(ctx, key, bytes, 0).Err()
}

func (a *adapter) Publish(domainEvent *domain.DomainEvent) error {
	err := a.redis.Publish(ctx, domainEvent.Topic, domainEvent).Err()
	return err
}

func (a *adapter) Subscribe(topic string, handler contracts.ReplyHandler) error {
	pubsub := a.redis.Subscribe(ctx, topic)
	a.pubsubs.Set(topic, pubsub)

	ch := pubsub.Channel()

	for msg := range ch {
		handler(msg.Payload)
	}

	return nil
}

func (a *adapter) Unsubscribe(topic string) error {
	if pubsub, ok := a.pubsubs.Get(topic); ok {
		if pubsub, ok := pubsub.(*redis.PubSub); ok {
			pubsub.Unsubscribe(ctx)
		}
	}

	return domain.ErrNotFound
}
