package user

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/opentracing/opentracing-go"
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
func (s *Postgres) FindAll(ctx context.Context) ([]*User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "user.Postgres.FinAll")
	defer span.Finish()

	findAll := `
SELECT
	id,
	email,
	username,
	name,
	created_at,
	updated_at
FROM users
ORDER BY name ASC;
`

	rows, err := s.db.QueryContext(ctx, findAll)
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
func (s *Postgres) Find(ctx context.Context, id string) (*User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "user.Postgres.Find")
	span.SetTag("id", id)
	defer span.Finish()

	findByID := `
SELECT
	username,
	email,
	name,
	password,
	created_at,
	updated_at
FROM users
WHERE id = $1
LIMIT 1;
`

	row := s.db.QueryRowContext(ctx, findByID, id)

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
	span, ctx := opentracing.StartSpanFromContext(ctx, "user.Postgres.FindByUsername")
	span.SetTag("username", username)
	defer span.Finish()

	findByUsername := `
SELECT
	id,
	email,
	name,
	password,
	created_at,
	updated_at
FROM users
WHERE username = $1
LIMIT 1;
`

	row := s.db.QueryRowContext(ctx, findByUsername, username)

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
func (s *Postgres) FindUserByEmail(ctx context.Context, email string) (*User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "user.Postgres.FindUserByEmail")
	span.SetTag("email", email)
	defer span.Finish()

	findByEmail := `
SELECT
  id,
  username,
  name,
  password,
  created_at,
  updated_at
FROM users
WHERE email = $1
LIMIT 1;
`

	row := s.db.QueryRowContext(ctx, findByEmail, email)

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

func (s *Postgres) FindRepositoryOwner(ctx context.Context, repositoryID string) (*User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "user.Postgres.FindRepositoryOwner")
	span.SetTag("repository", repositoryID)
	defer span.Finish()

	findRepositoryOwner := `
SELECT
  id,
  email,
  username,
  name,
  created_at,
  updated_at
FROM users
WHERE id = (SELECT owner_id
            FROM repositories
            WHERE repositories.id = $1)
`

	row := s.db.QueryRowContext(ctx, findRepositoryOwner, repositoryID)

	var id string
	var email string
	var username string
	var name string
	var created time.Time
	var updated time.Time
	if err := row.Scan(&id, &email, &username, &name, &created, &updated); err != nil {
		return nil, err
	}

	return &User{
		ID:       id,
		Email:    email,
		Username: username,
		Name:     name,
		Created:  created,
		Updated:  updated,
	}, nil
}

// Create a user in postgres and return it with the ID set in the store.
func (s *Postgres) Create(ctx context.Context, u *User) (*User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "user.Postgres.Create")
	span.SetTag("user_username", u.Username)
	span.SetTag("user_email", u.Email)
	defer span.Finish()

	pass, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	create := `
INSERT INTO users (email, username, name, password) VALUES ($1, $2, $3, $4)
RETURNING id, created_at, updated_at;
`

	err = s.db.QueryRowContext(
		ctx,
		create,
		u.Email, u.Username, u.Name, pass,
	).Scan(&u.ID, &u.Created, &u.Updated)

	return u, err
}

// Update a user by its username.
// TODO: Update users by their id?
func (s *Postgres) Update(ctx context.Context, user *User) (*User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "user.Postgres.Update")
	span.SetTag("id", user.ID)
	span.SetTag("username", user.Username)
	defer span.Finish()

	update := `
UPDATE users
SET username = $2, email = $3, name = $4, updated_at = now()
WHERE id = $1;
`

	stmt, err := s.db.PrepareContext(ctx, update)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, user.ID, user.Username, user.Email, user.Name)
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

	return s.Find(ctx, user.ID)
}

// Delete a user by its id.
func (s *Postgres) Delete(ctx context.Context, id string) error {
	panic("implement me")
}
