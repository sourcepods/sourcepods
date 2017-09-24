package repository

import (
	"context"
	"errors"

	"github.com/gitpods/gitpods/storage"
)

var (
	OwnerNotFoundError = errors.New("owner not found")
	AlreadyExistsError = errors.New("repository already exists")
)

type (
	// Store or retrieve repositories from some database.
	Store interface {
		List(ctx context.Context, owner *Owner) ([]*Repository, []*Stats, *Owner, error)
		Find(ctx context.Context, owner *Owner, name string) (*Repository, *Stats, *Owner, error)
		Create(ctx context.Context, owner *Owner, repository *Repository) (*Repository, error)
	}

	// Storage manages the git storage
	Storage interface {
		Create(ctx context.Context, owner string, name string) error
		SetDescription(ctx context.Context, owner, name, description string) error
		Tree(ctx context.Context, owner, name, branch string) ([]storage.TreeObject, error)
	}

	// Service to interact with repositories.
	Service interface {
		List(ctx context.Context, owner *Owner) ([]*Repository, []*Stats, *Owner, error)
		Find(ctx context.Context, owner *Owner, name string) (*Repository, *Stats, *Owner, error)
		Create(ctx context.Context, owner *Owner, repository *Repository) (*Repository, error)
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

func (s *service) List(ctx context.Context, owner *Owner) ([]*Repository, []*Stats, *Owner, error) {
	return s.repositories.List(ctx, owner)
}

func (s *service) Find(ctx context.Context, owner *Owner, name string) (*Repository, *Stats, *Owner, error) {
	_, err := s.storage.Tree(ctx, owner.Username, name, "master")
	if err != nil {
		return nil, nil, nil, err
	}

	return s.repositories.Find(ctx, owner, name)
}

func (s *service) Create(ctx context.Context, owner *Owner, repository *Repository) (*Repository, error) {
	if err := ValidateCreate(repository); err != nil {
		return nil, err
	}

	r, err := s.repositories.Create(ctx, owner, repository)
	if err != nil {
		return r, err
	}

	if err := s.storage.Create(ctx, owner.Username, r.Name); err != nil {
		return r, err
	}

	if err := s.storage.SetDescription(ctx, owner.Username, r.Name, r.Description); err != nil {
		return r, err
	}

	return r, nil
}
