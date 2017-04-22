package main

import (
	"log"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/urfave/cli"
)

const DefaultEnv = "production"

type Config struct {
	Addr string `arg:"env:ADDR" json:"addr"`
	Env  string `arg:"env:ENV" json:"env"`
}

func main() {
	config := Config{
		Addr: ":3000",
		Env:  DefaultEnv,
	}
	arg.MustParse(&config)

	app := cli.NewApp()
	app.Name = "gitloud"
	app.Usage = "git flying loudly in the cloud!"

	app.Action = ActionWeb(config)

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
