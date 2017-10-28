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
		List(ctx context.Context, owner string) ([]*Repository, []*Stats, string, error)
		Find(ctx context.Context, owner string, name string) (*Repository, *Stats, string, error)
		Create(ctx context.Context, owner string, repository *Repository) (*Repository, error)
	}

	// Storage manages the git storage
	Storage interface {
		Create(ctx context.Context, owner string, name string) error
		SetDescription(ctx context.Context, owner, name, description string) error
		Branches(ctx context.Context, owner string, name string) ([]storage.Branch, error)
		Commit(ctx context.Context, owner string, name string, rev string) (storage.Commit, error)
	}

	// Service to interact with repositories.
	Service interface {
		List(ctx context.Context, owner string) ([]*Repository, []*Stats, string, error)
		Find(ctx context.Context, owner string, name string) (*Repository, *Stats, string, error)
		Create(ctx context.Context, owner string, repository *Repository) (*Repository, error)
		Branches(ctx context.Context, owner string, name string) ([]*Branch, error)
		Commit(ctx context.Context, owner string, name string, rev string) (storage.Commit, error)
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

func (s *service) List(ctx context.Context, owner string) ([]*Repository, []*Stats, string, error) {
	return s.repositories.List(ctx, owner)
}

func (s *service) Find(ctx context.Context, owner string, name string) (*Repository, *Stats, string, error) {
	return s.repositories.Find(ctx, owner, name)
}

func (s *service) Create(ctx context.Context, owner string, repository *Repository) (*Repository, error) {
	if err := ValidateCreate(repository); err != nil {
		return nil, err
	}

	r, err := s.repositories.Create(ctx, owner, repository)
	if err != nil {
		return r, err
	}

	if err := s.storage.Create(ctx, owner, r.Name); err != nil {
		return r, err
	}

	if err := s.storage.SetDescription(ctx, owner, r.Name, r.Description); err != nil {
		return r, err
	}

	return r, nil
}

func (s *service) Branches(ctx context.Context, owner string, name string) ([]*Branch, error) {
	bs, err := s.storage.Branches(ctx, owner, name)
	if err != nil {
		return nil, err
	}

	var branches []*Branch
	for _, b := range bs {
		branches = append(branches, &Branch{
			Name: b.Name,
			Sha1: b.Sha1,
			Type: b.Type,
		})
	}

	return branches, nil
}

func (s *service) Commit(ctx context.Context, owner string, name string, rev string) (storage.Commit, error) {
	return s.storage.Commit(ctx, owner, name, rev)
}
