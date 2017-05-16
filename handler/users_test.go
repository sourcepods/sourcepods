package handler_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gitpods/gitpods"
	"github.com/gitpods/gitpods/handler"
	"github.com/gitpods/gitpods/store"
	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"
)

func newUsersAPI(usersStore handler.UsersStore) *handler.UsersAPI {
	if usersStore == nil {
		usersStore = store.NewUsersInMemory()
	}

	return &handler.UsersAPI{
		Logger: log.NewNopLogger(),
		Store:  usersStore,
	}
}

func TestUserList(t *testing.T) {
	r := newUsersAPI(nil)
	res, content, err := Request(r.Routes(), http.MethodGet, "/", nil)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("Content-Type"))
	assert.JSONEq(
		t,
		`[{"id":"25558000-2565-48dc-84eb-18754da2b0a2","username":"metalmatze","name":"Matthias Loibl","email":"metalmatze@example.com"},{"id":"911d24ae-ad9b-4e50-bf23-9dcbdc8134c6","username":"tboerger","name":"Thomas Boerger","email":"tboerger@example.com"}]`,
		string(content),
	)
}

func TestUser(t *testing.T) {
	r := newUsersAPI(nil)
	res, content, err := Request(r.Routes(), http.MethodGet, "/metalmatze", nil)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("Content-Type"))
	assert.JSONEq(
		t,
		`{"id":"25558000-2565-48dc-84eb-18754da2b0a2","username":"metalmatze","name":"Matthias Loibl","email":"metalmatze@example.com"}`,
		string(content),
	)
}

func TestUserNotFound(t *testing.T) {
	r := newUsersAPI(nil)
	res, content, err := Request(r.Routes(), http.MethodGet, "/foobar", nil)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("Content-Type"))
	assert.JSONEq(t, string(handler.JsonNotFound), string(content))
}

func TestUserCreate(t *testing.T) {
	usersStore := store.NewUsersInMemory()
	r := newUsersAPI(usersStore)

	payloadUser := gitpods.User{
		ID:       "28195928-2e77-431b-b1fc-43f543cfdc2a",
		Username: "foobar",
		Name:     "Foo Bar",
		Email:    "foobar@example.com",
	}

	payload, err := json.Marshal(payloadUser)
	assert.NoError(t, err)

	res, content, err := Request(r.Routes(), http.MethodPost, "/", payload)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, string(payload), string(content))

	user, err := usersStore.GetUser(payloadUser.Username)
	assert.NoError(t, err)
	assert.Equal(t, payloadUser, *user)
}

func TestUserUpdate(t *testing.T) {
	usersStore := store.NewUsersInMemory()
	r := newUsersAPI(usersStore)

	user, err := usersStore.GetUser("metalmatze")
	assert.NoError(t, err)

	newEmail := "matze@example.com"
	user.Email = newEmail

	payload, err := json.Marshal(user)
	assert.NoError(t, err)

	res, content, err := Request(r.Routes(), http.MethodPut, "/metalmatze", payload)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, string(payload), string(content))

	// Test if store was really updated
	user, err = usersStore.GetUser("metalmatze")
	assert.NoError(t, err)
	assert.Equal(t, newEmail, user.Email)
}

func TestUserUpdateBadRequest(t *testing.T) {
	r := newUsersAPI(nil)
	res, content, err := Request(r.Routes(), http.MethodPut, "/metalmatze", nil)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("Content-Type"))
	assert.JSONEq(t, `{"message":"failed to unmarshal user"}`, string(content))
}

func TestUserUpdateNotFound(t *testing.T) {
	usersStore := store.NewUsersInMemory()
	r := newUsersAPI(usersStore)

	user, err := usersStore.GetUser("metalmatze")
	assert.NoError(t, err)

	payload, err := json.Marshal(user)
	assert.NoError(t, err)

	res, content, err := Request(r.Routes(), http.MethodPut, "/foobar", payload)

	assert.NoError(t, err)
	assertNotFoundJson(t, res, content)
}

func TestUserDelete(t *testing.T) {
	usersStore := store.NewUsersInMemory()
	r := newUsersAPI(usersStore)

	_, err := usersStore.GetUser("metalmatze")
	assert.NoError(t, err)

	res, content, err := Request(r.Routes(), http.MethodDelete, "/metalmatze", nil)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "", string(content))

	// Test if store was really updated
	_, err = usersStore.GetUser("metalmatze")
	assert.Equal(t, err, store.UserNotFound)
}

func TestUserDeleteNotFound(t *testing.T) {
	r := DefaultTestAuthRouter()

	res, content, err := Request(r, http.MethodDelete, "/users/foobar", nil)
	assert.NoError(t, err)
	assertNotFoundJson(t, res, content)
}
