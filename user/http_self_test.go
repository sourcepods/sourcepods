package user

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gitpods/gitpods/session"
	"github.com/stretchr/testify/assert"
)

func TestHTTPSelf(t *testing.T) {
	s := &testService{}
	h := NewUserHandler(s)

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(t, err)
	ctx := context.WithValue(req.Context(), session.CookieUserID, "bb5e0c5f-73d9-4c9a-8c0d-8110e720e1b2")
	ctx = context.WithValue(ctx, session.CookieUserUsername, "username1")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	expected := `{"data":{"type":"users","id":"bb5e0c5f-73d9-4c9a-8c0d-8110e720e1b2","attributes":{"created_at":1257894000,"email":"email1@example.com","name":"Name 1","updated_at":1257895800,"username":"username1"}}}`

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expected, strings.TrimSpace(w.Body.String()))
}
