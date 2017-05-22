package authorization

import (
	"time"

	"github.com/gitpods/gitpods/session"
	"github.com/gitpods/gitpods/user"
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

func (s *loggingService) AuthenticateUser(email, password string) (*user.User, error) {
	start := time.Now()

	user, err := s.service.AuthenticateUser(email, password)

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

func (s *loggingService) CreateSession(userID, userUsername string) (*session.Session, error) {
	// Don't log anything here, it's done in the service being called.
	return s.service.CreateSession(userID, userUsername)
}
