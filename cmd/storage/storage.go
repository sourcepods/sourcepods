package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gitpods/gitpods/cmd"
	"github.com/gitpods/gitpods/storage"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/oklog/oklog/pkg/group"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
)

type storageConf struct {
	GRPCAddr string
	HTTPAddr string
	LogJSON  bool
	LogLevel string
	Root     string
}

var (
	storageConfig = storageConf{}

	storageFlags = []cli.Flag{
		cli.StringFlag{
			Name:        cmd.FlagGRPCAddr,
			EnvVar:      cmd.EnvGRPCAddr,
			Value:       ":3033",
			Destination: &storageConfig.GRPCAddr,
		},
		cli.StringFlag{
			Name:        cmd.FlagHTTPAddr,
			EnvVar:      cmd.EnvHTTPAddr,
			Value:       ":3030",
			Destination: &storageConfig.HTTPAddr,
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

	gitStorage, err := storage.NewStorage(storageConfig.Root)
	if err != nil {
		return err
	}

	lis, err := net.Listen("tcp", storageConfig.GRPCAddr)
	if err != nil {
		return fmt.Errorf("failed to create grpc listener: %v", err)
	}

	var gr group.Group
	{
		gr.Add(func() error {
			sig := make(chan os.Signal)
			signal.Notify(sig, os.Interrupt)
			<-sig
			return nil
		}, func(err error) {})
	}
	{
		gh := NewGitHTTP(storageConfig.Root)
		gh.Logger = logger

		server := &http.Server{
			Addr:    storageConfig.HTTPAddr,
			Handler: gh.Handler(),
		}

		gr.Add(func() error {
			level.Info(logger).Log(
				"msg", "starting gitpods storage http server",
				"addr", storageConfig.HTTPAddr,
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
	{
		gs := storage.NewStorageServer(grpc.NewServer(), gitStorage)
		gr.Add(func() error {
			level.Info(logger).Log(
				"msg", "starting gitpods storage grpc server",
				"addr", storageConfig.GRPCAddr,
			)
			return gs.Serve(lis)
		}, func(err error) {
			gs.GracefulStop()
			level.Info(logger).Log("msg", "grpc server shutdown gracefully")
		})
	}

	return gr.Run()
}
