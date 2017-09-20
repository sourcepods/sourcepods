package user

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
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
func (s *Postgres) FindAll() ([]*User, error) {
	rows, err := s.db.Query(`SELECT
	id,
	email,
	username,
	name,
	created_at,
	updated_at
FROM users
ORDER BY name ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		var id string
		var email string
		var username string
		var name string
		var created time.Time
		var updated time.Time
		rows.Scan(&id, &email, &username, &name, &created, &updated)

		users = append(users, &User{
			ID:       id,
			Email:    email,
			Username: username,
			Name:     name,
			Created:  created,
			Updated:  updated,
		})
	}

	return users, nil
}

// Find a user by its ID.
func (s *Postgres) Find(id string) (*User, error) {
	query := `SELECT
	username,
	email,
	name,
	password,
	created_at,
	updated_at
FROM users
WHERE id = $1
LIMIT 1`

	row := s.db.QueryRow(query, id)

	var username string
	var email string
	var name string
	var password string
	var created time.Time
	var updated time.Time
	if err := row.Scan(&username, &email, &name, &password, &created, &updated); err != nil {
		if err == sql.ErrNoRows {
			return nil, NotFoundError
		}
		return nil, err
	}

	return &User{
		ID:       id,
		Email:    email,
		Username: username,
		Name:     name,
		Password: password,
		Created:  created,
		Updated:  updated,
	}, nil
}

// FindByUsername finds a user by its username.
func (s *Postgres) FindByUsername(ctx context.Context, username string) (*User, error) {
	query := `SELECT
	id,
	email,
	name,
	password,
	created_at,
	updated_at
FROM users
WHERE username = $1
LIMIT 1`

	row := s.db.QueryRow(query, username)

	var id string
	var email string
	var name string
	var password string
	var created time.Time
	var updated time.Time
	if err := row.Scan(&id, &email, &name, &password, &created, &updated); err != nil {
		if err == sql.ErrNoRows {
			return nil, NotFoundError
		}
		return nil, err
	}

	return &User{
		ID:       id,
		Email:    email,
		Username: username,
		Name:     name,
		Password: password,
		Created:  created,
		Updated:  updated,
	}, nil
}

// FindUserByEmail by its email.
func (s *Postgres) FindUserByEmail(email string) (*User, error) {
	row := s.db.QueryRow(`SELECT
	id,
	username,
	name,
	password,
	created_at,
	updated_at
FROM users
WHERE email = $1
LIMIT 1`, email)

	var id string
	var username string
	var name string
	var password string
	var created time.Time
	var updated time.Time
	if err := row.Scan(&id, &username, &name, &password, &created, &updated); err != nil {
		return nil, err
	}

	return &User{
		ID:       id,
		Email:    email,
		Username: username,
		Name:     name,
		Password: password,
		Created:  created,
		Updated:  updated,
	}, nil
}

// Create a user in postgres and return it with the ID set in the store.
func (s *Postgres) Create(u *User) (*User, error) {
	pass, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	err = s.db.QueryRow(
		`INSERT INTO users (email, username, name, password) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`,
		u.Email, u.Username, u.Name, pass,
	).Scan(&u.ID, &u.Created, &u.Updated)

	return u, err
}

// Update a user by its username.
// TODO: Update users by their id?
func (s *Postgres) Update(user *User) (*User, error) {
	stmt, err := s.db.Prepare(`UPDATE users SET username = $2, email = $3, name = $4, updated_at = now() WHERE id = $1`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(user.ID, user.Username, user.Email, user.Name)
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

	return s.Find(user.ID)
}

// Delete a user by its id.
func (s *Postgres) Delete(id string) error {
	panic("implement me")
}
