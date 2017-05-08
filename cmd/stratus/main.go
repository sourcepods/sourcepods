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
		ui := &GitPodsRunner{}
		env := []string{
			fmt.Sprintf("GITPODS_ADDR=%s", uiAddrFlag),
			fmt.Sprintf("GITPODS_ADDR_API=%s", apiAddrFlag),
			fmt.Sprintf("GITPODS_ENV=%s", envFlag),
			fmt.Sprintf("GITPODS_LOGLEVEL=%s", loglevelFlag),
		}

		g.Add(func() error {
			log.Println("starting ui")
			return ui.Run("ui", env)
		}, func(err error) {
			log.Println("stopping ui")
			ui.Stop()
		})
	}
	{
		ui := &GitPodsRunner{}
		env := []string{
			fmt.Sprintf("GITPODS_ADDR=%s", apiAddrFlag),
			fmt.Sprintf("GITPODS_ENV=%s", envFlag),
			fmt.Sprintf("GITPODS_LOGLEVEL=%s", loglevelFlag),
		}

		g.Add(func() error {
			log.Println("starting api")
			return ui.Run("api", env)
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

type GitPodsRunner struct {
	cmd *exec.Cmd
}

func (r *GitPodsRunner) Run(name string, env []string) error {
	file := "./dist/" + name
	_, err := os.Stat(file)
	if err != nil {
		if err := build(name); err != nil {
			return err
		}
	}

	r.cmd = exec.Command("./dist/" + name)

	r.cmd.Env = env
	r.cmd.Stdin = os.Stdin
	r.cmd.Stdout = os.Stdout
	r.cmd.Stderr = os.Stderr
	return r.cmd.Run()
}

func (r *GitPodsRunner) Stop() {
	if r.cmd == nil || r.cmd.Process == nil {
		return
	}
	r.cmd.Process.Kill()
}

func build(name string) error {
	cmd := exec.Command("go", "build", "-v", "-i", "-o", "./dist/"+name, "./cmd/"+name)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

//// RunAPI runs a development server and restarts it with a new build if files change.
//func RunAPI(env []string) func() error {
//	return func() error {
//		builds := make(chan bool)
//
//		go BuildForever(builds)
//
//		go func() {
//			if err := build(); err == nil {
//				builds <- true
//			}
//		}()
//
//		var cmd *exec.Cmd
//		for {
//			<-builds
//			if cmd != nil {
//				cmd.Process.Kill()
//			}
//
//			cmd = exec.Command("./dist/gitpods")
//			go func() {
//				cmd.Env = env
//				cmd.Stdin = os.Stdin
//				cmd.Stdout = os.Stdout
//				cmd.Stderr = os.Stderr
//				if err := cmd.Run(); err != nil {
//					log.Println(err)
//					return
//				}
//			}()
//		}
//
//		return nil
//	}
//}
