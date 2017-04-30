package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

const ProductionEnv = "production"

func main() {
	app := cli.NewApp()
	app.Name = "gitpod"
	app.Usage = "git flying loudly in the cloud!"

	app.Action = ActionWeb
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "addr",
			EnvVar: "GITPOD_ADDR",
			Usage:  "The address gitpod runs on",
			Value:  ":3000",
		},
		cli.StringFlag{
			Name:   "env",
			EnvVar: "GITPOD_ENV",
			Usage:  "The environment gitpod should run in",
			Value:  ProductionEnv,
		},
		cli.StringFlag{
			Name:   "loglevel",
			EnvVar: "GITPOD_LOGLEVEL",
			Usage:  "The log level to filter logs with before printing",
			Value:  "info",
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
