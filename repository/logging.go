package repository

import (
	"context"
	"time"

	"github.com/gitpods/gitpods/storage"
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

func (s *loggingService) List(ctx context.Context, owner *Owner) ([]*Repository, []*Stats, *Owner, error) {
	start := time.Now()

	repositories, stats, owner, err := s.service.List(ctx, owner)

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

func (s *loggingService) Find(ctx context.Context, owner *Owner, name string) (*Repository, *Stats, *Owner, error) {
	start := time.Now()

	repository, stats, owner, err := s.service.Find(ctx, owner, name)

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

func (s *loggingService) Create(ctx context.Context, owner *Owner, repository *Repository) (*Repository, error) {
	start := time.Now()

	repository, err := s.service.Create(ctx, owner, repository)

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

func (s *loggingService) Tree(ctx context.Context, owner *Owner, name string) ([]storage.TreeObject, error) {
	start := time.Now()

	objects, err := s.service.Tree(ctx, owner, name)

	logger := log.With(s.logger,
		"method", "Tree",
		"owner", owner,
		"name", name,
		"duration", time.Since(start),
	)

	if err != nil {
		level.Warn(logger).Log(
			"msg", "failed to get tree for repository",
			"err", err,
		)
	} else {
		level.Debug(logger).Log()
	}

	return objects, err
}
