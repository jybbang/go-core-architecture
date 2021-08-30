package dapr

import (
	"context"
	"encoding/json"
	"errors"

	dapr "github.com/dapr/go-sdk/client"

	"github.com/jybbang/go-core-architecture/core"
)

type adapter struct {
	client   dapr.Client
	settings DaprSettings
}

type DaprSettings struct {
	StoreName  string
	PubsubName string
}

var daprClient dapr.Client

func init() {
	client, err := dapr.NewClient()

	if err != nil {
		panic(err)
	}

	daprClient = client
}

func NewDaprAdapter(settings DaprSettings) *adapter {
	return &adapter{
		settings: settings,
	}
}

func (a *adapter) IsConnected() bool {
	return a.client != nil
}

func (a *adapter) Connect(ctx context.Context) error {
	a.client = daprClient

	return nil
}

func (a *adapter) Disconnect() {}

func (a *adapter) Has(ctx context.Context, key string) bool {
	value, err := a.client.GetState(ctx, a.settings.StoreName, key)

	if err != nil {
		return false
	}

	return value.Value != nil
}

func (a *adapter) Get(ctx context.Context, key string, dest interface{}) error {
	value, err := a.client.GetState(ctx, a.settings.StoreName, key)

	if value == nil {
		return core.ErrNotFound
	}

	if err != nil {
		return err
	}

	return json.Unmarshal(value.Value, dest)
}

func (a *adapter) Set(ctx context.Context, key string, value interface{}) error {
	bytes, err := json.Marshal(value)

	if err != nil {
		return err
	}

	err = a.client.SaveState(ctx, a.settings.StoreName, key, bytes)

	if err != nil {
		return err
	}

	return nil
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
	return a.client.DeleteState(ctx, a.settings.StoreName, key)
}

func (a *adapter) Publish(ctx context.Context, coreEvent core.DomainEventer) error {
	bytes, err := json.Marshal(coreEvent)

	if err != nil {
		return err
	}

	err = a.client.PublishEvent(ctx, a.settings.PubsubName, coreEvent.GetTopic(), bytes)

	return err
}

func (a *adapter) Subscribe(ctx context.Context, topic string, handler core.ReplyHandler) error {
	return errors.New("not supported operation")
}

func (a *adapter) Unsubscribe(ctx context.Context, topic string) error {
	return errors.New("not supported operation")
}
