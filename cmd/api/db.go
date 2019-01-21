package main

import (
	"log"
	"path/filepath"

	"github.com/mattes/migrate"
	_ "github.com/mattes/migrate/database/postgres"
	_ "github.com/mattes/migrate/source/file"
	"github.com/sourcepods/sourcepods/cmd"
	"github.com/urfave/cli"
)

type dbConf struct {
	DatabaseDriver string
	DatabaseDSN    string
	MigrationsPath string
}

var (
	dbConfig = dbConf{}

	dbFlags = []cli.Flag{
		cli.StringFlag{
			Name:        cmd.FlagDatabaseDriver,
			EnvVar:      cmd.EnvDatabaseDriver,
			Usage:       "The database driver to use: memory & postgres",
			Value:       "postgres",
			Destination: &dbConfig.DatabaseDriver,
		},
		cli.StringFlag{
			Name:        cmd.FlagDatabaseDSN,
			EnvVar:      cmd.EnvDatabaseDSN,
			Usage:       "The database connection data",
			Destination: &dbConfig.DatabaseDSN,
		},
		cli.StringFlag{
			Name:        cmd.FlagMigrationsPath,
			EnvVar:      cmd.EnvMigrationsPath,
			Usage:       "The path to the folder containing all migrations",
			Destination: &dbConfig.MigrationsPath,
		},
	}
)

func dbMigrateAction(c *cli.Context) error {
	log.Println("start migration...")

	m, err := newMigrate(dbConfig.MigrationsPath, dbConfig.DatabaseDSN)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil {
		return err
	}

	log.Println("migration successful")
	return nil
}

func dbResetAction(c *cli.Context) error {
	log.Println("start resetting...")

	m, err := newMigrate(dbConfig.MigrationsPath, dbConfig.DatabaseDSN)
	if err != nil {
		return err
	}

	if err := m.Down(); err != nil {
		return err
	}

	if err := m.Up(); err != nil {
		return err
	}

	log.Println("reset successful")
	return nil
}

func newMigrate(path string, dsn string) (*migrate.Migrate, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	return migrate.New("file://"+path, dsn)
}
