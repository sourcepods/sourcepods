package store

import (
	"sync"

	"github.com/gitpods/gitpods"
	"github.com/go-errors/errors"
	"golang.org/x/crypto/bcrypt"
)

var UserNotFound = errors.New("user not found")

type UsersInMemory struct {
	mu    sync.RWMutex
	users []*gitpods.User
}

func NewUsersInMemory() *UsersInMemory {
	pass1, _ := bcrypt.GenerateFromPassword([]byte("kubernetes"), bcrypt.DefaultCost)
	pass2, _ := bcrypt.GenerateFromPassword([]byte("golang"), bcrypt.DefaultCost)

	return &UsersInMemory{
		users: []*gitpods.User{{
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

func (s *UsersInMemory) List() ([]*gitpods.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.users, nil
}

func (s *UsersInMemory) GetUser(username string) (*gitpods.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, user := range s.users {
		if user.Username == username {
			return user, nil
		}
	}

	return nil, UserNotFound
}

func (s *UsersInMemory) GetUserByEmail(email string) (*gitpods.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, user := range s.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, UserNotFound
}

func (s *UsersInMemory) CreateUser(user *gitpods.User) (*gitpods.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.users = append(s.users, user)

	return user, nil
}

func (s *UsersInMemory) UpdateUser(username string, updatedUser *gitpods.User) (*gitpods.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, user := range s.users {
		if user.Username == username {
			s.users[i].Username = updatedUser.Username
			s.users[i].Name = updatedUser.Name
			s.users[i].Email = updatedUser.Email
			return updatedUser, nil
		}
	}
	return updatedUser, UserNotFound
}

func (s *UsersInMemory) DeleteUser(username string) error {
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
