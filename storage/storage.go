package storage

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

type (
	Storage interface {
		Create(owner, name string) error
		Description(owner, name, description string) error
	}

	storage struct {
		git  string
		root string
	}
)

func NewStorage(root string) (Storage, error) {
	if err := os.MkdirAll(root, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage root: %s", root)
	}
	return &storage{
		git:  "/usr/bin/git",
		root: root,
	}, nil
}

func (s *storage) Create(owner, name string) error {
	dir := filepath.Join(s.root, owner, name)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to repository directory: %s", dir)
	}

	cmd := exec.CommandContext(context.Background(), s.git, "init", "--bare")
	cmd.Dir = dir
	return cmd.Run()
}

func (s *storage) Description(owner, name, description string) error {
	file := filepath.Join(s.root, owner, name, "description")
	return ioutil.WriteFile(file, []byte(description+"\n"), 0644)
}
