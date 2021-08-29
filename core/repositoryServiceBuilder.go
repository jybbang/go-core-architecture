package core

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/google/uuid"
)

// Builder Object for RepositoryService
type repositoryServiceBuilder struct {
	tableName         string
	userIdKey         string
	model             Entitier
	queryRepository   queryRepositoryAdapter
	commandRepository commandRepositoryAdapter
	cbSettings        CircuitBreakerSettings
}

// Constructor for RepositoryServiceBuilder
func NewRepositoryServiceBuilder(model Entitier, tableName string) *repositoryServiceBuilder {
	if model == nil {
		panic("model is required")
	}
	if strings.TrimSpace(tableName) == "" {
		panic("tableName is required")
	}

	o := new(repositoryServiceBuilder)
	o.tableName = tableName
	o.model = model
	o.model.SetID(uuid.Nil)
	o.userIdKey = http.CanonicalHeaderKey("Userid")
	o.cbSettings = CircuitBreakerSettings{Name: o.tableName + "-repository"}

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
		tableName:         b.tableName,
		userIdKey:         b.userIdKey,
		queryRepository:   b.queryRepository,
		commandRepository: b.commandRepository,
	}
	instance.initialize(b.cbSettings)

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

// Builder method to set the field queryRepository in RepositoryServiceBuilder
func (b *repositoryServiceBuilder) UserIdKeyInContext(useridKey string) *repositoryServiceBuilder {
	if strings.TrimSpace(useridKey) == "" {
		panic("useridKey is required")
	}

	b.userIdKey = http.CanonicalHeaderKey(useridKey)
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
	b.cbSettings = setting
	return b
}
