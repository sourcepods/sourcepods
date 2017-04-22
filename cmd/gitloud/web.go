package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gitloud/gitloud/handler"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/urfave/cli"
)

func ActionWeb(config Config) cli.ActionFunc {
	return func(c *cli.Context) error {
		logger := newLogger(config.Env)

		r := mux.NewRouter()

		// instantiate default middlewares
		middlewares := alice.New(LoggerMiddleware(logger))

		r.Handle("/", middlewares.ThenFunc(HomeHandler)).Methods(http.MethodGet)

		api := r.PathPrefix("/api").Subrouter()
		{
			api.Handle("/users", middlewares.ThenFunc(handler.UserList)).Methods(http.MethodGet)
			api.Handle("/users", middlewares.ThenFunc(handler.UserCreate)).Methods(http.MethodPost)
			api.Handle("/users/{id}", middlewares.ThenFunc(handler.User)).Methods(http.MethodGet)
			api.Handle("/users/{id}", middlewares.ThenFunc(handler.UserUpdate)).Methods(http.MethodPut)
			api.Handle("/users/{id}", middlewares.ThenFunc(handler.UserDelete)).Methods(http.MethodDelete)
		}

		r.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/"))).Methods(http.MethodGet)

		level.Info(logger).Log("msg", "starting gitloud", "addr", config.Addr)
		return http.ListenAndServe(config.Addr, r)
	}
}

func newLogger(env string) log.Logger {
	var logger log.Logger
	if env == DefaultEnv {
		logger = log.NewJSONLogger(log.NewSyncWriter(os.Stdout))
		logger = level.NewFilter(logger, level.AllowInfo())
	} else {
		logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
		logger = level.NewFilter(logger, level.AllowAll())
	}

	return log.With(logger, "ts", log.DefaultTimestampUTC)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "gitloud")
}

func LoggerMiddleware(logger log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			level.Debug(logger).Log(
				"duration", time.Since(start),
				"path", r.URL.Path,
				"method", r.Method,
			)
		})
	}
}
