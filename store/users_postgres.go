package store

import (
	"database/sql"

	"github.com/gitpods/gitpods"
)

type UsersPostgres struct {
	db *sql.DB
}

func NewUsersPostgres(db *sql.DB) *UsersPostgres {
	return &UsersPostgres{db: db}
}

func (s UsersPostgres) List() ([]*gitpods.User, error) {
	rows, err := s.db.Query(`SELECT * FROM users;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*gitpods.User
	for rows.Next() {
		var user *gitpods.User
		rows.Scan(user.ID, user.Email, user.Username, user.Name)
		users = append(users, user)
	}

	return users, nil
}

func (s *UsersPostgres) GetUserByUsername(username string) (*gitpods.User, error) {
	row := s.db.QueryRow(`SELECT * FROM users WHERE username = $1 LIMIT 1`, username)

	var user *gitpods.User
	if err := row.Scan(&user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UsersPostgres) GetUserByEmail(email string) (*gitpods.User, error) {
	row := s.db.QueryRow(`SELECT id, email, username, name, password FROM users WHERE email = $1 LIMIT 1`, email)

	var user gitpods.User
	if err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Name, &user.Password); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UsersPostgres) CreateUser(*gitpods.User) (*gitpods.User, error) {
	panic("implement me")
}

func (s *UsersPostgres) UpdateUser(username string, user *gitpods.User) (*gitpods.User, error) {
	panic("implement me")
}

func (s *UsersPostgres) DeleteUser(username string) error {
	panic("implement me")
}
