package user

import (
	"context"

	"github.com/opentracing/opentracing-go"
)

type tracingService struct {
	service Service
}

// NewTracingService wraps the Service and provides tracing for its methods.
func NewTracingService(s Service) Service {
	return &tracingService{s}
}

func (s *tracingService) FindAll(ctx context.Context) ([]*User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "user.Service.FindAll")
	defer span.Finish()

	return s.service.FindAll(ctx)
}

func (s *tracingService) Find(ctx context.Context, id string) (*User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "user.Service.Find")
	span.SetTag("id", id)
	defer span.Finish()

	return s.service.Find(ctx, id)
}

func (s *tracingService) FindByUsername(ctx context.Context, username string) (*User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "user.Service.FindByUsername")
	span.SetTag("username", username)
	defer span.Finish()

	return s.service.FindByUsername(ctx, username)
}

func (s *tracingService) FindRepositoryOwner(ctx context.Context, repositoryID string) (*User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "user.Service.FindRepositoryOwner")
	span.SetTag("repository", repositoryID)
	defer span.Finish()

	return s.service.FindRepositoryOwner(ctx, repositoryID)
}

func (s *tracingService) Create(ctx context.Context, user *User) (*User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "user.Service.Create")
	span.SetTag("username", user.Username)
	defer span.Finish()

	return s.service.Create(ctx, user)
}

func (s *tracingService) Update(ctx context.Context, user *User) (*User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "user.Service.Update")
	span.SetTag("id", user.ID)
	span.SetTag("username", user.Username)
	defer span.Finish()

	return s.service.Update(ctx, user)
}

func (s *tracingService) Delete(ctx context.Context, username string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "user.Service.Delete")
	span.SetTag("username", username)
	defer span.Finish()

	return s.service.Delete(ctx, username)
}
