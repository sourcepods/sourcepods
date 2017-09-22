package authorization

import (
	"context"

	"github.com/gitpods/gitpods/session"
	"github.com/gitpods/gitpods/user"
	"github.com/go-kit/kit/metrics"
)

type metricsService struct {
	loginAttempts metrics.Counter
	service       Service
}

func NewMetricsService(loginAttempts metrics.Counter, service Service) Service {
	// Initialize counters with 0
	loginAttempts.With("status", "failure").Add(0)
	loginAttempts.With("status", "success").Add(0)

	return &metricsService{loginAttempts: loginAttempts, service: service}
}

func (s *metricsService) AuthenticateUser(ctx context.Context, email, password string) (*user.User, error) {
	u, err := s.service.AuthenticateUser(ctx, email, password)

	if err != nil {
		s.loginAttempts.With("status", "failure").Add(1)
	} else {
		s.loginAttempts.With("status", "success").Add(1)
	}

	return u, err
}

func (s *metricsService) CreateSession(ctx context.Context, userID, userUsername string) (*session.Session, error) {
	// Don't do anything here, it's done in the service being called.
	return s.service.CreateSession(ctx, userID, userUsername)
}
