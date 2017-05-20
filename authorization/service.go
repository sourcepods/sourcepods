package authorization

import (
	"github.com/gitpods/gitpods"
	"github.com/gitpods/gitpods/session"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	AuthenticateUser(email, password string) (*gitpods.User, error)
	CreateSession(string, string) (*session.Session, error)
}

type Store interface {
	FindUserByEmail(string) (*gitpods.User, error)
}

func NewService(store Store, sessions session.Service) Service {
	return &service{store: store, sessions: sessions}
}

type service struct {
	store    Store
	sessions session.Service
}

// AuthenticateUser by querying the store for the hashed password
// and comparing it against the one passed.
func (s *service) AuthenticateUser(email, password string) (*gitpods.User, error) {
	user, err := s.store.FindUserByEmail(email)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) CreateSession(userID, userUsername string) (*session.Session, error) {
	return s.sessions.CreateSession(userID, userUsername)
}
