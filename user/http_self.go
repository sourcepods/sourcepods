package user

import (
	"net/http"

	"github.com/gitpods/gitpods/session"
	"github.com/go-chi/chi"
	"github.com/google/jsonapi"
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
		sessionUser := session.GetSessionUser(r.Context())

		user, err := s.FindByUsername(r.Context(), sessionUser.Username)
		if err != nil {
			return // TODO
		}

		res := &response{
			ID:       user.ID,
			Email:    user.Email,
			Username: user.Username,
			Name:     user.Name,
			Created:  user.Created,
			Updated:  user.Updated,
		}

		if err := jsonapi.MarshalPayload(w, res); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
