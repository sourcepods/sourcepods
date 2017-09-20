package user

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type loggingService struct {
	logger log.Logger
	Service
}

// NewLoggingService wraps the Service and provides logging for its methods.
func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{logger, s}
}

func (s *loggingService) FindAll() ([]*User, error) {
	start := time.Now()

	users, err := s.Service.FindAll()

	logger := log.With(s.logger,
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

func (s *loggingService) FindByUsername(ctx context.Context, username string) (*User, error) {
	start := time.Now()

	user, err := s.Service.FindByUsername(ctx, username)

	logger := log.With(s.logger,
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

func (s *loggingService) Create(user *User) (*User, error) {
	start := time.Now()

	user, err := s.Service.Create(user)

	logger := log.With(s.logger,
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

func (s *loggingService) Update(user *User) (*User, error) {
	start := time.Now()

	user, err := s.Service.Update(user)

	logger := log.With(s.logger,
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

func (s *loggingService) Delete(username string) error {
	start := time.Now()

	err := s.Service.Delete(username)

	logger := log.With(s.logger,
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
