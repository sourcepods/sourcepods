package repository

import (
	"database/sql"

	"github.com/gitpods/gitpods"
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
func (s *Postgres) ListByOwner(id string) ([]*gitpods.Repository, error) {
	rows, err := s.db.Query(`SELECT id, name, description, website, default_branch, private, bare FROM repositories WHERE owner_id = $1`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var repositories []*gitpods.Repository
	for rows.Next() {
		var id string
		var name string
		var description string
		var website string
		var defaultBranch string
		var private bool
		var bare bool
		rows.Scan(&id, &name, &description, &website, &defaultBranch, &private, &bare)

		repositories = append(repositories, &gitpods.Repository{
			ID:            id,
			Name:          name,
			Description:   description,
			Website:       website,
			DefaultBranch: defaultBranch,
			Private:       private,
			Bare:          bare,
		})
	}

	return repositories, nil
}
