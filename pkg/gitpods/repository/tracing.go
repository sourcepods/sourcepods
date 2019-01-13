package repository

import (
	"context"

	"github.com/gitpods/gitpods/pkg/storage"
	"github.com/opentracing/opentracing-go"
)

//TracingRequestID returns the request ID as string for tracing
type TracingRequestID func(context.Context) string

type tracingService struct {
	service   Service
	requestID TracingRequestID
}

// NewTracingService wraps the Service and provides tracing for its methods.
func NewTracingService(s Service, requestID TracingRequestID) Service {
	return &tracingService{service: s, requestID: requestID}
}

func (s *tracingService) List(ctx context.Context, owner string) ([]*Repository, string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.Service.List")
	span.SetTag("request", s.requestID(ctx))
	span.SetTag("owner", owner)
	defer span.Finish()

	return s.service.List(ctx, owner)
}

func (s *tracingService) Find(ctx context.Context, owner string, name string) (*Repository, string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.Service.Find")
	span.SetTag("request", s.requestID(ctx))
	span.SetTag("owner", owner)
	span.SetTag("name", name)
	defer span.Finish()

	return s.service.Find(ctx, owner, name)
}

func (s *tracingService) Create(ctx context.Context, owner string, repository *Repository) (*Repository, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.Service.Create")
	span.SetTag("request", s.requestID(ctx))
	span.SetTag("owner", owner)
	span.SetTag("name", repository.Name)
	defer span.Finish()

	return s.service.Create(ctx, owner, repository)
}

func (s *tracingService) Branches(ctx context.Context, owner string, name string) ([]*Branch, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.Service.Branches")
	span.SetTag("request", s.requestID(ctx))
	span.SetTag("owner", owner)
	span.SetTag("name", name)
	defer span.Finish()

	return s.service.Branches(ctx, owner, name)
}

func (s *tracingService) Commit(ctx context.Context, owner string, name string, rev string) (storage.Commit, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.Service.Commit")
	span.SetTag("request", s.requestID(ctx))
	span.SetTag("owner", owner)
	span.SetTag("name", name)
	span.SetTag("rev", rev)
	defer span.Finish()

	return s.service.Commit(ctx, owner, name, rev)
}
