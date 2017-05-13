package handler

import (
	"net/http"

	"github.com/gitpods/gitpods"
	"github.com/gitpods/gitpods/store"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type UserStore interface {
	GetUser(string) (*gitpods.User, error)
	GetUserRepositories(string) (*gitpods.User, []*gitpods.Repository, error)
}

func User(logger log.Logger, s UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.Context().Value(SessionUserUsername)
		if username == nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		user, err := s.GetUser(username.(string))
		if err == store.UserNotFound {
			jsonResponseBytes(w, JsonNotFound, http.StatusNotFound)
			return
		}
		if err != nil {
			msg := "failed to get user"
			level.Warn(logger).Log("msg", msg, "err", err)
			jsonResponse(w, map[string]string{"message": msg}, http.StatusInternalServerError)
			return
		}

		jsonResponse(w, user, http.StatusOK)
	}
}

func UserRepositories(logger log.Logger, store UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.Context().Value(SessionUserUsername)
		if username == nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		user, repositories, err := store.GetUserRepositories(username.(string))
		if err != nil {
			jsonResponseBytes(w, JsonNotFound, http.StatusNotFound)
			return
		}

		var resp struct {
			User         *gitpods.User         `json:"user"`
			Repositories []*gitpods.Repository `json:"repositories"`
		}

		resp.User = user
		resp.Repositories = repositories

		jsonResponse(w, resp, http.StatusOK)
	}
}
