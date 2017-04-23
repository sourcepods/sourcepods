package main

import (
	"net/http"
	"os"

	"github.com/gitloud/gitloud/handler"
	"github.com/gitloud/gitloud/store"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/urfave/cli"
)

func ActionWeb(config Config) cli.ActionFunc {
	return func(c *cli.Context) error {
		// Create the logger based on the environment: production/development/test
		logger := newLogger(config.Env)

		// Create a simple store running in memory for example purposes
		userStore := store.NewUserInMemory()

		// Create the http router and return it for use
		r := handler.NewRouter(logger, userStore)

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
