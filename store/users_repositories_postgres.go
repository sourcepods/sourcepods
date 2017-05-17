package store

import (
	"database/sql"

	"github.com/gitpods/gitpods"
)

type UsersRepositoriesPostgres struct {
	db *sql.DB
}

func NewUsersRepositoriesPostgres(db *sql.DB) *UsersRepositoriesPostgres {
	return &UsersRepositoriesPostgres{db: db}
}

func (s *UsersRepositoriesPostgres) List(username string) ([]*gitpods.Repository, error) {
	query := `
	SELECT id, name, description, website, default_branch, private, bare
		FROM repositories
		WHERE owner_id = (SELECT id
				FROM users
				WHERE username = $1)`

	rows, err := s.db.Query(query, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var repositories []*gitpods.Repository
	for rows.Next() {
		var repo gitpods.Repository
		rows.Scan(
			&repo.ID,
			&repo.Name,
			&repo.Description,
			&repo.Website,
			&repo.DefaultBranch,
			&repo.Private,
			&repo.Bare,
		)

		repositories = append(repositories, &repo)
	}

	return repositories, nil
}
