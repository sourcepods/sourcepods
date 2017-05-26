package session

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"strings"

	"github.com/stretchr/testify/assert"
)

var (
	uuid   = "2e075a73-98d3-4980-b0c7-ba06fbd2cc36"
	expiry = time.Date(2009, 11, 10, 23, 00, 00, 00, time.UTC)
)

type testService struct{}

func (s *testService) CreateSession(string, string) (*Session, error) {
	panic("implement me")
}

func (s *testService) FindSession(id string) (*Session, error) {
	if id == uuid {
		return &Session{
			ID:     uuid,
			Expiry: expiry,
			User: User{
				ID:       "ab2dfdfc-0603-4752-ad7f-0e57256feaa8",
				Username: "foobar",
			},
		}, nil
	}
	return nil, errors.New("session not found")
}

func (s *testService) ClearSessions() (int64, error) {
	panic("implement me")
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusTeapot)
	user := GetSessionUser(r)
	w.Write([]byte(fmt.Sprintf(`{"id":"%s","username":"%s"}`, user.ID, user.Username)))
}

func TestAuthorized(t *testing.T) {
	s := &testService{}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Cookie", fmt.Sprintf("%s=%s", CookieName, uuid))

	w := httptest.NewRecorder()
	h := Authorized(s)(http.HandlerFunc(testHandler))

	h.ServeHTTP(w, req)

	expected := `{"id":"ab2dfdfc-0603-4752-ad7f-0e57256feaa8","username":"foobar"}`
	assert.Equal(t, http.StatusTeapot, w.Code)
	assert.Equal(t, expected, w.Body.String())
}

func TestAuthorizedNoCookie(t *testing.T) {
	s := &testService{}

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	w := httptest.NewRecorder()
	h := Authorized(s)(http.HandlerFunc(testHandler))

	h.ServeHTTP(w, req)

	expected := `{"errors":[{"title":"Unauthorized","detail":"Your Cookie is not valid","status":"401"}]}`
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, expected, strings.TrimSpace(w.Body.String()))
}

func TestAuthorizedInvalidCookie(t *testing.T) {
	s := &testService{}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Cookie", fmt.Sprintf("%s=%s", CookieName, "blablabla"))

	w := httptest.NewRecorder()
	h := Authorized(s)(http.HandlerFunc(testHandler))

	h.ServeHTTP(w, req)

	expected := `{"errors":[{"title":"Unauthorized","detail":"Your Cookie is not valid","status":"401"}]}`
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, expected, strings.TrimSpace(w.Body.String()))
}

func TestGetSessionUser(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx := context.WithValue(req.Context(), cookieUserID, "id")
	ctx = context.WithValue(ctx, cookieUserUsername, "username")
	req = req.WithContext(ctx)

	user := GetSessionUser(req)
	assert.Equal(t, "id", user.ID)
	assert.Equal(t, "username", user.Username)
}
