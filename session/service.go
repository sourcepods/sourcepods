package session

import (
	"time"
)

type Service interface {
	CreateSession(string, string) (*Session, error)
	FindSession(string) (*Session, error)
	ClearSessions() (int64, error)
}

type Store interface {
	SaveSession(*Session) error
	FindSession(string) (*Session, error)
	ClearSessions() (int64, error)
}

func NewService(store Store) Service {
	return &service{store: store}
}

type service struct {
	store Store
}

func (s *service) CreateSession(userID, userUsername string) (*Session, error) {
	sess := &Session{
		Expiry: time.Now().Add(DefaultSessionDuration),
		User: SessionUser{
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
