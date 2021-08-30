package nats

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/nats-io/nats.go"
	cmap "github.com/orcaman/concurrent-map"

	"github.com/jybbang/go-core-architecture/core"
)

type adapter struct {
	client   *clientProxy
	settings NatsSettings
}

type clientProxy struct {
	nats        *nats.Conn
	pubsubs     cmap.ConcurrentMap
	handlers    cmap.ConcurrentMap
	isConnected bool
}

type clients struct {
	clients map[string]*clientProxy
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
					clients: make(map[string]*clientProxy),
				}
			})
	}
	return clientsInstance
}

func NewNatsAdapter(settings NatsSettings) *adapter {
	return &adapter{
		settings: settings,
	}
}

func (a *adapter) IsConnected() bool {
	return a.client.isConnected
}

func (a *adapter) Connect(ctx context.Context) error {
	clientsInstance := getClients()

	clientsInstance.mutex.Lock()
	defer clientsInstance.mutex.Unlock()

	url := a.settings.Url

	if strings.TrimSpace(url) == "" {
		return fmt.Errorf("url is required")
	}

	cli, ok := clientsInstance.clients[url]

	if !ok || !cli.isConnected {
		natsClient, err := nats.Connect(url)

		if err != nil {
			return err
		}

		// Check context cancellation
		if err := ctx.Err(); err != nil {
			return err
		}

		handlers := cli.handlers

		if handlers == nil {
			handlers = cmap.New()
		}

		clientsInstance.clients[url] = &clientProxy{
			nats:        natsClient,
			pubsubs:     cmap.New(),
			handlers:    handlers,
			isConnected: true,
		}
	}

	a.client = clientsInstance.clients[url]

	for _, k := range a.client.handlers.Keys() {
		v, _ := a.client.handlers.Get(k)

		a.Subscribe(context.Background(), k, v.(core.ReplyHandler))
	}

	return nil
}

func (a *adapter) Disconnect() {
	a.client.nats.Close()

	a.client.isConnected = false
}

func (a *adapter) Publish(ctx context.Context, coreEvent core.DomainEventer) error {
	bytes, err := json.Marshal(coreEvent)

	if err != nil {
		return err
	}

	return a.client.nats.Publish(coreEvent.GetTopic(), bytes)
}

func (a *adapter) Subscribe(ctx context.Context, topic string, handler core.ReplyHandler) error {
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
	if pubsub, ok := a.client.pubsubs.Get(topic); ok {
		if pubsub, ok := pubsub.(*nats.Subscription); ok {
			pubsub.Unsubscribe()
		}
	}

	a.client.pubsubs.Remove(topic)

	a.client.handlers.Remove(topic)

	return nil
}
