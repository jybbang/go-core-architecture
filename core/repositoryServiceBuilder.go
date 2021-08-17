package core

import (
	"reflect"

	"github.com/google/uuid"
	"github.com/sony/gobreaker"
)

// Builder Object for RepositoryService
type repositoryServiceBuilder struct {
	tableName         string
	model             Entitier
	queryRepository   queryRepositoryAdapter
	commandRepository commandRepositoryAdapter
	cb                *gobreaker.CircuitBreaker
}

// Constructor for RepositoryServiceBuilder
func NewRepositoryServiceBuilder(model Entitier, tableName string) *repositoryServiceBuilder {
	if model == nil {
		panic("model is required")
	}
	if tableName == "" {
		panic("tableName is required")
	}

	o := new(repositoryServiceBuilder)
	o.tableName = tableName
	o.model = model
	o.model.SetID(uuid.Nil)

	st := gobreaker.Settings{
		Name: tableName + "-repository",
	}
	o.cb = gobreaker.NewCircuitBreaker(st)

	return o
}

// Build Method which creates RepositoryService
func (b *repositoryServiceBuilder) Build() *repositoryService {
	typeOf := reflect.TypeOf(b.model)
	key := typeOf.Elem().Name()

	if repositories.Has(key) {
		panic("this repository service already created")
	}

	repository := b.Create()

	repositories.Set(key, repository)
	return repository
}

// Build Method which creates EventBus
func (b *repositoryServiceBuilder) Create() *repositoryService {
	if b.queryRepository == nil {
		panic("queryRepository adapter is required")
	}
	if b.commandRepository == nil {
		panic("commandRepository adapter is required")
	}

	instance := &repositoryService{
		queryRepository:   b.queryRepository,
		commandRepository: b.commandRepository,
		cb:                b.cb,
	}
	instance.initialize()

	return instance
}

// Builder method to set the field queryRepository in RepositoryServiceBuilder
func (b *repositoryServiceBuilder) QueryRepositoryAdapter(adapter queryRepositoryAdapter) *repositoryServiceBuilder {
	if adapter == nil {
		panic("adapter is required")
	}

	b.queryRepository = adapter
	b.queryRepository.SetModel(b.model, b.tableName)
	return b
}

// Builder method to set the field commandRepository in RepositoryServiceBuilder
func (b *repositoryServiceBuilder) CommandRepositoryAdapter(adapter commandRepositoryAdapter) *repositoryServiceBuilder {
	if adapter == nil {
		panic("adapter is required")
	}

	b.commandRepository = adapter
	b.commandRepository.SetModel(b.model, b.tableName)
	return b
}

// Builder method to set the field messaging in RepositoryServiceBuilder
func (b *repositoryServiceBuilder) CircuitBreaker(setting CircuitBreakerSettings) *repositoryServiceBuilder {
	b.cb = gobreaker.NewCircuitBreaker(setting.toGobreakerSettings(b.cb.Name()))
	return b
}
