package repository

import (
	"encoding/json"
	"net/http"

	"github.com/pressly/chi"
)

// NewHandler returns a RESTful http router interacting with the Service.
func NewHandler(s Service) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", listByOwner(s))

	return r
}

func listByOwner(s Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := chi.URLParam(r, "username")

		repositories, err := s.ListByOwnerUsername(username)
		if err != nil {
			return // TODO
		}

		data, err := json.Marshal(repositories)
		if err != nil {
			return // TODO
		}

		w.Write(data)
	}
}
