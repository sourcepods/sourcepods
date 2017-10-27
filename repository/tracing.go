package repository

import (
	"context"

	"github.com/gitpods/gitpods/storage"
	"github.com/opentracing/opentracing-go"
)

type tracingService struct {
	service Service
}

// NewLoggingService wraps the Service and provides tracing for its methods.
func NewTracingService(s Service) Service {
	return &tracingService{service: s}
}

func (s *tracingService) List(ctx context.Context, owner string) ([]*Repository, []*Stats, string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.Service.List")
	span.SetTag("owner", owner)
	defer span.Finish()

	return s.service.List(ctx, owner)
}

func (s *tracingService) Find(ctx context.Context, owner string, name string) (*Repository, *Stats, string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.Service.Find")
	span.SetTag("owner", owner)
	span.SetTag("name", name)
	defer span.Finish()

	return s.service.Find(ctx, owner, name)
}

func (s *tracingService) Create(ctx context.Context, owner string, repository *Repository) (*Repository, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.Service.Create")
	span.SetTag("owner", owner)
	span.SetTag("name", repository.Name)
	defer span.Finish()

	return s.service.Create(ctx, owner, repository)
}

func (s *tracingService) Branches(ctx context.Context, owner string, name string) ([]*Branch, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.Service.Branches")
	span.SetTag("owner", owner)
	span.SetTag("name", name)
	defer span.Finish()

	return s.service.Branches(ctx, owner, name)
}

func (s *tracingService) Tree(ctx context.Context, owner string, name string, recursive bool) ([]storage.TreeObject, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.Service.Tree")
	span.SetTag("owner", owner)
	span.SetTag("name", name)
	defer span.Finish()

	return s.service.Tree(ctx, owner, name, recursive)
}
