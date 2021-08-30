package core

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/google/uuid"
	"gopkg.in/jeevatkm/go-model.v1"
)

// Builder Object for RepositoryService
type repositoryServiceBuilder struct {
	tableName         string
	userIdKey         string
	model             Entitier
	queryRepository   queryRepositoryAdapter
	commandRepository commandRepositoryAdapter
	cbSettings        CircuitBreakerSettings
	settings          RepositoryServiceSettings
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
	o.userIdKey = http.CanonicalHeaderKey("Userid")
	o.cbSettings = CircuitBreakerSettings{
		AllowedRequestInHalfOpen: 1,
		DurationOfBreak:          time.Duration(60 * time.Second),
		SamplingDuration:         time.Duration(60 * time.Second),
		SamplingFailureCount:     5,
	}
	o.settings = RepositoryServiceSettings{
		ConnectionTimeout: time.Duration(10 * time.Second),
	}

	o.model.SetID(uuid.Nil)

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
		settings:          b.settings,
	}

	instance.cb = b.cbSettings.ToCircuitBreaker(b.tableName+"-repository", instance.onCircuitOpen)

	instance.initialize()

	return instance
}

// Builder method to set the field messaging in EventBusBuilder
func (b *repositoryServiceBuilder) Settings(settings RepositoryServiceSettings) *repositoryServiceBuilder {
	err := model.Copy(&b.settings, settings)

	if err != nil {
		panic(fmt.Errorf("settings mapping errors occurred: %v", err))
	}

	return b
}

// Builder method to set the field messaging in RepositoryServiceBuilder
func (b *repositoryServiceBuilder) CircuitBreaker(settings CircuitBreakerSettings) *repositoryServiceBuilder {
	err := model.Copy(&b.cbSettings, settings)

	if err != nil {
		panic(fmt.Errorf("cb settings mapping errors occurred: %v", err))
	}

	return b
}

// Builder method to set the field queryRepository in RepositoryServiceBuilder
func (b *repositoryServiceBuilder) UserIdKeyInContext(useridKey string) *repositoryServiceBuilder {
	if strings.TrimSpace(useridKey) == "" {
		panic("useridKey is required")
	}

	b.userIdKey = http.CanonicalHeaderKey(useridKey)

	return b
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
