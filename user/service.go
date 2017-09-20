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
	FindAll() ([]*User, error)
	Find(string) (*User, error)
	FindByUsername(context.Context, string) (*User, error)
	Create(*User) (*User, error)
	Update(*User) (*User, error)
	Delete(string) error
}

// Store users after manipulation or read them.
type Store interface {
	FindAll() ([]*User, error)
	Find(string) (*User, error)
	FindByUsername(context.Context, string) (*User, error)
	Create(*User) (*User, error)
	Update(*User) (*User, error)
	Delete(string) error
}

type service struct {
	users Store
}

// NewService returns a Service that handles all interactions with users.
func NewService(users Store) Service {
	return &service{users: users}
}

func (s *service) FindAll() ([]*User, error) {
	return s.users.FindAll()
}

func (s *service) Find(id string) (*User, error) {
	return s.users.Find(id)
}

func (s *service) FindByUsername(ctx context.Context, username string) (*User, error) {
	return s.users.FindByUsername(ctx, username)
}

func (s *service) Create(user *User) (*User, error) {
	panic("implement me")
}

func (s *service) Update(user *User) (*User, error) {
	errs := ValidateUpdate(user)
	if len(errs) > 0 {
		return nil, errs[0]
	}

	return s.users.Update(user)
}

func (s *service) Delete(username string) error {
	panic("implement me")
}
