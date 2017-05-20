package user

import "github.com/gitpods/gitpods"

// Service handles all interactions with users.
type Service interface {
	FindAll() ([]*gitpods.User, error)
	FindByUsername(string) (*gitpods.User, error)
	Create(*gitpods.User) (*gitpods.User, error)
	Update(string, *gitpods.User) (*gitpods.User, error)
	Delete(string) error
}

// Store users after manipulation or read them.
type Store interface {
	FindAll() ([]*gitpods.User, error)
	Find(string) (*gitpods.User, error)
	FindByUsername(string) (*gitpods.User, error)
	Create(*gitpods.User) (*gitpods.User, error)
	Update(string, *gitpods.User) (*gitpods.User, error)
	Delete(string) error
}

type service struct {
	users Store
}

// NewService returns a Service that handles all interactions with users.
func NewService(users Store) Service {
	return &service{users: users}
}

func (s *service) FindAll() ([]*gitpods.User, error) {
	return s.users.FindAll()
}

func (s *service) FindByUsername(username string) (*gitpods.User, error) {
	return s.users.FindByUsername(username)
}

func (s *service) Create(user *gitpods.User) (*gitpods.User, error) {
	return user, nil
}

func (s *service) Update(username string, user *gitpods.User) (*gitpods.User, error) {
	return user, nil
}

func (s *service) Delete(username string) error {
	return nil
}
