package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Entity struct {
	ID         uuid.UUID      `gorm:"primaryKey" bson:"_id,omitempty"`
	CreateUser string         `bson:"create_user,omitempty"`
	UpdateUser string         `bson:"update_user,omitempty"`
	CreatedAt  time.Time      `bson:"created_at,omitempty"`
	UpdatedAt  time.Time      `bson:"updated_at,omitempty"`
	DeletedAt  gorm.DeletedAt `gorm:"index" bson:"deleted_at,omitempty"`
}

type Entitier interface {
	GetID() uuid.UUID
	SetID(uuid.UUID)
	CopyWith(interface{}) bool
}

func (e *Entity) GetID() uuid.UUID {
	return e.ID
}

func (e *Entity) SetID(id uuid.UUID) {
	e.ID = id
}

func (e *Entity) CopyWith(src interface{}) bool {
	source, ok := src.(*Entity)
	e.ID = source.ID
	e.CreateUser = source.CreateUser
	e.UpdateUser = source.UpdateUser
	e.CreatedAt = source.CreatedAt
	e.UpdatedAt = source.UpdatedAt
	e.DeletedAt = source.DeletedAt
	return ok
}
