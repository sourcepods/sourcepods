package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gitloud/gitloud"
	"github.com/gorilla/mux"
)

type UserStore interface {
	List() ([]gitloud.User, error)
	GetUser(string) (gitloud.User, error)
	CreateUser(gitloud.User) error
	UpdateUser(string, gitloud.User) error
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
			http.Error(w, "failed to get user by username", http.StatusInternalServerError)
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

		if err := store.CreateUser(user); err != nil {
			log.Println(err)
			http.Error(w, "failed create user", http.StatusInternalServerError)
			return
		}
	}
}

func UserUpdate(store UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		username := vars["username"]

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

		if err := store.UpdateUser(username, user); err != nil {
			log.Println(err)
			http.Error(w, "failed create user", http.StatusInternalServerError)
			return
		}
	}
}

func UserDelete(store UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		username := vars["username"]

		if err := store.DeleteUser(username); err != nil {
			log.Println(err)
			http.Error(w, "failed to delete user", http.StatusBadRequest)
			return
		}
	}
}
