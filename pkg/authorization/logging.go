package authorization

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/sourcepods/sourcepods/pkg/gitpods/user"
	"github.com/sourcepods/sourcepods/pkg/session"
)

type loggingService struct {
	logger  log.Logger
	service Service
}

// NewLoggingService wraps the Service and provides logging for its methods.
func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{logger: logger, service: s}
}

func (s *loggingService) AuthenticateUser(ctx context.Context, email, password string) (*user.User, error) {
	start := time.Now()

	user, err := s.service.AuthenticateUser(ctx, email, password)

	logger := log.With(s.logger,
		"method", "AuthenticateUser",
		"duration", time.Since(start),
	)

	if err != nil {
		level.Warn(logger).Log("msg", "failed to authenticate user", "err", err)
	} else {
		level.Debug(logger).Log()
	}

	return user, err
}

func (s *loggingService) CreateSession(ctx context.Context, userID, userUsername string) (*session.Session, error) {
	// Don't log anything here, it's done in the service being called.
	return s.service.CreateSession(ctx, userID, userUsername)
}
