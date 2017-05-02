package handler_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gitpods/gitpod/handler"
	"github.com/gitpods/gitpod/store"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics/discard"
	"github.com/gobuffalo/packr"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

const NotFoundJsonString = `{"error":"Not Found"}`

var (
	box            = packr.NewBox("../public")
	httpTestClient = &http.Client{Timeout: 5 * time.Second}
)

func TestApiNotFound(t *testing.T) {
	r := DefaultTestRouter()
	res, content, err := Request(r, http.MethodGet, "/api/404", nil)
	assert.NoError(t, err)
	assertNotFoundJson(t, res, content)
}

// Helpers

func DiscardMetrics() handler.RouterMetrics {
	return handler.RouterMetrics{
		LoginAttempts: discard.NewCounter(),
	}
}

func DefaultTestRouter() *mux.Router {
	userStore := store.NewUserInMemory()
	return handler.NewRouter(log.NewNopLogger(), DiscardMetrics(), box, userStore)
}

func Request(r *mux.Router, method string, url string, payload []byte) (*http.Response, []byte, error) {
	ts := httptest.NewServer(r)
	defer ts.Close()

	req, err := http.NewRequest(method, ts.URL+url, bytes.NewReader(payload))
	if err != nil {
		return nil, nil, err
	}

	res, err := httpTestClient.Do(req)

	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	return res, content, nil
}

func assertNotFoundJson(t *testing.T, res *http.Response, content []byte) {
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))
	assert.JSONEq(t, NotFoundJsonString, string(content))
}
