package core

import (
	"time"

	"github.com/google/uuid"
)

type Entity struct {
	ID         uuid.UUID `gorm:"primaryKey" bson:"_id,omitempty"`
	CreateUser string
	UpdateUser string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Entitier interface {
	GetID() uuid.UUID
	SetID(uuid.UUID)
	SetCreatedAt(user string, timestamp time.Time)
	SetUpdatedAt(user string, timestamp time.Time)
}

func (e *Entity) GetID() uuid.UUID {
	return e.ID
}

func (e *Entity) SetID(id uuid.UUID) {
	e.ID = id
}

func (e *Entity) SetCreatedAt(user string, timestamp time.Time) {
	e.CreateUser = user
	e.CreatedAt = timestamp
}

func (e *Entity) SetUpdatedAt(user string, timestamp time.Time) {
	e.UpdateUser = user
	e.UpdatedAt = timestamp
}
