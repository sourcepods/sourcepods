package user

import (
	"database/sql"
	"errors"

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

// FindAll users.
func (s *Postgres) FindAll() ([]*gitpods.User, error) {
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

// Find a user by its ID.
func (s *Postgres) Find(id string) (*gitpods.User, error) {
	panic("implement me")
}

// FindByUsername finds a user by its username.
func (s *Postgres) FindByUsername(username string) (*gitpods.User, error) {
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

// FindUserByEmail by its email.
func (s *Postgres) FindUserByEmail(email string) (*gitpods.User, error) {
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

// Create a user in postgres and return it with the ID set in the store.
func (s *Postgres) Create(*gitpods.User) (*gitpods.User, error) {
	panic("implement me")
}

// Update a user by its username.
// TODO: Update users by their id?
func (s *Postgres) Update(username string, user *gitpods.User) (*gitpods.User, error) {
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

// Delete a user by its id.
func (s *Postgres) Delete(id string) error {
	panic("implement me")
}
