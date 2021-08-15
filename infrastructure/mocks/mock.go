package mocks

import (
	"context"

	"github.com/google/uuid"
	cmap "github.com/orcaman/concurrent-map"

	"github.com/jybbang/go-core-architecture/core"
	"gopkg.in/jeevatkm/go-model.v1"
)

type adapter struct {
	model   core.Entitier
	db      cmap.ConcurrentMap
	pubsubs cmap.ConcurrentMap
	states  cmap.ConcurrentMap
}

func NewMockAdapter() *adapter {
	return &adapter{
		db:      cmap.New(),
		pubsubs: cmap.New(),
		states:  cmap.New(),
	}
}

func (a *adapter) Has(ctx context.Context, key string) (ok bool, err error) {
	// Check context cancellation
	if err := ctx.Err(); err != nil {
		return false, err
	}

	core.Log.Info("mock has - {}", key)
	return a.states.Has(key), nil
}

func (a *adapter) Get(ctx context.Context, key string, dest interface{}) (ok bool, err error) {
	// Check context cancellation
	if err := ctx.Err(); err != nil {
		return false, err
	}

	core.Log.Info("mock get - {}", key)
	if resp, ok := a.states.Get(key); ok {
		model.Copy(dest, resp)
		return true, nil
	}
	return false, core.ErrNotFound
}

func (a *adapter) Set(ctx context.Context, key string, value interface{}) error {
	// Check context cancellation
	if err := ctx.Err(); err != nil {
		return err
	}

	core.Log.Info("mock set - {}", key, value)
	a.states.Set(key, value)
	return nil
}

func (a *adapter) Publish(ctx context.Context, coreEvent core.DomainEventer) error {
	// Check context cancellation
	if err := ctx.Err(); err != nil {
		return err
	}

	core.Log.Info("mock publish - {} {}", coreEvent.GetID(), coreEvent.GetTopic())
	return nil
}

func (a *adapter) Subscribe(ctx context.Context, topic string, handler core.ReplyHandler) error {
	// Check context cancellation
	if err := ctx.Err(); err != nil {
		return err
	}

	core.Log.Info("mock subscribe - {}", topic)
	a.pubsubs.Set(topic, handler)
	return nil
}

func (a *adapter) Unsubscribe(ctx context.Context, topic string) error {
	// Check context cancellation
	if err := ctx.Err(); err != nil {
		return err
	}

	core.Log.Info("mock unsubscribe - {}", topic)
	a.pubsubs.Remove(topic)
	return nil
}

func (a *adapter) FakeSend(topic string, receivedData interface{}) {
	core.Log.Info("mock fake send - {} {}", topic, receivedData)
	if resp, ok := a.pubsubs.Get(topic); ok {
		resp.(core.ReplyHandler)(receivedData)
	}
}

func (a *adapter) SetModel(model core.Entitier) {
	core.Log.Info("mock setmodel")
	a.model = model
}

func (a *adapter) Find(ctx context.Context, dest core.Entitier, id uuid.UUID) (ok bool, err error) {
	// Check context cancellation
	if err = ctx.Err(); err != nil {
		return false, err
	}

	core.Log.Info("mock find - {}", id)
	if resp, ok := a.db.Get(id.String()); ok {
		model.Copy(dest, resp)
		return true, nil
	}
	return false, core.ErrNotFound
}

func (a *adapter) Any(ctx context.Context) (ok bool, err error) {
	// Check context cancellation
	if err = ctx.Err(); err != nil {
		return false, err
	}

	core.Log.Info("mock any")
	count, err := a.Count(ctx)
	return count > 0, err
}

func (a *adapter) AnyWithFilter(ctx context.Context, query interface{}, args interface{}) (ok bool, err error) {
	// Check context cancellation
	if err = ctx.Err(); err != nil {
		return false, err
	}

	core.Log.Info("mock anywithfilter")
	count, err := a.CountWithFilter(ctx, query, args)
	return count > 0, err
}

func (a *adapter) Count(ctx context.Context) (count int64, err error) {
	// Check context cancellation
	if err = ctx.Err(); err != nil {
		return 0, err
	}

	core.Log.Info("mock count")
	resp := a.db.Count()
	return int64(resp), nil
}

func (a *adapter) CountWithFilter(ctx context.Context, query interface{}, args interface{}) (count int64, err error) {
	// Check context cancellation
	if err = ctx.Err(); err != nil {
		return 0, err
	}

	core.Log.Info("mock countwithfilter")
	resp := a.db.Count()
	return int64(resp), nil
}

func (a *adapter) List(ctx context.Context, dest []core.Entitier) error {
	// Check context cancellation
	if err := ctx.Err(); err != nil {
		return err
	}

	core.Log.Info("mock list")
	for _, v := range a.db.Items() {
		entity := v.(core.Entitier)
		dest = append(dest, entity)
	}

	return nil
}

func (a *adapter) ListWithFilter(ctx context.Context, dest []core.Entitier, query interface{}, args interface{}) error {
	// Check context cancellation
	if err := ctx.Err(); err != nil {
		return err
	}

	core.Log.Info("mock listwithfilter")
	for _, v := range a.db.Items() {
		entity := v.(core.Entitier)
		dest = append(dest, entity)
	}

	return nil
}

func (a *adapter) Remove(ctx context.Context, entity core.Entitier) error {
	// Check context cancellation
	if err := ctx.Err(); err != nil {
		return err
	}

	core.Log.Info("mock remove - {}", entity)
	a.db.Remove(entity.GetID().String())
	return nil
}

func (a *adapter) RemoveRange(ctx context.Context, entities []core.Entitier) error {
	// Check context cancellation
	if err := ctx.Err(); err != nil {
		return err
	}

	core.Log.Info("mock removerange")
	for _, v := range entities {
		a.Remove(ctx, v)
	}
	return nil
}

func (a *adapter) Add(ctx context.Context, entity core.Entitier) error {
	// Check context cancellation
	if err := ctx.Err(); err != nil {
		return err
	}

	core.Log.Info("mock add - ", entity)
	a.db.Set(entity.GetID().String(), entity)
	return nil
}

func (a *adapter) AddRange(ctx context.Context, entities []core.Entitier) error {
	// Check context cancellation
	if err := ctx.Err(); err != nil {
		return err
	}

	core.Log.Info("mock addrange")
	for _, v := range entities {
		a.Add(ctx, v)
	}
	return nil
}

func (a *adapter) Update(ctx context.Context, entity core.Entitier) error {
	// Check context cancellation
	if err := ctx.Err(); err != nil {
		return err
	}

	core.Log.Info("mock update - ", entity)
	a.db.Set(entity.GetID().String(), entity)
	return nil
}

func (a *adapter) UpdateRange(ctx context.Context, entities []core.Entitier) error {
	// Check context cancellation
	if err := ctx.Err(); err != nil {
		return err
	}

	core.Log.Info("mock updaterange")
	for _, v := range entities {
		a.Update(ctx, v)
	}
	return nil
}
