package handler

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gitloud/gitloud/store"
	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"
)

func TestApiNotFound(t *testing.T) {
	res, content := request(t, http.MethodGet, "/api/404", nil)

	assert.Equal(t, http.StatusNotFound, res.StatusCode)
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))
	assert.JSONEq(t, `{"error":"Not Found"}`, string(content))
}

func request(t *testing.T, method string, url string, payload []byte) (*http.Response, []byte) {
	userStore := store.NewUserInMemory()
	r := NewRouter(log.NewNopLogger(), box, userStore)

	return requestWithRouter(t, r, method, url, payload)
}

func requestWithRouter(t *testing.T, h http.Handler, method string, url string, payload []byte) (*http.Response, []byte) {
	ts := httptest.NewServer(h)
	defer ts.Close()

	req, err := http.NewRequest(method, ts.URL+url, bytes.NewReader(payload))
	assert.NoError(t, err)

	res, err := http.DefaultClient.Do(req)

	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	assert.NoError(t, err)

	return res, content
}
