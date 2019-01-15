package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/oklog/run"
	"github.com/sourcepods/sourcepods/cmd"
	"github.com/sourcepods/sourcepods/pkg/storage"
	jaeger "github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/urfave/cli"
)

type storageConf struct {
	GRPCAddr   string
	HTTPAddr   string
	LogJSON    bool
	LogLevel   string
	Root       string
	TracingURL string
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
		cli.StringFlag{
			Name:        cmd.FlagTracingURL,
			EnvVar:      cmd.EnvTracingURL,
			Usage:       "The url to send spans for tracing to",
			Destination: &storageConfig.TracingURL,
		},
	}
)

func storageAction(c *cli.Context) error {
	logger := cmd.NewLogger(storageConfig.LogJSON, storageConfig.LogLevel)
	logger = log.WithPrefix(logger, "app", c.App.Name)

	if storageConfig.TracingURL != "" {
		traceConfig := config.Configuration{
			Sampler: &config.SamplerConfig{
				Type:  jaeger.SamplerTypeConst,
				Param: 1,
			},
			Reporter: &config.ReporterConfig{
				LocalAgentHostPort: storageConfig.TracingURL,
			},
		}

		traceCloser, err := traceConfig.InitGlobalTracer(c.App.Name)
		if err != nil {
			return err
		}
		defer traceCloser.Close()

		level.Info(logger).Log(
			"msg", "tracing enabled",
			"addr", storageConfig.TracingURL,
		)
	} else {
		level.Info(logger).Log("msg", "tracing is disabled, no url given")
	}

	root := storageConfig.Root
	if root == "" {
		return errors.New("the root has to be a valid path")
	}

	if filepath.IsAbs(root) {
		root = filepath.Clean(root)
	} else {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		root = filepath.Join(wd, root)
	}

	gitStorage, err := storage.NewLocalStorage(storageConfig.Root)
	if err != nil {
		return err
	}

	lis, err := net.Listen("tcp", storageConfig.GRPCAddr)
	if err != nil {
		return fmt.Errorf("failed to create grpc listener: %v", err)
	}

	var gr run.Group
	{
		sig := make(chan os.Signal)
		gr.Add(func() error {
			signal.Notify(sig, os.Interrupt)
			<-sig
			return nil
		}, func(err error) {
			close(sig)
		})
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
				"msg", "starting SourcePods storage http server",
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
		gs := storage.NewStorageServer(gitStorage)
		gr.Add(func() error {
			level.Info(logger).Log(
				"msg", "starting SourcePods storage grpc server",
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
