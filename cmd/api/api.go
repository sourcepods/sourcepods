package main

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gitpods/gitpods/handler"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/metrics/prometheus"
	"github.com/oklog/oklog/pkg/group"
	"github.com/pressly/chi"
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/urfave/cli"
)

const (
	FlagAddr               = "addr"
	FlagDatabaseDriver     = "database-driver"
	FlagDatabaseDatasource = "database-datasource"
	FlagEnv                = "env"
	FlagLogLevel           = "loglevel"
	FlagSecret             = "secret"

	ProductionEnv = "production"
)

var FlagsAPI = []cli.Flag{
	cli.StringFlag{
		Name:   FlagAddr,
		EnvVar: "GITPODS_ADDR",
		Usage:  "The address gitpods API runs on",
		Value:  ":3010",
	},
	cli.StringFlag{
		Name:   FlagDatabaseDriver,
		EnvVar: "GITPODS_DATABASE_DRIVER",
		Usage:  "The database driver to use: memory & postgres",
		Value:  "postgres",
	},
	cli.StringFlag{
		Name:   FlagDatabaseDatasource,
		EnvVar: "GITPODS_DATABASE_DATASOURCE",
		Usage:  "The database connection data",
	},
	cli.StringFlag{
		Name:   FlagEnv,
		EnvVar: "GITPODS_ENV",
		Usage:  "The environment gitpods should run in",
		Value:  ProductionEnv,
	},
	cli.StringFlag{
		Name:   FlagLogLevel,
		EnvVar: "GITPODS_LOGLEVEL",
		Usage:  "The log level to filter logs with before printing",
		Value:  "info",
	},
	cli.StringFlag{
		Name:   FlagSecret,
		EnvVar: "GITPODS_SECRET",
		Usage:  "This secret is going to be used to generate cookies",
		Value:  "secret", // TODO: Remove this to force users to pass a real secret, no default
	},
}

type StoreCloser func() error

func ActionAPI(c *cli.Context) error {
	addr := c.String(FlagAddr)
	databaseDriver := c.String(FlagDatabaseDriver)
	databaseDSN := c.String(FlagDatabaseDatasource)
	env := c.String(FlagEnv)
	loglevel := c.String(FlagLogLevel)
	secret := c.String(FlagSecret)

	// Create the logger based on the environment: production/development/test
	logger := newLogger(env, loglevel)
	logger = log.WithPrefix(logger, "app", "api")

	store, dbCloser, err := NewRouterStore(databaseDriver, databaseDSN, []byte(secret))
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
		Addr:    addr,
		Handler: r,
	}

	var gr group.Group
	{
		gr.Add(func() error {
			level.Info(logger).Log("msg", "starting gitpods api", "addr", addr)
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

func newLogger(env string, loglevel string) log.Logger {
	var logger log.Logger

	if env == ProductionEnv {
		logger = log.NewJSONLogger(log.NewSyncWriter(os.Stdout))
	} else {
		logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	}

	switch strings.ToLower(loglevel) {
	case "debug":
		logger = level.NewFilter(logger, level.AllowDebug())
	case "warn":
		logger = level.NewFilter(logger, level.AllowWarn())
	case "error":
		logger = level.NewFilter(logger, level.AllowError())
	default:
		logger = level.NewFilter(logger, level.AllowInfo())
	}

	return log.With(logger, "ts", log.DefaultTimestampUTC)
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
