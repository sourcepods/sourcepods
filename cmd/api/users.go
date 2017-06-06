package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	"github.com/gitpods/gitpods/cmd"
	"github.com/gitpods/gitpods/user"
	"github.com/urfave/cli"
)

const (
	usersEmail    = "email"
	usersUsername = "username"
	usersName     = "name"
	usersPassword = "password"
)

type usersConf struct {
	DatabaseDriver string
	DatabaseDSN    string
}

var (
	usersConfig = &usersConf{}

	usersFlags = []cli.Flag{
		cli.StringFlag{
			Name:        cmd.FlagDatabaseDriver,
			EnvVar:      cmd.EnvDatabaseDriver,
			Usage:       "The database driver to use: postgres",
			Value:       "postgres",
			Destination: &usersConfig.DatabaseDriver,
		},
		cli.StringFlag{
			Name:        cmd.FlagDatabaseDSN,
			EnvVar:      cmd.EnvDatabaseDSN,
			Usage:       "The database connection data",
			Destination: &usersConfig.DatabaseDSN,
		},
	}
)

func usersCreateAction(c *cli.Context) error {
	email := c.String(usersEmail)
	username := c.String(usersUsername)
	name := c.String(usersName)
	password := c.String(usersPassword)

	//
	// Stores
	//
	var users user.Store

	switch apiConfig.DatabaseDriver {
	default:
		db, err := sql.Open("postgres", apiConfig.DatabaseDSN)
		if err != nil {
			return err
		}
		defer db.Close()

		users = user.NewPostgresStore(db)
	}

	u := &user.User{
		Email:    email,
		Username: username,
		Name:     name,
		Password: password,
	}

	errs := user.ValidateCreate(u)
	if len(errs) > 0 {
		for _, err := range errs {
			fmt.Println(err)
		}
		os.Exit(1)
	}

	u, err := users.Create(u)
	if err != nil {
		return err
	}

	u.Password = "hidden"

	data, err := json.MarshalIndent(u, "", "	")
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", string(data))

	return nil
}
