package core

import (
	"log"
	"reflect"

	"github.com/opentracing/opentracing-go"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/reporter"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"

	cmap "github.com/orcaman/concurrent-map"
)

var mediatorInstance *mediator

var eventBusInstance *eventbus

var statesInstance *stateService

var repositories cmap.ConcurrentMap = cmap.New()

var openTracerCloser reporter.Reporter

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

// reference:
// https://github.com/openzipkin/zipkin-go/blob/master/examples/httpserver_test.go
func UseTracing(settings TracingSettings) {
	// set up a span reporter
	reporter := zipkinhttp.NewReporter(settings.Endpoint)
	openTracerCloser = reporter

	// create our local service endpoint
	endpoint, err := zipkin.NewEndpoint(settings.ServiceName, "localhost")
	if err != nil {
		log.Fatalf("unable to create local endpoint: %+v\n", err)
	}

	// initialize our tracer
	nativeTracer, err := zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(endpoint))
	if err != nil {
		log.Fatalf("unable to create tracer: %+v\n", err)
	}

	// use zipkin-go-opentracing to wrap our tracer
	tracer := zipkinot.Wrap(nativeTracer)

	// optionally set as Global OpenTracing tracer instance
	opentracing.SetGlobalTracer(tracer)
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
