package store

import (
	"sync"

	"github.com/gitpods/gitpods"
	"github.com/go-errors/errors"
	"golang.org/x/crypto/bcrypt"
)

var UserNotFound = errors.New("user not found")

type UserInMemory struct {
	mu    sync.RWMutex
	users []gitpods.User
}

func NewUserInMemory() *UserInMemory {
	pass1, _ := bcrypt.GenerateFromPassword([]byte("kubernetes"), bcrypt.DefaultCost)
	pass2, _ := bcrypt.GenerateFromPassword([]byte("golang"), bcrypt.DefaultCost)

	return &UserInMemory{
		users: []gitpods.User{{
			ID:       "25558000-2565-48dc-84eb-18754da2b0a2",
			Username: "metalmatze",
			Name:     "Matthias Loibl",
			Email:    "metalmatze@example.com",
			Password: string(pass1),
		}, {
			ID:       "911d24ae-ad9b-4e50-bf23-9dcbdc8134c6",
			Username: "tboerger",
			Name:     "Thomas Boerger",
			Email:    "tboerger@example.com",
			Password: string(pass2),
		}},
	}
}

func (s *UserInMemory) List() ([]gitpods.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.users, nil
}

func (s *UserInMemory) GetUser(username string) (gitpods.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, user := range s.users {
		if user.Username == username {
			return user, nil
		}
	}

	return gitpods.User{}, UserNotFound
}

func (s *UserInMemory) GetUserByEmail(email string) (gitpods.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, user := range s.users {
		if user.Email == email {
			return user, nil
		}
	}
	return gitpods.User{}, UserNotFound
}

func (s *UserInMemory) CreateUser(user gitpods.User) (gitpods.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.users = append(s.users, user)

	return user, nil
}

func (s *UserInMemory) UpdateUser(username string, updatedUser gitpods.User) (gitpods.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, user := range s.users {
		if user.Username == username {
			s.users[i] = updatedUser
			return updatedUser, nil
		}
	}
	return updatedUser, UserNotFound
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
