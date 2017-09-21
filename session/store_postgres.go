package session

import (
	"context"
	"database/sql"

	"github.com/opentracing/opentracing-go"
)

// NewPostgresStore returns a Postgres implementation of the Store.
func NewPostgresStore(db *sql.DB) Store {
	return &Postgres{db: db}
}

// Postgres implementation of the Store.
type Postgres struct {
	db *sql.DB
}

// SaveSession in the store.
func (s *Postgres) SaveSession(ctx context.Context, session *Session) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "session.Postgres.SaveSession")
	defer span.Finish()

	return s.db.QueryRowContext(
		ctx,
		`INSERT INTO sessions(expires, owner_id) VALUES($1, $2) RETURNING id`,
		session.Expiry, session.User.ID,
	).Scan(&session.ID)
}

// FindSession that aren't expired
func (s *Postgres) FindSession(ctx context.Context, id string) (*Session, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "session.Postgres.FindSession")
	span.SetTag("id", id)
	defer span.Finish()

	query := `
SELECT
	sessions.id,
	sessions.expires,
	users.id       AS user_id,
	users.username AS user_username
FROM sessions
	JOIN users ON sessions.owner_id = users.id
WHERE sessions.id = $1`

	row := s.db.QueryRowContext(ctx, query, id)

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
func (s *Postgres) ClearSessions(ctx context.Context) (int64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "session.Postgres.ClearSessions")
	defer span.Finish()

	res, err := s.db.ExecContext(ctx, `DELETE FROM sessions WHERE expires < now()`)
	if err != nil {
		return 0, nil
	}
	return res.RowsAffected()
}
