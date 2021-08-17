package core

import (
	"net/http"
	"reflect"

	"github.com/opentracing/opentracing-go"
	cmap "github.com/orcaman/concurrent-map"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics/prometheus"
)

var mediatorInstance *mediator

var eventBusInstance *eventbus

var statesInstance *stateService

var repositories cmap.ConcurrentMap = cmap.New()

var openTracer opentracing.Tracer

func AddMetrics(settings MetricsSettings) {
	http.Handle(settings.Endpoint, promhttp.Handler())
}

func AddTracing(settings TracingSettings) {
	metricsFactory := prometheus.New()
	tracer, _, err := config.Configuration{
		ServiceName: settings.ServiceName,
	}.NewTracer(
		config.Metrics(metricsFactory),
	)
	if err != nil {
		panic(err)
	}
	openTracer = tracer
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
	if model == nil {
		panic("model is required")
	}

	typeOf := reflect.TypeOf(model)
	key := typeOf.Elem().Name()

	if !repositories.Has(key) {
		panic("you should create repository service before use it")
	}

	if value, ok := repositories.Get(key); ok {
		return value.(*repositoryService)
	}

	panic("repository not found exception")
}
