package mocks

import (
	"context"

	"github.com/google/uuid"
	cmap "github.com/orcaman/concurrent-map"
	"go.uber.org/zap"

	"github.com/jybbang/go-core-architecture/core"
	"gopkg.in/jeevatkm/go-model.v1"
)

type adapter struct {
	model   core.Entitier
	db      cmap.ConcurrentMap
	pubsubs cmap.ConcurrentMap
	states  cmap.ConcurrentMap
	setting MockSettings
}

type MockSettings struct {
	Log *zap.SugaredLogger
}

var mock = &adapter{
	db:      cmap.New(),
	pubsubs: cmap.New(),
	states:  cmap.New(),
}

func NewMockAdapter(setting MockSettings) *adapter {
	mock.setting = setting
	return mock
}

func (a *adapter) Has(ctx context.Context, key string) (ok bool, err error) {
	// Check context cancellation
	if err := ctx.Err(); err != nil {
		return false, err
	}

	a.setting.Log.Debugw("mock has", "key", key)
	return a.states.Has(key), nil
}

func (a *adapter) Get(ctx context.Context, key string, dest interface{}) (ok bool, err error) {
	// Check context cancellation
	if err := ctx.Err(); err != nil {
		return false, err
	}

	a.setting.Log.Debugw("mock get", "key", key)
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

	a.setting.Log.Debugw("mock set", "key", key, "value", value)
	a.states.Set(key, value)
	return nil
}

func (a *adapter) Delete(ctx context.Context, key string) error {
	// Check context cancellation
	if err := ctx.Err(); err != nil {
		return err
	}

	a.setting.Log.Debugw("mock delete", "key", key)
	a.states.Remove(key)
	return nil
}

func (a *adapter) Publish(ctx context.Context, coreEvent core.DomainEventer) error {
	// Check context cancellation
	if err := ctx.Err(); err != nil {
		return err
	}

	a.setting.Log.Debugw("mock publish", "id", coreEvent.GetID(), "topic", coreEvent.GetTopic())
	return nil
}

func (a *adapter) Subscribe(ctx context.Context, topic string, handler core.ReplyHandler) error {
	// Check context cancellation
	if err := ctx.Err(); err != nil {
		return err
	}

	a.setting.Log.Debugw("mock subscribe", "topic", topic)
	a.pubsubs.Set(topic, handler)
	return nil
}

func (a *adapter) Unsubscribe(ctx context.Context, topic string) error {
	// Check context cancellation
	if err := ctx.Err(); err != nil {
		return err
	}

	a.setting.Log.Debugw("mock unsubscribe", "topic", topic)
	a.pubsubs.Remove(topic)
	return nil
}

func (a *adapter) FakeSend(topic string, receivedData interface{}) {
	a.setting.Log.Debugw("mock fake send - {} {}", topic, receivedData)
	if resp, ok := a.pubsubs.Get(topic); ok {
		resp.(core.ReplyHandler)(receivedData)
	}
}

func (a *adapter) SetModel(model core.Entitier) {
	a.setting.Log.Debugw("mock setmodel")
	a.model = model
}

func (a *adapter) Find(ctx context.Context, dest core.Entitier, id uuid.UUID) (ok bool, err error) {
	// Check context cancellation
	if err = ctx.Err(); err != nil {
		return false, err
	}

	a.setting.Log.Debugw("mock find", "id", id)
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

	a.setting.Log.Debugw("mock any")
	count, err := a.Count(ctx)
	return count > 0, err
}

func (a *adapter) AnyWithFilter(ctx context.Context, query interface{}, args interface{}) (ok bool, err error) {
	// Check context cancellation
	if err = ctx.Err(); err != nil {
		return false, err
	}

	a.setting.Log.Debugw("mock anywithfilter")
	count, err := a.CountWithFilter(ctx, query, args)
	return count > 0, err
}

func (a *adapter) Count(ctx context.Context) (count int64, err error) {
	// Check context cancellation
	if err = ctx.Err(); err != nil {
		return 0, err
	}

	a.setting.Log.Debugw("mock count")
	resp := a.db.Count()
	return int64(resp), nil
}

func (a *adapter) CountWithFilter(ctx context.Context, query interface{}, args interface{}) (count int64, err error) {
	// Check context cancellation
	if err = ctx.Err(); err != nil {
		return 0, err
	}

	a.setting.Log.Debugw("mock countwithfilter")
	resp := a.db.Count()
	return int64(resp), nil
}

func (a *adapter) List(ctx context.Context, dest []core.Entitier) error {
	// Check context cancellation
	if err := ctx.Err(); err != nil {
		return err
	}

	a.setting.Log.Debugw("mock list")
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

	a.setting.Log.Debugw("mock listwithfilter")
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

	a.setting.Log.Debugw("mock remove", "entity", entity)
	a.db.Remove(entity.GetID().String())
	return nil
}

func (a *adapter) RemoveRange(ctx context.Context, entities []core.Entitier) error {
	// Check context cancellation
	if err := ctx.Err(); err != nil {
		return err
	}

	a.setting.Log.Debugw("mock removerange")
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

	a.setting.Log.Debugw("mock add", "entity", entity)
	a.db.Set(entity.GetID().String(), entity)
	return nil
}

func (a *adapter) AddRange(ctx context.Context, entities []core.Entitier) error {
	// Check context cancellation
	if err := ctx.Err(); err != nil {
		return err
	}

	a.setting.Log.Debugw("mock addrange")
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

	a.setting.Log.Debugw("mock update", "entity", entity)
	a.db.Set(entity.GetID().String(), entity)
	return nil
}

func (a *adapter) UpdateRange(ctx context.Context, entities []core.Entitier) error {
	// Check context cancellation
	if err := ctx.Err(); err != nil {
		return err
	}

	a.setting.Log.Debugw("mock updaterange")
	for _, v := range entities {
		a.Update(ctx, v)
	}
	return nil
}
