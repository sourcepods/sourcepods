package authorization

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gitpods/gitpods/session"
	"github.com/google/jsonapi"
	"github.com/pressly/chi"
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
		var form struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		badCredentials := []*jsonapi.ErrorObject{{
			Title:  http.StatusText(http.StatusBadRequest),
			Detail: "Bad Credentials",
			Status: fmt.Sprintf("%d", http.StatusBadRequest),
		}}

		if err := json.NewDecoder(io.LimitReader(r.Body, megabyte)).Decode(&form); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			jsonapi.MarshalErrors(w, badCredentials)
			return
		}

		user, err := s.AuthenticateUser(form.Email, form.Password)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			jsonapi.MarshalErrors(w, badCredentials)
			return
		}

		sess, err := s.CreateSession(user.ID, user.Username)
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
