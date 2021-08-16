package mocks

import (
	"context"
	"sync/atomic"

	"github.com/google/uuid"
	cmap "github.com/orcaman/concurrent-map"
	"go.uber.org/zap"
	"gopkg.in/jeevatkm/go-model.v1"

	"github.com/jybbang/go-core-architecture/core"
)

type adapter struct {
	model          core.Entitier
	db             cmap.ConcurrentMap
	pubsubs        cmap.ConcurrentMap
	states         cmap.ConcurrentMap
	setting        MockSettings
	publishedCount uint32
}

type MockSettings struct {
	Log *zap.SugaredLogger
}

func NewMockAdapter() *adapter {
	logger, _ := zap.NewDevelopment()

	return &adapter{
		db:      cmap.New(),
		pubsubs: cmap.New(),
		states:  cmap.New(),
		setting: MockSettings{
			Log: logger.Sugar(),
		},
	}
}

func NewMockAdapterWithSettings(setting MockSettings) *adapter {

	return &adapter{
		db:      cmap.New(),
		pubsubs: cmap.New(),
		states:  cmap.New(),
		setting: setting,
	}
}

func (a *adapter) GetPublishedCount() uint32 {
	return a.publishedCount
}

func (a *adapter) GetDbCount() int {
	return a.db.Count()
}

func (a *adapter) GetPubsubsCount() int {
	return a.pubsubs.Count()
}

func (a *adapter) GetStatesCount() int {
	return a.states.Count()
}

func (a *adapter) Has(ctx context.Context, key string) (ok bool, err error) {
	// Check context cancellation
	if err := ctx.Err(); err != nil {
		return false, err
	}

	ok = a.states.Has(key)

	defer a.setting.Log.Debugw("mock has", "key", key, "ok", ok)

	return ok, nil
}

func (a *adapter) Get(ctx context.Context, key string, dest interface{}) (err error) {
	// Check context cancellation
	if err := ctx.Err(); err != nil {
		return err
	}

	if resp, ok := a.states.Get(key); ok {
		model.Copy(dest, resp)
		defer a.setting.Log.Debugw("mock get", "key", key, "result", dest)
		return nil
	}

	defer a.setting.Log.Debugw("mock get", "key", key)

	return core.ErrNotFound
}

func (a *adapter) Set(ctx context.Context, key string, value interface{}) error {
	// Check context cancellation
	if err := ctx.Err(); err != nil {
		return err
	}

	defer a.setting.Log.Debugw("mock set", "key", key, "value", value)

	a.states.Set(key, value)
	return nil
}

func (a *adapter) Delete(ctx context.Context, key string) error {
	// Check context cancellation
	if err := ctx.Err(); err != nil {
		return err
	}

	defer a.setting.Log.Debugw("mock delete", "key", key)

	a.states.Remove(key)
	return nil
}

func (a *adapter) Publish(ctx context.Context, coreEvent core.DomainEventer) error {
	// Check context cancellation
	if err := ctx.Err(); err != nil {
		return err
	}

	defer a.setting.Log.Debugw("mock publish", "id", coreEvent.GetID(), "event", coreEvent)

	atomic.AddUint32(&a.publishedCount, 1)
	return nil
}

func (a *adapter) Subscribe(ctx context.Context, topic string, handler core.ReplyHandler) error {
	// Check context cancellation
	if err := ctx.Err(); err != nil {
		return err
	}

	defer a.setting.Log.Debugw("mock subscribe", "topic", topic)

	a.pubsubs.Set(topic, handler)
	return nil
}

func (a *adapter) Unsubscribe(ctx context.Context, topic string) error {
	// Check context cancellation
	if err := ctx.Err(); err != nil {
		return err
	}

	defer a.setting.Log.Debugw("mock unsubscribe", "topic", topic)

	a.pubsubs.Remove(topic)
	return nil
}

func (a *adapter) FakeSend(topic string, receivedData interface{}) {
	defer a.setting.Log.Debugw("mock fake send", "topic", topic, "data", receivedData)

	if resp, ok := a.pubsubs.Get(topic); ok {
		resp.(core.ReplyHandler)(receivedData)
	}
}

