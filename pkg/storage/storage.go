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
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/sourcepods/sourcepods/pkg/command"
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

	// StorageOption is for injecting configuration into LocalStorage
	StorageOption func(Storage)

	// LocalStorage implements Storage for Local disk-access
	LocalStorage struct {
		git    string
		root   string
		logger log.Logger
	}

	// Repository is the interface for manipulating repos
	Repository interface {
		GetID() string
		SetDescription(ctx context.Context, description string) error
		ListBranches(ctx context.Context) ([]Branch, error)
		GetCommit(ctx context.Context, ref string) (Commit, error)
		Tree(ctx context.Context, ref, path string) ([]TreeEntry, error)
		UploadPack(ctx context.Context, stdin io.Reader, stdout, stderr io.Writer) (int32, error)
		ReceivePack(ctx context.Context, stdin io.Reader, stdout, stderr io.Writer) (int32, error)
	}

	// LocalRepository implements Repository for Local disk-access
	LocalRepository struct {
		git    string
		path   string
		id     string
		logger log.Logger
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
	return filepath.Join(s.root, id[0:2], id[2:4], id[4:])
}

// Create a new repository
func (s *LocalStorage) Create(ctx context.Context, id string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storage.LocalStorage.Create")
	span.SetTag("repo_path", id)
	defer span.Finish()
	dir := s.repoPath(id)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to repository directory: %s", dir)
	}

	errBuf := &bytes.Buffer{}
	cmd, err := command.New(ctx, dir, s.git, []string{"init", "--bare"}, command.StderrWriter(errBuf))
	if err != nil {
		injectError(span, err, "")
	}

	err = cmd.Wait()
	if err != nil {
		injectError(span, err, errBuf.String())
	}
	return err
}

// GetRepository from Storage
// TODO: Cache these somehow?
func (s *LocalStorage) GetRepository(ctx context.Context, repoPath string) (Repository, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storage.LocalStorage.GetRepository")
	span.SetTag("repo_path", repoPath)
	defer span.Finish()
	dir := s.repoPath(repoPath)

	out, err := command.NewSimple(ctx, dir, s.git, "config", "--null", "core.repositoryformatversion")
	if err != nil {
		injectError(span, err, out)
		return nil, ErrRepoNotValid
	}

	if strings.TrimSuffix(out, "\x00") != "0" {
		injectError(span, ErrRepoNotValid, out)
		return nil, ErrRepoNotValid
	}

	return &LocalRepository{git: s.git, path: dir, id: repoPath, logger: s.logger}, nil
}

