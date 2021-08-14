package core

import (
	"time"

	"github.com/google/uuid"
)

type Entity struct {
	ID         uuid.UUID `gorm:"primaryKey" bson:"_id,omitempty"`
	CreateUser string    `bson:"create_user,omitempty"`
	UpdateUser string    `bson:"update_user,omitempty"`
	CreatedAt  time.Time `bson:"created_at,omitempty"`
	UpdatedAt  time.Time `bson:"updated_at,omitempty"`
}

type Entitier interface {
	GetID() uuid.UUID
	SetID(uuid.UUID)
	SetCreatedAt(time.Time)
	SetUpdatedAt(time.Time)
}

func (e *Entity) GetID() uuid.UUID {
	return e.ID
}

func (e *Entity) SetID(id uuid.UUID) {
	e.ID = id
}

func (e *Entity) SetCreatedAt(timestamp time.Time) {
	e.CreatedAt = timestamp
}

func (e *Entity) SetUpdatedAt(timestamp time.Time) {
	e.UpdatedAt = timestamp
}
