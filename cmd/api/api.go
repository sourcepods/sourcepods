package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gitpods/gitpods/cmd"
	"github.com/gitpods/gitpods/handler"
	"github.com/gitpods/gitpods/user"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/metrics/prometheus"
	_ "github.com/lib/pq"
	"github.com/oklog/oklog/pkg/group"
	"github.com/pressly/chi"
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/urfave/cli"
)

type apiConf struct {
	Addr           string
	DatabaseDriver string
	DatabaseDSN    string
	LogJson        bool
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
			Name:        cmd.FlagLogJson,
			EnvVar:      cmd.EnvLogJson,
			Usage:       "The logger will log json lines",
			Destination: &apiConfig.LogJson,
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

	logger := cmd.NewLogger(apiConfig.LogJson, apiConfig.LogLevel)
	logger = log.WithPrefix(logger, "app", "api")

	//
	// Stores
	//
	var (
		users user.Store
	)

	switch apiConfig.DatabaseDriver {
	case "memory":
		users = user.NewMemoryRepository()
	default:
		db, err := sql.Open("postgres", apiConfig.DatabaseDSN)
		if err != nil {
			return err
		}

		users = user.NewPostgresRepository(db)
	}
	//
	// Services
	//
	var us user.Service
	us = user.NewService(users)
	us = user.NewLoggingService(log.WithPrefix(logger, "service", "user"), us)
	//
	//
	//

	httpLogger := log.With(logger, "component", "http")

	router := chi.NewRouter()
	router.Use(handler.LoggerMiddleware(httpLogger))
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "hi")
	})

	router.Mount("/users", user.NewHandler(us))

	http.ListenAndServe(":3020", router)

	store, dbCloser, err := NewRouterStore(apiConfig.DatabaseDriver, apiConfig.DatabaseDSN, []byte(apiConfig.Secret))
	if err != nil {
		level.Error(logger).Log(
			"msg", "failed to initialize store",
			"err", err,
		)
		os.Exit(1)
	}

	r := chi.NewRouter()
	r.Use(handler.LoggerMiddleware(logger))

	r.Mount("/", handler.NewRouter(logger, prometheusMetrics(), store))

	server := &http.Server{
		Addr:    apiConfig.Addr,
		Handler: r,
	}

	var gr group.Group
	{
		gr.Add(func() error {
			level.Info(logger).Log("msg", "starting gitpods api", "addr", apiConfig.Addr)
			return server.ListenAndServe()
		}, func(err error) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			log.With(logger, "grouperr", err)

			if err := server.Shutdown(ctx); err != nil {
				level.Error(logger).Log("msg", "failed to shutdown http server gracefully", "err", err)
				return
			}
			level.Info(logger).Log("msg", "http server shutdown gracefully")
		})
	}
	{
		gr.Add(func() error {
			select {}
		}, func(err error) {
			dbCloser()
			level.Info(logger).Log("msg", "database shutdown gracefully")
		})
	}

	return gr.Run()
}

func prometheusMetrics() handler.RouterMetrics {
	return handler.RouterMetrics{
		LoginAttempts: prometheus.NewCounterFrom(prom.CounterOpts{
			Namespace: "gitpods",
			Name:      "login_attempts_total",
			Help:      "Number of attempts to login and their status",
		}, []string{"status"}),
	}
}
