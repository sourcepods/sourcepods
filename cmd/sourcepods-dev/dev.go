package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/koding/websocketproxy"
	"github.com/oklog/run"
	"github.com/sourcepods/sourcepods/cmd"
	"github.com/urfave/cli"
)

var (
	devFlags = []cli.Flag{
		// Global
		cli.StringFlag{
			Name:  "database-dsn",
			Usage: "The database connection data",
			Value: "postgres://root@localhost:26257/sourcepods?sslmode=disable",
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

	for _, dir := range []string{"./dev", "./dev/keys"} {
		exists, err := exists(dir)
		if err != nil {
			return err
		}
		if !exists {
			color.HiRed("Development folder %s doesn't exists. Run `sourcepods-dev setup` first", dir)
			return nil
		}
	}

	uiRunner := NewRunner("ui", []string{
		fmt.Sprintf("%s=%s", cmd.EnvHTTPAddr, uiAddrFlag),
		fmt.Sprintf("%s=%s", cmd.EnvAPIURL, "http://localhost:3000/api"), // TODO
		fmt.Sprintf("%s=%s", cmd.EnvLogLevel, loglevelFlag),
		fmt.Sprintf("%s=%v", cmd.EnvLogJSON, logJSONFlag),
		fmt.Sprintf("%s=%v", cmd.EnvTracingURL, tracingURL),
	})

	apiRunner := NewRunner("api", []string{
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

	storageRunner := NewRunner("storage", []string{
		fmt.Sprintf("%s=%s", cmd.EnvHTTPAddr, storageAddrFlag),
		fmt.Sprintf("%s=%s", cmd.EnvLogLevel, loglevelFlag),
		fmt.Sprintf("%s=%v", cmd.EnvLogJSON, logJSONFlag),
		fmt.Sprintf("%s=%s", cmd.EnvRoot, storageRootFlag),
		fmt.Sprintf("%s=%v", cmd.EnvTracingURL, tracingURL),
	})

	sshRunner := NewRunner("ssh", []string{
		fmt.Sprintf("%s=%d", cmd.EnvSSHPort, 2222),
		fmt.Sprintf("%s=%s", cmd.EnvSSHHostKeyPath, "./dev/keys/"),
		fmt.Sprintf("%s=%s", cmd.EnvStorageGRPCURL, "localhost:3033"),
		fmt.Sprintf("%s=%s", cmd.EnvLogLevel, loglevelFlag),
		fmt.Sprintf("%s=%v", cmd.EnvLogJSON, logJSONFlag),
		fmt.Sprintf("%s=%v", cmd.EnvTracingURL, tracingURL),
	})

	caddy := CaddyRunner{}

	if watchFlag {
		watcher := &FileWatcher{}
		watcher.Add(uiRunner, apiRunner, storageRunner, sshRunner)

		go watcher.Watch()
	}

	var g run.Group
	{
		stop := make(chan os.Signal, 1)
		g.Add(func() error {
			signal.Notify(stop, os.Interrupt)
			<-stop
			return nil
		}, func(err error) {
			close(stop)
		})
	}
	{
		g.Add(func() error {
			color.HiGreen("starting api")
			return apiRunner.Run()
		}, func(err error) {
			color.HiYellow("stopping api")
			apiRunner.Shutdown()
		})
	}
	{
		g.Add(func() error {
			color.HiGreen("starting storage")
			return storageRunner.Run()
		}, func(err error) {
			color.HiYellow("stopping storage")
			storageRunner.Shutdown()
		})
	}
	{
		g.Add(func() error {
			color.HiGreen("starting ssh")
			return sshRunner.Run()
		}, func(err error) {
			color.HiYellow("stopping ssh")
			sshRunner.Shutdown()
		})
	}
	{
		if uiModeFlag == "binary" {
			g.Add(func() error {
				color.HiGreen("starting ui")
				return uiRunner.Run()
			}, func(err error) {
				color.HiYellow("stopping ui")
				uiRunner.Shutdown()
			})
		}
	}
	{
		g.Add(func() error {
			color.HiGreen("starting caddy")
			return caddy.Run()
		}, func(err error) {
			color.HiYellow("stopping caddy")
			caddy.Stop()
		})
	}

	if uiModeFlag == "docker" {
		err := ensureUIContainer()
		if err != nil {
			return fmt.Errorf("failed to start UI container: %v", err)
		}
	}

	if uiModeFlag == "dart" {
		{
			c := exec.Command("webdev", "serve", "--hot-reload", "web:3011")
			g.Add(func() error {
				color.HiGreen("starting dart")

				c.Dir = "ui"
				stdout, err := c.StdoutPipe()
				if err != nil {
					return err
				}
				stderr, err := c.StderrPipe()
				if err != nil {
					return err
				}

				multi := io.MultiReader(stdout, stderr)

				if err := c.Start(); err != nil {
					return err
				}

				scanner := bufio.NewScanner(multi)
				for scanner.Scan() {
					fmt.Printf("%s\t%s\n", color.HiBlueString("dart"), scanner.Text())
				}

				return c.Wait()
			}, func(err error) {
				color.HiYellow("stopping dart")
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

			addr := "localhost:3011"

			director := func(r *http.Request) {
				if redirect(r.URL.Path) {
					r.URL.Path = "/"
				}
				r.URL.Scheme = "http"
				r.URL.Host = addr
			}

			wsURL, _ := url.Parse("ws://" + addr)
			ws := websocketproxy.NewProxy(wsURL)

			server := &http.Server{
				Addr: ":3010",
				Handler: proxy{
					reverse: &httputil.ReverseProxy{
						Director: director,
					},
					websocket: ws,
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

type proxy struct {
	reverse   *httputil.ReverseProxy
	websocket *websocketproxy.WebsocketProxy
}

func (p proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Upgrade") == "websocket" {
		r.Header.Del("Origin")
		p.websocket.ServeHTTP(w, r)
		return
	}

	p.reverse.ServeHTTP(w, r)
}

func ensureUIContainer() error {
	name := "sourcepods-ui"
	args := []string{"run", "-d", "-p=3010:3010", "--name=" + name, "quay.io/sourcepods/ui:latest"}
	return ensureContainer(name, args)
}

func ensureContainer(name string, args []string) error {
	docker, err := exec.LookPath("docker")
	if err != nil {
		return err
	}

	cmd := exec.Command(docker, []string{
		"ps",
		"--filter", "name=" + name,
	}...)
	output, err := cmd.Output()
	if err != nil {
		return err
	}

	// If more than 2 lines, container exists
	lines := strings.Split(string(output), "\n")
	if len(lines) > 2 {
		return nil
	}

	color.Blue("docker %s", strings.Join(args, " "))

	cmd = exec.Command(docker, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Run(); err != nil {
		return err
	}

	return nil
}
