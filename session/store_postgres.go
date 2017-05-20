package session

import "database/sql"

// NewPostgresStore returns a Postgres implementation of the Store.
func NewPostgresStore(db *sql.DB) Store {
	return &Postgres{db: db}
}

// Postgres implementation of the Store.
type Postgres struct {
	db *sql.DB
}

// SaveSession in the store.
func (s *Postgres) SaveSession(session *Session) error {
	return s.db.QueryRow(
		`INSERT INTO sessions(expires, owner_id) VALUES($1, $2) RETURNING id`,
		session.Expiry, session.User.ID,
	).Scan(&session.ID)
}

// FindSession that aren't expired
func (s *Postgres) FindSession(id string) (*Session, error) {
	query := `
SELECT
	sessions.id,
	sessions.expires,
	users.id       AS user_id,
	users.username AS user_username
FROM sessions
	JOIN users ON sessions.owner_id = users.id
WHERE sessions.id = $1`

	row := s.db.QueryRow(query, id)

	session := Session{
		User: User{},
	}

	err := row.Scan(&session.ID, &session.Expiry, &session.User.ID, &session.User.Username)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

// ClearSessions that are expired.
func (s *Postgres) ClearSessions() (int64, error) {
	res, err := s.db.Exec(`DELETE FROM sessions WHERE expires < now()`)
	if err != nil {
		return 0, nil
	}
	return res.RowsAffected()
}
