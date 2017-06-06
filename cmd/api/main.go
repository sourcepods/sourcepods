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

	app.Commands = []cli.Command{{
		Name:  "db",
		Usage: "Run actions on the database",
		Subcommands: []cli.Command{{
			Name:   "migrate",
			Flags:  dbFlags,
			Action: dbMigrateAction,
		}, {
			Name:   "reset",
			Flags:  dbFlags,
			Action: dbResetAction,
		}},
	}}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
