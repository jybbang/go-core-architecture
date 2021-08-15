package core

import (
	"time"

	"github.com/google/uuid"
)

type Entity struct {
	ID         uuid.UUID `validate:"required" gorm:"primaryKey" bson:"_id,omitempty"`
	CreateUser string    `bson:"create_user,omitempty"`
	UpdateUser string    `bson:"update_user,omitempty"`
	CreatedAt  time.Time `bson:"created_at,omitempty"`
	UpdatedAt  time.Time `bson:"updated_at,omitempty"`
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
	// TODO add create user from context
	e.CreatedAt = timestamp
}

func (e *Entity) SetUpdatedAt(user string, timestamp time.Time) {
	// TODO add update user from context
	e.UpdatedAt = timestamp
}
