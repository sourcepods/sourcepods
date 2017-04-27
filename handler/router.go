package handler

import (
	"net/http"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/gobuffalo/packr"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewRouter(logger log.Logger, box packr.Box, userStore UserStore) *mux.Router {
	r := mux.NewRouter()

	// instantiate default middlewares
	middlewares := alice.New(LoggerMiddleware(logger))

	r.Handle("/", middlewares.ThenFunc(HomeHandler(box))).Methods(http.MethodGet)

	api := r.PathPrefix("/api").Subrouter()
	{
		api.Handle("/users", middlewares.ThenFunc(UserList(userStore))).Methods(http.MethodGet)
		api.Handle("/users", middlewares.ThenFunc(UserCreate(userStore))).Methods(http.MethodPost)
		api.Handle("/users/{username}", middlewares.ThenFunc(User(userStore))).Methods(http.MethodGet)
		api.Handle("/users/{username}", middlewares.ThenFunc(UserUpdate(userStore))).Methods(http.MethodPut)
		api.Handle("/users/{username}", middlewares.ThenFunc(UserDelete(userStore))).Methods(http.MethodDelete)

		api.NotFoundHandler = middlewares.ThenFunc(func(w http.ResponseWriter, r *http.Request) {
			WriteJson(w, map[string]string{"error": http.StatusText(http.StatusNotFound)}, http.StatusNotFound)
		})
	}

	r.Handle("/metrics", promhttp.Handler()).Methods(http.MethodGet)

	r.PathPrefix("/js").Handler(middlewares.Then(http.FileServer(box))).Methods(http.MethodGet)

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
