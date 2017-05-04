package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gitpods/gitpod"
	"github.com/gitpods/gitpod/store"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/metrics"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

const (
	loginAttemptFailed  = "failed"
	loginAttemptSuccess = "success"
)

var (
	BadCredentialsJson = map[string]string{"message": "Bad credentials"}
	cookieStore        = sessions.NewCookieStore([]byte("secret"))
)

type LoginStore interface {
	GetUserByEmail(string) (gitpod.User, error)
}

type loginForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Authorize(logger log.Logger, loginAttempts metrics.Counter, s LoginStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var form loginForm

		if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
			level.Warn(logger).Log("msg", "failed to unmarshal form", "err", err)
			loginAttempts.With("status", loginAttemptFailed).Add(1)
			jsonResponse(w, map[string]string{"message": "failed to unmarshal form"}, http.StatusBadRequest)
			return
		}

		user, err := s.GetUserByEmail(form.Email)
		if err == store.UserNotFound {
			level.Debug(logger).Log("msg", "user by email doesn't exist", "email", form.Email)
			loginAttempts.With("status", loginAttemptFailed).Add(1)
			jsonResponse(w, BadCredentialsJson, http.StatusUnauthorized)
			return
		}
		if err != nil {
			level.Warn(logger).Log("msg", "failed to get user by email", "err", err)
			loginAttempts.With("status", loginAttemptFailed).Add(1)
			jsonResponse(w, BadCredentialsJson, http.StatusUnauthorized)
			return
		}

		// TODO: Move this into some kind of service?
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password)); err != nil {
			level.Debug(logger).Log("msg", "login password doesn't match", "err", err)
			loginAttempts.With("status", loginAttemptFailed).Add(1)
			jsonResponse(w, BadCredentialsJson, http.StatusUnauthorized)
			return
		}

		loginAttempts.With("status", loginAttemptSuccess).Add(1)

		jsonResponse(w, user, http.StatusOK)
	}
}

func AuthorizedUser(logger log.Logger, s LoginStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := s.GetUserByEmail("metalmatze@example.com") // TODO: Get username or id from some token
		if err == store.UserNotFound {
			jsonResponse(w, NotFoundJson, http.StatusNotFound)
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
