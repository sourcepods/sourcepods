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
		Tree(ctx context.Context, owner, name, branch string, recursive bool) ([]TreeObject, error)
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

func (s *storage) Tree(ctx context.Context, owner, name, branch string, recursive bool) ([]TreeObject, error) {
	var objects []TreeObject

	args := []string{"ls-tree", branch}
	if recursive {
		args = []string{"ls-tree", "-r", branch}
	}

	path := filepath.Join(s.root, owner, name)
	cmd := exec.CommandContext(ctx, s.git, args...)
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

	indexes := make(chan int, 64)
	done := make(chan int, 64)

	for w := 0; w < 4; w++ {
		go func() {
			for index := range indexes {
				commitHash, err := s.commitHash(ctx, owner, name, branch, objects[index].File)
				if err != nil {
					log.Println(err)
					return
				}

				commit, err := s.commit(ctx, owner, name, commitHash)
				if err != nil {
					log.Println(err)
					return
				}

				objects[index].Commit = commit
				done <- index
			}
		}()
	}

	for i := 0; i < len(objects); i++ {
		indexes <- i
	}
	close(indexes)

	for i := 0; i < len(objects); i++ {
		<-done
	}
	close(done)

	return objects, nil
}

func (s *storage) commitHash(ctx context.Context, owner, name, rev, file string) (string, error) {
	args := []string{"log", "-1", "--pretty=%H", rev, "--", file}
	cmd := exec.CommandContext(ctx, s.git, args...)
	cmd.Dir = filepath.Join(s.root, owner, name)
	out, err := cmd.Output()

	return strings.TrimSpace(string(out)), err
}

func (s *storage) commit(ctx context.Context, owner, name, hash string) (Commit, error) {
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
