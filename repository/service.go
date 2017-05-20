package repository

import (
	"github.com/gitpods/gitpods"
)

type (
	Store interface {
		ListByOwner(string) ([]*gitpods.Repository, error)
	}

	UserStore interface {
		FindByUsername(string) (*gitpods.User, error)
	}

	Service interface {
		ListByOwnerUsername(string) ([]*gitpods.Repository, error)
	}

	service struct {
		repositories Store
		users        UserStore
	}
)

func NewService(users UserStore, repositories Store) Service {
	return &service{
		users:        users,
		repositories: repositories,
	}
}

func (s *service) ListByOwnerUsername(username string) ([]*gitpods.Repository, error) {
	u, err := s.users.FindByUsername(username)
	if err != nil {
		return nil, err
	}

	repositories, err := s.repositories.ListByOwner(u.ID)
	if err != nil {
		return nil, err
	}

	return repositories, err
}

//type UsersRepositoriesStore interface {
//	List(username string) ([]*gitpods.Store, error)
//}
//
//type UsersRepositoriesAPI struct {
//	logger log.Logger
//	store  UsersRepositoriesStore
//}
//
//func (a UsersRepositoriesAPI) Routes() *chi.Mux {
//	r := chi.NewRouter()
//	r.Get("/", a.List)
//
//	return r
//}
//
//func (a UsersRepositoriesAPI) List(w http.ResponseWriter, r *http.Request) {
//	username := chi.URLParam(r, "username")
//
//	repositories, err := a.store.List(username)
//	if err != nil {
//		msg := "failed to list user's repositories"
//		level.Warn(a.logger).Log(
//			"msg", msg,
//			"err", err,
//		)
//		jsonResponse(w, msg, http.StatusInternalServerError)
//		return
//	}
//
//	jsonResponse(w, repositories, http.StatusOK)
//}
