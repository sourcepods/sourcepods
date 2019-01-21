package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

// Runner builds, runs, stops and restarts SourcePods.
type Runner struct {
	name    string
	env     []string
	cmd     *exec.Cmd
	restart chan bool
}

func NewRunner(name string, env []string) *Runner {
	return &Runner{
		name:    name,
		env:     env,
		restart: make(chan bool, 16),
	}
}

func (r *Runner) Name() string {
	return r.name
}

func (r *Runner) Run() error {
	if err := r.Build(); err == nil {
		r.restart <- true
	}

	for {
		_, more := <-r.restart
		if more {
			if r.cmd != nil {
				r.Stop()
			}

			go func() {
				r.cmd = exec.Command("./dev/" + r.name)
				r.cmd.Env = r.env
				stdout, err := r.cmd.StdoutPipe()
				if err != nil {
					return
				}
				stderr, err := r.cmd.StderrPipe()
				if err != nil {
					return
				}

				multi := io.MultiReader(stdout, stderr)

				if r.cmd.Start() != nil {
					return
				}

				scanner := bufio.NewScanner(multi)

				for scanner.Scan() {
					fmt.Printf("%s\t%s\n", color.HiBlueString(r.name), scanner.Text())
				}

				if err = r.cmd.Wait(); err != nil {
					return
				}
			}()
		} else {
			return nil
		}
	}
}

func (r *Runner) Stop() {
	if r.cmd == nil || r.cmd.Process == nil {
		return
	}
	r.cmd.Process.Kill()
}

func (r *Runner) Build() error {
	cmd := exec.Command("make", "dev/"+r.name)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println(strings.Join(cmd.Args, " "))

	return cmd.Run()
}

func (r *Runner) Restart() {
	r.restart <- true
}

func (r *Runner) Shutdown() {
	close(r.restart)
	r.Stop()
}

// CaddyRunner runs caddy
type CaddyRunner struct {
	cmd *exec.Cmd
}

func (r *CaddyRunner) Run() error {
	r.cmd = exec.Command(filepath.Join(".", "dev", "caddy"), "-conf", "./dev/Caddyfile")
	stdout, err := r.cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := r.cmd.StderrPipe()
	if err != nil {
		return err
	}

	multi := io.MultiReader(stdout, stderr)

	if err := r.cmd.Start(); err != nil {
		return err
	}

	scanner := bufio.NewScanner(multi)

	for scanner.Scan() {
		fmt.Printf("%s\t%s\n", color.HiBlueString("caddy"), scanner.Text())
	}

	return r.cmd.Wait()
}

func (r *CaddyRunner) Stop() {
	if r.cmd == nil || r.cmd.Process == nil {
		return
	}
	r.cmd.Process.Kill()
}
