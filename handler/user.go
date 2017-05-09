package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gitpods/gitpods"
	"github.com/gitpods/gitpods/store"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/gorilla/mux"
)

type UserStore interface {
	List() ([]gitpods.User, error)
	GetUser(string) (gitpods.User, error)
	CreateUser(gitpods.User) (gitpods.User, error)
	UpdateUser(string, gitpods.User) (gitpods.User, error)
	DeleteUser(string) error
}

type UsersAPI struct {
	logger log.Logger
	store  UserStore
}

func (a *UsersAPI) List(w http.ResponseWriter, r *http.Request) {
	users, err := a.store.List()
	if err != nil {
		http.Error(w, "failed to list users", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, users, http.StatusOK)
}

func (a *UsersAPI) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	user, err := a.store.GetUser(username)
	if err != nil {
		jsonResponseBytes(w, JsonNotFound, http.StatusNotFound)
		return
	}

	jsonResponse(w, user, http.StatusOK)
}

func (a *UsersAPI) Create(w http.ResponseWriter, r *http.Request) {
	a.logger = log.With(a.logger, "handler", "UserCreate")

	var user gitpods.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		msg := "failed to unmarshal user"
		level.Warn(a.logger).Log("msg", msg, "err", err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	if err := user.Validate(); err != nil {
		level.Debug(a.logger).Log("msg", "user invalid", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := a.store.CreateUser(user)
	if err != nil {
		msg := "failed create user"
		level.Warn(a.logger).Log("msg", msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	jsonResponse(w, user, http.StatusOK)
}

func (a *UsersAPI) Update(w http.ResponseWriter, r *http.Request) {
	logger := log.With(a.logger, "handler", "UserUpdate")

	vars := mux.Vars(r)
	username := vars["username"]

	var user gitpods.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		msg := "failed to unmarshal user"
		level.Warn(a.logger).Log("msg", msg, "err", err)
		jsonResponse(w, map[string]string{"message": msg}, http.StatusBadRequest)
		return
	}

	if err := user.Validate(); err != nil {
		level.Debug(logger).Log("msg", "user invalid", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := a.store.UpdateUser(username, user)
	if err == store.UserNotFound {
		jsonResponseBytes(w, JsonNotFound, http.StatusNotFound)
		return
	}

	jsonResponse(w, user, http.StatusOK)
}

func (a *UsersAPI) Delete(w http.ResponseWriter, r *http.Request) {
	logger := log.With(a.logger, "handler", "UserDelete")

	vars := mux.Vars(r)
	username := vars["username"]

	err := a.store.DeleteUser(username)
	if err == store.UserNotFound {
		jsonResponseBytes(w, JsonNotFound, http.StatusNotFound)
		return
	}
	if err != nil {
		msg := "failed to delete user"
		level.Warn(logger).Log("msg", msg, "err", err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
}
