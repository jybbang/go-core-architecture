package core

import (
	"io"
	"reflect"

	"github.com/opentracing/opentracing-go"
	cmap "github.com/orcaman/concurrent-map"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics/prometheus"
)

var mediatorInstance *mediator

var eventBusInstance *eventbus

var statesInstance *stateService

var repositories cmap.ConcurrentMap = cmap.New()

var openTracer opentracing.Tracer
var openTracerCloser io.Closer

func Close() {
	if openTracerCloser != nil {
		openTracerCloser.Close()
	}
	GetEventbus().close()
	GetStateService().close()
	for _, r := range repositories.Items() {
		r.(*repositoryService).close()
	}
}

func UseTracing(settings TracingSettings) {
	metricsFactory := prometheus.New()
	tracer, closer, err := config.Configuration{
		ServiceName: settings.ServiceName,
		Reporter:    &config.ReporterConfig{CollectorEndpoint: settings.Endpoint},
	}.NewTracer(
		config.Metrics(metricsFactory),
	)
	if err != nil {
		panic(err)
	}
	openTracer = tracer
	openTracerCloser = closer
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
