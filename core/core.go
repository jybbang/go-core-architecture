package core

import (
	"reflect"

	cmap "github.com/orcaman/concurrent-map"
	"github.com/sony/gobreaker"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var mediatorInstance *mediator

var eventBusInstance *eventbus

var statesInstance *stateService

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

func SetLogger(logger *zap.Logger) {
	Log = logger.Sugar()
}

func OnCbStateChange(name string, from gobreaker.State, to gobreaker.State) {
	Log.Infow("circuit breaker state changed", "name", name, "from", from.String(), "to", to.String())
}

func GetMediator() *mediator {
	if mediatorInstance == nil {
		panic("you should create mediator before use it")
	}
	return mediatorInstance
}

func Geteventbus() *eventbus {
	if eventBusInstance == nil {
		panic("you should create event bus before use it")
	}
	return eventBusInstance
}

func GetStateService() *stateService {
	if statesInstance == nil {
		panic("you should create state service before use it")
	}
	return statesInstance
}

func GetRepositoryService(model Entitier) *repositoryService {
	valueOf := reflect.ValueOf(model)
	key := valueOf.Type().Name()

	if !repositories.Has(key) {
		panic("you should create repository service before use it")
	}

	if value, ok := repositories.Get(key); ok {
		return value.(*repositoryService)
	}

	panic("repository not found exception")
}
