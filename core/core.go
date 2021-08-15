package core

import (
	"net/http"
	"reflect"
	"time"

	cmap "github.com/orcaman/concurrent-map"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sony/gobreaker"
	"go.uber.org/zap"
)

type MetricsSettings struct {
	Endpoint   string
	ListenAddr string
}

var defaultMetricsSettings = MetricsSettings{
	Endpoint: "/metrics",
}

var mediatorInstance *mediator

var eventBusInstance *eventbus

var statesInstance *stateService

var repositories cmap.ConcurrentMap = cmap.New()

var Log *zap.SugaredLogger

const cbDefaultTimeout = time.Duration(30 * time.Second)

const cbDefaultAllowedRequests = 3

func init() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	Log = logger.Sugar()
}

func onCbStateChange(name string, from gobreaker.State, to gobreaker.State) {
	Log.Infow("circuit breaker state changed", "name", name, "from", from.String(), "to", to.String())
}

func SetLogger(logger *zap.Logger) {
	Log = logger.Sugar()
}

func AddMetrics(settings *MetricsSettings) {
	if settings == nil {
		settings = &defaultMetricsSettings
	}
	http.Handle(settings.Endpoint, promhttp.Handler())
}

func GetMediator() *mediator {
	if mediatorInstance == nil {
		panic("you should create mediator before use it")
	}
	return mediatorInstance
}

func GetEventbus() *eventbus {
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
