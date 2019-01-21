package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Name = "sourcepods-dev"
	app.Usage = "Runs SourcePods on you local development machine"
	app.Flags = devFlags
	app.Action = devAction

	app.Commands = []cli.Command{{
		Name:   "setup",
		Usage:  "Sets up all the things necessary for developing SourcePods",
		Action: devSetupAction,
	}}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
