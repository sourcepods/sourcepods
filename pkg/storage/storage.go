package storage

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

var (
	// ErrRepoNotValid is returned for invalid repositories
	ErrRepoNotValid = fmt.Errorf("not a valid repository")
)

type (
	// Storage TODO: is something that should be split up
	Storage interface {
		Create(ctx context.Context, id string) error
		GetRepository(ctx context.Context, id string) (Repository, error)
	}

	StorageOption func(Storage)

	// LocalStorage implements Storage for Local disk-access
	LocalStorage struct {
		git    string
		root   string
		logger log.Logger
	}

	// Repository is the interface for manipulating repos
	Repository interface {
		SetDescription(ctx context.Context, description string) error
		ListBranches(ctx context.Context) ([]Branch, error)
		GetCommit(ctx context.Context, ref string) (Commit, error)
		Tree(ctx context.Context, ref, path string) ([]TreeEntry, error)
	}

	// LocalRepository implements Repository for Local disk-access
	LocalRepository struct {
		git  string
		path string
	}
)

// LoggerOption injects a logger into LocalStorage
func LoggerOption(logger log.Logger) StorageOption {
	return func(s Storage) {
		ls, ok := s.(*LocalStorage)
		if !ok {
			return
		}
		ls.logger = logger
	}
}

// NewLocalStorage returns a LocalStorage in the given `root`
func NewLocalStorage(root string, opts ...StorageOption) (*LocalStorage, error) {
	if err := os.MkdirAll(root, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage root: %s", root)
	}
	ls := &LocalStorage{
		git:    "/usr/bin/git",
		root:   root,
		logger: log.NewNopLogger(),
	}

	for _, opt := range opts {
		opt(ls)
	}
	return ls, nil
}

func (s *LocalStorage) repoPath(id string) string {
	id = strings.Replace(id, "-", "", -1)
	return filepath.Join(s.root, id[:2], id[2:2], id[4:])
}

// Create a new repository
func (s *LocalStorage) Create(ctx context.Context, id string) error {
	dir := s.repoPath(id)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to repository directory: %s", dir)
	}

	cmd := exec.CommandContext(ctx, s.git, "init", "--bare")
	cmd.Dir = dir

	//TODO: Don't throw away stdout/err...
	return cmd.Run()
}

// GetRepository from Storage
// TODO: Cache these somehow?
func (s *LocalStorage) GetRepository(ctx context.Context, repoPath string) (Repository, error) {
	dir := s.repoPath(repoPath)

	cmd := exec.CommandContext(ctx, s.git, "config", "--null", "core.repositoryformatversion")
	cmd.Dir = dir

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, ErrRepoNotValid
	}
	if strings.TrimSuffix(string(output), "\x00") != "0" {
		return nil, ErrRepoNotValid
	}

	return &LocalRepository{git: s.git, path: dir}, nil
}

// Branch of a repository
type Branch struct {
	Name string
	Sha1 string
	Type string
}

// SetDescription of repository
func (r *LocalRepository) SetDescription(ctx context.Context, description string) error {
	file := filepath.Join(r.path, "description")
	return ioutil.WriteFile(file, []byte(description+"\n"), 0644)
}

// ListBranches returns all branches of a given repository
func (r *LocalRepository) ListBranches(ctx context.Context) ([]Branch, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	args := []string{"for-each-ref", "--format=%(objectname) %(objecttype) %(refname)", "refs/heads"}
	cmd := exec.CommandContext(ctx, r.git, args...)
	cmd.Dir = r.path
	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err = cmd.Start(); err != nil {
		return nil, err
	}

	var bs []Branch
	scanner := bufio.NewScanner(out)
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), " ")

		bs = append(bs, Branch{
			Name: strings.TrimPrefix(s[2], "refs/heads/"),
			Sha1: s[0],
			Type: s[1],
		})
	}

	if err := cmd.Wait(); err != nil {
		return nil, err
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
func (r *LocalRepository) GetCommit(ctx context.Context, ref string) (Commit, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	args := []string{"cat-file", "-p", ref}
	cmd := exec.CommandContext(ctx, r.git, args...)
	cmd.Dir = r.path
	out, err := cmd.StdoutPipe()
	if err != nil {
		return Commit{}, err
	}
	if err = cmd.Start(); err != nil {
		return Commit{}, err
	}

	commit, err := parseCommit(out, ref)
	if err != nil {
		return Commit{}, err
	}

	if err := cmd.Wait(); err != nil {
		return Commit{}, err
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

//TreeEntry is a file or folder at a given path in a repository
type TreeEntry struct {
	Mode   string
	Type   string
	Object string
	Path   string
}

//Tree returns the files and folders at a given ref at a path in a repository
func (r *LocalRepository) Tree(ctx context.Context, ref, path string) ([]TreeEntry, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	span, ctx := opentracing.StartSpanFromContext(ctx, "storage.Server.Tree")
	span.SetTag("ref", ref)
	span.SetTag("path", path)
	defer span.Finish()

	if !strings.HasSuffix(path, "/") {
		path += "/"
	}

	treeEntries, err := r.tree(ctx, ref, path)
	if err != nil {
		return nil, err
	}

	return treeEntries, nil
}

func (r *LocalRepository) tree(ctx context.Context, ref, path string) ([]TreeEntry, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storage.Server.tree")
	span.SetTag("ref", ref)
	span.SetTag("path", path)
	defer span.Finish()

	args := []string{"ls-tree", ref, path}
	cmd := exec.CommandContext(ctx, r.git, args...)
	cmd.Dir = r.path
	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, errors.Wrap(err, "failed to run git ls-tree")
	}
	if err = cmd.Start(); err != nil {
		return nil, errors.Wrap(err, "failed to run git ls-tree")
	}

	// TODO: The scanner takes unbounded inputs. This could cause OOMs

	var treeEntries []TreeEntry
	scanner := bufio.NewScanner(out)
	for scanner.Scan() {
		line := scanner.Text()
		te, err := parseTreeEntry(line)
		if err != nil {
			return treeEntries, errors.Wrap(err, "unable to parse tree entry line")
		}
		treeEntries = append(treeEntries, te)
	}

	if err := cmd.Wait(); err != nil {
		return nil, errors.Wrap(err, "failed to wait for command to finish")
	}

	return treeEntries, nil
}

func parseTreeEntry(s string) (TreeEntry, error) {
	tabs := strings.Split(s, "\t")
	if len(tabs) != 2 {
		return TreeEntry{}, errors.New("expected 2 tab separated inputs")
	}
	spaces := strings.Split(tabs[0], " ")
	if len(spaces) != 3 {
		return TreeEntry{}, errors.New("expected 3 space separated inputs")
	}

	return TreeEntry{
		Mode:   spaces[0],
		Type:   spaces[1],
		Object: spaces[2],
		Path:   tabs[1],
	}, nil
}
