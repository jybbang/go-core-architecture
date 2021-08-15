package nats

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/nats-io/nats.go"
	cmap "github.com/orcaman/concurrent-map"

	"github.com/jybbang/go-core-architecture/core"
)

type adapter struct {
	nats     *nats.Conn
	pubsubs  cmap.ConcurrentMap
	settings NatsSettings
}

type clients struct {
	clients map[string]*nats.Conn
	pubsubs map[string]cmap.ConcurrentMap
	mutex   sync.Mutex
}

type NatsSettings struct {
	Url string
}

var clientsSync sync.Once

var clientsInstance *clients

func getClients() *clients {
	if clientsInstance == nil {
		clientsSync.Do(
			func() {
				clientsInstance = &clients{
					clients: make(map[string]*nats.Conn),
				}
			})
	}
	return clientsInstance
}

func getNatsClient(ctx context.Context, settings NatsSettings) (*nats.Conn, cmap.ConcurrentMap) {
	clientsInstance := getClients()

	clientsInstance.mutex.Lock()
	defer clientsInstance.mutex.Unlock()

	url := settings.Url
	_, ok := clientsInstance.clients[url]
	if !ok {
		natsClient, err := nats.Connect(url)
		if err != nil {
			panic(err)
		}
		// Check context cancellation
		if err := ctx.Err(); err != nil {
			panic(err)
		}

		core.Log.Infow("natsClient created")
		clientsInstance.clients[url] = natsClient
		clientsInstance.pubsubs[url] = cmap.New()
	}

	client := clientsInstance.clients[url]
	pubsub := clientsInstance.pubsubs[url]
	return client, pubsub
}

func NewNatsAdapter(ctx context.Context, settings NatsSettings) *adapter {
	client, pubsub := getNatsClient(ctx, settings)
	natsService := &adapter{
		nats:     client,
		pubsubs:  pubsub,
		settings: settings,
	}

	return natsService
}

func (a *adapter) Publish(ctx context.Context, coreEvent core.DomainEventer) error {
	bytes, err := json.Marshal(coreEvent)
	if err != nil {
		return err
	}
	return a.nats.Publish(coreEvent.GetTopic(), bytes)
}

func (a *adapter) Subscribe(ctx context.Context, topic string, handler core.ReplyHandler) error {
	pubsub, err := a.nats.Subscribe(topic, func(m *nats.Msg) {
		handler(m.Data)
	})
	if err != nil {
		return err
	}

	a.pubsubs.Set(topic, pubsub)

	return nil
}

func (a *adapter) Unsubscribe(ctx context.Context, topic string) error {
	if pubsub, ok := a.pubsubs.Get(topic); ok {
		if pubsub, ok := pubsub.(*nats.Subscription); ok {
			pubsub.Unsubscribe()
		}
	}

	return core.ErrNotFound
}
