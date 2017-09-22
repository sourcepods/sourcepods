package user

import (
	"context"
	"errors"
)

var (
	NotFoundError = errors.New("user not found")
)

// Service handles all interactions with users.
type Service interface {
	FindAll(context.Context) ([]*User, error)
	Find(context.Context, string) (*User, error)
	FindByUsername(context.Context, string) (*User, error)
	Create(context.Context, *User) (*User, error)
	Update(context.Context, *User) (*User, error)
	Delete(context.Context, string) error
}

// Store users after manipulation or read them.
type Store interface {
	FindAll(context.Context) ([]*User, error)
	Find(context.Context, string) (*User, error)
	FindByUsername(context.Context, string) (*User, error)
	Create(context.Context, *User) (*User, error)
	Update(context.Context, *User) (*User, error)
	Delete(context.Context, string) error
}

type service struct {
	users Store
}

// NewService returns a Service that handles all interactions with users.
func NewService(users Store) Service {
	return &service{users: users}
}

func (s *service) FindAll(ctx context.Context) ([]*User, error) {
	return s.users.FindAll(ctx)
}

func (s *service) Find(ctx context.Context, id string) (*User, error) {
	return s.users.Find(ctx, id)
}

func (s *service) FindByUsername(ctx context.Context, username string) (*User, error) {
	return s.users.FindByUsername(ctx, username)
}

func (s *service) Create(ctx context.Context, user *User) (*User, error) {
	panic("implement me")
}

func (s *service) Update(ctx context.Context, user *User) (*User, error) {
	errs := ValidateUpdate(user)
	if len(errs) > 0 {
		return nil, errs[0]
	}

	return s.users.Update(ctx, user)
}

func (s *service) Delete(ctx context.Context, username string) error {
	panic("implement me")
}
