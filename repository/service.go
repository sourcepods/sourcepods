package repository

import (
	"errors"
)

var (
	OwnerNotFoundError = errors.New("owner not found")
	AlreadyExistsError = errors.New("repository already exists")
)

type (
	// Store or retrieve repositories from some database.
	Store interface {
		ListByOwnerUsername(string) ([]*Repository, []*Stats, *Owner, error)
		Find(string, string) (*Repository, *Stats, *Owner, error)
		Create(ownerID string, repository *Repository) (*Repository, error)
	}

	// Service to interact with repositories.
	Service interface {
		ListByOwnerUsername(string) ([]*Repository, []*Stats, *Owner, error)
		Find(string, string) (*Repository, *Stats, *Owner, error)
		Create(ownerID string, repository *Repository) (*Repository, error)
	}

	service struct {
		repositories Store
	}
)

// NewService to interact with repositories.
func NewService(repositories Store) Service {
	return &service{
		repositories: repositories,
	}
}

func (s *service) ListByOwnerUsername(username string) ([]*Repository, []*Stats, *Owner, error) {
	return s.repositories.ListByOwnerUsername(username)
}

func (s *service) Find(owner string, name string) (*Repository, *Stats, *Owner, error) {
	return s.repositories.Find(owner, name)
}

func (s *service) Create(ownerID string, repository *Repository) (*Repository, error) {
	if err := ValidateCreate(repository); err != nil {
		return nil, err
	}

	return s.repositories.Create(ownerID, repository)
}
