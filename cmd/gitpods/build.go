package main

import (
	"github.com/urfave/cli"
)

func buildAction(c *cli.Context) error {
	apiRunner := NewGitPodsRunner("api", []string{})
	uiRunner := NewGitPodsRunner("ui", []string{})
	webpack := &WebpackRunner{}

	if err := apiRunner.Build(); err != nil {
		return err
	}

	if err := uiRunner.Build(); err != nil {
		return err
	}

	if err := webpack.Run(false); err != nil {
		return err
	}

	return nil
}
