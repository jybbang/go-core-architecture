package mocks

import (
	"log"

	"github.com/google/uuid"
	cmap "github.com/orcaman/concurrent-map"

	"github.com/jybbang/core-architecture/application/contracts"
	"github.com/jybbang/core-architecture/domain"
)

type adapter struct {
	model   domain.Entity
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
	log.Println("mock has", key)
	return a.states.Has(key), nil
}

func (a *adapter) Get(key string, dest domain.Copyable) (bool, error) {
	log.Println("mock get", key)
	if resp, ok := a.states.Get(key); ok {
		dest.CopyWith(resp)
		return true, nil
	}
	return false, domain.ErrNotFound
}

func (a *adapter) Set(key string, value interface{}) error {
	log.Println("mock set", key, value)
	a.states.Set(key, value)
	return nil
}

func (a *adapter) Publish(domainEvent *domain.DomainEvent) error {
	log.Println("mock publish", domainEvent.ID, domainEvent.Topic)
	return nil
}

func (a *adapter) Subscribe(topic string, handler contracts.ReplyHandler) error {
	log.Println("mock subscribe", topic)
	a.pubsubs.Set(topic, handler)
	return nil
}

func (a *adapter) Unsubscribe(topic string) error {
	log.Println("mock unsubscribe", topic)
	a.pubsubs.Remove(topic)
	return nil
}

func (a *adapter) SetModel(model domain.Entity) {
	log.Println("mock setmodel")
	a.model = model
}

func (a *adapter) Find(dto domain.Copyable, id uuid.UUID) error {
	log.Println("mock find", id)
	if resp, ok := a.db.Get(id.String()); ok {
		dto.CopyWith(resp)
		return nil
	}
	return domain.ErrNotFound
}

func (a *adapter) Any() (bool, error) {
	log.Println("mock any")
	count, err := a.Count()
	return count > 0, err
}

func (a *adapter) AnyWithFilter(query interface{}, args interface{}) (bool, error) {
	log.Println("mock anywithfilter")
	count, err := a.CountWithFilter(query, args)
	return count > 0, err
}

func (a *adapter) Count() (int64, error) {
	log.Println("mock count")
	count := a.db.Count()
	return int64(count), nil
}

func (a *adapter) CountWithFilter(query interface{}, args interface{}) (int64, error) {
	log.Println("mock countwithfilter")
	count := a.db.Count()
	return int64(count), nil
}

func (a *adapter) List(dtos []domain.Copyable) error {
	log.Println("mock list")
	for _, v := range a.db.Items() {
		entity := v.(domain.Copyable)
		dtos = append(dtos, entity)
	}

	return nil
}

func (a *adapter) ListWithFilter(dtos []domain.Copyable, query interface{}, args interface{}) error {
	log.Println("mock listwithfilter")
	for _, v := range a.db.Items() {
		entity := v.(domain.Copyable)
		dtos = append(dtos, entity)
	}

	return nil
}

func (a *adapter) Remove(entity *domain.Entity) error {
	log.Println("mock remove", entity)
	a.db.Remove(entity.ID.String())
	return nil
}

func (a *adapter) RemoveRange(entities []*domain.Entity) error {
	log.Println("mock removerange")
	for _, v := range entities {
		a.Remove(v)
	}
	return nil
}

func (a *adapter) Add(entity *domain.Entity) error {
	log.Println("mock add", entity)
	a.db.Set(entity.ID.String(), entity)
	return nil
}

func (a *adapter) AddRange(entities []*domain.Entity) error {
	log.Println("mock addrange")
	for _, v := range entities {
		a.Add(v)
	}
	return nil
}

func (a *adapter) Update(entity *domain.Entity) error {
	log.Println("mock update", entity)
	return a.Add(entity)
}

func (a *adapter) UpdateRange(entities []*domain.Entity) error {
	log.Println("mock updaterange")
	return a.AddRange(entities)
}
