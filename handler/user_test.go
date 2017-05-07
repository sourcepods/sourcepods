package handler_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gitpods/gitpods"
	"github.com/gitpods/gitpods/handler"
	"github.com/gitpods/gitpods/store"
	"github.com/stretchr/testify/assert"
)

func TestUserList(t *testing.T) {
	routerStore := DefaultRouterStore()
	r := DefaultTestAuthRouterWithStore(routerStore)

	res, content, err := Request(r, http.MethodGet, "/api/users", nil)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("Content-Type"))
	assert.JSONEq(
		t,
		`[{"id":"25558000-2565-48dc-84eb-18754da2b0a2","username":"metalmatze","name":"Matthias Loibl","email":"metalmatze@example.com"},{"id":"911d24ae-ad9b-4e50-bf23-9dcbdc8134c6","username":"tboerger","name":"Thomas Boerger","email":"tboerger@example.com"}]`,
		string(content),
	)
}

func TestUserListUnauthorized(t *testing.T) {
	r := DefaultTestRouter()
	res, content, err := Request(r, http.MethodGet, "/api/users", nil)

	assert.NoError(t, err)
	assertUnauthorized(t, res, content)
}

func TestUser(t *testing.T) {
	r := DefaultTestAuthRouter()
	res, content, err := Request(r, http.MethodGet, "/api/users/metalmatze", nil)

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
	r := DefaultTestAuthRouter()
	res, content, err := Request(r, http.MethodGet, "/api/users/foobar", nil)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("Content-Type"))
	assert.JSONEq(t, string(handler.JsonNotFound), string(content))
}

func TestUserCreate(t *testing.T) {
	routerStore := DefaultRouterStore()
	r := DefaultTestAuthRouterWithStore(routerStore)

	payloadUser := gitpod.User{
		ID:       "28195928-2e77-431b-b1fc-43f543cfdc2a",
		Username: "foobar",
		Name:     "Foo Bar",
		Email:    "foobar@example.com",
	}

	payload, err := json.Marshal(payloadUser)
	assert.NoError(t, err)

	res, content, err := Request(r, http.MethodPost, "/api/users", payload)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, string(payload), string(content))

	user, err := routerStore.UserStore.GetUser(payloadUser.Username)
	assert.NoError(t, err)
	assert.Equal(t, payloadUser, user)
}

func TestUserUpdate(t *testing.T) {
	routerStore := DefaultRouterStore()
	r := DefaultTestAuthRouterWithStore(routerStore)

	user, err := routerStore.UserStore.GetUser("metalmatze")
	assert.NoError(t, err)

	newEmail := "matze@example.com"
	user.Email = newEmail

	payload, err := json.Marshal(user)
	assert.NoError(t, err)

	res, content, err := Request(r, http.MethodPut, "/api/users/metalmatze", payload)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, string(payload), string(content))

	// Test if store was really updated
	user, err = routerStore.UserStore.GetUser("metalmatze")
	assert.NoError(t, err)
	assert.Equal(t, newEmail, user.Email)
}

func TestUserUpdateBadRequest(t *testing.T) {
	r := DefaultTestAuthRouter()
	res, content, err := Request(r, http.MethodPut, "/api/users/metalmatze", nil)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("Content-Type"))
	assert.JSONEq(t, `{"message":"failed to unmarshal user"}`, string(content))
}

func TestUserUpdateNotFound(t *testing.T) {
	routerStore := DefaultRouterStore()
	r := DefaultTestAuthRouterWithStore(routerStore)

	user, err := routerStore.UserStore.GetUser("metalmatze")
	assert.NoError(t, err)

	payload, err := json.Marshal(user)
	assert.NoError(t, err)

	res, content, err := Request(r, http.MethodPut, "/api/users/foobar", payload)

	assert.NoError(t, err)
	assertNotFoundJson(t, res, content)
}

func TestUserDelete(t *testing.T) {
	routerStore := DefaultRouterStore()
	r := DefaultTestAuthRouterWithStore(routerStore)

	_, err := routerStore.UserStore.GetUser("metalmatze")
	assert.NoError(t, err)

	res, content, err := Request(r, http.MethodDelete, "/api/users/metalmatze", nil)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "", string(content))

	// Test if store was really updated
	_, err = routerStore.UserStore.GetUser("metalmatze")
	assert.Equal(t, err, store.UserNotFound)
}

func TestUserDeleteNotFound(t *testing.T) {
	r := DefaultTestAuthRouter()

	res, content, err := Request(r, http.MethodDelete, "/api/users/foobar", nil)
	assert.NoError(t, err)
	assertNotFoundJson(t, res, content)
}
