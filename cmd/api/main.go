package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "gitpods-api"
	app.Usage = "git flying loudly in the cloud!"

	app.Action = ActionAPI
	app.Flags = FlagsAPI

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
