package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "sourcepods-ssh"
	app.Usage = "git in the cloud!"

	app.Action = sshAction
	app.Flags = sshFlags

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
