package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/oklog/oklog/pkg/group"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Name = "stratus"

	app.Commands = []cli.Command{{
		Name:   "dev",
		Usage:  "Runs gitpod in development mode",
		Action: actionDev,
		Flags: []cli.Flag{
			cli.IntFlag{Name: "port", Usage: "The port to run gitpod on", Value: 3000},
			cli.StringFlag{Name: "env", Usage: "Set the env gitpod runs in", Value: "development"},
		},
	}}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func actionDev(c *cli.Context) error {
	portFlag := c.Int("port")
	envFlag := c.Bool("env")

	env := []string{
		fmt.Sprintf("GITPOD_ADDR=:%d", portFlag),
		fmt.Sprintf("GITPOD_ENV=%s", envFlag),
	}

	var g group.Group
	g.Add(RunGitpod(env), func(err error) {
		if err != nil {
			log.Fatal(err)
			return
		}
	})

	g.Add(RunWebpack, func(err error) {
		if err != nil {
			log.Fatal(err)
			return
		}
	})

	return g.Run()
}

// RunGitpod runs a development server and restarts it with a new build if files change.
func RunGitpod(env []string) func() error {
	return func() error {
		builds := make(chan bool)

		go BuildForever(builds)

		go func() {
			if err := build(); err == nil {
				builds <- true
			}
		}()

		var cmd *exec.Cmd
		for {
			<-builds
			if cmd != nil {
				cmd.Process.Kill()
			}

			cmd = exec.Command("./dist/gitpod")
			go func() {
				cmd.Env = env
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				if err := cmd.Run(); err != nil {
					log.Println(err)
					return
				}
			}()
		}

		return nil
	}
}

// BuildForever watches the filesystem and builds a new binary if something changes.
// It notifies a channel that a build was created
func BuildForever(builds chan<- bool) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println(err)
		return
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op != fsnotify.Chmod && event.Name != "" {
					start := time.Now()
					if err := build(); err == nil { // only notify and log if binary was created successfully.
						log.Println("built a new binary in", time.Since(start))
						builds <- true // notify that a new build was successfully created
					}
					watcher.Remove(event.Name)
					watcher.Add(event.Name)
				}
			case err := <-watcher.Errors:
				log.Println(err)
			}
		}
	}()

	files, err := findGoFiles()
	if err != nil {
		log.Println(err)
	}

	for _, file := range files {
		if err := watcher.Add(filepath.Join(".", file)); err != nil {
			log.Println(err)
		}
	}

	select {}
}

func findGoFiles() ([]string, error) {
	var files []string
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if strings.HasPrefix(path, "cmd/stratus") { // don't watch stratus itself
			return nil
		}
		if strings.HasSuffix(path, ".go") {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func build() error {
	cmd := exec.Command("go", "build", "-v", "-i", "-o", "./dist/gitpod", "./cmd/gitpod")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

// RunWebpack starts webpack in the background and watches for changes
func RunWebpack() error {
	file := "./webpack.config.js"
	_, err := os.Stat(file)
	if err != nil {
		// webpack config not found
		return nil
	}

	cmd := exec.Command(filepath.Join("node_modules", ".bin", "webpack"), "--watch")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
