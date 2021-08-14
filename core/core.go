package core

import (
	"reflect"
	"sync"

	cmap "github.com/orcaman/concurrent-map"
	"github.com/sony/gobreaker"
	"go.uber.org/zap"
)

type singletons struct {
	mediator sync.Once
	eventBus sync.Once
	states   sync.Once
}

var syncs singletons

var mediatorInstance *Mediator

var eventBusInstance *EventBus

var statesInstance *StateService

var repositories cmap.ConcurrentMap

var Log *zap.SugaredLogger

func init() {
	logger, _ := zap.NewProduction()
	Log = logger.Sugar()

	repositories = cmap.New()
}

func GetMediator() *Mediator {
	if mediatorInstance == nil {
		syncs.mediator.Do(
			func() {
				mediatorInstance = &Mediator{
					requestHandlers:      cmap.New(),
					notificationHandlers: cmap.New(),
				}
				Log.Info("mediator created")
			})
	}
	return mediatorInstance
}

func GetEventBus() *EventBus {
	if eventBusInstance == nil {
		syncs.eventBus.Do(
			func() {
				var st gobreaker.Settings
				st.Name = "eventbus"

				eventBusInstance = &EventBus{
					mediator:     GetMediator(),
					domainEvents: make([]DomainEventer, 0),
					cb:           gobreaker.NewCircuitBreaker(st),
				}
				Log.Info("eventbus created")
			})
	}
	return eventBusInstance
}

func GetStateService() *StateService {
	if statesInstance == nil {
		syncs.states.Do(
			func() {
				var st gobreaker.Settings
				st.Name = "state"

				statesInstance = &StateService{
					cb: gobreaker.NewCircuitBreaker(st),
				}
				Log.Info("state created")
			})
	}
	return statesInstance
}

func GetRepositoryService(model Entitier) *RepositoryService {
	valueOf := reflect.ValueOf(model)
	key := valueOf.Type().Name()

	if !repositories.Has(key) {
		var st gobreaker.Settings
		st.Name = key + "repository"
		value := &RepositoryService{
			model: model,
			cb:    gobreaker.NewCircuitBreaker(st),
		}
		repositories.Set(key, value)
	}

	if value, ok := repositories.Get(key); ok {
		return value.(*RepositoryService)
	}

	panic("repository not found exception")
}
