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

func (s *loggingService) ListByOwnerUsername(username string) ([]*Repository, []*Stats, *Owner, error) {
	start := time.Now()

	repositories, stats, owner, err := s.service.ListByOwnerUsername(username)

	logger := log.With(s.logger,
		"method", "ListByOwnerUsername",
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

	return repositories, stats, owner, err

}

func (s *loggingService) Find(ownerUsername string, name string) (*Repository, *Stats, *Owner, error) {
	start := time.Now()

	repository, stats, owner, err := s.service.Find(ownerUsername, name)

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

	return repository, stats, owner, err
}
