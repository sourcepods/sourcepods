package authorization

import (
	"context"
	"testing"
	"time"

	"github.com/gitpods/gitpods/session"
	"github.com/gitpods/gitpods/user"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

var (
	u1 = user.User{
		ID:       "12b5e0b0-f8c4-4b31-bf4a-8fb77a9ca89e",
		Email:    "foobar@example.com",
		Username: "foobar",
		Name:     "Foo Bar",
		Password: "$2y$10$o/x4Dnb/7wOAFTlWEwRBpuYQNg51v3gfl4v9hD0Hs3cgQrBfghCpy",
	}
	expiry = time.Date(2009, 11, 10, 23, 00, 00, 00, time.UTC)
)

type testStore struct{}

func (*testStore) FindUserByEmail(context.Context, string) (*user.User, error) {
	return &u1, nil
}

type sessionService struct{}

func (s sessionService) Create(ctx context.Context, id string, username string) (*session.Session, error) {
	return &session.Session{
		ID:     "410f59a5-75e6-4332-a0d3-ef06a0bfb2a5",
		Expiry: expiry,
		User: session.User{
			ID:       id,
			Username: username,
		},
	}, nil
}

func (s sessionService) Find(context.Context, string) (*session.Session, error) {
	// We don't need this for these tests.
	panic("implement me")
}

func (s sessionService) Delete(context.Context, string) error {
	// We don't need this for these tests.
	panic("implement me")
}

func (s sessionService) DeleteExpired(context.Context) (int64, error) {
	// We don't need this for these tests.
	panic("implement me")
}

func TestService_AuthenticateUser(t *testing.T) {
	store := &testStore{}
	ss := &sessionService{}
	s := NewService(store, ss)

	u, err := s.AuthenticateUser(context.Background(), "foobar@example.com", "bar")
	assert.Equal(t, bcrypt.ErrMismatchedHashAndPassword, err)
	assert.Nil(t, u)

	u, err = s.AuthenticateUser(context.Background(), "foobar@example.com", "baz")
	assert.NoError(t, err)
	assert.Equal(t, &u1, u)
}

func TestService_CreateSession(t *testing.T) {
	store := &testStore{}
	ss := &sessionService{}
	s := NewService(store, ss)

	expected := session.Session{
		ID:     "410f59a5-75e6-4332-a0d3-ef06a0bfb2a5",
		Expiry: expiry,
		User: session.User{
			ID:       u1.ID,
			Username: u1.Username,
		},
	}

	sess, err := s.CreateSession(context.Background(), u1.ID, u1.Username)
	assert.NoError(t, err)
	assert.Equal(t, &expected, sess)
}
