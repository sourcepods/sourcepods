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
	}, {
		Name:  "users",
		Usage: "Manage users",
		Flags: usersFlags,
		Subcommands: []cli.Command{{
			Name:   "create",
			Action: usersCreateAction,
			Flags: []cli.Flag{
				cli.StringFlag{Name: usersEmail},
				cli.StringFlag{Name: usersUsername},
				cli.StringFlag{Name: usersName},
				cli.StringFlag{Name: usersPassword},
			},
		}},
	}}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
