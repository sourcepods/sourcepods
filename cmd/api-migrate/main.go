package main

import (
	"log"
	"os"

	"github.com/gitpods/gitpods/cmd"
	"github.com/mattes/migrate"
	_ "github.com/mattes/migrate/database/postgres"
	_ "github.com/mattes/migrate/source/file"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "gitpods-api"
	app.Usage = "git in the cloud!"

	app.Action = migrateAction
	app.Flags = migrateFlags

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

type migrateConf struct {
	DatabaseDriver string
	DatabaseDSN    string
	MigrationsPath string
}

var (
	migrateConfig = migrateConf{}

	migrateFlags = []cli.Flag{
		cli.StringFlag{
			Name:        cmd.FlagDatabaseDriver,
			EnvVar:      cmd.EnvDatabaseDriver,
			Usage:       "The database driver to use: memory & postgres",
			Value:       "postgres",
			Destination: &migrateConfig.DatabaseDriver,
		},
		cli.StringFlag{
			Name:        cmd.FlagDatabaseDSN,
			EnvVar:      cmd.EnvDatabaseDSN,
			Usage:       "The database connection data",
			Destination: &migrateConfig.DatabaseDSN,
		},
		cli.StringFlag{
			Name:        cmd.FlagMigrationsPath,
			EnvVar:      cmd.EnvMigrationsPath,
			Usage:       "The path to the folder containing all migrations",
			Destination: &migrateConfig.MigrationsPath,
		},
	}
)

func migrateAction(c *cli.Context) error {
	m, err := migrate.New("file://"+migrateConfig.MigrationsPath, migrateConfig.DatabaseDSN)
	if err != nil {
		return err
	}

	return m.Up()
}
