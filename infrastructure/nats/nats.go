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
	client   *client
	settings NatsSettings
	mutex    sync.Mutex
}

type client struct {
	nats     *nats.Conn
	pubsubs  cmap.ConcurrentMap
	handlers cmap.ConcurrentMap
	isOpened bool
}

type clients struct {
	clients map[string]*client
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
					clients: make(map[string]*client),
				}
			})
	}
	return clientsInstance
}

func NewNatsAdapter(ctx context.Context, settings NatsSettings) *adapter {
	natsService := &adapter{
		settings: settings,
	}

	natsService.setClient(ctx)
	return natsService
}

func (a *adapter) setClient(ctx context.Context) {
	clientsInstance := getClients()

	clientsInstance.mutex.Lock()
	defer clientsInstance.mutex.Unlock()

	url := a.settings.Url

	if strings.TrimSpace(url) == "" {
		panic("url is required")
	}

	cli, ok := clientsInstance.clients[url]
	if !ok || !cli.isOpened {
		natsClient, err := nats.Connect(url)
		if err != nil {
			panic(err)
		}
		// Check context cancellation
		if err := ctx.Err(); err != nil {
			panic(err)
		}

		handlers := cli.handlers
		if handlers == nil {
			handlers = cmap.New()
		}

		clientsInstance.clients[url] = &client{
			nats:     natsClient,
			pubsubs:  cmap.New(),
			handlers: handlers,
			isOpened: true,
		}
	}

	client := clientsInstance.clients[url]

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
	a.client.nats.Close()
}

func (a *adapter) Publish(ctx context.Context, coreEvent core.DomainEventer) error {
	if !a.client.isOpened {
		a.Open()
	}

	bytes, err := json.Marshal(coreEvent)
	if err != nil {
		return err
	}
	return a.client.nats.Publish(coreEvent.GetTopic(), bytes)
}

func (a *adapter) Subscribe(ctx context.Context, topic string, handler core.ReplyHandler) error {
	if !a.client.isOpened {
		a.Open()
	}

	pubsub, err := a.client.nats.Subscribe(topic, func(m *nats.Msg) {
		handler(m.Data)
	})
	if err != nil {
		return err
	}

	a.client.pubsubs.Set(topic, pubsub)
	a.client.handlers.Set(topic, handler)
	return nil
}

func (a *adapter) Unsubscribe(ctx context.Context, topic string) error {
	if !a.client.isOpened {
		a.Open()
	}

	if pubsub, ok := a.client.pubsubs.Get(topic); ok {
		if pubsub, ok := pubsub.(*nats.Subscription); ok {
			pubsub.Unsubscribe()
		}
	}

	a.client.pubsubs.Remove(topic)
	a.client.handlers.Remove(topic)
	return nil
}
