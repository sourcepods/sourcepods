package repository

import (
	"database/sql"
)

type postgres struct {
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) *postgres {
	return &postgres{db: db}
}

func (s *postgres) ListByOwner(id string) ([]*Repository, error) {
	rows, err := s.db.Query(`SELECT id, name, description, website, default_branch, private, bare FROM repositories WHERE owner_id = $1`, id)
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
		rows.Scan(&id, &name, &description, &website, &defaultBranch, &private, &bare)

		repositories = append(repositories, &Repository{
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
