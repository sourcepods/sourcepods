package repository

type (
	// Store or retrieve repositories from some database.
	Store interface {
		ListAggregateByOwnerUsername(string) ([]*RepositoryAggregate, error)
	}

	// Service to interact with repositories.
	Service interface {
		ListAggregateByOwnerUsername(string) ([]*RepositoryAggregate, error)
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

func (s *service) ListAggregateByOwnerUsername(username string) ([]*RepositoryAggregate, error) {
	return s.repositories.ListAggregateByOwnerUsername(username)
}
