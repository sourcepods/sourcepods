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
		cli.StringFlag{
			Name:  "ui-addr",
			Usage: "The address to run the UI on",
			Value: ":3010",
		},
		cli.StringFlag{
			Name:  "api-addr",
			Usage: "The address to run the API on",
			Value: ":3020",
		},
		cli.BoolFlag{
			Name:  "dart",
			Usage: "Run pub serve as a development server for dart",
		},
		cli.StringFlag{
			Name:  "database-driver",
			Usage: "The database driver to use: memory & postgres",
			Value: "postgres",
		},
		cli.StringFlag{
			Name:  "database-dsn",
			Usage: "The database connection data",
			Value: "postgres://postgres:postgres@localhost:5432?sslmode=disable",
		},
		cli.StringFlag{
			Name:  "log-level",
			Usage: "The log level to filter logs with before printing",
			Value: "debug",
		},
		cli.BoolFlag{
			Name:  "log-json",
			Usage: "Log json instead of key-value pairs",
		},
		cli.StringFlag{
			Name:  "root",
			Usage: "Storage's root to write to",
			Value: "./dev/storage-data",
		},
		cli.StringFlag{
			Name:  "storage-addr",
			Usage: "The address to run the storage on",
			Value: ":3030",
		},
		cli.BoolTFlag{
			Name:  "tracing",
			Usage: "Enable tracing",
		},
		cli.BoolFlag{
			Name:  "watch,w",
			Usage: "Watch files in this project and rebuild binaries if something changes",
		},
	}
)

func devAction(c *cli.Context) error {
	apiAddrFlag := c.String("api-addr")
	dart := c.Bool("dart")
	databaseDriver := c.String("database-driver")
	databaseDSN := c.String("database-dsn")
	logJSONFlag := c.Bool("log-json")
	loglevelFlag := c.String("log-level")
	rootFlag := c.String("root")
	storageAddrFlag := c.String("storage-addr")
	tracingFlag := c.BoolT("tracing")
	uiAddrFlag := c.String("ui-addr")
	watch := c.Bool("watch")

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
		fmt.Sprintf("%s=%s", cmd.EnvDatabaseDriver, databaseDriver),
		fmt.Sprintf("%s=%s", cmd.EnvDatabaseDSN, databaseDSN),
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
		fmt.Sprintf("%s=%s", cmd.EnvRoot, rootFlag),
		fmt.Sprintf("%s=%v", cmd.EnvTracingURL, tracingURL),
	})

	caddy := CaddyRunner{}

	if watch {
		watcher := &FileWatcher{}
		watcher.Add(uiRunner, apiRunner)

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
		g.Add(func() error {
			log.Println("starting caddy")
			return caddy.Run()
		}, func(err error) {
			log.Println("stopping caddy")
			caddy.Stop()
		})
	}

	if dart {
		{
			c := exec.Command("pub", "serve", "--port=3011")
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
				if path == "/main.dart" {
					return false
				}
				if path == "/main.template.dart" {
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
	} else {
		{
			g.Add(func() error {
				log.Println("starting ui")
				return uiRunner.Run()
			}, func(err error) {
				log.Println("stopping ui")
				uiRunner.Shutdown()
			})
		}
	}

	return g.Run()
}
