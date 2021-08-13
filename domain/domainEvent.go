package domain

import (
	"time"

	"github.com/google/uuid"
)

type DomainEvent struct {
	ID                      uuid.UUID
	EventID                 uuid.UUID
	Topic                   string
	CanPublishToEventsource bool
	IsPublished             bool
	CreatedAt               time.Time
	PublishedAt             time.Time
}

type DomainEventer interface {
	GetID() uuid.UUID
	GetEventID() uuid.UUID
	GetTopic() string
	GetCanPublishToEventsource() bool
	AddingEvent()
	PublishingEvent(time.Time)
}

func (e *DomainEvent) GetID() uuid.UUID {
	return e.ID
}

func (e *DomainEvent) GetEventID() uuid.UUID {
	return e.EventID
}

func (e *DomainEvent) GetTopic() string {
	return e.Topic
}

func (e *DomainEvent) GetCanPublishToEventsource() bool {
	return e.CanPublishToEventsource
}

func (e *DomainEvent) AddingEvent() {
	e.EventID = uuid.New()
	e.CreatedAt = time.Now()
}

func (e *DomainEvent) PublishingEvent(requestAt time.Time) {
	e.IsPublished = true
	e.PublishedAt = requestAt
}
