package authorization

import (
	"context"

	"github.com/gitpods/gitpods/session"
	"github.com/gitpods/gitpods/user"
	"github.com/opentracing/opentracing-go"
)

type tracingService struct {
	service Service
}

// NewTracingService wraps the Service and provides tracing for its methods.
func NewTracingService(s Service) Service {
	return &tracingService{service: s}
}

func (s *tracingService) AuthenticateUser(ctx context.Context, email, password string) (*user.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authorization.Service.AuthenticateUser")
	span.SetTag("email", email)
	defer span.Finish()

	return s.service.AuthenticateUser(ctx, email, password)
}

func (s *tracingService) CreateSession(ctx context.Context, userID string, userUsername string) (*session.Session, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authorization.Service.CreateSession")
	span.SetTag("user_id", userID)
	span.SetTag("user_username", userUsername)
	defer span.Finish()

	return s.service.CreateSession(ctx, userID, userUsername)
}
