package user

import (
	"database/sql"
	"errors"
)

type postgres struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *postgres {
	return &postgres{db: db}
}

func (r *postgres) FindAll() ([]*User, error) {
	rows, err := r.db.Query(`SELECT id, email, username, name FROM users`)
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
		rows.Scan(&id, &email, &username, &name)

		users = append(users, &User{
			ID:       id,
			Email:    email,
			Username: username,
			Name:     name,
		})
	}

	return users, nil
}

func (r *postgres) Find(id string) (*User, error) {
	panic("implement me")
}

func (r *postgres) FindByUsername(username string) (*User, error) {
	row := r.db.QueryRow(`SELECT id, email, username, name, password FROM users WHERE username = $1 LIMIT 1`, username)

	var id string
	var email string
	var uusername string
	var name string
	var password string
	if err := row.Scan(&id, &email, &uusername, &name, &password); err != nil {
		return nil, err
	}

	return &User{
		ID:       id,
		Email:    email,
		Username: username,
		Name:     name,
		Password: password,
	}, nil
}

func (r *postgres) Create(*User) (*User, error) {
	panic("implement me")
}

func (r *postgres) Update(username string, user *User) (*User, error) {
	stmt, err := r.db.Prepare(`UPDATE users SET username=$1, email=$2, name=$3 WHERE username=$1`)
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

	return r.FindByUsername(username)
}

func (r *postgres) Delete(string) error {
	panic("implement me")
}
