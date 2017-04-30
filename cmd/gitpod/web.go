package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gitpods/gitpod/handler"
	"github.com/gitpods/gitpod/store"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/gobuffalo/packr"
	"github.com/oklog/oklog/pkg/group"
	"github.com/urfave/cli"
)

func ActionWeb(c *cli.Context) error {
	addr := c.String("addr")
	env := c.String("env")

	// Create the logger based on the environment: production/development/test
	logger := newLogger(env)

	var server *http.Server

	var gr group.Group
	gr.Add(func() error {
		// Create FileServer handler with buffalo's packr to serve file from disk or from within the binary.
		// The path is relative to this file.
		box := packr.NewBox("../../public")

		// Create a simple store running in memory for example purposes
		userStore := store.NewUserInMemory()

		// Create the http router and return it for use
		r := handler.NewRouter(logger, box, userStore)

		level.Info(logger).Log("msg", "starting gitpod", "addr", addr)

		server := &http.Server{Addr: addr, Handler: r}
		return server.ListenAndServe()
	}, func(err error) {
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		server.Shutdown(ctx)
		level.Info(logger).Log(
			"msg", "http server shutdown gracefully",
			"err", err,
		)
	})

	return gr.Run()
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
