package middlewares

import (
	"context"
	"time"

	"github.com/jybbang/go-core-architecture/core"
	"go.uber.org/zap"
)

type performanceMiddleware struct {
	core.Middleware
	log           *zap.Logger
	limitDuration time.Duration
}

func NewPerformanceMiddleware(logger *zap.Logger, limitDuration time.Duration) *performanceMiddleware {
	return &performanceMiddleware{
		log:           logger,
		limitDuration: limitDuration,
	}
}

func (m *performanceMiddleware) Run(ctx context.Context, request core.Request) core.Result {
	defer m.timeMeasurement(time.Now(), request)
	return m.Next()
}

func (m *performanceMiddleware) timeMeasurement(start time.Time, request core.Request) {
	elapsed := time.Since(start)
	if elapsed > time.Duration(m.limitDuration) {
		m.log.Warn("send request long running", zap.Reflect("request", request), zap.Duration("measure", elapsed))
	}
}
