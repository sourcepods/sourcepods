package user

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/jsonapi"
)

type response struct {
	ID       string    `jsonapi:"primary,users"`
	Email    string    `jsonapi:"attr,email"`
	Username string    `jsonapi:"attr,username"`
	Name     string    `jsonapi:"attr,name"`
	Created  time.Time `jsonapi:"attr,created_at"`
	Updated  time.Time `jsonapi:"attr,updated_at"`
}

// NewUsersHandler returns a RESTful http router interacting with the Service.
func NewUsersHandler(s Service) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", list(s))
	r.Get("/{username}", get(s))
	r.Put("/{username}", update(s))

	return r
}

func list(s Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := s.FindAll()
		if err != nil {
			jsonapi.MarshalErrors(w, []*jsonapi.ErrorObject{{
				Title:  http.StatusText(http.StatusInternalServerError),
				Status: fmt.Sprintf("%d", http.StatusInternalServerError),
			}})
			return
		}

		res := make([]interface{}, len(users))
		for i, user := range users {
			res[i] = &response{
				ID:       user.ID,
				Email:    user.Email,
				Username: user.Username,
				Name:     user.Name,
				Created:  user.Created,
				Updated:  user.Updated,
			}
		}

		if err := jsonapi.MarshalManyPayload(w, res); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func get(s Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := chi.URLParam(r, "username")

		user, err := s.FindByUsername(username)
		if err != nil {
			jsonapi.MarshalErrors(w, []*jsonapi.ErrorObject{{
				Title:  http.StatusText(http.StatusNotFound),
				Detail: "Can't find user with this username",
				Status: fmt.Sprintf("%d", http.StatusNotFound),
			}})
			return
		}

		res := &response{
			ID:       user.ID,
			Email:    user.Email,
			Username: user.Username,
			Name:     user.Name,
			Created:  user.Created,
			Updated:  user.Updated,
		}

		if err := jsonapi.MarshalOnePayload(w, res); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func update(s Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//username := chi.URLParam(r, "username")

		var user *User
		if err := json.NewDecoder(io.LimitReader(r.Body, 5242880)).Decode(&user); err != nil {
			return // TODO
		}

		user, err := s.Update(user)
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
