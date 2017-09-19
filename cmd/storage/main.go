package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "gitpods-storage"
	app.Usage = "git in the cloud!"

	app.Action = storageAction
	app.Flags = storageFlags

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
