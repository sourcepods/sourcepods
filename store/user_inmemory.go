package store

import (
	"math/rand"
	"sync"
	"time"

	"github.com/gitloud/gitloud"
	"github.com/go-errors/errors"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var UserNotFound = errors.New("user not found")

type UserInMemory struct {
	mu    sync.RWMutex
	users []gitloud.User
}

func NewUserInMemory() *UserInMemory {
	return &UserInMemory{
		users: []gitloud.User{{
			ID:       "abcd-efgh-1234-5678",
			Username: "metalmatze",
			Name:     "Matthias Loibl",
			Email:    "mail@matthiasloibl.com",
			Password: "encrypted with bcrypt",
		}, {
			ID:       "bcde-fghi-2345-6789",
			Username: "tboerger",
			Name:     "Thomas Boerger",
			Email:    "thomas@webhippie.de",
			Password: "encrypted with bcrypt",
		}},
	}
}

func (s *UserInMemory) List() ([]gitloud.User, error) {
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)

	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.users, nil
}

func (s *UserInMemory) GetUser(username string) (gitloud.User, error) {
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)

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
	time.Sleep(time.Duration(rand.Intn(1500)) * time.Millisecond)

	s.mu.Lock()
	defer s.mu.Unlock()
	s.users = append(s.users, user)

	return nil
}

func (s *UserInMemory) UpdateUser(username string, updateUser gitloud.User) error {
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)

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
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)

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
