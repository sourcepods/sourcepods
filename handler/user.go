package handler

import (
	"net/http"

	"github.com/gitpods/gitpods/store"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

func User(logger log.Logger, s UsersStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.Context().Value(SessionUserUsername)
		if username == nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		user, err := s.GetUserByUsername(username.(string))
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

func UserRepositories(logger log.Logger, store UsersRepositoriesStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.Context().Value(SessionUserUsername)
		if username == nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		repositories, err := store.List(username.(string))
		if err != nil {
			jsonResponseBytes(w, JsonNotFound, http.StatusNotFound)
			return
		}

		jsonResponse(w, repositories, http.StatusOK)
	}
}
