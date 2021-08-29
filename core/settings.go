package core

import (
	"fmt"
	"strings"
	"time"

	"github.com/sony/gobreaker"
	"gopkg.in/jeevatkm/go-model.v1"
)

type TracingSettings struct {
	ServiceName string
	Endpoint    string
}

type EventBusSettings struct {
	BufferedEventBufferCount int           `model:",omitempty"`
	BufferedEventBufferTime  time.Duration `model:",omitempty"`
	BufferedEventTimeout     time.Duration `model:",omitempty"`
}

type CircuitBreakerSettings struct {
	Name                     string
	AllowedRequestInHalfOpen int           `model:",omitempty"`
	DurationOfBreak          time.Duration `model:",omitempty"`
	SamplingDuration         time.Duration `model:",omitempty"`
	SamplingFailureCount     int           `model:",omitempty"`
	SamplingFailureRatio     float64
	OnStateChange            func(name string, from string, to string)
}

func (s *CircuitBreakerSettings) ToCircuitBreaker(defaultName string, onCircuitOpen func()) *gobreaker.CircuitBreaker {
	ss := &CircuitBreakerSettings{
		AllowedRequestInHalfOpen: 1,
		DurationOfBreak:          time.Duration(60 * time.Second),
		SamplingDuration:         time.Duration(60 * time.Second),
		SamplingFailureCount:     5,
	}

	err := model.Copy(ss, s)
	if err != nil {
		panic(fmt.Errorf("mapping errors occurred: %v", err))
	}

	if strings.TrimSpace(s.Name) == "" {
		ss.Name = defaultName
	}
	if strings.TrimSpace(ss.Name) == "" {
		panic("name is required")
	}

	settings := gobreaker.Settings{
		Name:        ss.Name,
		MaxRequests: uint32(ss.AllowedRequestInHalfOpen),
		Interval:    ss.SamplingDuration,
		Timeout:     ss.DurationOfBreak,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			if ss.SamplingFailureRatio > 0 {
				failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
				return counts.Requests >= uint32(ss.SamplingFailureCount) && failureRatio >= ss.SamplingFailureRatio
			}
			return counts.TotalFailures >= uint32(ss.SamplingFailureCount)
		},
		OnStateChange: func(name string, from, to gobreaker.State) {
			if ss.OnStateChange != nil {
				ss.OnStateChange(name, from.String(), to.String())
			}
			onCircuitOpen()
		},
	}

	return gobreaker.NewCircuitBreaker(settings)
}
