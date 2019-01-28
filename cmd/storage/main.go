package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "sourcepods-storage"
	app.Usage = "git in the cloud!"

	app.Action = storageAction
	app.Flags = storageFlags

	app.Commands = cli.Commands{
		{
			Name:        "tree",
			Usage:       "Print the tree of a repository",
			Description: "Print the tree of a repository at a given ref with a path\n\tExample: tree sourcepods sourcepods master pkg/api",
			ArgsUsage:   "owner name ref path",
			Action:      treeAction,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
