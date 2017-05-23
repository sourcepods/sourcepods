package repository

import (
	"database/sql"
	"time"
)

// Postgres implementation of the Store.
type Postgres struct {
	db *sql.DB
}

// NewPostgresStore returns a Postgres implementation of the Store.
func NewPostgresStore(db *sql.DB) *Postgres {
	return &Postgres{db: db}
}

// ListAggregateByOwnerUsername retrieves a list of repositories based on their ownership.
func (s *Postgres) ListAggregateByOwnerUsername(username string) ([]*RepositoryAggregate, error) {
	query := `
SELECT
	id,
	name,
	description,
	website,
	default_branch,
	private,
	bare,
	created_at,
	updated_at,
	(SELECT 42) AS stars,
	(SELECT 23) AS forks
FROM repositories
WHERE owner_id = (SELECT id
		  FROM users
		  WHERE username = $1)
ORDER BY updated_at DESC`

	rows, err := s.db.Query(query, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var repositories []*RepositoryAggregate
	for rows.Next() {
		var id string
		var name string
		var description sql.NullString
		var website sql.NullString
		var defaultBranch string
		var private bool
		var bare bool
		var created time.Time
		var updated time.Time
		var stars int
		var forks int

		rows.Scan(
			&id,
			&name,
			&description,
			&website,
			&defaultBranch,
			&private,
			&bare,
			&created,
			&updated,
			&stars,
			&forks,
		)

		repositories = append(repositories, &RepositoryAggregate{
			Repository: &Repository{
				ID:            id,
				Name:          name,
				Description:   description.String,
				Website:       website.String,
				DefaultBranch: defaultBranch,
				Private:       private,
				Bare:          bare,
				Created:       created,
				Updated:       updated,
			},
			Stars: stars,
			Forks: forks,
		})
	}

	return repositories, nil
}
