package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
	"github.com/opentracing/opentracing-go"
)

// Postgres implementation of the Store.
type Postgres struct {
	db *sql.DB
}

// NewPostgresStore returns a Postgres implementation of the Store.
func NewPostgresStore(db *sql.DB) *Postgres {
	return &Postgres{db: db}
}

// List retrieves a list of repositories based on their ownership.
func (s *Postgres) List(ctx context.Context, owner string) ([]*Repository, string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.Postgres.List")
	span.SetTag("owner", owner)
	defer span.Finish()

	listByOwnerID := `
SELECT
	id,
	name,
	description,
	website,
	default_branch,
	created_at,
	updated_at,
	(SELECT 42) AS stars,
	(SELECT 23) AS forks
FROM repositories
WHERE owner_id = (SELECT id FROM users WHERE username = $1 LIMIT 1)
ORDER BY updated_at DESC;
`

	rows, err := s.db.QueryContext(ctx, listByOwnerID, owner)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	var repositories []*Repository

	for rows.Next() {
		var id string
		var name string
		var description sql.NullString
		var website sql.NullString
		var defaultBranch string
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
			&created,
			&updated,
			&stars,
			&forks,
		)

		repositories = append(repositories, &Repository{
			ID:            id,
			Name:          name,
			Description:   description.String,
			Website:       website.String,
			DefaultBranch: defaultBranch,
			Created:       created,
			Updated:       updated,
		})
	}

	return repositories, owner, nil
}

func (s *Postgres) Find(ctx context.Context, owner string, name string) (*Repository, string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.Postgres.Find")
	span.SetTag("owner", owner)
	defer span.Finish()

	findByOwnerAndName := `
SELECT
	id,
	description,
	website,
	default_branch,
	created_at,
	updated_at,
	owner_id
FROM repositories
WHERE
	name = $2 AND
	owner_id = (SELECT id FROM users WHERE username = $1);`

	row := s.db.QueryRowContext(ctx, findByOwnerAndName, owner, name)

	var id string
	var description sql.NullString
	var website sql.NullString
	var defaultBranch string
	var created time.Time
	var updated time.Time
	var ownerID string

	if err := row.Scan(
		&id,
		&description,
		&website,
		&defaultBranch,
		&created,
		&updated,
		&ownerID,
	); err != nil {
		return nil, "", err
	}

	return &Repository{
			ID:            id,
			Name:          name,
			Description:   description.String,
			Website:       website.String,
			DefaultBranch: defaultBranch,
			Created:       created,
			Updated:       updated,
		},
		owner,
		nil
}

func (s *Postgres) Create(ctx context.Context, owner string, r *Repository) (*Repository, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.Postgres.Create")
	span.SetTag("owner", owner)
	span.SetTag("name", r.Name)
	defer span.Finish()

	var description *string
	if r.Description != "" {
		description = &r.Description
	}

	var website *string
	if r.Website != "" {
		website = &r.Website
	}

	create := `
INSERT INTO repositories (owner_id, name, description, website, default_branch)
VALUES ((SELECT id FROM users WHERE username = $1 LIMIT 1), $2, $3, $4, $5)
RETURNING id, created_at, updated_at;
`

	row := s.db.QueryRowContext(ctx, create,
		owner,
		r.Name,
		description,
		website,
		r.DefaultBranch,
	)

	if err := row.Scan(&r.ID, &r.Created, &r.Updated); err != nil {
		if err, ok := err.(*pq.Error); ok {
			if err.Code == pq.ErrorCode("23505") {
				return nil, ErrAlreadyExists
			}
		}
		return nil, err
	}

	return r, nil
}
