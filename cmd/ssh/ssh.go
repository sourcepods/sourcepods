package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/oklog/run"
	"github.com/sourcepods/sourcepods/cmd"
	"github.com/sourcepods/sourcepods/pkg/ssh"
	"github.com/sourcepods/sourcepods/pkg/storage"
	jaeger "github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/urfave/cli"
)

type sshConf struct {
	HostKeyPath    string
	LogJSON        bool
	LogLevel       string
	SSHAddr        string
	StorageGRPCURL string
	TracingURL     string
}

var (
	sshConfig = sshConf{}
	sshFlags  = []cli.Flag{
		cli.StringFlag{
			Name:        cmd.FlagSSHAddr,
			Value:       ":3022",
			Usage:       "The SSH address to listen on",
			Destination: &sshConfig.SSHAddr,
		},
		cli.StringFlag{
			Name:        cmd.FlagSSHHostKeyPath,
			Value:       "/etc/ssh/",
			Usage:       "The path to looks for ssh host-keys in",
			Destination: &sshConfig.HostKeyPath,
		},
		cli.StringFlag{
			Name:        cmd.FlagStorageGRPCURL,
			Usage:       "The storage's grpc url to connect with",
			Destination: &sshConfig.StorageGRPCURL,
		},
		cli.BoolFlag{
			Name:        cmd.FlagLogJSON,
			Usage:       "The logger will log json lines",
			Destination: &sshConfig.LogJSON,
		},
		cli.StringFlag{
			Name:        cmd.FlagLogLevel,
			Usage:       "The log level to filter logs with before printing",
			Value:       "info",
			Destination: &sshConfig.LogLevel,
		},
		cli.StringFlag{
			Name:        cmd.FlagTracingURL,
			Usage:       "The url to send spans for tracing to",
			Destination: &sshConfig.TracingURL,
		},
	}
)

func sshAction(c *cli.Context) error {
	if sshConfig.StorageGRPCURL == "" {
		return errors.New("the storage grpc url can not be empty")
	}

	logger := cmd.NewLogger(sshConfig.LogJSON, sshConfig.LogLevel)
	logger = log.WithPrefix(logger, "app", c.App.Name)

	// TODO: Metrics FFS...
	//apiMetrics := apiMetrics()

	if sshConfig.TracingURL != "" {
		traceConfig := config.Configuration{
			Sampler: &config.SamplerConfig{
				Type:  jaeger.SamplerTypeConst,
				Param: 1,
			},
			Reporter: &config.ReporterConfig{
				LocalAgentHostPort: sshConfig.TracingURL,
			},
		}

		traceCloser, err := traceConfig.InitGlobalTracer(c.App.Name)
		if err != nil {
			return err
		}
		defer traceCloser.Close()

		level.Info(logger).Log(
			"msg", "tracing enabled",
			"addr", sshConfig.TracingURL,
		)
	} else {
		level.Info(logger).Log("msg", "tracing is disabled, no url given")
	}

	//
	// Storage
	//
	storageClient, err := storage.NewClient(sshConfig.StorageGRPCURL)
	if err != nil {
		return err
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
		ss := ssh.NewServer(sshConfig.SSHAddr, sshConfig.HostKeyPath, logger, storageClient)
		gr.Add(func() error {
			level.Info(logger).Log(
				"msg", "starting SourcePods git-ssh server",
				"addr", sshConfig.SSHAddr,
			)
			return ss.ListenAndServe()
		}, func(err error) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			ss.Shutdown(ctx)
			level.Info(logger).Log("msg", "grpc server shutdown gracefully")
		})
	}

	return gr.Run()
}
