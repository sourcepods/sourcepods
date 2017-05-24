package repository

import (
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type loggingService struct {
	logger  log.Logger
	service Service
}

// NewLoggingService wraps the Service and provides logging for its methods.
func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{logger: logger, service: s}
}

func (s *loggingService) ListAggregateByOwnerUsername(username string) ([]*Repository, []*Stats, error) {
	start := time.Now()

	repositories, stats, err := s.service.ListAggregateByOwnerUsername(username)

	logger := log.With(s.logger,
		"method", "ListAggregateByOwnerUsername",
		"duration", time.Since(start),
	)

	if err != nil {
		level.Warn(logger).Log(
			"msg", "failed to list repositories by owner's username",
			"err", err,
		)
	} else {
		level.Debug(logger).Log()
	}

	return repositories, stats, err

}

func (s *loggingService) Find(owner string, name string) (*Repository, *Stats, error) {
	start := time.Now()

	repository, stats, err := s.service.Find(owner, name)

	logger := log.With(s.logger,
		"method", "Find",
		"duration", time.Since(start),
	)

	if err != nil {
		level.Warn(logger).Log(
			"msg", "failed to find repository by owner & name",
			"err", err,
		)
	} else {
		level.Debug(logger).Log()
	}

	return repository, stats, err
}
