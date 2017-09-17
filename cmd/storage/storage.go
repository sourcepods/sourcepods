package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gitpods/gitpods/cmd"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/oklog/oklog/pkg/group"
	"github.com/urfave/cli"
)

type storageConf struct {
	Addr     string
	LogJSON  bool
	LogLevel string
	Root     string
}

var (
	storageConfig = storageConf{}

	storageFlags = []cli.Flag{
		cli.StringFlag{
			Name:        cmd.FlagHTTPAddr,
			EnvVar:      cmd.EnvHTTPAddr,
			Value:       ":3030",
			Destination: &storageConfig.Addr,
		},
		cli.BoolFlag{
			Name:        cmd.FlagLogJSON,
			EnvVar:      cmd.EnvLogJSON,
			Usage:       "The logger will log json lines",
			Destination: &storageConfig.LogJSON,
		},
		cli.StringFlag{
			Name:        cmd.FlagLogLevel,
			EnvVar:      cmd.EnvLogLevel,
			Usage:       "The log level to filter logs with before printing",
			Value:       "info",
			Destination: &storageConfig.LogLevel,
		},
		cli.StringFlag{
			Name:        cmd.FlagRoot,
			EnvVar:      cmd.EnvRoot,
			Usage:       "The root folder to store all git repositories in",
			Destination: &storageConfig.Root,
		},
	}
)

func storageAction(c *cli.Context) error {
	logger := cmd.NewLogger(storageConfig.LogJSON, storageConfig.LogLevel)
	logger = log.WithPrefix(logger, "app", "storage")

	if storageConfig.Root == "" {
		return errors.New("the root has to be a valid path")
	}

	if err := os.MkdirAll(storageConfig.Root, 0755); err != nil {
		return fmt.Errorf("failed to create storage root: %s", storageConfig.Root)
	}

	gh := NewGitHTTP(storageConfig.Root)
	gh.Logger = logger

	server := &http.Server{
		Addr:    storageConfig.Addr,
		Handler: gh.Handler(),
	}

	var gr group.Group
	{
		gr.Add(func() error {
			sig := make(chan os.Signal)
			signal.Notify(sig, os.Interrupt)
			<-sig
			return nil
		}, func(err error) {

		})
	}
	{
		gr.Add(func() error {
			level.Info(logger).Log(
				"msg", "starting gitpods storage",
				"addr", storageConfig.Addr,
			)
			return server.ListenAndServe()
		}, func(err error) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			if err := server.Shutdown(ctx); err != nil {
				level.Error(logger).Log(
					"msg", "failed to shutdown http server gracefully",
					"err", err,
				)
				return
			}
			level.Info(logger).Log("msg", "http server shutdown gracefully")
		})
	}

	return gr.Run()
}
