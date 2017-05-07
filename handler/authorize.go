package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gitpods/gitpods"
	"github.com/gitpods/gitpods/store"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/metrics"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

const (
	// SessionName is the name used to store the cookie on the client
	SessionName         = "_gitpods_session"
	SessionUserID       = "user_id"
	SessionUserUsername = "user_username"

	loginAttemptFailed  = "failed"
	loginAttemptSuccess = "success"
)

type LoginStore interface {
	GetUser(string) (gitpods.User, error)
	GetUserByEmail(string) (gitpods.User, error)
}

// Authorized users will have a user information in the next handlers.
func Authorized(logger log.Logger, cookies sessions.Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := cookies.Get(r, SessionName)
			if err != nil {
				level.Warn(logger).Log("msg", "can't get session from cookie store", "err", err)
				jsonResponseBytes(w, JsonBadCredentials, http.StatusInternalServerError)
				return
			}

			id, found := session.Values[SessionUserID]
			if !found || id == "" {
				level.Debug(logger).Log("msg", "user's id can't be found in session")
				jsonResponseBytes(w, JsonUnauthorized, http.StatusUnauthorized)
				return
			}

			username, found := session.Values[SessionUserUsername]
			if !found || username == "" {
				level.Debug(logger).Log("msg", "user's username can't be found in session")
				jsonResponseBytes(w, JsonUnauthorized, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), SessionUserID, id)
			ctx = context.WithValue(r.Context(), SessionUserUsername, username)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

func Authorize(logger log.Logger, loginAttempts metrics.Counter, cookies sessions.Store, s LoginStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := cookies.Get(r, SessionName)
		if err != nil {
			level.Warn(logger).Log("msg", "can't get session from cookie store", "err", err)
			loginAttempts.With("status", loginAttemptFailed).Add(1)
			jsonResponseBytes(w, JsonBadCredentials, http.StatusInternalServerError)
			return
		}

		var form struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

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
			jsonResponseBytes(w, JsonBadCredentials, http.StatusUnauthorized)
			return
		}
		if err != nil {
			level.Warn(logger).Log("msg", "failed to get user by email", "err", err)
			loginAttempts.With("status", loginAttemptFailed).Add(1)
			jsonResponseBytes(w, JsonBadCredentials, http.StatusUnauthorized)
			return
		}

		// TODO: Move this into some kind of service?
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password)); err != nil {
			level.Debug(logger).Log("msg", "login password doesn't match", "err", err)
			loginAttempts.With("status", loginAttemptFailed).Add(1)
			jsonResponseBytes(w, JsonBadCredentials, http.StatusUnauthorized)
			return
		}

		session.Values[SessionUserID] = user.ID
		session.Values[SessionUserUsername] = user.Username

		if err := session.Save(r, w); err != nil {
			level.Warn(logger).Log("msg", "can't save session to cookie store", "err", err)
			loginAttempts.With("status", loginAttemptFailed).Add(1)
			jsonResponseBytes(w, JsonBadCredentials, http.StatusInternalServerError)
			return
		}

		loginAttempts.With("status", loginAttemptSuccess).Add(1)
		jsonResponse(w, user, http.StatusOK)
	}
}

func AuthorizedUser(logger log.Logger, s LoginStore) http.HandlerFunc {
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
