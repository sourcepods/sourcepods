package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gitpods/gitpods/authorization"
	"github.com/gitpods/gitpods/cmd"
	"github.com/gitpods/gitpods/repository"
	"github.com/gitpods/gitpods/session"
	"github.com/gitpods/gitpods/user"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	_ "github.com/lib/pq"
	"github.com/oklog/oklog/pkg/group"
	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	"github.com/urfave/cli"
)

type apiConf struct {
	Addr           string
	DatabaseDriver string
	DatabaseDSN    string
	LogJSON        bool
	LogLevel       string
	Secret         string
}

var (
	apiConfig = apiConf{}

	apiFlags = []cli.Flag{
		cli.StringFlag{
			Name:        cmd.FlagAddr,
			EnvVar:      cmd.EnvAddr,
			Usage:       "The address gitpods API runs on",
			Value:       ":3010",
			Destination: &apiConfig.Addr,
		},
		cli.StringFlag{
			Name:        cmd.FlagDatabaseDriver,
			EnvVar:      cmd.EnvDatabaseDriver,
			Usage:       "The database driver to use: memory & postgres",
			Value:       "postgres",
			Destination: &apiConfig.DatabaseDriver,
		},
		cli.StringFlag{
			Name:        cmd.FlagDatabaseDSN,
			EnvVar:      cmd.EnvDatabaseDSN,
			Usage:       "The database connection data",
			Destination: &apiConfig.DatabaseDSN,
		},
		cli.StringFlag{
			Name:        cmd.FlagLogLevel,
			EnvVar:      cmd.EnvLogLevel,
			Usage:       "The log level to filter logs with before printing",
			Value:       "info",
			Destination: &apiConfig.LogLevel,
		},
		cli.BoolFlag{
			Name:        cmd.FlagLogJSON,
			EnvVar:      cmd.EnvLogJSON,
			Usage:       "The logger will log json lines",
			Destination: &apiConfig.LogJSON,
		},
		cli.StringFlag{
			Name:        cmd.FlagSecret,
			EnvVar:      cmd.EnvSecret,
			Usage:       "This secret is going to be used to generate cookies",
			Destination: &apiConfig.Secret,
		},
	}
)

func apiAction(c *cli.Context) error {
	if apiConfig.Secret == "" {
		return errors.New("the secret for the api can't be empty")
	}

	logger := cmd.NewLogger(apiConfig.LogJSON, apiConfig.LogLevel)
	logger = log.WithPrefix(logger, "app", "api")

	//
	// Stores
	//
	var (
		repositories repository.Store
		sessions     session.Store
		users        user.Store
	)

	switch apiConfig.DatabaseDriver {
	default:
		db, err := sql.Open("postgres", apiConfig.DatabaseDSN)
		if err != nil {
			return err
		}
		defer db.Close()

		users = user.NewPostgresStore(db)
		sessions = session.NewPostgresStore(db)
		repositories = repository.NewPostgresStore(db)
	}

	//
	// Services
	//
	var ss session.Service
	ss = session.NewService(sessions)

	var as authorization.Service
	as = authorization.NewService(users.(authorization.Store), ss)
	as = authorization.NewLoggingService(log.WithPrefix(logger, "service", "authorization"), as)

	var us user.Service
	us = user.NewService(users)
	us = user.NewLoggingService(log.WithPrefix(logger, "service", "user"), us)

	var rs repository.Service
	rs = repository.NewService(users, repositories)

	//
	// Router
	//
	router := chi.NewRouter()

	//httpLogger := log.With(logger, "component", "http")
	//router.Use(handler.LoggerMiddleware(httpLogger))
	router.Use(middleware.Logger) // TODO

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "hi")
	})

	router.Mount("/authorize", authorization.NewHandler(as))

	router.Group(func(router chi.Router) {
		router.Use(session.Authorized(ss))

		router.Mount("/user", user.NewUserHandler(us))
		router.Mount("/users", user.NewUsersHandler(us))
		router.Mount("/users/:username/repositories", repository.NewHandler(rs))
	})

	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error":"Not Found"}`))
	})

	server := &http.Server{
		Addr:    apiConfig.Addr,
		Handler: router,
	}

	var gr group.Group
	{
		gr.Add(func() error {
			dur := time.Minute
			level.Info(logger).Log("msg", "starting session cleaner", "interval", dur)
			for {
				if _, err := ss.ClearSessions(); err != nil {
					return err
				}
				time.Sleep(dur)
			}
		}, func(err error) {
		})
	}
	{
		gr.Add(func() error {
			level.Info(logger).Log(
				"msg", "starting gitpods api",
				"addr", apiConfig.Addr,
			)
			return server.ListenAndServe()
		}, func(err error) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			log.With(logger, "err", err)

			if err := server.Shutdown(ctx); err != nil {
				level.Error(logger).Log("msg", "failed to shutdown http server gracefully", "err", err)
				return
			}
			level.Info(logger).Log("msg", "http server shutdown gracefully")
		})
	}

	return gr.Run()
}
