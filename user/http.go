package user

import (
	"encoding/json"
	"net/http"

	"github.com/pressly/chi"
)

func NewHandler(s Service) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", list(s))
	r.Get("/:username", get(s))

	return r
}

func list(s Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := s.FindAll()
		if err != nil {
			return // TODO
		}

		data, err := json.Marshal(users)
		if err != nil {
			return // TODO
		}

		w.Write(data)
	}
}

func get(s Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := chi.URLParam(r, "username")

		user, err := s.FindByUsername(Username(username))
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
