package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"time"

	"github.com/gitpods/gitpods/cmd"
	"github.com/oklog/run"
	"github.com/urfave/cli"
)

var (
	devFlags = []cli.Flag{
		// Global
		cli.StringFlag{
			Name:  "database-dsn",
			Usage: "The database connection data",
			Value: "postgres://root@localhost:26257/gitpods?sslmode=disable",
		},
		cli.BoolFlag{
			Name:  "log-json",
			Usage: "Log json instead of key-value pairs",
		},
		cli.StringFlag{
			Name:  "log-level",
			Usage: "The log level to filter logs with before printing",
			Value: "debug",
		},
		cli.BoolTFlag{
			Name:  "tracing",
			Usage: "Enable tracing",
		},
		cli.BoolFlag{
			Name:  "watch,w",
			Usage: "Watch files in this project and rebuild binaries if something changes",
		},
		// API
		cli.StringFlag{
			Name:  "api-addr",
			Usage: "The address to run the API on",
			Value: ":3020",
		},
		// Storage
		cli.StringFlag{
			Name:  "storage-addr",
			Usage: "The address to run the storage on",
			Value: ":3030",
		},
		cli.StringFlag{
			Name:  "storage-root",
			Usage: "Storage's root to write to",
			Value: "./dev/storage-data",
		},
		// UI
		cli.StringFlag{
			Name:  "ui",
			Usage: "How to run the UI. Run docker container or compile and run a binary. Run Dart Dev server",
			Value: "docker",
		},
		cli.StringFlag{
			Name:  "ui-addr",
			Usage: "The address to run the UI on",
			Value: ":3010",
		},
	}
)

func devAction(c *cli.Context) error {
	// Global
	databaseDSNFlag := c.String("database-dsn")
	logJSONFlag := c.Bool("log-json")
	loglevelFlag := c.String("log-level")
	tracingFlag := c.BoolT("tracing")
	watchFlag := c.Bool("watch")

	// API
	apiAddrFlag := c.String("api-addr")

	// Storage
	storageAddrFlag := c.String("storage-addr")
	storageRootFlag := c.String("storage-root")

	// UI
	uiModeFlag := c.String("ui")
	uiAddrFlag := c.String("ui-addr")

	tracingURL := ""
	if tracingFlag {
		tracingURL = "localhost:6831"
	}

	uiRunner := NewGitPodsRunner("ui", []string{
		fmt.Sprintf("%s=%s", cmd.EnvHTTPAddr, uiAddrFlag),
		fmt.Sprintf("%s=%s", cmd.EnvAPIURL, "http://localhost:3000/api"), // TODO
		fmt.Sprintf("%s=%s", cmd.EnvLogLevel, loglevelFlag),
		fmt.Sprintf("%s=%v", cmd.EnvLogJSON, logJSONFlag),
		fmt.Sprintf("%s=%v", cmd.EnvTracingURL, tracingURL),
	})

	apiRunner := NewGitPodsRunner("api", []string{
		fmt.Sprintf("%s=%s", cmd.EnvHTTPAddr, apiAddrFlag),
		fmt.Sprintf("%s=%s", cmd.EnvDatabaseDSN, databaseDSNFlag),
		fmt.Sprintf("%s=%s", cmd.EnvMigrationsPath, "./schema/postgres"),
		fmt.Sprintf("%s=%s", cmd.EnvLogLevel, loglevelFlag),
		fmt.Sprintf("%s=%v", cmd.EnvLogJSON, logJSONFlag),
		fmt.Sprintf("%s=%s", cmd.EnvSecret, "secret"),
		fmt.Sprintf("%s=%s", cmd.EnvStorageGRPCURL, "localhost:3033"),
		fmt.Sprintf("%s=%s", cmd.EnvStorageHTTPURL, "http://localhost:3030"),
		fmt.Sprintf("%s=%v", cmd.EnvTracingURL, tracingURL),
	})

	storageRunner := NewGitPodsRunner("storage", []string{
		fmt.Sprintf("%s=%s", cmd.EnvHTTPAddr, storageAddrFlag),
		fmt.Sprintf("%s=%s", cmd.EnvLogLevel, loglevelFlag),
		fmt.Sprintf("%s=%v", cmd.EnvLogJSON, logJSONFlag),
		fmt.Sprintf("%s=%s", cmd.EnvRoot, storageRootFlag),
		fmt.Sprintf("%s=%v", cmd.EnvTracingURL, tracingURL),
	})

	caddy := CaddyRunner{}

	if watchFlag {
		watcher := &FileWatcher{}
		watcher.Add(uiRunner, apiRunner, storageRunner)

		go watcher.Watch()
	}

	var g run.Group
	{
		stop := make(chan os.Signal, 1)
		g.Add(func() error {
			log.Println("waiting for interrupt")
			signal.Notify(stop, os.Interrupt)
			<-stop
			return nil
		}, func(err error) {
			close(stop)
		})
	}
	{
		g.Add(func() error {
			log.Println("starting api")
			return apiRunner.Run()
		}, func(err error) {
			log.Println("stopping api")
			apiRunner.Shutdown()
		})
	}
	{
		g.Add(func() error {
			log.Println("starting storage")
			return storageRunner.Run()
		}, func(err error) {
			log.Println("stopping storage")
			storageRunner.Shutdown()
		})
	}
	{
		if uiModeFlag == "binary" {
			g.Add(func() error {
				log.Println("starting ui")
				return uiRunner.Run()
			}, func(err error) {
				log.Println("stopping ui")
				uiRunner.Shutdown()
			})
		}
	}
	{
		g.Add(func() error {
			log.Println("starting caddy")
			return caddy.Run()
		}, func(err error) {
			log.Println("stopping caddy")
			caddy.Stop()
		})
	}

		{
			c := exec.Command("webdev", "serve", "--hot-reload", "web:3011")
			g.Add(func() error {
				c.Dir = "ui"
				c.Stdout = os.Stdout
				c.Stderr = os.Stderr
				c.Stdin = os.Stdin
				return c.Run()
			}, func(err error) {
				if c == nil || c.Process == nil {
					return
				}
				c.Process.Kill()
			})
		}
		{
			redirect := func(path string) bool {
				if path == "/main.dart.js" {
					return false
				}
				if path == "/main.dart.bootstrap.js" {
					return false
				}
				if path == "/main.ddc.js" {
					return false
				}
				if path == "/main.ddc.js.map" {
					return false
				}
				if path == "/$assetDigests" {
					return false
				}
				if strings.HasPrefix(path, "/components") {
					return false
				}
				if strings.HasPrefix(path, "/img") {
					return false
				}
				if strings.HasPrefix(path, "/packages") {
					return false
				}
				return true
			}

			director := func(r *http.Request) {
				if redirect(r.URL.Path) {
					r.URL.Path = "/"
				}
				r.URL.Scheme = "http"
				r.URL.Host = "localhost:3011"
			}

			server := &http.Server{
				Addr: ":3010",
				Handler: &httputil.ReverseProxy{
					Director: director,
				},
			}

			g.Add(func() error {
				return server.ListenAndServe()
			}, func(err error) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				server.Shutdown(ctx)
			})
		}
	}

	return g.Run()
}
