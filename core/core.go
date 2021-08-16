package core

import (
	"net/http"
	"reflect"
	"time"

	"github.com/opentracing/opentracing-go"
	cmap "github.com/orcaman/concurrent-map"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sony/gobreaker"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics/prometheus"
)

type MetricsSettings struct {
	Endpoint string
}

type TracingSettings struct {
	ServiceName string
}

type CircuitBreakerSettings struct {
	Name                     string
	AllowedRequestInHalfOpen int
	SamplingFailureCount     int
	SamplingFailureRatio     float64
	SamplingDuration         time.Duration
	DurationOfBreak          time.Duration
	OnStateChange            func(name string, from string, to string)
}

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

func (s *CircuitBreakerSettings) toGobreakerSettings() gobreaker.Settings {
	if s.AllowedRequestInHalfOpen < 1 {
		s.AllowedRequestInHalfOpen = 1
	}
	if s.SamplingDuration <= time.Duration(0) {
		s.SamplingDuration = time.Duration(60 * time.Second)
	}
	if s.DurationOfBreak <= time.Duration(0) {
		s.DurationOfBreak = time.Duration(60 * time.Second)
	}
	if s.SamplingFailureCount < 1 {
		s.SamplingFailureCount = 5
	}

	return gobreaker.Settings{
		Name:        s.Name,
		MaxRequests: uint32(s.AllowedRequestInHalfOpen),
		Interval:    s.SamplingDuration,
		Timeout:     s.DurationOfBreak,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			if s.SamplingFailureRatio > 0 {
				failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
				return counts.Requests >= uint32(s.SamplingFailureCount) && failureRatio >= s.SamplingFailureRatio
			}
			return counts.TotalFailures >= uint32(s.SamplingFailureCount)
		},
		OnStateChange: func(name string, from, to gobreaker.State) {
			if s.OnStateChange != nil {
				s.OnStateChange(name, from.String(), to.String())
			}
		},
	}
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
	typeOf := reflect.TypeOf(model)
	key := typeOf.Elem().Name()

	if key == "" {
		panic("key is required")
	}

	if !repositories.Has(key) {
		panic("you should create repository service before use it")
	}

	if value, ok := repositories.Get(key); ok {
		return value.(*repositoryService)
	}

	panic("repository not found exception")
}
