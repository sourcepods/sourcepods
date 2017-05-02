package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gitpods/gitpod"
	"github.com/gitpods/gitpod/store"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/gorilla/mux"
)

type UserStore interface {
	List() ([]gitpod.User, error)
	GetUser(string) (gitpod.User, error)
	CreateUser(gitpod.User) (gitpod.User, error)
	UpdateUser(string, gitpod.User) (gitpod.User, error)
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

func UserList(logger log.Logger, store UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := store.List()
		if err != nil {
			http.Error(w, "failed to list users", http.StatusInternalServerError)
			return
		}

		WriteJson(w, users, http.StatusOK)
	}
}

func User(logger log.Logger, store UserStore) http.HandlerFunc {
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

func UserCreate(logger log.Logger, store UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger = log.With(logger, "handler", "UserCreate")

		var user gitpod.User
		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			msg := "failed to unmarshal user"
			level.Warn(logger).Log("msg", msg, "err", err)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		if err := user.Validate(); err != nil {
			level.Debug(logger).Log("msg", "user invalid", "err", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user, err := store.CreateUser(user)
		if err != nil {
			msg := "failed create user"
			level.Warn(logger).Log("msg", msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		WriteJson(w, user, http.StatusOK)
	}
}

func UserUpdate(logger log.Logger, s UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger = log.With(logger, "handler", "UserUpdate")

		vars := mux.Vars(r)
		username := vars["username"]

		var user gitpod.User

		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			msg := "failed to unmarshal user"
			level.Warn(logger).Log("msg", msg, "err", err)
			WriteJson(w, map[string]string{"message": msg}, http.StatusBadRequest)
			return
		}

		if err := user.Validate(); err != nil {
			level.Debug(logger).Log("msg", "user invalid", "err", err)
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

func UserDelete(logger log.Logger, s UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger = log.With(logger, "handler", "UserDelete")

		vars := mux.Vars(r)
		username := vars["username"]

		err := s.DeleteUser(username)
		if err == store.UserNotFound {
			WriteJson(w, NotFoundJson, http.StatusNotFound)
			return
		}
		if err != nil {
			msg := "failed to delete user"
			level.Warn(logger).Log("msg", msg, "err", err)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}
	}
}
