package session

import (
	"context"

	"github.com/opentracing/opentracing-go"
)

type tracingService struct {
	service Service
}

// NewLoggingService wraps the Service and provides tracing for its methods.
func NewTracingService(s Service) Service {
	return &tracingService{service: s}
}

func (s *tracingService) CreateSession(ctx context.Context, userID, userUsername string) (*Session, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "session.Service.CreateSession")
	span.SetTag("user_id", userID)
	span.SetTag("user_username", userUsername)
	defer span.Finish()

	return s.service.CreateSession(ctx, userID, userUsername)
}

func (s *tracingService) FindSession(ctx context.Context, id string) (*Session, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "session.Service.FindSession")
	span.SetTag("id", id)
	defer span.Finish()

	return s.service.FindSession(ctx, id)
}

func (s *tracingService) ClearSessions(ctx context.Context) (int64, error) {
	return s.service.ClearSessions(ctx)
}
