package storage

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type (
	Storage interface {
		Create(ctx context.Context, owner, name string) error
		SetDescription(ctx context.Context, owner, name, description string) error
		Tree(ctx context.Context, owner, name, branch string) ([]TreeObject, error)
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

func (s *storage) Create(ctx context.Context, owner, name string) error {
	dir := filepath.Join(s.root, owner, name)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to repository directory: %s", dir)
	}

	cmd := exec.CommandContext(ctx, s.git, "init", "--bare")
	cmd.Dir = dir
	return cmd.Run()
}

func (s *storage) SetDescription(ctx context.Context, owner, name, description string) error {
	file := filepath.Join(s.root, owner, name, "description")
	return ioutil.WriteFile(file, []byte(description+"\n"), 0644)
}

type TreeObject struct {
	Mode   string
	Type   string
	Object string
	File   string
}

func (s *storage) Tree(ctx context.Context, owner, name, branch string) ([]TreeObject, error) {
	var objects []TreeObject

	path := filepath.Join(s.root, owner, name)
	cmd := exec.CommandContext(ctx, s.git, "ls-tree", branch)
	cmd.Dir = path
	out, err := cmd.Output()
	if err != nil {
		return objects, err
	}

	scanner := bufio.NewScanner(bytes.NewBuffer(out))
	for scanner.Scan() {
		tabs := strings.Split(scanner.Text(), "\t")
		metas, file := tabs[0], tabs[1]
		meta := strings.Split(metas, " ")

		mode, typ, object := meta[0], meta[1], meta[2]

		objects = append(objects, TreeObject{
			Mode:   mode,
			Type:   typ,
			Object: object,
			File:   file,
		})
	}

	return objects, nil
}
