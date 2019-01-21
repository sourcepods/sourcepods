package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
)

type BuildRestarter interface {
	Name() string
	Build() error
	Restart()
}

type FileWatcher struct {
	Restarters []BuildRestarter
}

func (w *FileWatcher) Add(r ...BuildRestarter) {
	w.Restarters = append(w.Restarters, r...)
}

func (w *FileWatcher) Watch() {
	// TODO: Find new/deleted Go files not only on start, but also while running
	files, err := w.findGoFiles()
	if err != nil {
		log.Println(err)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println(err)
	}
	defer watcher.Close()

	for _, file := range files {
		if err := watcher.Add(filepath.Join(".", file)); err != nil {
			log.Println(err)
		}
	}

	for {
		select {
		case event := <-watcher.Events:
			if event.Op != fsnotify.Chmod && event.Name != "" {
				for _, restarter := range w.Restarters {
					color.HiYellow("rebuilding %s\n", restarter.Name())
					if err := restarter.Build(); err == nil { // only notify and log if binary was created successfully.
						color.HiYellow("restarting %s", restarter.Name())
						restarter.Restart()
					}
				}
				watcher.Remove(event.Name)
				watcher.Add(event.Name)
			}
		case err := <-watcher.Errors:
			log.Println(err)
		}
	}
}

func (w *FileWatcher) findGoFiles() ([]string, error) {
	var files []string
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if strings.HasPrefix(path, "cmd/sourcepods-dev") { // don't watch gitpods itself
			return nil
		}
		if strings.HasPrefix(path, "vendor") {
			return nil
		}
		if strings.HasSuffix(path, ".go") {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}
