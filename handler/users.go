package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gitpods/gitpods"
	"github.com/gitpods/gitpods/store"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pressly/chi"
)

type UsersStore interface {
	List() ([]*gitpods.User, error)
	GetUserByUsername(string) (*gitpods.User, error)
	CreateUser(*gitpods.User) (*gitpods.User, error)
	UpdateUser(string, *gitpods.User) (*gitpods.User, error)
	DeleteUser(string) error
}

type UsersAPI struct {
	Logger log.Logger
	Store  UsersStore
}

func (a *UsersAPI) Routes() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", a.List)
	r.Post("/", a.Create)
	r.Get("/:username", a.Get)
	r.Put("/:username", a.Update)
	r.Delete("/:username", a.Delete)

	return r
}

func (a *UsersAPI) List(w http.ResponseWriter, r *http.Request) {
	users, err := a.Store.List()
	if err != nil {
		http.Error(w, "failed to list users", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, users, http.StatusOK)
}

func (a *UsersAPI) Get(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")

	user, err := a.Store.GetUserByUsername(username)
	if err != nil {
		jsonResponseBytes(w, JsonNotFound, http.StatusNotFound)
		return
	}

	jsonResponse(w, user, http.StatusOK)
}

func (a *UsersAPI) Create(w http.ResponseWriter, r *http.Request) {
	a.Logger = log.With(a.Logger, "handler", "UserCreate")

	var user *gitpods.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		msg := "failed to unmarshal user"
		level.Warn(a.Logger).Log("msg", msg, "err", err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	if err := user.Validate(); err != nil {
		level.Debug(a.Logger).Log("msg", "user invalid", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := a.Store.CreateUser(user)
	if err != nil {
		msg := "failed create user"
		level.Warn(a.Logger).Log("msg", msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	jsonResponse(w, user, http.StatusOK)
}

func (a *UsersAPI) Update(w http.ResponseWriter, r *http.Request) {
	logger := log.With(a.Logger, "handler", "UserUpdate")
	username := chi.URLParam(r, "username")

	var user *gitpods.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		msg := "failed to unmarshal user"
		level.Warn(a.Logger).Log("msg", msg, "err", err)
		jsonResponse(w, map[string]string{"message": msg}, http.StatusBadRequest)
		return
	}

	if err := user.Validate(); err != nil {
		level.Debug(logger).Log("msg", "user invalid", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := a.Store.UpdateUser(username, user)
	if err == store.UserNotFound {
		jsonResponseBytes(w, JsonNotFound, http.StatusNotFound)
		return
	}

	jsonResponse(w, user, http.StatusOK)
}

func (a *UsersAPI) Delete(w http.ResponseWriter, r *http.Request) {
	logger := log.With(a.Logger, "handler", "UserDelete")
	username := chi.URLParam(r, "username")

	err := a.Store.DeleteUser(username)
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
