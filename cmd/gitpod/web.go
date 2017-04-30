package main

import (
	"context"
	"net/http"
	"os"
	"strings"
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
	loglevel := c.String("loglevel")

	// Create the logger based on the environment: production/development/test
	logger := newLogger(env, loglevel)

	// Create FileServer handler with buffalo's packr to serve file from disk or from within the binary.
	// The path is relative to this file.
	box := packr.NewBox("../../public")

	// Create a simple store running in memory for example purposes
	userStore := store.NewUserInMemory()

	// Create the http router and return it for use
	r := handler.NewRouter(logger, box, userStore)

	server := &http.Server{Addr: addr, Handler: r}

	var gr group.Group
	{
		gr.Add(func() error {
			level.Info(logger).Log("msg", "starting gitpod", "addr", addr)
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
