package dapr

import (
	"context"
	"encoding/json"
	"errors"
	"sync"

	dapr "github.com/dapr/go-sdk/client"

	"github.com/jybbang/go-core-architecture/core"
)

type adapter struct {
	dapr     dapr.Client
	settings DaprSettings
}

type DaprSettings struct {
	StoreName  string
	PubsubName string
}

var daprClient dapr.Client

var clientsSync sync.Once

func getClient(ctx context.Context, settings DaprSettings) dapr.Client {
	if daprClient == nil {
		clientsSync.Do(
			func() {
				client, err := dapr.NewClient()
				if err != nil {
					panic(err)
				}
				// Check context cancellation
				if err := ctx.Err(); err != nil {
					panic(err)
				}
				daprClient = client
			})
	}
	return daprClient
}

func NewDaprAdapter(ctx context.Context, settings DaprSettings) *adapter {
	client := getClient(ctx, settings)
	daprService := &adapter{
		dapr:     client,
		settings: settings,
	}

	return daprService
}

func (a *adapter) Has(ctx context.Context, key string) (ok bool, err error) {
	value, err := a.dapr.GetState(ctx, a.settings.StoreName, key)
	if err != nil {
		return false, err
	}
	return value.Value != nil, err
}

func (a *adapter) Get(ctx context.Context, key string, dest interface{}) (err error) {
	value, err := a.dapr.GetState(ctx, a.settings.StoreName, key)
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
	err = a.dapr.SaveState(ctx, a.settings.StoreName, key, bytes)
	if err != nil {
		return err
	}
	return nil
}

func (a *adapter) Delete(ctx context.Context, key string) error {
	return a.dapr.DeleteState(ctx, a.settings.StoreName, key)
}

func (a *adapter) Publish(ctx context.Context, coreEvent core.DomainEventer) error {
	bytes, err := json.Marshal(coreEvent)
	if err != nil {
		return err
	}
	err = a.dapr.PublishEvent(ctx, a.settings.PubsubName, coreEvent.GetTopic(), bytes)
	return err
}

func (a *adapter) Subscribe(ctx context.Context, topic string, handler core.ReplyHandler) error {
	return errors.New("not supported operation")
}

func (a *adapter) Unsubscribe(ctx context.Context, topic string) error {
	return errors.New("not supported operation")
}
