package repository

type (
	// Store or retrieve repositories from some database.
	Store interface {
		ListAggregateByOwnerUsername(string) ([]*RepositoryAggregate, error)
		Find(string, string) (*Repository, *Stats, error)
	}

	// Service to interact with repositories.
	Service interface {
		ListAggregateByOwnerUsername(string) ([]*RepositoryAggregate, error)
		Find(string, string) (*Repository, *Stats, error)
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

func (s *service) Find(owner string, name string) (*Repository, *Stats, error) {
	return s.repositories.Find(owner, name)
}
