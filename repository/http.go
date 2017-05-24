package repository

import (
	"encoding/json"
	"net/http"

	"github.com/pressly/chi"
)

// NewUsersHandler returns a RESTful http router interacting with the Service.
func NewUsersHandler(s Service) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", listByOwner(s))

	return r
}

func listByOwner(s Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := chi.URLParam(r, "username")

		repositories, err := s.ListAggregateByOwnerUsername(username)
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

func NewHandler(s Service) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", get(s))

	return r
}

type responseWithStats struct {
	*Repository
	*Stats
}

func get(s Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		owner := chi.URLParam(r, "owner")
		name := chi.URLParam(r, "name")

		repository, stats, err := s.Find(owner, name)
		if err != nil {
			return // TODO
		}

		res := responseWithStats{
			repository,
			stats,
		}

		data, err := json.Marshal(res)
		if err != nil {
			return // TODO
		}

		w.Write(data)
	}
}