func (a *adapter) SetModel(model core.Entitier) {
	defer a.setting.Log.Debugw("mock setmodel", "model", model)

	a.model = model
}

func (a *adapter) Find(ctx context.Context, id uuid.UUID, dest core.Entitier) (err error) {
	// Check context cancellation
	if err = ctx.Err(); err != nil {
		return err
	}

	if resp, ok := a.db.Get(id.String()); ok {
		defer a.setting.Log.Debugw("mock find", "id", id, "result", resp)

		model.Copy(dest, resp)
		return nil
	}

	defer a.setting.Log.Debugw("mock find", "id", id)

	return core.ErrNotFound
}

func (a *adapter) Any(ctx context.Context) (ok bool, err error) {
	// Check context cancellation
	if err = ctx.Err(); err != nil {
		return false, err
	}

	count, err := a.Count(ctx)

	ok = count > 0

	defer a.setting.Log.Debugw("mock any", "ok", ok)

	return ok, err
}

func (a *adapter) AnyWithFilter(ctx context.Context, query interface{}, args interface{}) (ok bool, err error) {
	// Check context cancellation
	if err = ctx.Err(); err != nil {
		return false, err
	}

	count, err := a.CountWithFilter(ctx, query, args)

	ok = count > 0

	defer a.setting.Log.Debugw("mock anywithfilter", "ok", ok)

	return ok, err
}

func (a *adapter) Count(ctx context.Context) (count int64, err error) {
	// Check context cancellation
	if err = ctx.Err(); err != nil {
		return 0, err
	}

	resp := a.db.Count()

	defer a.setting.Log.Debugw("mock count", "count", resp)

	return int64(resp), nil
}

func (a *adapter) CountWithFilter(ctx context.Context, query interface{}, args interface{}) (count int64, err error) {
	// Check context cancellation
	if err = ctx.Err(); err != nil {
		return 0, err
	}

	resp := a.db.Count()

	defer a.setting.Log.Debugw("mock countwithfilter", "count", resp, "query", query, "args", args)

	return int64(resp), nil
}

func (a *adapter) List(ctx context.Context) (result []core.Entitier, err error) {
	// Check context cancellation
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	for _, v := range a.db.Items() {
		entity := v.(core.Entitier)
		result = append(result, entity)
	}

	defer a.setting.Log.Debugw("mock list", "result", result)

	return result, nil
}

func (a *adapter) ListWithFilter(ctx context.Context, query interface{}, args interface{}) (result []core.Entitier, err error) {
	// Check context cancellation
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	for _, v := range a.db.Items() {
		entity := v.(core.Entitier)
		result = append(result, entity)
	}

	defer a.setting.Log.Debugw("mock listwithfilter", "result", result, "query", query, "args", args)

	return result, nil
}

func (a *adapter) Remove(ctx context.Context, entity core.Entitier) error {
	// Check context cancellation
	if err := ctx.Err(); err != nil {
		return err
	}

	defer a.setting.Log.Debugw("mock remove", "entity", entity)

	a.db.Remove(entity.GetID().String())
	return nil
}

func (a *adapter) RemoveRange(ctx context.Context, entities []core.Entitier) error {
	// Check context cancellation
	if err := ctx.Err(); err != nil {
		return err
	}

	defer a.setting.Log.Debugw("mock removerange")

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

	defer a.setting.Log.Debugw("mock add", "entity", entity)

	a.db.Set(entity.GetID().String(), entity)
	return nil
}

func (a *adapter) AddRange(ctx context.Context, entities []core.Entitier) error {
	// Check context cancellation
	if err := ctx.Err(); err != nil {
		return err
	}

	defer a.setting.Log.Debugw("mock addrange")

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

	defer a.setting.Log.Debugw("mock update", "entity", entity)

	a.db.Set(entity.GetID().String(), entity)
	return nil
}

func (a *adapter) UpdateRange(ctx context.Context, entities []core.Entitier) error {
	// Check context cancellation
	if err := ctx.Err(); err != nil {
		return err
	}

	defer a.setting.Log.Debugw("mock updaterange")

	for _, v := range entities {
		a.Update(ctx, v)
	}
	return nil
}
