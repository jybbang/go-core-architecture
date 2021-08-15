package core

import (
	"reflect"

	"github.com/google/uuid"
	"github.com/sony/gobreaker"
)

// Builder Object for RepositoryService
type repositoryServiceBuilder struct {
	key               string
	model             Entitier
	queryRepository   queryRepositoryAdapter
	commandRepository commandRepositoryAdapter
	cb                *gobreaker.CircuitBreaker
}

// Constructor for RepositoryServiceBuilder
func NewRepositoryServiceBuilder(model Entitier) *repositoryServiceBuilder {
	valueOf := reflect.ValueOf(model)
	key := valueOf.Type().Name()

	o := new(repositoryServiceBuilder)
	o.key = key
	o.model = model

	st := gobreaker.Settings{
		Name:          key + "-repository",
		OnStateChange: onCbStateChange,
	}
	o.cb = gobreaker.NewCircuitBreaker(st)

	return o
}

// Build Method which creates RepositoryService
func (b *repositoryServiceBuilder) Build() *repositoryService {
	if repositories.Has(b.key) {
		panic("this repository service already created")
	}

	repository := &repositoryService{
		model:             b.model,
		queryRepository:   b.queryRepository,
		commandRepository: b.commandRepository,
		cb:                b.cb,
	}
	repository.initialize()

	repositories.Set(b.key, repository)
	return repository
}

// Builder method to set the field queryRepository in RepositoryServiceBuilder
func (b *repositoryServiceBuilder) QueryRepositoryAdapter(adapter queryRepositoryAdapter) *repositoryServiceBuilder {
	b.queryRepository = adapter
	b.model.SetID(uuid.Nil)
	b.queryRepository.SetModel(b.model)
	return b
}

// Builder method to set the field commandRepository in RepositoryServiceBuilder
func (b *repositoryServiceBuilder) CommandRepositoryAdapter(adapter commandRepositoryAdapter) *repositoryServiceBuilder {
	b.commandRepository = adapter
	b.model.SetID(uuid.Nil)
	b.commandRepository.SetModel(b.model)
	return b
}

// Builder method to set the field messaging in RepositoryServiceBuilder
func (b *repositoryServiceBuilder) CircuitBreaker(setting gobreaker.Settings) *repositoryServiceBuilder {
	setting.Name = b.cb.Name()
	setting.OnStateChange = onCbStateChange
	b.cb = gobreaker.NewCircuitBreaker(setting)
	return b
}
