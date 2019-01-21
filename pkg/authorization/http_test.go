package authorization

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sourcepods/sourcepods/pkg/gitpods/user"
	"github.com/sourcepods/sourcepods/pkg/session"
	"github.com/stretchr/testify/assert"
)

type testService struct{}

func (s *testService) AuthenticateUser(ctx context.Context, email, password string) (*user.User, error) {
	if email == "foobar@example.com" && password == "baz" {
		return &u1, nil
	}
	return nil, errors.New("bad credentials")
}

func (s *testService) CreateSession(ctx context.Context, id string, username string) (*session.Session, error) {
	return &session.Session{
		ID:     "410f59a5-75e6-4332-a0d3-ef06a0bfb2a5",
		Expiry: expiry,
		User: session.User{
			ID:       id,
			Username: username,
		},
	}, nil
}

func TestHTTPAuthorize(t *testing.T) {
	s := &testService{}
	h := NewHandler(s)

	payload := strings.NewReader(`{"email": "foobar@example.com","password": "baz"}`)
	req, err := http.NewRequest(http.MethodPost, "/", payload)
	assert.NoError(t, err)

	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	cookie := "_sourcepods_session=410f59a5-75e6-4332-a0d3-ef06a0bfb2a5; Path=/; Expires=Tue, 10 Nov 2009 23:00:00 GMT"

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, cookie, w.Header().Get("Set-Cookie"))
	assert.Equal(t, "", w.Body.String())
}

func TestHTTPAuthorizeBadCredentials(t *testing.T) {
	s := &testService{}
	h := NewHandler(s)

	payload := strings.NewReader(`{"email": "foobar@example.com","password": "bla"}`)
	req, err := http.NewRequest(http.MethodPost, "/", payload)
	assert.NoError(t, err)

	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	badCredentials := `{"errors":[{"title":"Bad Request","detail":"Incorrect email or password","status":"400"}]}`

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "", w.Header().Get("Set-Cookie"))
	assert.Equal(t, badCredentials, strings.TrimSpace(w.Body.String()))
}
