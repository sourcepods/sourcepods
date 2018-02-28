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
func (s *Postgres) Save(ctx context.Context, session *Session) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "session.Postgres.SaveSession")
	defer span.Finish()

	save := `INSERT INTO sessions (expires, owner_id) VALUES ($1, $2) RETURNING id;`

	return s.db.QueryRowContext(
		ctx,
		save,
		session.Expiry, session.User.ID,
	).Scan(&session.ID)
}

// Find that aren't expired
func (s *Postgres) Find(ctx context.Context, id string) (*Session, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "session.Postgres.Find")
	span.SetTag("id", id)
	defer span.Finish()

	findByID := `
SELECT
	sessions.id,
	sessions.expires,
	users.id       AS user_id,
	users.username AS user_username
FROM sessions
	JOIN users ON sessions.owner_id = users.id
WHERE sessions.id = $1;
`

	row := s.db.QueryRowContext(ctx, findByID, id)

	session := Session{
		User: User{},
	}

	err := row.Scan(&session.ID, &session.Expiry, &session.User.ID, &session.User.Username)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

// DeleteExpired sessions that are expired.
func (s *Postgres) DeleteExpired(ctx context.Context) (int64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "session.Postgres.DeleteExpired")
	defer span.Finish()

	clearExpired := `DELETE FROM sessions WHERE expires < now();`

	res, err := s.db.ExecContext(ctx, clearExpired)
	if err != nil {
		return 0, nil
	}
	return res.RowsAffected()
}
