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
func (s *Postgres) List(ctx context.Context, owner string) ([]*Repository, []*Stats, string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.Postgres.List")
	span.SetTag("owner", owner)
	defer span.Finish()

	row := s.db.QueryRowContext(ctx, userIdByUsername, owner)

	var owner_id string
	row.Scan(&owner_id)

	if owner_id == "" {
		return nil, nil, "", OwnerNotFoundError
	}

	span.SetTag("owner_id", owner_id)
	span.SetTag("owner_username", owner)

	rows, err := s.db.QueryContext(ctx, listByOwnerId, owner_id)
	if err != nil {
		return nil, nil, "", err
	}
	defer rows.Close()

	var repositories []*Repository
	var stats []*Stats

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

		repositories = append(repositories, &Repository{
			ID:            id,
			Name:          name,
			Description:   description.String,
			Website:       website.String,
			DefaultBranch: defaultBranch,
			Private:       private,
			Bare:          bare,
			Created:       created,
			Updated:       updated,
		})

		stats = append(stats, &Stats{
			Stars:                  stars,
			Forks:                  forks,
			IssueTotalCount:        66,
			IssueOpenCount:         35,
			IssueClosedCount:       31,
			PullRequestTotalCount:  20,
			PullRequestOpenCount:   4,
			PullRequestClosedCount: 16,
		})
	}

	return repositories, stats, owner, nil
}

func (s *Postgres) Find(ctx context.Context, owner string, name string) (*Repository, *Stats, string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.Postgres.Find")
	span.SetTag("owner_username", owner)
	defer span.Finish()

	row := s.db.QueryRowContext(ctx, findByOwnerAndName, owner, name)

	var id string
	var description sql.NullString
	var website sql.NullString
	var defaultBranch string
	var private bool
	var bare bool
	var created time.Time
	var updated time.Time
	var ownerID string

	if err := row.Scan(
		&id,
		&description,
		&website,
		&defaultBranch,
		&private,
		&bare,
		&created,
		&updated,
		&ownerID,
	); err != nil {
		return nil, nil, "", err
	}

	return &Repository{
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
		&Stats{
			Stars:                  42,
			Forks:                  23,
			IssueTotalCount:        66,
			IssueOpenCount:         13,
			IssueClosedCount:       53,
			PullRequestTotalCount:  20,
			PullRequestOpenCount:   2,
			PullRequestClosedCount: 18,
		},
		owner,
		nil
}

func (s *Postgres) Create(ctx context.Context, owner string, r *Repository) (*Repository, error) {
	row := s.db.QueryRowContext(ctx, userIdByUsername, owner)

	var owner_id string
	row.Scan(&owner_id)

	if owner_id == "" {
		return nil, OwnerNotFoundError
	}

	var description *string
	if r.Description != "" {
		description = &r.Description
	}

	var website *string
	if r.Website != "" {
		website = &r.Website
	}

	row = s.db.QueryRowContext(ctx, create,
		owner_id,
		r.Name,
		description,
		website,
		r.DefaultBranch,
		r.Private,
		r.Bare,
	)

	if err := row.Scan(&r.ID, &r.Created, &r.Updated); err != nil {
		if err, ok := err.(*pq.Error); ok {
			if err.Code == pq.ErrorCode("23505") {
				return nil, AlreadyExistsError
			}
		}
		return nil, err
	}

	return r, nil
}
