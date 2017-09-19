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

func (s *loggingService) List(owner *Owner) ([]*Repository, []*Stats, *Owner, error) {
	start := time.Now()

	repositories, stats, owner, err := s.service.List(owner)

	logger := log.With(s.logger,
		"method", "List",
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

func (s *loggingService) Find(owner *Owner, name string) (*Repository, *Stats, *Owner, error) {
	start := time.Now()

	repository, stats, owner, err := s.service.Find(owner, name)

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

func (s *loggingService) Create(owner *Owner, repository *Repository) (*Repository, error) {
	start := time.Now()

	repository, err := s.service.Create(owner, repository)

	logger := log.With(s.logger,
		"method", "Create",
		"duration", time.Since(start),
	)

	if err != nil {
		if err != AlreadyExistsError {
			level.Warn(logger).Log(
				"msg", "failed to create repository",
				"err", err,
			)
		}
	} else {
		level.Debug(logger).Log(
			"owner", owner.Username,
			"name", repository.Name,
		)
	}

	return repository, err
}
