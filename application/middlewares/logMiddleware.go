package middlewares

import (
	"github.com/jybbang/go-core-architecture/application"
	"go.uber.org/zap"
)

type logMiddleware struct {
	application.Middleware
	log *zap.SugaredLogger
}

func NewLogMiddleware() *logMiddleware {
	logger, _ := zap.NewProduction()

	return &logMiddleware{
		log: logger.Sugar(),
	}
}

func (m *logMiddleware) Run(request application.Request) (bool, error) {
	m.log.Info("mediator log")
	return true, nil
}
