package user

import (
	"encoding/json"
	"net/http"

	"github.com/gitpods/gitpods/session"
	"github.com/pressly/chi"
)

// NewUserHandler returns a RESTful http router interacting with the Service
// and the authenticated user set as the username.
func NewUserHandler(s Service) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", self(s))

	return r
}

func self(s Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionUser := session.GetSessionUser(r)

		user, err := s.FindByUsername(sessionUser.Username)
		if err != nil {
			return // TODO
		}

		data, err := json.Marshal(user)
		if err != nil {
			return // TODO
		}

		w.Write(data)
	}
}
