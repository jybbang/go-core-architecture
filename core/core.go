package core

import (
	"reflect"
	"sync"

	cmap "github.com/orcaman/concurrent-map"
	"github.com/reactivex/rxgo/v2"
	"github.com/sony/gobreaker"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type singletons struct {
	mediator sync.Once
	eventBus sync.Once
	states   sync.Once
	logger   sync.Mutex
}

var syncs singletons

var mediatorInstance *Mediator

var eventBusInstance *EventBus

var statesInstance *StateService

var repositories cmap.ConcurrentMap = cmap.New()

var Log *zap.SugaredLogger

func init() {
	logger, err := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(zapcore.DebugLevel),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",

			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			TimeKey:    "time",
			EncodeTime: zapcore.ISO8601TimeEncoder,

			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}.Build()
	if err != nil {
		panic(err)
	}
	Log = logger.Sugar()
}

func SetupLogger(logger *zap.Logger) {
	syncs.logger.Lock()
	defer syncs.logger.Unlock()
	Log = logger.Sugar()
}

func OnCbStateChange(name string, from gobreaker.State, to gobreaker.State) {
	Log.Infow("circuit breaker state changed", "name", name, "from", from.String(), "to", to.String())
}

func GetMediator() *Mediator {
	if mediatorInstance == nil {
		syncs.mediator.Do(
			func() {
				mediatorInstance = &Mediator{
					requestHandlers:      cmap.New(),
					notificationHandlers: cmap.New(),
				}
				mediatorInstance.Setup()

				Log.Infow("mediator created")
			})
	}
	return mediatorInstance
}

func GetEventBus() *EventBus {
	if eventBusInstance == nil {
		syncs.eventBus.Do(
			func() {
				st := gobreaker.Settings{
					Name:          "eventbus",
					MaxRequests:   3,
					OnStateChange: OnCbStateChange,
				}

				eventBusInstance = &EventBus{
					mediator:     GetMediator(),
					domainEvents: make([]DomainEventer, 0),
					ch:           make(chan rxgo.Item, 1),
					cb:           gobreaker.NewCircuitBreaker(st),
				}
				eventBusInstance.Initialize()

				Log.Infow("eventbus created")
			})
	}
	return eventBusInstance
}

func GetStateService() *StateService {
	if statesInstance == nil {
		syncs.states.Do(
			func() {
				st := gobreaker.Settings{
					Name:          "state-service",
					MaxRequests:   3,
					OnStateChange: OnCbStateChange,
				}

				statesInstance = &StateService{
					cb: gobreaker.NewCircuitBreaker(st),
				}
				statesInstance.Initialize()

				Log.Infow("state created")
			})
	}
	return statesInstance
}

func GetRepositoryService(model Entitier) *RepositoryService {
	valueOf := reflect.ValueOf(model)
	key := valueOf.Type().Name()

	if !repositories.Has(key) {
		st := gobreaker.Settings{
			Name:          key + "-repository",
			MaxRequests:   3,
			OnStateChange: OnCbStateChange,
		}

		repository := &RepositoryService{
			model: model,
			cb:    gobreaker.NewCircuitBreaker(st),
		}
		repository.Initialize()

		repositories.Set(key, repository)
		Log.Infow("repository created")
	}

	if value, ok := repositories.Get(key); ok {
		return value.(*RepositoryService)
	}

	panic("repository not found exception")
}
