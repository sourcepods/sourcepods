package repository

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
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
func (s *Postgres) List(owner *Owner) ([]*Repository, []*Stats, *Owner, error) {
	if owner.ID == "" && owner.Username != "" {
		query := `SELECT id FROM users WHERE username = $1`
		row := s.db.QueryRow(query, owner.Username)

		var id string
		row.Scan(&id)

		if id == "" {
			return nil, nil, nil, OwnerNotFoundError
		}

		owner.ID = id
	}

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

	rows, err := s.db.Query(query, owner.Username)
	if err != nil {
		return nil, nil, nil, err
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

func (s *Postgres) Find(owner *Owner, name string) (*Repository, *Stats, *Owner, error) {
	query := `
SELECT
	id,
	description,
	website,
	default_branch,
	private,
	bare,
	created_at,
	updated_at,
	owner_id
FROM repositories
WHERE
	name = $2 AND
	owner_id = (SELECT id FROM users WHERE username = $1) `

	row := s.db.QueryRow(query, owner.Username, name)

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
		return nil, nil, nil, err
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

func (s *Postgres) Create(owner *Owner, r *Repository) (*Repository, error) {
	query := `
INSERT INTO repositories (owner_id, name, description, website, default_branch, private, bare)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, created_at, updated_at`

	var description *string
	if r.Description != "" {
		description = &r.Description
	}

	var website *string
	if r.Website != "" {
		website = &r.Website
	}

	row := s.db.QueryRow(query,
		owner.ID,
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
