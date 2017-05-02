package handler

import (
	"net/http"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/metrics"
	"github.com/gobuffalo/packr"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	NotFoundJson = map[string]string{"error": http.StatusText(http.StatusNotFound)}
)

// TODO: Refactor this to possibly a struct with instances of interfaces
type Store interface {
	LoginStore
	UserStore
}

type RouterMetrics struct {
	LoginAttempts metrics.Counter
}

func NewRouter(logger log.Logger, metrics RouterMetrics, box packr.Box, store Store) *mux.Router {
	r := mux.NewRouter().StrictSlash(true)

	// instantiate default middlewares
	middlewares := alice.New(LoggerMiddleware(logger))

	r.Handle("/", middlewares.ThenFunc(HomeHandler(box))).Methods(http.MethodGet)
	r.Handle("/favicon.ico", middlewares.Then(http.FileServer(box))).Methods(http.MethodGet)
	r.Handle("/favicon.png", middlewares.Then(http.FileServer(box))).Methods(http.MethodGet)
	r.PathPrefix("/js").Handler(middlewares.Then(http.FileServer(box))).Methods(http.MethodGet)
	r.PathPrefix("/css").Handler(middlewares.Then(http.FileServer(box))).Methods(http.MethodGet)
	r.PathPrefix("/img").Handler(middlewares.Then(http.FileServer(box))).Methods(http.MethodGet)

	r.Handle("/metrics", promhttp.Handler()).Methods(http.MethodGet)

	api := r.PathPrefix("/api").Subrouter()
	{
		api.Handle("/authorize", middlewares.ThenFunc(Authorize(logger, metrics.LoginAttempts, store))).Methods(http.MethodPost)

		api.Handle("/users", middlewares.ThenFunc(UserList(store))).Methods(http.MethodGet)
		api.Handle("/users", middlewares.ThenFunc(UserCreate(store))).Methods(http.MethodPost)
		api.Handle("/users/{username}", middlewares.ThenFunc(User(store))).Methods(http.MethodGet)
		api.Handle("/users/{username}", middlewares.ThenFunc(UserUpdate(store))).Methods(http.MethodPut)
		api.Handle("/users/{username}", middlewares.ThenFunc(UserDelete(store))).Methods(http.MethodDelete)

		api.NotFoundHandler = middlewares.ThenFunc(func(w http.ResponseWriter, r *http.Request) {
			WriteJson(w, NotFoundJson, http.StatusNotFound)
		})
	}

	r.NotFoundHandler = HomeHandler(box)

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
