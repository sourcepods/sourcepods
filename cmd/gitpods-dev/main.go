package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Name = "gitpods-dev"
	app.Usage = "Runs gitpods on you local development machine"
	app.Flags = devFlags
	app.Action = devAction

	app.Commands = []cli.Command{{
		Name:   "setup",
		Usage:  "Sets up all the things necessary for developing gitpods",
		Action: devSetupAction,
	}}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
