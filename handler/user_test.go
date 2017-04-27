package handler

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gitloud/gitloud"
	"github.com/gitloud/gitloud/store"
	"github.com/go-kit/kit/log"
	"github.com/gobuffalo/packr"
	"github.com/stretchr/testify/assert"
)

var (
	box = packr.NewBox("../public")
)

func TestUserList(t *testing.T) {
	res, content := request(t, http.MethodGet, "/api/users", nil)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))
	assert.JSONEq(
		t,
		`[{"id":"25558000-2565-48dc-84eb-18754da2b0a2","username":"metalmatze","name":"Matthias Loibl","email":"metalmatze@example.com"},{"id":"911d24ae-ad9b-4e50-bf23-9dcbdc8134c6","username":"tboerger","name":"Thomas Boerger","email":"tboerger@example.com"}]`,
		string(content),
	)
}

func TestUser(t *testing.T) {
	res, content := request(t, http.MethodGet, "/api/users/metalmatze", nil)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))
	assert.JSONEq(
		t,
		`{"id":"25558000-2565-48dc-84eb-18754da2b0a2","username":"metalmatze","name":"Matthias Loibl","email":"metalmatze@example.com"}`,
		string(content),
	)
}

func TestUserCreate(t *testing.T) {
	userStore := store.NewUserInMemory()
	r := NewRouter(log.NewNopLogger(), box, userStore)

	payloadUser := gitloud.User{
		ID:       "28195928-2e77-431b-b1fc-43f543cfdc2a",
		Username: "foobar",
		Name:     "Foo Bar",
		Email:    "foobar@example.com",
	}

	payload, err := json.Marshal(payloadUser)
	assert.NoError(t, err)

	res, content := requestWithRouter(t, r, http.MethodPost, "/api/users", payload)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "", string(content))

	user, err := userStore.GetUser(payloadUser.Username)
	assert.NoError(t, err)
	assert.Equal(t, payloadUser, user)
}

func TestUserUpdate(t *testing.T) {
	userStore := store.NewUserInMemory()
	r := NewRouter(log.NewNopLogger(), box, userStore)

	user, err := userStore.GetUser("metalmatze")
	assert.NoError(t, err)

	newEmail := "matze@example.com"
	user.Email = newEmail

	payload, err := json.Marshal(user)
	assert.NoError(t, err)

	res, content := requestWithRouter(t, r, http.MethodPut, "/api/users/metalmatze", payload)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "", string(content))

	// Test if store was really updated
	user, err = userStore.GetUser("metalmatze")
	assert.NoError(t, err)
	assert.Equal(t, newEmail, user.Email)
}

func TestUserDelete(t *testing.T) {
	userStore := store.NewUserInMemory()
	r := NewRouter(log.NewNopLogger(), box, userStore)

	_, err := userStore.GetUser("metalmatze")
	assert.NoError(t, err)

	res, content := requestWithRouter(t, r, http.MethodDelete, "/api/users/metalmatze", nil)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "", string(content))

	// Test if store was really updated
	_, err = userStore.GetUser("metalmatze")
	assert.Equal(t, err, store.UserNotFound)
}
