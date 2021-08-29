package nats

import (
	"context"
	"encoding/json"
	"strings"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	cmap "github.com/orcaman/concurrent-map"

	"github.com/jybbang/go-core-architecture/core"
)

type adapter struct {
	nats     *nats.Conn
	pubsubs  cmap.ConcurrentMap
	handlers cmap.ConcurrentMap
	settings NatsSettings
	mutex    sync.Mutex
}

type clients struct {
	clients  map[string]*nats.Conn
	pubsubs  map[string]cmap.ConcurrentMap
	handlers map[string]cmap.ConcurrentMap
	mutex    sync.Mutex
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
					pubsubs: make(map[string]cmap.ConcurrentMap),
				}
			})
	}
	return clientsInstance
}

func NewNatsAdapter(ctx context.Context, settings NatsSettings) *adapter {
	natsService := &adapter{
		settings: settings,
	}

	natsService.open(ctx)
	return natsService
}

func (a *adapter) open(ctx context.Context) {
	clientsInstance := getClients()

	clientsInstance.mutex.Lock()
	defer clientsInstance.mutex.Unlock()

	url := a.settings.Url

	if strings.TrimSpace(url) == "" {
		panic("url is required")
	}

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

		clientsInstance.clients[url] = natsClient
		clientsInstance.pubsubs[url] = cmap.New()

		if _, ok := clientsInstance.handlers[url]; !ok {
			clientsInstance.handlers[url] = cmap.New()
		}
	}

	client := clientsInstance.clients[url]
	pubsubs := clientsInstance.pubsubs[url]
	handlers := clientsInstance.handlers[url]

	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.nats = client
	a.pubsubs = pubsubs
	a.handlers = handlers

	for _, k := range handlers.Keys() {
		if v, ok := handlers.Get(k); ok {
			go a.Subscribe(context.Background(), k, v.(core.ReplyHandler))
		}
	}
}

func (a *adapter) OnCircuitOpen() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	a.open(ctx)
}

func (a *adapter) Close() {
	a.nats.Close()
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
	a.handlers.Set(topic, handler)

	return nil
}

func (a *adapter) Unsubscribe(ctx context.Context, topic string) error {
	if pubsub, ok := a.pubsubs.Get(topic); ok {
		if pubsub, ok := pubsub.(*nats.Subscription); ok {
			pubsub.Unsubscribe()
		}
	}

	a.pubsubs.Remove(topic)
	a.handlers.Remove(topic)
	return nil
}
