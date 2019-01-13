package storage

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
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
		Branches(ctx context.Context, owner string, name string) ([]Branch, error)
		Commit(ctx context.Context, owner string, name string, rev string) (Commit, error)
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

// Branches
type Branch struct {
	Name string
	Sha1 string
	Type string
}

func (s *storage) Branches(ctx context.Context, owner string, name string) ([]Branch, error) {
	args := []string{"for-each-ref", "--format=%(objectname) %(objecttype) %(refname)", "refs/heads"}
	cmd := exec.CommandContext(ctx, s.git, args...)
	cmd.Dir = filepath.Join(s.root, owner, name)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var bs []Branch
	scanner := bufio.NewScanner(bytes.NewBuffer(out))
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), " ")

		bs = append(bs, Branch{
			Name: strings.TrimPrefix(s[2], "refs/heads/"),
			Sha1: s[0],
			Type: s[1],
		})
	}

	return bs, nil
}

type Commit struct {
	Hash    string
	Tree    string
	Parent  string
	Message string

	Author      string
	AuthorEmail string
	AuthorDate  time.Time

	Committer      string
	CommitterEmail string
	CommitterDate  time.Time
}

func (s *storage) Commit(ctx context.Context, owner string, name string, rev string) (Commit, error) {
	args := []string{"cat-file", "-p", rev}
	cmd := exec.CommandContext(ctx, s.git, args...)
	cmd.Dir = filepath.Join(s.root, owner, name)
	out, err := cmd.Output()
	if err != nil {
		return Commit{}, err
	}

	commit, err := parseCommit(bytes.NewBuffer(out), rev)
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
			c.Message = line
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
