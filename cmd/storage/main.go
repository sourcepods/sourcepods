package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "gitpoids-storage"
	app.Usage = "git int the cloud!"

	app.Action = storageAction
	app.Flags = storageFlags

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
