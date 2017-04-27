package store

import (
	"sync"

	"github.com/gitloud/gitloud"
	"github.com/go-errors/errors"
)

var UserNotFound = errors.New("user not found")

type UserInMemory struct {
	mu    sync.RWMutex
	users []gitloud.User
}

func NewUserInMemory() *UserInMemory {
	return &UserInMemory{
		users: []gitloud.User{{
			ID:       "25558000-2565-48dc-84eb-18754da2b0a2",
			Username: "metalmatze",
			Name:     "Matthias Loibl",
			Email:    "metalmatze@example.com",
			Password: "encrypted with bcrypt",
		}, {
			ID:       "911d24ae-ad9b-4e50-bf23-9dcbdc8134c6",
			Username: "tboerger",
			Name:     "Thomas Boerger",
			Email:    "tboerger@example.com",
			Password: "encrypted with bcrypt",
		}},
	}
}

func (s *UserInMemory) List() ([]gitloud.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.users, nil
}

func (s *UserInMemory) GetUser(username string) (gitloud.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, user := range s.users {
		if user.Username == username {
			return user, nil
		}
	}

	return gitloud.User{}, UserNotFound
}

func (s *UserInMemory) CreateUser(user gitloud.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.users = append(s.users, user)

	return nil
}

func (s *UserInMemory) UpdateUser(username string, updateUser gitloud.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, user := range s.users {
		if user.Username == username {
			s.users[i] = updateUser
			return nil
		}
	}
	return UserNotFound
}

func (s *UserInMemory) DeleteUser(username string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, user := range s.users {
		if user.Username == username {
			s.users = append(s.users[:i], s.users[i+1:]...)
			return nil
		}
	}
	return UserNotFound
}
