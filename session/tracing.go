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

func (s *tracingService) Create(ctx context.Context, userID, userUsername string) (*Session, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "session.Service.Create")
	span.SetTag("user_id", userID)
	span.SetTag("user_username", userUsername)
	defer span.Finish()

	return s.service.Create(ctx, userID, userUsername)
}

func (s *tracingService) Find(ctx context.Context, id string) (*Session, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "session.Service.Find")
	span.SetTag("id", id)
	defer span.Finish()

	return s.service.Find(ctx, id)
}

func (s *tracingService) DeleteExpired(ctx context.Context) (int64, error) {
	return s.service.DeleteExpired(ctx)
}
