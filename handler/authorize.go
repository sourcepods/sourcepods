package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gitpods/gitpod"
	"github.com/gitpods/gitpod/store"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"golang.org/x/crypto/bcrypt"
)

const (
	loginAttemptFailed  = "failed"
	loginAttemptSuccess = "success"
)

var (
	BadCredentialsJson = map[string]string{"message": "Bad credentials"}
)

type LoginStore interface {
	GetUserByEmail(string) (gitpod.User, error)
}

type loginForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Authorize(logger log.Logger, s LoginStore) http.HandlerFunc {
	// TODO: This registers this collector globally, which causes them to register twice in tests and thus fails.
	//loginAttempts := prometheus.NewCounterVec(prometheus.CounterOpts{
	//	Namespace: "gitpod",
	//	Name:      "login_attempts",
	//	Help:      "Number of attempted logins and their success",
	//}, []string{"status"})
	//prometheus.MustRegister(loginAttempts)

	return func(w http.ResponseWriter, r *http.Request) {
		var form loginForm

		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
			level.Warn(logger).Log("msg", "failed to unmarshal form", "err", err)
			//loginAttempts.WithLabelValues(loginAttemptFailed).Inc()
			WriteJson(w, map[string]string{"message": "failed to unmarshal form"}, http.StatusBadRequest)
			return
		}

		user, err := s.GetUserByEmail(form.Email)
		if err == store.UserNotFound {
			level.Debug(logger).Log("msg", "user by email doesn't exist", "email", form.Email)
			//loginAttempts.WithLabelValues(loginAttemptFailed).Inc()
			WriteJson(w, BadCredentialsJson, http.StatusUnauthorized)
			return
		}
		if err != nil {
			level.Warn(logger).Log("msg", "failed to get user by email", "err", err)
			//loginAttempts.WithLabelValues(loginAttemptFailed).Inc()
			WriteJson(w, BadCredentialsJson, http.StatusUnauthorized)
			return
		}

		// TODO: Move this into some kind of service?
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password)); err != nil {
			level.Debug(logger).Log("msg", "login password doesn't match", "err", err)
			//loginAttempts.WithLabelValues(loginAttemptFailed).Inc()
			WriteJson(w, BadCredentialsJson, http.StatusUnauthorized)
			return
		}

		//loginAttempts.WithLabelValues(loginAttemptSuccess).Inc()

		WriteJson(w, user, http.StatusOK)
	}
}
