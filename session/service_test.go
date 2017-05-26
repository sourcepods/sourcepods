package session

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testStore struct{}

func (s *testStore) SaveSession(sess *Session) error {
	sess.ID = "6ae1485d-13e8-4535-ba93-1d497f1b809c"
	return nil
}

func (s *testStore) FindSession(id string) (*Session, error) {
	if id == "6ae1485d-13e8-4535-ba93-1d497f1b809c" {
		return &Session{
			ID:     "6ae1485d-13e8-4535-ba93-1d497f1b809c",
			Expiry: time.Now().Add(defaultExpiry),
			User: User{
				ID:       "9749ca6a-82b2-41b5-882b-e89df9e56a2e",
				Username: "foobar",
			},
		}, nil
	}
	return nil, errors.New("session not found")
}

func (s *testStore) ClearSessions() (int64, error) {
	panic("implement me")
}

func TestService_CreateSession(t *testing.T) {
	store := &testStore{}
	s := NewService(store)

	sess, err := s.CreateSession("9749ca6a-82b2-41b5-882b-e89df9e56a2e", "foobar")
	assert.NoError(t, err)
	assert.Len(t, sess.ID, 36)
	assert.WithinDuration(t, time.Now().Add(defaultExpiry), sess.Expiry, time.Second)
	assert.Equal(t, "9749ca6a-82b2-41b5-882b-e89df9e56a2e", sess.User.ID)
	assert.Equal(t, "foobar", sess.User.Username)
}

func TestService_FindSession(t *testing.T) {
	store := &testStore{}
	s := NewService(store)

	sess, err := s.FindSession("nope")
	assert.Error(t, err)
	assert.Nil(t, sess)

	sess, err = s.FindSession("6ae1485d-13e8-4535-ba93-1d497f1b809c")
	assert.NoError(t, err)
	assert.Len(t, sess.ID, 36)
	assert.WithinDuration(t, time.Now().Add(defaultExpiry), sess.Expiry, time.Second)
	assert.Equal(t, "9749ca6a-82b2-41b5-882b-e89df9e56a2e", sess.User.ID)
	assert.Equal(t, "foobar", sess.User.Username)
}
