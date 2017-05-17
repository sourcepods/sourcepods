package user

import (
	"sync"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

var NotFound = errors.New("user not found")

type memory struct {
	mu    sync.RWMutex
	users []*User
}

func NewMemoryRepository() *memory {
	pass1, _ := bcrypt.GenerateFromPassword([]byte("kubernetes"), bcrypt.DefaultCost)
	pass2, _ := bcrypt.GenerateFromPassword([]byte("golang"), bcrypt.DefaultCost)

	return &memory{
		users: []*User{{
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

func (r *memory) FindAll() ([]*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.users, nil
}

func (r *memory) Find(id ID) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, user := range r.users {
		if user.ID == id {
			return user, nil
		}
	}

	return nil, NotFound
}

func (r *memory) FindByUsername(username Username) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, user := range r.users {
		if user.Username == username {
			return user, nil
		}
	}

	return nil, NotFound
}

func (r *memory) Create(user *User) (*User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users = append(r.users, user)

	return user, nil
}

func (r *memory) Update(username Username, updated *User) (*User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, user := range r.users {
		if user.Username == username {
			r.users[i].Username = updated.Username
			r.users[i].Name = updated.Name
			r.users[i].Email = updated.Email
			return r.users[i], nil
		}
	}

	return nil, NotFound
}

func (r *memory) Delete(username Username) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, user := range r.users {
		if user.Username == username {
			r.users = append(r.users[:i], r.users[i+1:]...)
			return nil
		}
	}

	return NotFound
}
