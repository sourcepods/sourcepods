package user

import (
	"database/sql"
	"errors"

	"github.com/gitpods/gitpods"
)

type postgres struct {
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) *postgres {
	return &postgres{db: db}
}

func (s *postgres) FindAll() ([]*gitpods.User, error) {
	rows, err := s.db.Query(`SELECT id, email, username, name FROM users ORDER BY name ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*gitpods.User
	for rows.Next() {
		var id string
		var email string
		var username string
		var name string
		rows.Scan(&id, &email, &username, &name)

		users = append(users, &gitpods.User{
			ID:       id,
			Email:    email,
			Username: username,
			Name:     name,
		})
	}

	return users, nil
}

func (s *postgres) Find(id string) (*gitpods.User, error) {
	panic("implement me")
}

func (s *postgres) FindByUsername(username string) (*gitpods.User, error) {
	row := s.db.QueryRow(`SELECT id, email, username, name, password FROM users WHERE username = $1 LIMIT 1`, username)

	var id string
	var email string
	var uusername string
	var name string
	var password string
	if err := row.Scan(&id, &email, &uusername, &name, &password); err != nil {
		return nil, err
	}

	return &gitpods.User{
		ID:       id,
		Email:    email,
		Username: username,
		Name:     name,
		Password: password,
	}, nil
}

func (s *postgres) FindUserByEmail(email string) (*gitpods.User, error) {
	row := s.db.QueryRow(`SELECT id, username, name, password FROM users WHERE email = $1 LIMIT 1`, email)

	var id string
	var username string
	var name string
	var password string
	if err := row.Scan(&id, &username, &name, &password); err != nil {
		return nil, err
	}

	return &gitpods.User{
		ID:       id,
		Email:    email,
		Username: username,
		Name:     name,
		Password: password,
	}, nil
}

func (s *postgres) Create(*gitpods.User) (*gitpods.User, error) {
	panic("implement me")
}

func (s *postgres) Update(username string, user *gitpods.User) (*gitpods.User, error) {
	stmt, err := s.db.Prepare(`UPDATE users SET username=$1, email=$2, name=$3 WHERE username=$1`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(string(user.Username), user.Email, user.Name)
	if err != nil {
		return nil, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if affected != 1 {
		return nil, errors.New("no rows updated")
	}

	return s.FindByUsername(username)
}

func (s *postgres) Delete(string) error {
	panic("implement me")
}
