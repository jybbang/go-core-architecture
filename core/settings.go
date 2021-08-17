package core

import (
	"strings"
	"time"

	"github.com/sony/gobreaker"
)

type MetricsSettings struct {
	Endpoint string
}

type TracingSettings struct {
	ServiceName string
}

type EventbusSettings struct {
	BufferedEventBufferTime  time.Duration
	BufferedEventBufferCount int
	BufferedEventTimeout     time.Duration
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

func (s *CircuitBreakerSettings) toGobreakerSettings(defaultName string) gobreaker.Settings {
	if strings.TrimSpace(s.Name) == "" {
		s.Name = defaultName
	}
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
