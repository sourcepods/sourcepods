package repository

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

func (s *tracingService) List(ctx context.Context, owner *Owner) ([]*Repository, []*Stats, *Owner, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.Service.List")
	span.SetTag("owner_id", owner.ID)
	span.SetTag("owner_username", owner.Username)
	defer span.Finish()

	return s.service.List(ctx, owner)
}

func (s *tracingService) Find(ctx context.Context, owner *Owner, name string) (*Repository, *Stats, *Owner, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.Service.Find")
	span.SetTag("owner_id", owner.ID)
	span.SetTag("owner_username", owner.Username)
	span.SetTag("name", name)
	defer span.Finish()

	return s.service.Find(ctx, owner, name)
}

func (s *tracingService) Create(ctx context.Context, owner *Owner, repository *Repository) (*Repository, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.Service.Create")
	span.SetTag("owner_id", owner.ID)
	span.SetTag("owner_username", owner.Username)
	span.SetTag("name", repository.Name)
	defer span.Finish()

	return s.service.Create(ctx, owner, repository)
}
