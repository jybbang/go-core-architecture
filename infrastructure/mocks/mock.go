package mocks

import (
	"context"
	"reflect"
	"sync/atomic"

	"github.com/google/uuid"
	cmap "github.com/orcaman/concurrent-map"
	"go.uber.org/zap"
	"gopkg.in/jeevatkm/go-model.v1"

	"github.com/jybbang/go-core-architecture/core"
)

type adapter struct {
	tableName      string
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
	logger, _ := zap.NewProduction()

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

func (a *adapter) IsConnected() bool {
	return true
}

func (a *adapter) Connect(ctx context.Context) error {
	defer a.setting.Log.Debug("mock connect")

	return nil
}

func (a *adapter) Disconnect() {
	defer a.setting.Log.Debug("mock disconnect")

	a.db.Clear()

	a.pubsubs.Clear()

	a.states.Clear()
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

func (a *adapter) Has(ctx context.Context, key string) bool {
	// Check context cancellation
	if err := ctx.Err(); err != nil {
		return false
	}

	defer a.setting.Log.Debugw("mock has", "key", key)

	return a.states.Has(key)
}

func (a *adapter) Get(ctx context.Context, key string, dest interface{}) error {
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

func (a *adapter) SetModel(model core.Entitier, tableName string) {
	defer a.setting.Log.Debugw("mock setmodel", "model", model, "tableName", tableName)

	a.model = model
	a.tableName = tableName
}

func (a *adapter) Find(ctx context.Context, id uuid.UUID, dest core.Entitier) (err error) {
	// Check context cancellation
	if err = ctx.Err(); err != nil {
		return err
	}

	if resp, ok := a.db.Get(id.String()); ok {
		model.Copy(dest, resp)

		defer a.setting.Log.Debugw("mock find", "id", id, "dest", dest)
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

	defer a.setting.Log.Debugw("mock any with filter", "ok", ok)

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

	defer a.setting.Log.Debugw("mock count with filter", "count", resp, "query", query, "args", args)

	return int64(resp), nil
}

func (a *adapter) List(ctx context.Context, dest interface{}) (err error) {
	// Check context cancellation
	if err := ctx.Err(); err != nil {
		return err
	}

	resultsVal := reflect.ValueOf(dest)

	sliceVal := resultsVal.Elem()

	if sliceVal.Kind() == reflect.Interface {
		sliceVal = sliceVal.Elem()
	}

	for _, v := range a.db.Items() {
		entityVal := reflect.ValueOf(v)

		sliceVal = reflect.Append(sliceVal, entityVal)
	}

	resultsVal.Elem().Set(sliceVal)

	defer a.setting.Log.Debugw("mock list", "dest", dest)

	return nil
}

func (a *adapter) ListWithFilter(ctx context.Context, query interface{}, args interface{}, dest interface{}) (err error) {
	// Check context cancellation
	if err := ctx.Err(); err != nil {
		return err
	}

	resultsVal := reflect.ValueOf(dest)

	sliceVal := resultsVal.Elem()

	if sliceVal.Kind() == reflect.Interface {
		sliceVal = sliceVal.Elem()
	}

	for _, v := range a.db.Items() {
		entityVal := reflect.ValueOf(v)

		sliceVal = reflect.Append(sliceVal, entityVal)
	}

	resultsVal.Elem().Set(sliceVal)

	defer a.setting.Log.Debugw("mock list with filter", "dest", dest, "query", query, "args", args)

	return nil
}

func (a *adapter) Remove(ctx context.Context, id uuid.UUID) error {
	// Check context cancellation
	if err := ctx.Err(); err != nil {
		return err
	}

	defer a.setting.Log.Debugw("mock remove", "id", id)

	a.db.Remove(id.String())

	return nil
}

func (a *adapter) RemoveRange(ctx context.Context, ids []uuid.UUID) error {
	// Check context cancellation
	if err := ctx.Err(); err != nil {
		return err
	}

	defer a.setting.Log.Debugw("mock remove range")

	for _, id := range ids {
		a.Remove(ctx, id)
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

	defer a.setting.Log.Debugw("mock add range")

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

	defer a.setting.Log.Debugw("mock update range")

	for _, v := range entities {
		a.Update(ctx, v)
	}

	return nil
}
