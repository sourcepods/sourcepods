package handler

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gitloud/gitloud"
	"github.com/gitloud/gitloud/store"
	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"
)

func TestUserList(t *testing.T) {
	userStore := store.NewUserInMemory()
	r := NewRouter(log.NewNopLogger(), userStore)

	ts := httptest.NewServer(r)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/api/users")
	assert.NoError(t, err)

	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))
	assert.JSONEq(
		t,
		`[{"id":"abcd-efgh-1234-5678","username":"metalmatze","name":"Matthias Loibl","email":"mail@matthiasloibl.com"},{"id":"bcde-fghi-2345-6789","username":"tboerger","name":"Thomas Boerger","email":"thomas@webhippie.de"}]`,
		string(content),
	)
}

func TestUser(t *testing.T) {
	userStore := store.NewUserInMemory()
	r := NewRouter(log.NewNopLogger(), userStore)

	ts := httptest.NewServer(r)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/api/users/metalmatze")
	assert.NoError(t, err)

	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))
	assert.JSONEq(
		t,
		`{"id":"abcd-efgh-1234-5678","username":"metalmatze","name":"Matthias Loibl","email":"mail@matthiasloibl.com"}`,
		string(content),
	)
}

func TestUserCreate(t *testing.T) {
	userStore := store.NewUserInMemory()
	r := NewRouter(log.NewNopLogger(), userStore)

	ts := httptest.NewServer(r)
	defer ts.Close()

	payloadUser := gitloud.User{
		ID:       "cdef-ghij-3456-7890",
		Username: "foobar",
		Name:     "Foo Bar",
		Email:    "foobar@example.com",
	}

	payload, err := json.Marshal(payloadUser)
	assert.NoError(t, err)

	res, err := http.Post(ts.URL+"/api/users", "application/json", bytes.NewReader(payload))
	assert.NoError(t, err)

	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "", string(content))

	user, err := userStore.GetUser(payloadUser.Username)
	assert.NoError(t, err)
	assert.Equal(t, payloadUser, user)
}

func TestUserUpdate(t *testing.T) {
	userStore := store.NewUserInMemory()
	r := NewRouter(log.NewNopLogger(), userStore)

	ts := httptest.NewServer(r)
	defer ts.Close()

	user, err := userStore.GetUser("metalmatze")
	assert.NoError(t, err)

	newEmail := "me@metalmatze.de"
	user.Email = newEmail

	payload, err := json.Marshal(user)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPut, ts.URL+"/api/users/metalmatze", bytes.NewReader(payload))
	assert.NoError(t, err)

	res, err := http.DefaultClient.Do(req)

	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "", string(content))

	// Test if store was really updated
	user, err = userStore.GetUser("metalmatze")
	assert.NoError(t, err)
	assert.Equal(t, newEmail, user.Email)
}

func TestUserDelete(t *testing.T) {
	userStore := store.NewUserInMemory()
	r := NewRouter(log.NewNopLogger(), userStore)

	ts := httptest.NewServer(r)
	defer ts.Close()

	_, err := userStore.GetUser("metalmatze")
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodDelete, ts.URL+"/api/users/metalmatze", nil)
	assert.NoError(t, err)

	res, err := http.DefaultClient.Do(req)

	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "", string(content))

	// Test if store was really updated
	_, err = userStore.GetUser("metalmatze")
	assert.Equal(t, err, store.UserNotFound)
}
