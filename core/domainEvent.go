package core

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type DomainEvent struct {
	ID                         uuid.UUID
	EventID                    uuid.UUID
	Topic                      string
	CanNotPublishToEventsource bool
	IsPublished                bool
	CanBuffered                bool
	CreatedAt                  time.Time
	PublishedAt                time.Time
}

type DomainEventer interface {
	GetID() uuid.UUID
	GetEventID() uuid.UUID
	GetTopic() string
	GetCanNotPublishToEventsource() bool
	GetCanBuffered() bool
	SetAddingEvent()
	SetPublishingEvent(ctx context.Context, publishedAt time.Time)
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

func (e *DomainEvent) GetCanNotPublishToEventsource() bool {
	return e.CanNotPublishToEventsource
}

func (e *DomainEvent) GetCanBuffered() bool {
	return e.CanBuffered
}

func (e *DomainEvent) SetAddingEvent() {
	e.EventID = uuid.New()
	e.CreatedAt = time.Now()
}

func (e *DomainEvent) SetPublishingEvent(ctx context.Context, publishedAt time.Time) {
	e.IsPublished = true
	e.PublishedAt = publishedAt
}
