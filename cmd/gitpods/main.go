package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/gitpods/gitpods/cmd"
	"github.com/oklog/oklog/pkg/group"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Name = "gitpods"

	app.Commands = []cli.Command{{
		Name:   "dev",
		Usage:  "Runs gitpods on you local development machine",
		Action: actionDev,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "addr-ui",
				Usage: "The address to run the UI on",
				Value: ":3010",
			},
			cli.StringFlag{
				Name:  "addr-api",
				Usage: "The address to run the API on",
				Value: ":3020",
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
				Name:  "setup",
				Usage: "Setup all dependencies needed for local development",
			},
			cli.BoolFlag{
				Name:  "watch,w",
				Usage: "Watch files in this project and rebuild binaries if something changes",
			},
		},
	}}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func actionDev(c *cli.Context) error {
	if c.Bool("setup") {
		return ActionDevSetup(c)
	}

	uiAddrFlag := c.String("addr-ui")
	apiAddrFlag := c.String("addr-api")
	databaseDriver := c.String("database-driver")
	databaseDSN := c.String("database-dsn")
	loglevelFlag := c.String("log-level")
	watch := c.Bool("watch")

	uiRunner := NewGitPodsRunner("ui", []string{
		fmt.Sprintf("%s=%s", cmd.EnvAddr, uiAddrFlag),
		fmt.Sprintf("%s=%s", cmd.EnvAddrAPI, "http://localhost:3000/api"), // TODO
		fmt.Sprintf("%s=%s", cmd.EnvLogLevel, loglevelFlag),
	})

	apiRunner := NewGitPodsRunner("api", []string{
		fmt.Sprintf("%s=%s", cmd.EnvAddr, apiAddrFlag),
		fmt.Sprintf("%s=%s", cmd.EnvDatabaseDriver, databaseDriver),
		fmt.Sprintf("%s=%s", cmd.EnvDatabaseDSN, databaseDSN),
		fmt.Sprintf("%s=%s", cmd.EnvLogLevel, loglevelFlag),
		fmt.Sprintf("%s=%s", cmd.EnvSecret, "secret"),
	})

	caddy := CaddyRunner{}

	if watch {
		watcher := &FileWatcher{}
		watcher.Add(uiRunner, apiRunner)

		go watcher.Watch()
	}

	var g group.Group
	{
		g.Add(func() error {
			log.Println("starting ui")
			return uiRunner.Run()
		}, func(err error) {
			log.Println("stopping ui")
			uiRunner.Shutdown()
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
			log.Println("starting caddy")
			return caddy.Run()
		}, func(err error) {
			log.Println("stopping caddy")
			caddy.Stop()
		})
	}
	{
		g.Add(func() error {
			stop := make(chan os.Signal, 1)
			signal.Notify(stop, os.Interrupt)
			<-stop
			return nil
		}, func(err error) {
		})
	}

	webpack := &WebpackRunner{}
	if watch {
		g.Add(func() error {
			log.Println("starting webpack")
			return webpack.Run(true)
		}, func(err error) {
			log.Println("stopping webpack")
			webpack.Stop()
		})
	} else {
		webpack.Run(false)
	}

	return g.Run()
}
