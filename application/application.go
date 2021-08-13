package application

import (
	"reflect"
	"sync"

	"github.com/jybbang/go-core-architecture/domain"
	cmap "github.com/orcaman/concurrent-map"
	"github.com/sony/gobreaker"
	"go.uber.org/zap"
)

var log *zap.SugaredLogger

func init() {
	logger, _ := zap.NewProduction()
	log = logger.Sugar()
}

type singletons struct {
	mediator sync.Once
	eventBus sync.Once
	states   sync.Once
}

var syncs singletons

var mediatorInstance *mediator

var eventBusInstance *eventBus

var statesInstance *stateService

var repositories cmap.ConcurrentMap

func GetMediator() *mediator {
	if mediatorInstance == nil {
		syncs.mediator.Do(
			func() {
				mediatorInstance = &mediator{
					requestHandlers:      cmap.New(),
					notificationHandlers: cmap.New(),
				}
				log.Info("mediator created")
			})
	}
	return mediatorInstance
}

func GetEventBus() *eventBus {
	if eventBusInstance == nil {
		syncs.eventBus.Do(
			func() {
				var st gobreaker.Settings
				st.Name = "eventbus"

				eventBusInstance = &eventBus{
					mediator:     GetMediator(),
					domainEvents: make([]*domain.DomainEvent, 0),
					cb:           gobreaker.NewCircuitBreaker(st),
				}
				log.Info("eventbus created")
			})
	}
	return eventBusInstance
}

func GetStateService() *stateService {
	if statesInstance == nil {
		syncs.states.Do(
			func() {
				var st gobreaker.Settings
				st.Name = "state"

				statesInstance = &stateService{
					cb: gobreaker.NewCircuitBreaker(st),
				}
				log.Info("state created")
			})
	}
	return statesInstance
}

func GetRepositoryService(model *domain.Entity) *repositoryService {
	valueOf := reflect.ValueOf(model)
	key := valueOf.Type().Name()

	if !repositories.Has(key) {
		var st gobreaker.Settings
		st.Name = key + "repository"
		value := &repositoryService{
			model: model,
			cb:    gobreaker.NewCircuitBreaker(st),
		}
		repositories.Set(key, value)
	}

	if value, ok := repositories.Get(key); ok {
		return value.(*repositoryService)
	}

	panic("repository not found exception")
}
