package handler

import (
	"net/http"

	"github.com/gitpods/gitpods"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pressly/chi"
)

type UsersRepositoriesStore interface {
	List(username string) ([]*gitpods.Repository, error)
}

type UsersRepositoriesAPI struct {
	logger log.Logger
	store  UsersRepositoriesStore
}

func (a UsersRepositoriesAPI) Routes() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/", a.List)

	return r
}

func (a UsersRepositoriesAPI) List(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")

	repositories, err := a.store.List(username)
	if err != nil {
		msg := "failed to list user's repositories"
		level.Warn(a.logger).Log(
			"msg", msg,
			"err", err,
		)
		jsonResponse(w, msg, http.StatusInternalServerError)
		return
	}

	jsonResponse(w, repositories, http.StatusOK)
}
