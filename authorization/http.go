package authorization

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gitpods/gitpods/session"
	"github.com/pressly/chi"
)

const megabyte = 1024 * 1024 * 1024

func NewHandler(s Service) *chi.Mux {
	r := chi.NewRouter()

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		var form struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(io.LimitReader(r.Body, megabyte)).Decode(&form); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return // TODO
		}

		user, err := s.AuthenticateUser(form.Email, form.Password)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return // TODO
		}

		sess, err := s.CreateSession(user.ID, user.Username)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return // TODO
		}

		cookie := &http.Cookie{
			Name:    session.CookieName,
			Value:   sess.ID,
			Path:    "/",
			Expires: sess.Expiry,
		}

		http.SetCookie(w, cookie)
	})

	return r
}
