package handler

import (
	"net/http"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/metrics"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var (
	JsonUnauthorized   = []byte(`{"message":"Unauthorized"}`) // 401
	JsonNotFound       = []byte(`{"error":"Not Found"}`)      // 404
	JsonBadCredentials = []byte(`{"message":"Bad credentials"}`)
)

type RouterStore struct {
	LoginStore  LoginStore
	UserStore   UserStore
	CookieStore sessions.Store
}

type RouterMetrics struct {
	LoginAttempts metrics.Counter
}

func NewRouter(logger log.Logger, metrics RouterMetrics, store RouterStore) *mux.Router {
	r := mux.NewRouter().StrictSlash(true)

	r.Path("/authorize").Methods(http.MethodPost).Handler(Authorize(logger, metrics.LoginAttempts, store.CookieStore, store.LoginStore))

	apiAuthRouter := NewAuthRouter(logger, metrics, store)
	r.PathPrefix("/").Handler(Authorized(logger, store.CookieStore)(apiAuthRouter))

	return r
}

func NewAuthRouter(logger log.Logger, metrics RouterMetrics, store RouterStore) *mux.Router {
	r := mux.NewRouter().StrictSlash(true)

	r.Path("/user").Methods(http.MethodGet).Handler(AuthorizedUser(logger, store.LoginStore))

	r.Path("/users").Methods(http.MethodGet).Handler(UserList(logger, store.UserStore))
	r.Path("/users").Methods(http.MethodPost).Handler(UserCreate(logger, store.UserStore))
	r.Path("/users/{username}").Methods(http.MethodGet).Handler(User(logger, store.UserStore))
	r.Path("/users/{username}").Methods(http.MethodPut).Handler(UserUpdate(logger, store.UserStore))
	r.Path("/users/{username}").Methods(http.MethodDelete).Handler(UserDelete(logger, store.UserStore))

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponseBytes(w, JsonNotFound, http.StatusNotFound)
	})

	return r
}

func LoggerMiddleware(logger log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			level.Debug(logger).Log(
				"duration", time.Since(start),
				"method", r.Method,
				"path", r.URL.Path,
			)
		})
	}
}
