package v1

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gitpods/gitpods/pkg/gitpods/user"

	"github.com/gitpods/gitpods/pkg/gitpods/repository"
	"github.com/gitpods/gitpods/pkg/storage"
)

type repositoryTestService struct{}

func (repositoryTestService) List(ctx context.Context, owner string) ([]*repository.Repository, string, error) {
	panic("implement me")
}

func (repositoryTestService) Find(ctx context.Context, owner string, name string) (*repository.Repository, string, error) {
	panic("implement me")
}

func (repositoryTestService) Create(ctx context.Context, owner string, repository *repository.Repository) (*repository.Repository, error) {
	panic("implement me")
}

func (repositoryTestService) Branches(ctx context.Context, owner string, name string) ([]*repository.Branch, error) {
	panic("implement me")
}

func (repositoryTestService) Commit(ctx context.Context, owner string, name string, rev string) (storage.Commit, error) {
	panic("implement me")
}

type userTestService struct {
	FinAll func(context.Context) ([]*user.User, error)
}

func (u userTestService) FindAll(ctx context.Context) ([]*user.User, error) {
	return u.FinAll(ctx)
}

func (u userTestService) Find(context.Context, string) (*user.User, error) {
	panic("implement me")
}

func (u userTestService) FindByUsername(context.Context, string) (*user.User, error) {
	panic("implement me")
}

func (u userTestService) FindRepositoryOwner(ctx context.Context, repositoryID string) (*user.User, error) {
	panic("implement me")
}

func (u userTestService) Create(context.Context, *user.User) (*user.User, error) {
	panic("implement me")
}

func (u userTestService) Update(context.Context, *user.User) (*user.User, error) {
	panic("implement me")
}

func (u userTestService) Delete(context.Context, string) error { panic("implement me") }

func TestUsersListUsersHandler(t *testing.T) {
	findAll := func(ctx context.Context) ([]*user.User, error) {
		return []*user.User{{
			ID:       "2849392e-6eca-43f0-9bec-b16beac5c2b1",
			Email:    "mail@example.com",
			Username: "username",
			Name:     "User Name",
			Password: "secret",
		}}, nil
	}

	api, err := New(repositoryTestService{}, userTestService{FinAll: findAll})
	assert.NoError(t, err)

	ts := httptest.NewServer(api.Handler)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/v1/users")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	defer res.Body.Close()

	expectedUser := `[{
		"id": "2849392e-6eca-43f0-9bec-b16beac5c2b1",
		"email": "mail@example.com",
		"username": "username",
		"name": "User Name",
		"created_at": "0001-01-01T00:00:00.000Z",
		"updated_at": "0001-01-01T00:00:00.000Z"
	}]`

	body, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.JSONEq(t, expectedUser, string(body))
}
