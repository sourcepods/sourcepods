package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Name = "sourcepods-dev"
	app.Usage = "Runs sourcepods on you local development machine"
	app.Flags = devFlags
	app.Action = devAction

	app.Commands = []cli.Command{{
		Name:   "setup",
		Usage:  "Sets up all the things necessary for developing sourcepods",
		Action: devSetupAction,
	}}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
