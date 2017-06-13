package user

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	created = time.Date(2009, 11, 10, 23, 00, 00, 00, time.UTC)
	updated = time.Date(2009, 11, 10, 23, 30, 00, 00, time.UTC)

	u1 = User{
		ID:       "bb5e0c5f-73d9-4c9a-8c0d-8110e720e1b2",
		Email:    "email1@example.com",
		Username: "username1",
		Name:     "Name 1",
		Password: "password 1",
		Created:  created,
		Updated:  updated,
	}
	u2 = User{
		ID:       "e845eb7a-60c0-42c1-af53-5c328779efb8",
		Email:    "email2@example.com",
		Username: "username2",
		Name:     "Name 2",
		Password: "password 2",
		Created:  created,
		Updated:  updated,
	}
)

type testService struct{}

func (s *testService) FindAll() ([]*User, error) {
	return []*User{&u1, &u2}, nil
}

func (s *testService) Find(string) (*User, error) {
	panic("implement me")
}

func (s *testService) FindByUsername(username string) (*User, error) {
	switch username {
	case u1.Username:
		return &u1, nil
	case u2.Username:
		return &u2, nil
	default:
		return nil, errors.New("user not found")
	}
}

func (s *testService) Create(*User) (*User, error) {
	panic("implement me")
}

func (s *testService) Update(*User) (*User, error) {
	panic("implement me")
}

func (s *testService) Delete(string) error {
	panic("implement me")
}

func TestHTTPList(t *testing.T) {
	s := &testService{}

	h := NewUsersHandler(s)
	ts := httptest.NewServer(h)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/")
	assert.NoError(t, err)

	body, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)
	defer res.Body.Close()

	expected := `{"data":[{"type":"users","id":"bb5e0c5f-73d9-4c9a-8c0d-8110e720e1b2","attributes":{"created_at":1257894000,"email":"email1@example.com","name":"Name 1","updated_at":1257895800,"username":"username1"}},{"type":"users","id":"e845eb7a-60c0-42c1-af53-5c328779efb8","attributes":{"created_at":1257894000,"email":"email2@example.com","name":"Name 2","updated_at":1257895800,"username":"username2"}}]}`

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, expected, strings.TrimSpace(string(body)))
}

func TestHTTPUser(t *testing.T) {
	s := &testService{}

	h := NewUsersHandler(s)
	ts := httptest.NewServer(h)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/username1")
	assert.NoError(t, err)

	body, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)
	defer res.Body.Close()

	expected := `{"data":{"type":"users","id":"bb5e0c5f-73d9-4c9a-8c0d-8110e720e1b2","attributes":{"created_at":1257894000,"email":"email1@example.com","name":"Name 1","updated_at":1257895800,"username":"username1"}}}`

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, expected, strings.TrimSpace(string(body)))
}

func TestHTTPUserNotFound(t *testing.T) {
	s := &testService{}

	h := NewUsersHandler(s)
	ts := httptest.NewServer(h)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/foobar")
	assert.NoError(t, err)

	body, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)
	defer res.Body.Close()

	expected := `{"errors":[{"title":"Not Found","detail":"Can't find user with this username","status":"404"}]}`

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, expected, strings.TrimSpace(string(body)))
}
