package store

import (
	"errors"
	"math/rand"
	"sync"
	"time"

	"fmt"

	"github.com/gitpods/gitpods"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var RepositoryNotFound = errors.New("repository not found")

type RepositoriesInMemory struct {
	mu           sync.RWMutex
	repositories []*gitpods.Repository
}

func NewRepositoriesInMemory(usersStore *UsersInMemory) *RepositoriesInMemory {
	var repositories []*gitpods.Repository

	users, _ := usersStore.List()

	for i := 0; i < 20; i++ {
		repo := &gitpods.Repository{
			ID:            fmt.Sprintf("25558000-2565-48dc-84eb-18754da2b0a%d", i),
			Name:          fmt.Sprintf("Project %d", i),
			Description:   fmt.Sprintf("Description for project %d", i),
			DefaultBranch: "master",
			Private:       true,
			Bare:          true,
			Created:       time.Now(),
			Updated:       time.Now(),
			Owner:         users[rand.Intn(len(users))],
		}
		repositories = append(repositories, repo)
	}

	return &RepositoriesInMemory{repositories: repositories}
}

func (s *RepositoriesInMemory) List() ([]*gitpods.Repository, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.repositories, nil
}
