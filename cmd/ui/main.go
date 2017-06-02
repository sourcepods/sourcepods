package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "gitpods-ui"
	app.Usage = "git in the cloud!"

	app.Action = ActionUI
	app.Flags = uiFlags

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
