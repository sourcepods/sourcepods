package main

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gitpods/gitpods/handler"
	"github.com/gitpods/gitpods/store"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/metrics/prometheus"
	"github.com/gorilla/sessions"
	"github.com/oklog/oklog/pkg/group"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/urfave/cli"
)

const (
	FlagAddr     = "addr"
	FlagEnv      = "env"
	FlagLogLevel = "loglevel"
	FlagSecret   = "secret"

	ProductionEnv = "production"
)

var FlagsAPI = []cli.Flag{
	cli.StringFlag{
		Name:   FlagAddr,
		EnvVar: "GITPODS_ADDR",
		Usage:  "The address gitpods API runs on",
		Value:  ":3000",
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

func ActionAPI(c *cli.Context) error {
	addr := c.String(FlagAddr)
	env := c.String(FlagEnv)
	loglevel := c.String(FlagLogLevel)
	secret := c.String(FlagSecret)

	// Create the logger based on the environment: production/development/test
	logger := newLogger(env, loglevel)

	cookieStore := sessions.NewFilesystemStore("/tmp/gitpods_sessions", []byte(secret))

	// Create a simple store running in memory for example purposes
	userStore := store.NewUserInMemory()

	// Create a routerStore by passing concrete implementations to interfaces for the router.
	routerStore := handler.RouterStore{
		CookieStore: cookieStore,
		UserStore:   userStore,
		LoginStore:  userStore,
	}

	// Create the http router and return it for use
	r := handler.NewRouter(logger, prometheusMetrics(), routerStore)

	server := &http.Server{Addr: addr, Handler: r}

	var gr group.Group
	{
		gr.Add(func() error {
			level.Info(logger).Log("msg", "starting gitpods API", "addr", addr)
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
		LoginAttempts: prometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "gitpods",
			Name:      "login_attempts_total",
			Help:      "Number of attempts to login and their status",
		}, []string{"status"}),
	}
}
