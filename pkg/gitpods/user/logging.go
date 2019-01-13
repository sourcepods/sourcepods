package user

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

//LoggingRequestID returns the request ID as string for logging
type LoggingRequestID func(context.Context) string

type loggingService struct {
	service   Service
	requestID LoggingRequestID
	logger    log.Logger
}

// NewLoggingService wraps the Service and provides logging for its methods.
func NewLoggingService(s Service, requestID LoggingRequestID, logger log.Logger) Service {
	return &loggingService{service: s, requestID: requestID, logger: logger}
}

func (s *loggingService) FindAll(ctx context.Context) ([]*User, error) {
	start := time.Now()

	users, err := s.service.FindAll(ctx)

	logger := log.With(s.logger,
		"request", s.requestID(ctx),
		"method", "FindAll",
		"duration", time.Since(start),
	)

	if err != nil {
		level.Warn(logger).Log("msg", "failed to find all users", "err", err)
	} else {
		level.Debug(logger).Log()
	}

	return users, err
}

func (s *loggingService) Find(context.Context, string) (*User, error) {
	panic("implement me")
}

func (s *loggingService) FindByUsername(ctx context.Context, username string) (*User, error) {
	start := time.Now()

	user, err := s.service.FindByUsername(ctx, username)

	logger := log.With(s.logger,
		"request", s.requestID(ctx),
		"method", "FindByUsername",
		"duration", time.Since(start),
		"username", username,
	)

	if err != nil && err != NotFoundError {
		level.Warn(logger).Log("msg", "failed to find user by username", "err", err)
	} else {
		level.Debug(logger).Log()
	}

	return user, err
}

func (s *loggingService) FindRepositoryOwner(ctx context.Context, repositoryID string) (*User, error) {
	panic("implement me")
}

func (s *loggingService) Create(ctx context.Context, user *User) (*User, error) {
	start := time.Now()

	user, err := s.service.Create(ctx, user)

	logger := log.With(s.logger,
		"request", s.requestID(ctx),
		"method", "Create",
		"duration", time.Since(start),
		"username", user.Username,
	)

	if err != nil && err != NotFoundError {
		level.Warn(logger).Log("msg", "failed to create user", "err", err)
	} else {
		level.Debug(logger).Log()
	}

	return user, err
}

func (s *loggingService) Update(ctx context.Context, user *User) (*User, error) {
	start := time.Now()

	user, err := s.service.Update(ctx, user)

	logger := log.With(s.logger,
		"request", s.requestID(ctx),
		"method", "Update",
		"duration", time.Since(start),
		"username", user.Username,
	)

	if err != nil {
		level.Warn(logger).Log("msg", "failed to update user", "err", err)
	} else {
		level.Debug(logger).Log()
	}

	return user, err
}

func (s *loggingService) Delete(ctx context.Context, username string) error {
	start := time.Now()

	err := s.service.Delete(ctx, username)

	logger := log.With(s.logger,
		"request", s.requestID(ctx),
		"method", "Delete",
		"duration", time.Since(start),
		"username", username,
	)

	if err != nil {
		level.Warn(logger).Log("msg", "failed to update user", "err", err)
	} else {
		level.Debug(logger).Log()
	}

	return err
}
