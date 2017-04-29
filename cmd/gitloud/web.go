package main

import (
	"net/http"
	"os"

	"github.com/gitpods/gitpod/handler"
	"github.com/gitpods/gitpod/store"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/gobuffalo/packr"
	"github.com/urfave/cli"
)

func ActionWeb(c *cli.Context) error {
	addr := c.String("addr")
	env := c.String("env")

	// Create the logger based on the environment: production/development/test
	logger := newLogger(env)

	// Create FileServer handler with buffalo's packr to serve file from disk or from within the binary.
	// The path is relative to this file.
	box := packr.NewBox("../../public")

	// Create a simple store running in memory for example purposes
	userStore := store.NewUserInMemory()

	// Create the http router and return it for use
	r := handler.NewRouter(logger, box, userStore)

	level.Info(logger).Log("msg", "starting gitloud", "addr", addr)
	return http.ListenAndServe(addr, r)
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
