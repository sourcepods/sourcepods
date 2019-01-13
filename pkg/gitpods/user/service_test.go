package user

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type store struct {
	users []*User
}

var testUsers = []*User{{
	ID:       "12b5e0b0-f8c4-4b32-bf4a-8fb77a9ca19e",
	Email:    "user1@example.com",
	Username: "user1",
	Name:     "User 1",
	Password: "password1",
	Created:  time.Now(),
	Updated:  time.Now(),
}, {
	ID:       "12b5e0b0-f8c4-4b21-bf4a-8fb77a7ca89e",
	Email:    "user2@example.com",
	Username: "user2",
	Name:     "User 2",
	Password: "password2",
	Created:  time.Now(),
	Updated:  time.Now(),
}}

func (s *store) FindAll(ctx context.Context) ([]*User, error) {
	return s.users, nil
}

func (s *store) Find(ctx context.Context, id string) (*User, error) {
	for _, u := range s.users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, NotFoundError
}

func (s *store) FindByUsername(ctx context.Context, username string) (*User, error) {
	for _, u := range s.users {
		if u.Username == username {
			return u, nil
		}
	}
	return nil, NotFoundError
}

func (s *store) FindRepositoryOwner(ctx context.Context, repositoryID string) (*User, error) {
	panic("implement me")
}

func (s *store) Create(context.Context, *User) (*User, error) {
	panic("implement me")
}

func (s *store) Update(ctx context.Context, u *User) (*User, error) {
	// Normally merge the current user with the new fields here.
	return s.users[0], nil
}

func (s *store) Delete(context.Context, string) error {
	panic("implement me")
}

func TestService_FindAll(t *testing.T) {
	service := NewService(&store{testUsers})

	users, err := service.FindAll(context.Background())
	assert.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, "user1", users[0].Username)
}

func TestService_Find(t *testing.T) {
	service := NewService(&store{testUsers})

	u1, err := service.Find(context.Background(), "12b5e0b0-f8c4-4b32-bf4a-8fb77a9ca19e")
	assert.NoError(t, err)
	assert.Equal(t, "12b5e0b0-f8c4-4b32-bf4a-8fb77a9ca19e", u1.ID)

	u2, err := service.Find(context.Background(), "12b5e0b0-f8c4-4b21-bf4a-8fb77a7ca89e")
	assert.NoError(t, err)
	assert.Equal(t, "12b5e0b0-f8c4-4b21-bf4a-8fb77a7ca89e", u2.ID)

	_, err = service.Find(context.Background(), "foobar")
	assert.Error(t, err)
	assert.Equal(t, NotFoundError, err)
}

func TestService_FindByUsername(t *testing.T) {
	service := NewService(&store{testUsers})

	u1, err := service.FindByUsername(context.Background(), "user1")
	assert.NoError(t, err)
	assert.Equal(t, "user1", u1.Username)

	u2, err := service.FindByUsername(context.Background(), "user2")
	assert.NoError(t, err)
	assert.Equal(t, "user2", u2.Username)

	_, err = service.FindByUsername(context.Background(), "foobar")
	assert.Error(t, err)
	assert.Equal(t, NotFoundError, err)
}

func TestService_Update(t *testing.T) {
	service := NewService(&store{testUsers})

	user, err := service.Update(context.Background(), testUsers[0])
	assert.NoError(t, err)
	assert.Equal(t, "12b5e0b0-f8c4-4b32-bf4a-8fb77a9ca19e", user.ID)
	assert.Equal(t, "user1@example.com", user.Email)
	assert.Equal(t, "user1", user.Username)
	assert.Equal(t, "User 1", user.Name)
	assert.Equal(t, "password1", user.Password)

	// Add new fields to this user until it equals testUser[0]
	newUser := &User{}

	user, err = service.Update(context.Background(), newUser)
	assert.Nil(t, user)
	assert.Error(t, err)
	assert.Equal(t, "id is not a valid uuid v4", err.Error())
	newUser.ID = "12b5e0b0-f8c4-4b32-bf4a-8fb77a9ca19e"

	user, err = service.Update(context.Background(), newUser)
	assert.Nil(t, user)
	assert.Error(t, err)
	assert.Equal(t, "email is not valid", err.Error())
	newUser.Email = "user1@example.com"

	user, err = service.Update(context.Background(), newUser)
	assert.Nil(t, user)
	assert.Error(t, err)
	assert.Equal(t, "username is not between 4 and 32 characters long", err.Error())
	newUser.Username = "user1"

	user, err = service.Update(context.Background(), newUser)
	assert.Nil(t, user)
	assert.Error(t, err)
	assert.Equal(t, "name is not between 2 and 64 characters long", err.Error())
	newUser.Name = "User 1"

	user, err = service.Update(context.Background(), newUser)
	assert.Equal(t, testUsers[0], user)
	assert.NoError(t, err)
}
