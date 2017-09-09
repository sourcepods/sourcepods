package main

import (
	"os"
	"os/exec"
	"path/filepath"
)

// GitPodsRunner builds, runs, stops and restarts GitPods.
type GitPodsRunner struct {
	name    string
	env     []string
	cmd     *exec.Cmd
	restart chan bool
}

func NewGitPodsRunner(name string, env []string) *GitPodsRunner {
	return &GitPodsRunner{
		name:    name,
		env:     env,
		restart: make(chan bool, 16),
	}
}

func (r *GitPodsRunner) Name() string {
	return r.name
}

func (r *GitPodsRunner) Run() error {
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
				r.cmd.Stdin = os.Stdin
				r.cmd.Stdout = os.Stdout
				r.cmd.Stderr = os.Stderr
				r.cmd.Run()
			}()
		} else {
			return nil
		}
	}
	return nil
}

func (r *GitPodsRunner) Stop() {
	if r.cmd == nil || r.cmd.Process == nil {
		return
	}
	r.cmd.Process.Kill()
}

func (r *GitPodsRunner) Build() error {
	cmd := exec.Command("go", "build", "-v", "-i", "-o", "./dev/"+r.name, "./cmd/"+r.name)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (r *GitPodsRunner) Restart() {
	r.restart <- true
}

func (r *GitPodsRunner) Shutdown() {
	close(r.restart)
	r.Stop()
}

// CaddyRunner runs caddy
type CaddyRunner struct {
	cmd *exec.Cmd
}

func (r *CaddyRunner) Run() error {
	r.cmd = exec.Command(filepath.Join(".", "dev", "caddy"), "-conf", "./dev/Caddyfile")
	r.cmd.Stdin = os.Stdin
	r.cmd.Stdout = os.Stdout
	r.cmd.Stderr = os.Stderr

	return r.cmd.Run()
}

func (r *CaddyRunner) Stop() {
	if r.cmd == nil || r.cmd.Process == nil {
		return
	}
	r.cmd.Process.Kill()
}
