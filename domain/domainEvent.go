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

func (e *DomainEvent) AddingEvent() {
	e.EventID = uuid.New()
	e.CreatedAt = time.Now()
}

func (e *DomainEvent) PublishingEvent(requestAt time.Time) {
	e.IsPublished = true
	e.PublishedAt = requestAt
}
