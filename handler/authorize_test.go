package handler_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthorize(t *testing.T) {
	r := DefaultTestRouter()

	payload := `{"email":"metalmatze@example.com", "password":"kubernetes"}`

	res, content, err := Request(r, http.MethodPost, "/authorize", []byte(payload))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("Content-Type"))
	assert.Contains(t, res.Header.Get("Set-Cookie"), "_gitpods_session=")
	assert.Equal(t,
		`{"id":"25558000-2565-48dc-84eb-18754da2b0a2","email":"metalmatze@example.com","username":"metalmatze","name":"Matthias Loibl"}`,
		string(content),
	)
}

func TestAuthorizeBadRequest(t *testing.T) {
	r := DefaultTestRouter()

	res, content, err := Request(r, http.MethodPost, "/authorize", nil)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("Content-Type"))
	assert.Equal(t, `{"message":"failed to unmarshal form"}`, string(content))
}

func TestAuthorizeWrongEmail(t *testing.T) {
	r := DefaultTestRouter()

	payload := `{"email":"foobar@example.com", "password":"kubernetes"}`

	res, content, err := Request(r, http.MethodPost, "/authorize", []byte(payload))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("Content-Type"))
	assert.Equal(t, `{"message":"Bad credentials"}`, string(content))
}

func TestAuthorizeWrongPassword(t *testing.T) {
	r := DefaultTestRouter()

	payload := `{"email":"metalmatze@example.com", "password":"foobar"}`

	res, content, err := Request(r, http.MethodPost, "/authorize", []byte(payload))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("Content-Type"))
	assert.Equal(t, `{"message":"Bad credentials"}`, string(content))
}
