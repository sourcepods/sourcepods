package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

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
			cli.StringFlag{Name: "addr-ui", Usage: "The address to run the UI on", Value: ":3010"},
			cli.StringFlag{Name: "addr-api", Usage: "The address to run the API on", Value: ":3020"},
			cli.StringFlag{Name: "env", Usage: "Set the env gitpods runs in", Value: "development"},
			cli.StringFlag{Name: "log-level", Usage: "The log level to filter logs with before printing", Value: "debug"},
			cli.BoolFlag{Name: "setup", Usage: "Setup all dependencies needed for local development"},
			cli.BoolFlag{Name: "watch,w", Usage: "Watch files in this project and rebuild binaries if something changes"},
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
	envFlag := c.String("env")
	loglevelFlag := c.String("log-level")
	watch := c.Bool("watch")

	uiRunner := NewGitPodsRunner("ui", []string{
		fmt.Sprintf("GITPODS_ADDR=%s", uiAddrFlag),
		fmt.Sprintf("GITPODS_ADDR_API=%s", "http://localhost:3000/api"), // TODO
		fmt.Sprintf("GITPODS_ENV=%s", envFlag),
		fmt.Sprintf("GITPODS_LOGLEVEL=%s", loglevelFlag),
	})

	apiRunner := NewGitPodsRunner("api", []string{
		fmt.Sprintf("GITPODS_ADDR=%s", apiAddrFlag),
		fmt.Sprintf("GITPODS_ENV=%s", envFlag),
		fmt.Sprintf("GITPODS_LOGLEVEL=%s", loglevelFlag),
	})

	caddy := CaddyRunner{}

	//if watch {
	//	watcher := &FileWatcher{}
	//	watcher.Add(uiRunner, apiRunner)
	//
	//	go watcher.Watch()
	//}

	var g group.Group
	{
		g.Add(func() error {
			log.Println("starting ui")
			return uiRunner.Run()
		}, func(err error) {
			log.Println("stopping ui")
			uiRunner.Stop()
		})
	}
	{
		g.Add(func() error {
			log.Println("starting api")
			return apiRunner.Run()
		}, func(err error) {
			log.Println("stopping api")
			apiRunner.Stop()
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
