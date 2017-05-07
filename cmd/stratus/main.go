package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/oklog/oklog/pkg/group"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Name = "stratus"

	app.Commands = []cli.Command{{
		Name:   "dev",
		Usage:  "Runs gitpods on you local development machine",
		Action: actionDev,
		Flags: []cli.Flag{
			cli.StringFlag{Name: "addr-ui", Usage: "The address to run the UI on", Value: ":3000"},
			cli.StringFlag{Name: "addr-api", Usage: "The address to run the API on", Value: ":3010"},
			cli.StringFlag{Name: "env", Usage: "Set the env gitpods runs in", Value: "development"},
			cli.StringFlag{Name: "log-level", Usage: "The log level to filter logs with before printing", Value: "debug"},
		},
	}}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func actionDev(c *cli.Context) error {
	uiAddrFlag := c.String("addr-ui")
	apiAddrFlag := c.String("addr-api")
	envFlag := c.String("env")
	loglevelFlag := c.String("log-level")

	var g group.Group
	{
		webpack := &WebpackRunner{}
		g.Add(func() error {
			log.Println("starting webpack")
			return webpack.Run()
		}, func(err error) {
			log.Println("stopping webpack")
			webpack.Stop()
		})
	}
	{
		ui := &UIRunner{}
		env := []string{
			fmt.Sprintf("GITPODS_ADDR=%s", uiAddrFlag),
			fmt.Sprintf("GITPODS_ENV=%s", envFlag),
			fmt.Sprintf("GITPODS_LOGLEVEL=%s", loglevelFlag),
		}

		g.Add(func() error {
			log.Println("starting ui")
			return ui.Run(env)
		}, func(err error) {
			log.Println("stopping ui")
			ui.Stop()
		})
	}
	{
		ui := &APIRunner{}
		env := []string{
			fmt.Sprintf("GITPODS_ADDR=%s", apiAddrFlag),
			fmt.Sprintf("GITPODS_ENV=%s", envFlag),
			fmt.Sprintf("GITPODS_LOGLEVEL=%s", loglevelFlag),
		}

		g.Add(func() error {
			log.Println("starting api")
			return ui.Run(env)
		}, func(err error) {
			log.Println("stopping api")
			ui.Stop()
		})
	}

	//TODO: Add actor to the group that listens for system call to stop the group gracefully

	return g.Run()
}

type WebpackRunner struct {
	cmd *exec.Cmd
}

func (r *WebpackRunner) Run() error {
	file := "./webpack.config.js"
	_, err := os.Stat(file)
	if err != nil {
		// webpack config not found
		return nil
	}

	r.cmd = exec.Command(filepath.Join("node_modules", ".bin", "webpack"), "--watch")
	r.cmd.Stdin = os.Stdin
	r.cmd.Stdout = os.Stdout
	r.cmd.Stderr = os.Stderr

	return r.cmd.Run()
}

func (r *WebpackRunner) Stop() {
	if r.cmd == nil || r.cmd.Process == nil {
		return
	}
	r.cmd.Process.Kill()
}

type UIRunner struct {
	cmd *exec.Cmd
}

func (r *UIRunner) Run(env []string) error {
	//TODO: run build command if no binary available
	r.cmd = exec.Command("./dist/ui")

	fmt.Printf("api env: %+v\n", env)
	r.cmd.Env = env
	r.cmd.Stdin = os.Stdin
	r.cmd.Stdout = os.Stdout
	r.cmd.Stderr = os.Stderr
	return r.cmd.Run()
}

func (r *UIRunner) Stop() {
	if r.cmd == nil || r.cmd.Process == nil {
		return
	}
	r.cmd.Process.Kill()
}

type APIRunner struct {
	cmd *exec.Cmd
}

func (r *APIRunner) Run(env []string) error {
	//TODO: run build command if no binary available
	r.cmd = exec.Command("./dist/api")

	fmt.Printf("api env: %+v\n", env)
	r.cmd.Env = env
	r.cmd.Stdin = os.Stdin
	r.cmd.Stdout = os.Stdout
	r.cmd.Stderr = os.Stderr
	return r.cmd.Run()
}

func (r *APIRunner) Stop() {
	if r.cmd == nil || r.cmd.Process == nil {
		return
	}
	r.cmd.Process.Kill()
}
