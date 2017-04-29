package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

const DefaultEnv = "production"

func main() {
	app := cli.NewApp()
	app.Name = "gitloud"
	app.Usage = "git flying loudly in the cloud!"

	app.Action = ActionWeb
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "addr",
			EnvVar: "GITLOUD_ADDR",
			Usage:  "The address gitloud runs on",
			Value:  ":3000",
		},
		cli.StringFlag{
			Name:   "env",
			EnvVar: "GITLOUD_ENV",
			Usage:  "The environment gitloud should run in",
			Value:  DefaultEnv,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
