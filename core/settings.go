package core

import (
	"strings"
	"time"

	"github.com/sony/gobreaker"
)

type TracingSettings struct {
	ServiceName string
	Endpoint    string
}

type EventBusSettings struct {
	BufferedEventBufferCount int           `model:",omitempty"`
	BufferedEventBufferTime  time.Duration `model:",omitempty"`
	BufferedEventTimeout     time.Duration `model:",omitempty"`
	ConnectionTimeout        time.Duration `model:",omitempty"`
}

type StateServiceSettings struct {
	ConnectionTimeout time.Duration `model:",omitempty"`
}

type RepositoryServiceSettings struct {
	ConnectionTimeout time.Duration `model:",omitempty"`
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
	if strings.TrimSpace(s.Name) == "" {
		s.Name = defaultName
	}

	if strings.TrimSpace(s.Name) == "" {
		panic("name is required")
	}

	settings := gobreaker.Settings{
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

			onCircuitOpen()
		},
	}

	return gobreaker.NewCircuitBreaker(settings)
}
