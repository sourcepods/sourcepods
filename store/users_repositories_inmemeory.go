package store

import (
	"sync"

	"github.com/gitpods/gitpods"
)

type UsersRepositoriesInMemory struct {
	mu           sync.RWMutex
	users        *UsersInMemory
	repositories *RepositoriesInMemory
}

func NewUsersRepositoriesInMemory(users *UsersInMemory, repositories *RepositoriesInMemory) *UsersRepositoriesInMemory {
	return &UsersRepositoriesInMemory{
		users:        users,
		repositories: repositories,
	}
}

func (s *UsersRepositoriesInMemory) List(username string) ([]*gitpods.Repository, error) {
	user, err := s.users.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	repositories, err := s.repositories.List()
	if err != nil {
		return nil, err
	}

	var userRepos []*gitpods.Repository
	for _, repo := range repositories {
		if repo.Owner == user {
			userRepos = append(userRepos, &gitpods.Repository{
				ID:            repo.ID,
				Name:          repo.Name,
				Description:   repo.Description,
				Website:       repo.Website,
				DefaultBranch: repo.DefaultBranch,
				Private:       repo.Private,
				Bare:          repo.Bare,
				Created:       repo.Created,
				Updated:       repo.Updated,

				OwnerID: repo.Owner.ID,
			})
		}
	}

	return userRepos, nil
}
