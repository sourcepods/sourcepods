package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gitpods/gitpod"
	"github.com/gitpods/gitpod/store"
	"github.com/gorilla/mux"
)

type UserStore interface {
	List() ([]gitloud.User, error)
	GetUser(string) (gitloud.User, error)
	CreateUser(gitloud.User) (gitloud.User, error)
	UpdateUser(string, gitloud.User) (gitloud.User, error)
	DeleteUser(string) error
}

func WriteJson(w http.ResponseWriter, v interface{}, code int) {
	data, err := json.Marshal(v)
	if err != nil {
		http.Error(w, "failed to marshal to json", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func UserList(store UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := store.List()
		if err != nil {
			http.Error(w, "failed to list users", http.StatusInternalServerError)
			return
		}

		WriteJson(w, users, http.StatusOK)
	}
}

func User(store UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		username := vars["username"]

		user, err := store.GetUser(username)
		if err != nil {
			WriteJson(w, NotFoundJson, http.StatusNotFound)
			return
		}

		WriteJson(w, user, http.StatusOK)
	}
}

func UserCreate(store UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user gitloud.User

		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			log.Println(err)
			http.Error(w, "failed to unmarshal user", http.StatusBadRequest)
			return
		}

		if err := user.Validate(); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user, err := store.CreateUser(user)
		if err != nil {
			log.Println(err)
			http.Error(w, "failed create user", http.StatusInternalServerError)
			return
		}

		WriteJson(w, user, http.StatusOK)
	}
}

func UserUpdate(s UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		username := vars["username"]

		var user gitloud.User

		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			log.Println(err)
			WriteJson(w, map[string]string{"message": "failed to unmarshal user"}, http.StatusBadRequest)
			return
		}

		if err := user.Validate(); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user, err := s.UpdateUser(username, user)
		if err == store.UserNotFound {
			WriteJson(w, NotFoundJson, http.StatusNotFound)
			return
		}

		WriteJson(w, user, http.StatusOK)
	}
}

func UserDelete(s UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		username := vars["username"]

		err := s.DeleteUser(username)
		if err == store.UserNotFound {
			WriteJson(w, NotFoundJson, http.StatusNotFound)
			return
		}
		if err != nil {
			log.Println(err)
			http.Error(w, "failed to delete user", http.StatusBadRequest)
			return
		}
	}
}