// GetID returns the repos ID
func (r *LocalRepository) GetID() string {
	return r.id
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "storage.LocalRepository.ListBranches")
	defer span.Finish()

	errBuf := &bytes.Buffer{}
	args := []string{"for-each-ref", "--format=%(objectname) %(objecttype) %(refname)", "refs/heads"}
	cmd, err := command.New(ctx, r.path, r.git, args, command.StderrWriter(errBuf), command.StdoutPipe)
	if err != nil {
		injectError(span, err, "")
		return nil, err
	}

	var bs []Branch
	scanner := bufio.NewScanner(cmd.Stdout())
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), " ")

		bs = append(bs, Branch{
			Name: strings.TrimPrefix(s[2], "refs/heads/"),
			Sha1: s[0],
			Type: s[1],
		})
	}

	if err := cmd.Wait(); err != nil {
		injectError(span, err, errBuf.String())
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

func injectError(span opentracing.Span, err error, stderr string) {
	span.SetTag("error", true)
	span.LogKV("event", "error", "message", err, "stderr", stderr)
}

// GetCommit returns a single Commit from a Repository
func (r *LocalRepository) GetCommit(ctx context.Context, ref string) (Commit, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storage.LocalRepository.GetCommit")
	span.SetTag("ref", ref)
	defer span.Finish()

	errBuf := &bytes.Buffer{}
	args := []string{"cat-file", "-p", ref}
	cmd, err := command.New(ctx, r.path, r.git, args, command.StderrWriter(errBuf), command.StdoutPipe)
	if err != nil {
		injectError(span, err, errBuf.String())
		return Commit{}, err
	}
	defer cmd.Finish()

	commit, err := parseCommit(cmd.Stdout(), ref)
	if err != nil {
		injectError(span, err, errBuf.String())
		return Commit{}, err
	}

	if err := cmd.Wait(); err != nil {
		injectError(span, err, errBuf.String())
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "storage.LocalRepository.Tree")
	span.SetTag("ref", ref)
	span.SetTag("path", path)
	defer span.Finish()

	if !strings.HasSuffix(path, "/") {
		path += "/"
	}

	treeEntries, err := r.tree(ctx, ref, path)
	if err != nil {
		injectError(span, err, "")
		return nil, err
	}

	return treeEntries, nil
}

func (r *LocalRepository) tree(ctx context.Context, ref, path string) ([]TreeEntry, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storage.LocalRepository.tree")
	span.SetTag("ref", ref)
	span.SetTag("path", path)
	defer span.Finish()

	errBuf := &bytes.Buffer{}
	args := []string{"ls-tree", ref, path}
	cmd, err := command.New(ctx, r.path, r.git, args, command.StderrWriter(errBuf), command.StdoutPipe)
	if err != nil {
		injectError(span, err, errBuf.String())
		return nil, errors.Wrap(err, "failed to run git ls-tree")
	}
	defer cmd.Finish()

	// TODO: The scanner takes unbounded inputs. This could cause OOMs

	var treeEntries []TreeEntry
	scanner := bufio.NewScanner(cmd.Stdout())
	for scanner.Scan() {
		line := scanner.Text()
		te, err := parseTreeEntry(line)
		if err != nil {
			injectError(span, err, "parseTreeEntry")
			return treeEntries, errors.Wrap(err, "unable to parse tree entry line")
		}
		treeEntries = append(treeEntries, te)
	}

	if err := cmd.Wait(); err != nil {
		injectError(span, err, errBuf.String())
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

// UploadPack is a hack because we need r.path
//  int32 exitCode - The commands exit-code
//  error internalError - And internal error occured
func (r *LocalRepository) UploadPack(ctx context.Context, stdin io.Reader, stdout, stderr io.Writer) (int32, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storage.Repository.UploadPack")
	span.SetTag("repo_path", r.path)
	defer span.Finish()

	cmd, err := command.New(ctx, r.path, r.git, []string{"upload-pack", "--strict", "."},
		command.StdinWriter(stdin),
		command.StdoutWriter(stdout),
		command.StderrWriter(stderr),
	)
	if err != nil {
		return 0, errors.Wrap(err, "command failed")
	}

	return exitStatus(cmd.Wait())
}

// ReceivePack is a hack because we need r.path
//  int32 exitCode - The commands exit-code
//  error internalError - And internal error occured
func (r *LocalRepository) ReceivePack(ctx context.Context, stdin io.Reader, stdout, stderr io.Writer) (int32, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storage.Repository.ReceivePack")
	span.SetTag("repo_path", r.path)
	defer span.Finish()

	cmd, err := command.New(ctx, r.path, r.git, []string{"receive-pack", "."},
		command.StdinWriter(stdin),
		command.StdoutWriter(stdout),
		command.StderrWriter(stderr),
	)
	if err != nil {
		return 0, errors.Wrap(err, "command failed")
	}

	return exitStatus(cmd.Wait())
}

// Thankfully borrowed from https://github.com/gliderlabs/sshfront/blob/ff9cab19386c1b3bcdf1d574c5cbaf8bd046fc12/handlers.go#L25-L37
func exitStatus(err error) (int32, error) {
	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			// There is no platform independent way to retrieve
			// the exit code, but the following will work on Unix
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				return int32(status.ExitStatus()), nil
			}
		}
		return 0, err
	}
	return 0, nil
}
