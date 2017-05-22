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

// ListByOwner retrieves a list of repositories based on their ownership.
func (s *Postgres) ListByOwnerUsername(username string) ([]*Repository, error) {
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
	updated_at
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

	var repositories []*Repository
	for rows.Next() {
		var id string
		var name string
		var description string
		var website string
		var defaultBranch string
		var private bool
		var bare bool
		var created time.Time
		var updated time.Time
		rows.Scan(&id, &name, &description, &website, &defaultBranch, &private, &bare, &created, &updated)

		repositories = append(repositories, &Repository{
			ID:            id,
			Name:          name,
			Description:   description,
			Website:       website,
			DefaultBranch: defaultBranch,
			Private:       private,
			Bare:          bare,
			Created:       created,
			Updated:       updated,
		})
	}

	return repositories, nil
}
