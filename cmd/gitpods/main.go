package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "gitpod"
	app.Usage = "git flying loudly in the cloud!"

	app.Action = ActionWeb
	app.Flags = FlagsWeb

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
