package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "gitpods-api"
	app.Usage = "git in the cloud!"

	app.Action = apiAction
	app.Flags = apiFlags

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
