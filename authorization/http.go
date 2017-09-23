package authorization

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gitpods/gitpods/session"
	"github.com/go-chi/chi"
	"github.com/google/jsonapi"
	opentracing "github.com/opentracing/opentracing-go"
)

const megabyte = 1024 * 1024 * 1024

// NewHandler returns a RESTful http router interacting with the Service.
func NewHandler(s Service) *chi.Mux {
	r := chi.NewRouter()

	r.Post("/", authorize(s))

	return r
}

func authorize(s Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		span, ctx := opentracing.StartSpanFromContext(r.Context(), "authorization.Handler.authorize")
		defer span.Finish()

		var form struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		badCredentials := []*jsonapi.ErrorObject{{
			Title:  http.StatusText(http.StatusBadRequest),
			Detail: "Incorrect email or password",
			Status: fmt.Sprintf("%d", http.StatusBadRequest),
		}}

		if err := json.NewDecoder(io.LimitReader(r.Body, megabyte)).Decode(&form); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			jsonapi.MarshalErrors(w, badCredentials)
			return
		}

		span.SetTag("email", form.Email)

		user, err := s.AuthenticateUser(ctx, form.Email, form.Password)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			jsonapi.MarshalErrors(w, badCredentials)
			return
		}

		sess, err := s.CreateSession(ctx, user.ID, user.Username)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			jsonapi.MarshalErrors(w, badCredentials)
			return
		}

		cookie := &http.Cookie{
			Name:    session.CookieName,
			Value:   sess.ID,
			Path:    "/",
			Expires: sess.Expiry,
		}

		http.SetCookie(w, cookie)
	}
}
