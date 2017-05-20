package session

import (
	"time"
)

// Service creates, finds and clears sessions from their store.
type Service interface {
	CreateSession(string, string) (*Session, error)
	FindSession(string) (*Session, error)
	ClearSessions() (int64, error)
}

// Store session in a database.
type Store interface {
	SaveSession(*Session) error
	FindSession(string) (*Session, error)
	ClearSessions() (int64, error)
}

// NewService that talks to the store and returns sessions.
func NewService(store Store) Service {
	return &service{store: store}
}

type service struct {
	store Store
}

func (s *service) CreateSession(userID, userUsername string) (*Session, error) {
	sess := &Session{
		Expiry: time.Now().Add(defaultExpiry),
		User: User{
			ID:       userID,
			Username: userUsername,
		},
	}

	if err := s.store.SaveSession(sess); err != nil {
		return nil, err
	}

	return sess, nil
}

func (s *service) FindSession(id string) (*Session, error) {
	return s.store.FindSession(id)
}

func (s *service) ClearSessions() (int64, error) {
	return s.store.ClearSessions()
}
