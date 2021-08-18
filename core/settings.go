package core

import (
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

func (s *CircuitBreakerSettings) ToGobreakerSettings(defaultName string) gobreaker.Settings {
	settings := &CircuitBreakerSettings{
		AllowedRequestInHalfOpen: 1,
		DurationOfBreak:          time.Duration(60 * time.Second),
		SamplingDuration:         time.Duration(60 * time.Second),
		SamplingFailureCount:     5,
	}

	err := model.Copy(settings, s)
	if err != nil {
		panic(err)
	}

	if strings.TrimSpace(s.Name) == "" {
		settings.Name = defaultName
	}
	if strings.TrimSpace(settings.Name) == "" {
		panic("name is required")
	}

	return gobreaker.Settings{
		Name:        settings.Name,
		MaxRequests: uint32(settings.AllowedRequestInHalfOpen),
		Interval:    settings.SamplingDuration,
		Timeout:     settings.DurationOfBreak,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			if settings.SamplingFailureRatio > 0 {
				failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
				return counts.Requests >= uint32(settings.SamplingFailureCount) && failureRatio >= settings.SamplingFailureRatio
			}
			return counts.TotalFailures >= uint32(settings.SamplingFailureCount)
		},
		OnStateChange: func(name string, from, to gobreaker.State) {
			if settings.OnStateChange != nil {
				settings.OnStateChange(name, from.String(), to.String())
			}
		},
	}
}
