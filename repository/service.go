package repository

import (
	"github.com/gitpods/gitpods"
)

type (
	// Store or retrieve repositories from some database.
	Store interface {
		ListByOwner(string) ([]*gitpods.Repository, error)
	}

	// UserStore is used to find a user by its username.
	UserStore interface {
		FindByUsername(string) (*gitpods.User, error)
	}

	// Service to interact with repositories.
	Service interface {
		ListByOwnerUsername(string) ([]*gitpods.Repository, error)
	}

	service struct {
		repositories Store
		users        UserStore
	}
)

// NewService to interact with repositories.
func NewService(users UserStore, repositories Store) Service {
	return &service{
		users:        users,
		repositories: repositories,
	}
}

func (s *service) ListByOwnerUsername(username string) ([]*gitpods.Repository, error) {
	u, err := s.users.FindByUsername(username)
	if err != nil {
		return nil, err
	}

	repositories, err := s.repositories.ListByOwner(u.ID)
	if err != nil {
		return nil, err
	}

	return repositories, err
}
