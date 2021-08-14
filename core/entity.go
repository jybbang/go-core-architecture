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
	SetCreatedAt(string, time.Time)
	SetUpdatedAt(string, time.Time)
	CopyWith(interface{}) bool
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

func (e *Entity) CopyWith(src interface{}) bool {
	source, ok := src.(*Entity)
	e.ID = source.ID
	e.CreateUser = source.CreateUser
	e.UpdateUser = source.UpdateUser
	e.CreatedAt = source.CreatedAt
	e.UpdatedAt = source.UpdatedAt
	return ok
}
