package mocks

import (
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

func (a *adapter) Has(key string) (bool, error) {
	core.Log.Info("mock has - ", key)
	return a.states.Has(key), nil
}

func (a *adapter) Get(key string, dest core.Entitier) (bool, error) {
	core.Log.Info("mock get - ", key)
	if resp, ok := a.states.Get(key); ok {
		model.Copy(dest, resp)
		return true, nil
	}
	return false, core.ErrNotFound
}

func (a *adapter) Set(key string, value interface{}) error {
	core.Log.Info("mock set - ", key, value)
	a.states.Set(key, value)
	return nil
}

func (a *adapter) Publish(coreEvent core.DomainEventer) error {
	core.Log.Info("mock publish - ", coreEvent.GetID(), coreEvent.GetTopic())
	return nil
}

func (a *adapter) Subscribe(topic string, handler core.ReplyHandler) error {
	core.Log.Info("mock subscribe - ", topic)
	a.pubsubs.Set(topic, handler)
	return nil
}

func (a *adapter) FakeSend(topic string, data interface{}) {
	core.Log.Info("mock fake send - ", topic, data)
	if resp, ok := a.pubsubs.Get(topic); ok {
		resp.(core.ReplyHandler)(data)
	}
}

func (a *adapter) Unsubscribe(topic string) error {
	core.Log.Info("mock unsubscribe - ", topic)
	a.pubsubs.Remove(topic)
	return nil
}

func (a *adapter) SetModel(model core.Entitier) {
	core.Log.Info("mock setmodel")
	a.model = model
}

func (a *adapter) Find(dto core.Entitier, id uuid.UUID) error {
	core.Log.Info("mock find - ", id)
	if resp, ok := a.db.Get(id.String()); ok {
		model.Copy(dto, resp)
		return nil
	}
	return core.ErrNotFound
}

func (a *adapter) Any() (bool, error) {
	core.Log.Info("mock any")
	count, err := a.Count()
	return count > 0, err
}

func (a *adapter) AnyWithFilter(query interface{}, args interface{}) (bool, error) {
	core.Log.Info("mock anywithfilter")
	count, err := a.CountWithFilter(query, args)
	return count > 0, err
}

func (a *adapter) Count() (int64, error) {
	core.Log.Info("mock count")
	count := a.db.Count()
	return int64(count), nil
}

func (a *adapter) CountWithFilter(query interface{}, args interface{}) (int64, error) {
	core.Log.Info("mock countwithfilter")
	count := a.db.Count()
	return int64(count), nil
}

func (a *adapter) List(dtos []core.Entitier) error {
	core.Log.Info("mock list")
	for _, v := range a.db.Items() {
		entity := v.(core.Entitier)
		dtos = append(dtos, entity)
	}

	return nil
}

func (a *adapter) ListWithFilter(dtos []core.Entitier, query interface{}, args interface{}) error {
	core.Log.Info("mock listwithfilter")
	for _, v := range a.db.Items() {
		entity := v.(core.Entitier)
		dtos = append(dtos, entity)
	}

	return nil
}

func (a *adapter) Remove(entity core.Entitier) error {
	core.Log.Info("mock remove - ", entity)
	a.db.Remove(entity.GetID().String())
	return nil
}

func (a *adapter) RemoveRange(entities []core.Entitier) error {
	core.Log.Info("mock removerange")
	for _, v := range entities {
		a.Remove(v)
	}
	return nil
}

func (a *adapter) Add(entity core.Entitier) error {
	core.Log.Info("mock add - ", entity)
	a.db.Set(entity.GetID().String(), entity)
	return nil
}

func (a *adapter) AddRange(entities []core.Entitier) error {
	core.Log.Info("mock addrange")
	for _, v := range entities {
		a.Add(v)
	}
	return nil
}

func (a *adapter) Update(entity core.Entitier) error {
	core.Log.Info("mock update - ", entity)
	a.db.Set(entity.GetID().String(), entity)
	return nil
}

func (a *adapter) UpdateRange(entities []core.Entitier) error {
	core.Log.Info("mock updaterange")
	for _, v := range entities {
		a.Update(v)
	}
	return nil
}
