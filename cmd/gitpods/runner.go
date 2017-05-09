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
		restart: make(chan bool, 1),
	}
}

func (r *GitPodsRunner) Name() string {
	return r.name
}

func (r *GitPodsRunner) Run() error {
	file := "./dev/" + r.name
	_, err := os.Stat(file)
	if err != nil {
		if err := r.Build(); err == nil {
			r.restart <- true
		}
	}

	//// Enter the first for iteration to start the services
	//r.restart <- true

	var cmd *exec.Cmd
	for {
		//<-r.restart

		if cmd != nil {
			r.Stop()
		}

		r.cmd = exec.Command("./dev/" + r.name)
		r.cmd.Env = r.env
		r.cmd.Stdin = os.Stdin
		r.cmd.Stdout = os.Stdout
		r.cmd.Stderr = os.Stderr
		return r.cmd.Run()
	}

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

func (r GitPodsRunner) Restart() {
	r.restart <- true
}

// WebpackRunner runs webpack either ones or watches the files.
type WebpackRunner struct {
	cmd *exec.Cmd
}

func (r *WebpackRunner) Run(watch bool) error {
	file := "./webpack.config.js"
	_, err := os.Stat(file)
	if err != nil {
		// webpack config not found
		return nil
	}

	args := []string{}
	if watch {
		args = []string{"--watch"}
	}

	r.cmd = exec.Command(filepath.Join("node_modules", ".bin", "webpack"), args...)
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
