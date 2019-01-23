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
	"strings"
)

var (
	// ErrRepoNotValid is returned for invalid repositories
	ErrRepoNotValid = fmt.Errorf("not a valid repository")
)

type (
	// Storage TODO: is something that should be split up
	Storage interface {
		Create(ctx context.Context, owner, name string) error
		GetRepository(ctx context.Context, owner, name string) (Repository, error)
		//Branches(ctx context.Context, owner string, name string) ([]Branch, error)
		//Commit(ctx context.Context, owner string, name string, rev string) (Commit, error)
	}

	// Repository is the interface for manipulating repos
	Repository interface {
		SetDescription(ctx context.Context, description string) error
		ListBranches(ctx context.Context) ([]Branch, error)
		GetCommit(ctx context.Context, rev string) (Commit, error)
	}

	// LocalRepository implements Repository for Local disk-access
	LocalRepository struct {
		git  string
		path string
	}

	// LocalStorage implements Storage for Local disk-access
	LocalStorage struct {
		git  string
		root string
	}
)

// SetDescription of repository
func (r *LocalRepository) SetDescription(ctx context.Context, description string) error {
	file := filepath.Join(r.path, "description")
	return ioutil.WriteFile(file, []byte(description+"\n"), 0644)
}

// NewLocalStorage returns a LocalStorage in the given `root`
func NewLocalStorage(root string) (*LocalStorage, error) {
	if err := os.MkdirAll(root, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage root: %s", root)
	}
	return &LocalStorage{
		git:  "/usr/bin/git",
		root: root,
	}, nil
}

// GetRepository from Storage
// TODO: Cache these somehow?
func (s *LocalStorage) GetRepository(ctx context.Context, owner, name string) (Repository, error) {
	dir := filepath.Join(s.root, owner, name)

	cmd := exec.CommandContext(ctx, s.git, "config", "--null", "core.repositoryformatversion")
	cmd.Dir = dir

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, ErrRepoNotValid
	}
	if strings.TrimSuffix(string(output), "\x00") != "0" {
		return nil, ErrRepoNotValid
	}
	// TODO: return an actual Repository...
	return &LocalRepository{git: s.git, path: dir}, nil
}

// Create a new repository
func (s *LocalStorage) Create(ctx context.Context, owner, name string) error {
	dir := filepath.Join(s.root, owner, name)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to repository directory: %s", dir)
	}

	cmd := exec.CommandContext(ctx, s.git, "init", "--bare")
	cmd.Dir = dir
	return cmd.Run()
}

// Branch of a repository
type Branch struct {
	Name string
	Sha1 string
	Type string
}

// ListBranches returns all branches of a given repository
func (r *LocalRepository) ListBranches(ctx context.Context) ([]Branch, error) {
	args := []string{"for-each-ref", "--format=%(objectname) %(objecttype) %(refname)", "refs/heads"}
	cmd := exec.CommandContext(ctx, r.git, args...)
	cmd.Dir = r.path
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

// Commit holds a commit
type Commit struct {
	Hash    string
	Tree    string
	Parent  string
	Message string
	Body    string

	Author Signature

	Committer Signature
}

// GetCommit returns a single Commit from a Repository
func (r *LocalRepository) GetCommit(ctx context.Context, rev string) (Commit, error) {
	args := []string{"cat-file", "-p", rev}
	cmd := exec.CommandContext(ctx, r.git, args...)
	cmd.Dir = r.path
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

		c.Author, err = parseSignature(line)
		if err != nil {
			// This should probably just error out, and not return a partial commit...
			return false, err
		}
		return false, nil
	}

	if strings.HasPrefix(line, committerPrefix) {
		var err error
		line := strings.TrimPrefix(line, committerPrefix)

		c.Committer, err = parseSignature(line)
		if err != nil {
			// This should probably just error out, and not return a partial commit...
			return false, err
		}

		return false, nil
	}

	// skip any excessive header-lines
	return false, nil
}
