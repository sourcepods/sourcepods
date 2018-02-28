package session

import (
	"context"
	"time"
)

// Service creates, finds and clears sessions from their store.
type Service interface {
	Create(ctx context.Context, userID string, userUsername string) (*Session, error)
	Find(ctx context.Context, id string) (*Session, error)
	Delete(ctx context.Context, id string) error
	DeleteExpired(ctx context.Context) (int64, error)
}

// Store session in a database.
type Store interface {
	Save(ctx context.Context, s *Session) error
	Find(ctx context.Context, id string) (*Session, error)
	Delete(ctx context.Context, id string) error
	DeleteExpired(ctx context.Context) (int64, error)
}

// NewService that talks to the store and returns sessions.
func NewService(store Store) Service {
	return &service{store: store}
}

type service struct {
	store Store
}

func (s *service) Create(ctx context.Context, userID, userUsername string) (*Session, error) {
	sess := &Session{
		Expiry: time.Now().Add(defaultExpiry),
		User: User{
			ID:       userID,
			Username: userUsername,
		},
	}

	if err := s.store.Save(ctx, sess); err != nil {
		return nil, err
	}

	return sess, nil
}

func (s *service) Find(ctx context.Context, id string) (*Session, error) {
	return s.store.Find(ctx, id)
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.store.Delete(ctx, id)
}

func (s *service) DeleteExpired(ctx context.Context) (int64, error) {
	return s.store.DeleteExpired(ctx)
}
