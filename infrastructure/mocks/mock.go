package mocks

import (
	"github.com/google/uuid"
	cmap "github.com/orcaman/concurrent-map"

	"github.com/jybbang/go-core-architecture/application"
	"github.com/jybbang/go-core-architecture/application/contracts"
	"github.com/jybbang/go-core-architecture/domain"
	"github.com/jybbang/go-core-architecture/infrastructure"
)

type adapter struct {
	model   domain.Entitier
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
	infrastructure.Log.Debug("mock has", key)
	return a.states.Has(key), nil
}

func (a *adapter) Get(key string, dest domain.Entitier) (bool, error) {
	infrastructure.Log.Debug("mock get", key)
	if resp, ok := a.states.Get(key); ok {
		dest.CopyWith(resp)
		return true, nil
	}
	return false, application.ErrNotFound
}

func (a *adapter) Set(key string, value interface{}) error {
	infrastructure.Log.Debug("mock set", key, value)
	a.states.Set(key, value)
	return nil
}

func (a *adapter) Publish(domainEvent domain.DomainEventer) error {
	infrastructure.Log.Debug("mock publish", domainEvent.GetID(), domainEvent.GetTopic())
	return nil
}

func (a *adapter) Subscribe(topic string, handler contracts.ReplyHandler) error {
	infrastructure.Log.Debug("mock subscribe", topic)
	a.pubsubs.Set(topic, handler)
	return nil
}

func (a *adapter) Unsubscribe(topic string) error {
	infrastructure.Log.Debug("mock unsubscribe", topic)
	a.pubsubs.Remove(topic)
	return nil
}

func (a *adapter) SetModel(model domain.Entitier) {
	infrastructure.Log.Debug("mock setmodel")
	a.model = model
}

func (a *adapter) Find(dto domain.Entitier, id uuid.UUID) error {
	infrastructure.Log.Debug("mock find", id)
	if resp, ok := a.db.Get(id.String()); ok {
		dto.CopyWith(resp)
		return nil
	}
	return application.ErrNotFound
}

func (a *adapter) Any() (bool, error) {
	infrastructure.Log.Debug("mock any")
	count, err := a.Count()
	return count > 0, err
}

func (a *adapter) AnyWithFilter(query interface{}, args interface{}) (bool, error) {
	infrastructure.Log.Debug("mock anywithfilter")
	count, err := a.CountWithFilter(query, args)
	return count > 0, err
}

func (a *adapter) Count() (int64, error) {
	infrastructure.Log.Debug("mock count")
	count := a.db.Count()
	return int64(count), nil
}

func (a *adapter) CountWithFilter(query interface{}, args interface{}) (int64, error) {
	infrastructure.Log.Debug("mock countwithfilter")
	count := a.db.Count()
	return int64(count), nil
}

func (a *adapter) List(dtos []domain.Entitier) error {
	infrastructure.Log.Debug("mock list")
	for _, v := range a.db.Items() {
		entity := v.(domain.Entitier)
		dtos = append(dtos, entity)
	}

	return nil
}

func (a *adapter) ListWithFilter(dtos []domain.Entitier, query interface{}, args interface{}) error {
	infrastructure.Log.Debug("mock listwithfilter")
	for _, v := range a.db.Items() {
		entity := v.(domain.Entitier)
		dtos = append(dtos, entity)
	}

	return nil
}

func (a *adapter) Remove(entity domain.Entitier) error {
	infrastructure.Log.Debug("mock remove", entity)
	a.db.Remove(entity.GetID().String())
	return nil
}

func (a *adapter) RemoveRange(entities []domain.Entitier) error {
	infrastructure.Log.Debug("mock removerange")
	for _, v := range entities {
		a.Remove(v)
	}
	return nil
}

func (a *adapter) Add(entity domain.Entitier) error {
	infrastructure.Log.Debug("mock add", entity)
	a.db.Set(entity.GetID().String(), entity)
	return nil
}

func (a *adapter) AddRange(entities []domain.Entitier) error {
	infrastructure.Log.Debug("mock addrange")
	for _, v := range entities {
		a.Add(v)
	}
	return nil
}

func (a *adapter) Update(entity domain.Entitier) error {
	infrastructure.Log.Debug("mock update", entity)
	return a.Add(entity)
}

func (a *adapter) UpdateRange(entities []domain.Entitier) error {
	infrastructure.Log.Debug("mock updaterange")
	return a.AddRange(entities)
}
