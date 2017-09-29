package storage

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	authorLine    = regexp.MustCompile(`(.*)\s\<(.*)\>\s(\d*)`)
	committerLine = authorLine
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

	Commit Commit
}

type Commit struct {
	Hash    string
	Tree    string
	Parent  string
	Subject string

	Author      string
	AuthorEmail string
	AuthorDate  time.Time

	Committer      string
	CommitterEmail string
	CommitterDate  time.Time
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

	wg := sync.WaitGroup{}
	for i := range objects {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			commitHash, err := s.commitHash(ctx, owner, name, branch, objects[i].File)
			if err != nil {
				log.Println(err)
				return
			}

			commit, err := s.commit(ctx, owner, name, commitHash)
			if err != nil {
				log.Println(err)
				return
			}

			objects[i].Commit = commit
		}(i)
	}
	wg.Wait()

	return objects, nil
}

func (s *storage) commitHash(ctx context.Context, owner, name, rev, file string) (string, error) {
	//span, ctx := opentracing.StartSpanFromContext(ctx, "storage.Storage.commitHash")
	//span.SetTag("owner", owner)
	//span.SetTag("name", name)
	//span.SetTag("rev", rev)
	//span.SetTag("file", file)
	//defer span.Finish()

	args := []string{"log", "-1", "--pretty=%H", rev, "--", file}
	cmd := exec.CommandContext(ctx, s.git, args...)
	cmd.Dir = filepath.Join(s.root, owner, name)
	out, err := cmd.Output()

	return strings.TrimSpace(string(out)), err
}

func (s *storage) commit(ctx context.Context, owner, name, hash string) (Commit, error) {
	//span, ctx := opentracing.StartSpanFromContext(ctx, "storage.Storage.commit")
	//span.SetTag("owner", owner)
	//span.SetTag("name", name)
	//span.SetTag("hash", hash)
	//defer span.Finish()

	args := []string{"cat-file", "-p", hash}
	cmd := exec.CommandContext(ctx, s.git, args...)
	cmd.Dir = filepath.Join(s.root, owner, name)
	out, err := cmd.Output()
	if err != nil {
		return Commit{}, err
	}

	commit, err := parseCommit(bytes.NewBuffer(out), hash)
	if err != nil {
		return commit, err
	}

	return commit, nil
}

func parseCommit(r io.Reader, hash string) (Commit, error) {
	scanner := bufio.NewScanner(r)

	c := Commit{
		Hash: hash,
	}

	// if true then the following lines are the subject
	subject := false

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			subject = true
			continue
		}

		if subject {
			c.Subject = line
		}

		if strings.HasPrefix(line, "tree ") {
			c.Tree = strings.TrimPrefix(line, "tree ")
		} else if strings.HasPrefix(line, "parent ") {
			c.Parent = strings.TrimPrefix(line, "parent ")
		} else if strings.HasPrefix(line, "author ") {
			line := strings.TrimPrefix(line, "author ")
			author := authorLine.FindStringSubmatch(line)

			t, err := strconv.ParseInt(author[3], 10, 64)
			if err != nil {
				return c, err
			}

			c.Author = author[1]
			c.AuthorEmail = author[2]
			c.AuthorDate = time.Unix(t, 0)
		} else if strings.HasPrefix(line, "committer ") {
			line := strings.TrimPrefix(line, "committer ")
			committer := committerLine.FindStringSubmatch(line)

			t, err := strconv.ParseInt(committer[3], 10, 64)
			if err != nil {
				return c, err
			}

			c.Committer = committer[1]
			c.CommitterEmail = committer[2]
			c.CommitterDate = time.Unix(t, 0)
		}
	}

	return c, nil
}
