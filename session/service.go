package session

import (
	"context"
	"time"
)

// Service creates, finds and clears sessions from their store.
type Service interface {
	CreateSession(context.Context, string, string) (*Session, error)
	FindSession(context.Context, string) (*Session, error)
	ClearSessions(context.Context) (int64, error)
}

// Store session in a database.
type Store interface {
	SaveSession(context.Context, *Session) error
	FindSession(context.Context, string) (*Session, error)
	ClearSessions(context.Context) (int64, error)
}

// NewService that talks to the store and returns sessions.
func NewService(store Store) Service {
	return &service{store: store}
}

type service struct {
	store Store
}

func (s *service) CreateSession(ctx context.Context, userID, userUsername string) (*Session, error) {
	sess := &Session{
		Expiry: time.Now().Add(defaultExpiry),
		User: User{
			ID:       userID,
			Username: userUsername,
		},
	}

	if err := s.store.SaveSession(ctx, sess); err != nil {
		return nil, err
	}

	return sess, nil
}

func (s *service) FindSession(ctx context.Context, id string) (*Session, error) {
	return s.store.FindSession(ctx, id)
}

func (s *service) ClearSessions(ctx context.Context) (int64, error) {
	return s.store.ClearSessions(ctx)
}
