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
		List(owner *Owner) ([]*Repository, []*Stats, *Owner, error)
		Find(owner *Owner, name string) (*Repository, *Stats, *Owner, error)
		Create(owner *Owner, repository *Repository) (*Repository, error)
	}

	// Storage manages the git storage
	Storage interface {
		Create(owner string, name string) error
		Description(owner, name, description string) error
	}

	// Service to interact with repositories.
	Service interface {
		List(owner *Owner) ([]*Repository, []*Stats, *Owner, error)
		Find(owner *Owner, name string) (*Repository, *Stats, *Owner, error)
		Create(owner *Owner, repository *Repository) (*Repository, error)
	}

	service struct {
		repositories Store
		storage      Storage
	}
)

// NewService to interact with repositories.
func NewService(repositories Store, storage Storage) Service {
	return &service{
		repositories: repositories,
		storage:      storage,
	}
}

func (s *service) List(owner *Owner) ([]*Repository, []*Stats, *Owner, error) {
	return s.repositories.List(owner)
}

func (s *service) Find(owner *Owner, name string) (*Repository, *Stats, *Owner, error) {
	return s.repositories.Find(owner, name)
}

func (s *service) Create(owner *Owner, repository *Repository) (*Repository, error) {
	if err := ValidateCreate(repository); err != nil {
		return nil, err
	}

	r, err := s.repositories.Create(owner, repository)
	if err != nil {
		return r, err
	}

	if err := s.storage.Create(owner.Username, r.Name); err != nil {
		return r, err
	}

	if err := s.storage.Description(owner.Username, r.Name, r.Description); err != nil {
		return r, err
	}

	return r, nil
}
