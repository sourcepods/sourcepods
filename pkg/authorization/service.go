package authorization

import (
	"context"

	"github.com/gitpods/gitpods/pkg/gitpods/user"
	"github.com/gitpods/gitpods/pkg/session"
	"github.com/opentracing/opentracing-go"
	"golang.org/x/crypto/bcrypt"
)

// Service authenticates users and creates sessions for them.
type Service interface {
	AuthenticateUser(ctx context.Context, email, password string) (*user.User, error)
	CreateSession(context.Context, string, string) (*session.Session, error)
}

// Store finds users by emails.
type Store interface {
	FindUserByEmail(context.Context, string) (*user.User, error)
}

// NewService takes a store to find users by their email and
// takes a session service to create sessions for them once authenticated.
func NewService(store Store, sessions session.Service) Service {
	return &service{store: store, sessions: sessions}
}

type service struct {
	store    Store
	sessions session.Service
}

// AuthenticateUser by hashing the given password an comparing it against the one stored.
func (s *service) AuthenticateUser(ctx context.Context, email, password string) (*user.User, error) {
	u, err := s.store.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	span, ctx := opentracing.StartSpanFromContext(ctx, "authorization.Service.CompareHashAndPassword")
	defer span.Finish()

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return nil, err
	}

	return u, nil
}

func (s *service) CreateSession(ctx context.Context, userID, userUsername string) (*session.Session, error) {
	return s.sessions.Create(ctx, userID, userUsername)
}
