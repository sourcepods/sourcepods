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
	authorLine = regexp.MustCompile(`(.*)\s\<(.*)\>\s(\d*)`)
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

//Branch of a repository
type Branch struct {
	Name string
	Sha1 string
	Type string
}

//Branches returns all branches of a given repository
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

type Author struct {
	Name  string
	Email string
	Date  time.Time
}

func parseAuthor(line string) (Author, error) {
	committer := authorLine.FindStringSubmatch(line)

	t, err := strconv.ParseInt(committer[3], 10, 64)
	if err != nil {
		return Author{}, nil
	}

	return Author{
		Name:  committer[1],
		Email: committer[2],
		Date:  time.Unix(t, 0),
	}, nil
}

type Commit struct {
	Hash    string
	Tree    string
	Parent  string
	Message string
	Body    string

	Author Author

	Committer Author
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
	hasHeader := false

	for scanner.Scan() {
		line := scanner.Text()

		if !hasHeader {
			var err error
			hasHeader, err = parseCommitHeader(&c, line)
			if err != nil {
				return c, err
			}
			continue
		}

		if c.Message == "" {
			c.Message = line
			continue
		}

		// TODO: use a string.Builder instead?
		// TODO: handle commit signatures
		c.Body = fmt.Sprintf("%s\n%s", c.Body, line)
	}

	// Trim excessive stringing new lines
	c.Body = strings.TrimLeft(c.Body, "\n")

	return c, nil
}

// returns true when it's passed the header
func parseCommitHeader(c *Commit, line string) (bool, error) {
	const (
		treePrefix      = "tree "
		parentPrefix    = "parent "
		authorPrefix    = "author "
		committerPrefix = "committer "
	)

	if line == "" {
		return true, nil
	}
	if strings.HasPrefix(line, treePrefix) {
		c.Tree = strings.TrimPrefix(line, treePrefix)
		return false, nil
	}

	if strings.HasPrefix(line, parentPrefix) {
		c.Parent = strings.TrimPrefix(line, parentPrefix)
		return false, nil
	}

	if strings.HasPrefix(line, authorPrefix) {
		var err error
		line := strings.TrimPrefix(line, authorPrefix)

		c.Author, err = parseAuthor(line)
		if err != nil {
			// This should probably just error out, and not return a partial commit...
			return false, err
		}
		return false, nil
	}

	if strings.HasPrefix(line, committerPrefix) {
		var err error
		line := strings.TrimPrefix(line, committerPrefix)

		c.Committer, err = parseAuthor(line)
		if err != nil {
			// This should probably just error out, and not return a partial commit...
			return false, err
		}

		return false, nil
	}

	// skip any excessive header-lines
	return false, nil
}
