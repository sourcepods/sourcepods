package repository

type (
	// Store or retrieve repositories from some database.
	Store interface {
		ListByOwnerUsername(string) ([]*Repository, error)
	}

	// Service to interact with repositories.
	Service interface {
		ListByOwnerUsername(string) ([]*Repository, error)
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

func (s *service) ListByOwnerUsername(username string) ([]*Repository, error) {
	repositories, err := s.repositories.ListByOwnerUsername(username)
	if err != nil {
		return nil, err
	}

	return repositories, err
}
