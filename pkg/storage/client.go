package storage

import (
	"context"
	"io"
	"time"

	grpcopentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"gitlab.com/gitlab-org/gitaly/streamio"
	"google.golang.org/grpc"
)

// Client holds the gRPC-connection to the storage-server
type Client struct {
	repos    RepositoryClient
	branches BranchClient
	commits  CommitClient
	ssh      SSHClient
}

// NewClient returns a new Storage client.
func NewClient(storageAddr string) (*Client, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithUnaryInterceptor(grpcopentracing.UnaryClientInterceptor()))
	opts = append(opts, grpc.WithStreamInterceptor(grpcopentracing.StreamClientInterceptor()))

	conn, err := grpc.DialContext(context.Background(), storageAddr, opts...)
	if err != nil {
		return nil, err
	}

	return &Client{
		repos:    NewRepositoryClient(conn),
		branches: NewBranchClient(conn),
		commits:  NewCommitClient(conn),
		ssh:      NewSSHClient(conn),
	}, nil
}

// Create a repository
func (c *Client) Create(ctx context.Context, id string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storage.Client.Create")
	span.SetTag("id", id)
	defer span.Finish()

	_, err := c.repos.Create(ctx, &CreateRequest{Id: id})
	return err
}

// SetDescription of a repository
func (c *Client) SetDescription(ctx context.Context, id, description string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storage.Client.Description")
	span.SetTag("id", id)
	span.SetTag("description", description)
	defer span.Finish()

	_, err := c.repos.SetDescriptions(ctx, &SetDescriptionRequest{
		Id:          id,
		Description: description,
	})
	return err
}

// Branches returns all branches of a repository
func (c *Client) Branches(ctx context.Context, id string) ([]Branch, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storage.Client.Branches")
	span.SetTag("id", id)
	defer span.Finish()

	res, err := c.branches.List(ctx, &BranchesRequest{Id: id})

	var branches []Branch
	for _, b := range res.Branch {
		branches = append(branches, Branch{
			Name: b.Name,
			Sha1: b.Sha1,
			Type: b.Type,
		})
	}

	return branches, err
}

// CommitCount returns the number of commits that any given `ref` points to
func (c *Client) CommitCount(ctx context.Context, id, ref string) (int64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storage.Client.CommitCount")
	span.SetTag("id", id)
	span.SetTag("ref", ref)
	defer span.Finish()

	req := &CommitCountRequest{
		Id:  id,
		Ref: ref,
	}

	resp, err := c.commits.Count(ctx, req)

	return resp.GetCount(), err
}

// ListCommits returns all commits pointerd to by any given `ref`.
// `limit` sets the max amount of commits to be returned. Can not be 0.
// `skip` N*`limit` commits. Useful for pagination
func (c *Client) ListCommits(ctx context.Context, id, ref string, limit, skip int64) ([]Commit, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storage.Client.ListCommits")
	span.SetTag("id", id)
	span.SetTag("ref", ref)
	span.SetTag("limit", limit)
	span.SetTag("skip", skip)
	defer span.Finish()

	req := &CommitListRequest{
		Id:    id,
		Ref:   ref,
		Limit: limit,
		Skip:  skip,
	}

	var commits []Commit

	stream, err := c.commits.List(ctx, req)
	if err != nil {
		return nil, err
	}
	var resp *CommitListResponse
	for err == nil {
		resp, err = stream.Recv()
		if err != nil {
			break
		}
		for _, commit := range resp.GetCommits() {
			commits = append(commits, commitFromResponse(commit))
		}
	}

	if err == io.EOF {
		err = nil
	}
	return commits, err
}

// Commit returns a single commit from a given repository
func (c *Client) Commit(ctx context.Context, id, ref string) (Commit, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storage.Client.Commit")
	span.SetTag("id", id)
	span.SetTag("ref", ref)
	defer span.Finish()

	req := &CommitRequest{
		Id:  id,
		Ref: ref,
	}

	res, err := c.commits.Get(ctx, req)
	if err != nil {
		return Commit{}, err
	}

	return commitFromResponse(res), nil
}

func commitFromResponse(resp *CommitResponse) Commit {
	return Commit{
		Hash:    resp.GetHash(),
		Tree:    resp.GetTree(),
		Parent:  resp.GetParent(),
		Message: resp.GetMessage(),
		Author: Signature{
			Name:  resp.GetAuthor(),
			Email: resp.GetAuthorEmail(),
			Date:  time.Unix(resp.GetAuthorDate(), 0),
		},
		Committer: Signature{
			Name:  resp.GetCommitter(),
			Email: resp.GetCommitterEmail(),
			Date:  time.Unix(resp.GetCommitterDate(), 0),
		},
	}
}

//Tree returns the files and folders at a given ref at a path in a repository
func (c *Client) Tree(ctx context.Context, id, ref, path string) ([]TreeEntry, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storage.Client.Tree")
	span.SetTag("repo_path", id)
	span.SetTag("ref", ref)
	span.SetTag("path", path)
	defer span.Finish()

	req := &TreeRequest{
		Id:   id,
		Ref:  ref,
		Path: path,
	}

	res, err := c.repos.Tree(ctx, req)
	if err != nil {
		return nil, err
	}

	var treeEntries []TreeEntry
	for _, te := range res.TreeEntries {
		treeEntries = append(treeEntries, TreeEntry{
			Mode:   te.Mode,
			Type:   te.Type,
			Object: te.Object,
			Path:   te.Path,
		})
	}

	return treeEntries, nil
}

// UploadPack to a git-repo
func (c *Client) UploadPack(ctx context.Context, id string, stdin io.Reader, stdout, stderr io.Writer) (int32, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storage.Client.UploadPack")
	span.SetTag("repo_path", id)
	defer span.Finish()

	req := &GRERequest{Id: id}

	stream, err := c.ssh.UploadPack(ctx)
	if err != nil {
		return 0, err
	}

	if err := stream.Send(req); err != nil {
		return 0, nil
	}

	errC := make(chan error, 1)
	// Go-routine for sending stdin
	go func(errC chan error) {
		in := streamio.NewWriter(func(p []byte) error {
			return stream.Send(&GRERequest{Stdin: p})
		})
		if _, err := io.Copy(in, stdin); err != nil {
			errC <- err
		}
		if err := stream.CloseSend(); err != nil {
			errC <- err
		}
		close(errC)
	}(errC)

	var resp *GREResponse
	for ; err == nil; resp, err = stream.Recv() {
		if resp.GetExitCode() != nil {
			return resp.GetExitCode().GetExitCode(), nil
		}
		if len(resp.GetStderr()) > 0 {
			stderr.Write(resp.GetStderr())
		}
		if len(resp.GetStdout()) > 0 {
			stdout.Write(resp.GetStdout())
		}
	}
	if err == io.EOF {
		err = nil
	}

	for errIn := range errC {
		if errIn == io.EOF {
			continue
		}
		span.SetTag("error", true)
		span.LogKV("event", "error", "message", errIn)
	}

	return 0, err
}

// ReceivePack from a git-repo
func (c *Client) ReceivePack(ctx context.Context, id string, stdin io.Reader, stdout, stderr io.Writer) (int32, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storage.Client.ReceivePack")
	span.SetTag("repo_path", id)
	defer span.Finish()

	req := &GRERequest{Id: id}

	stream, err := c.ssh.ReceivePack(ctx)
	if err != nil {
		return 0, err
	}

	if err := stream.Send(req); err != nil {
		return 0, nil
	}

	errC := make(chan error, 1)
	// Go-routine for sending stdin
	go func(errC chan error) {
		in := streamio.NewWriter(func(p []byte) error {
			return stream.Send(&GRERequest{Stdin: p})
		})
		if _, err := io.Copy(in, stdin); err != nil {
			errC <- err
		}
		if err := stream.CloseSend(); err != nil {
			errC <- err
		}
		close(errC)
	}(errC)

	var resp *GREResponse
	for ; err == nil; resp, err = stream.Recv() {
		if resp.GetExitCode() != nil {
			return resp.GetExitCode().GetExitCode(), nil
		}
		if len(resp.GetStderr()) > 0 {
			stderr.Write(resp.GetStderr())
		}
		if len(resp.GetStdout()) > 0 {
			stdout.Write(resp.GetStdout())
		}
	}
	if err == io.EOF {
		err = nil
	}

	for errIn := range errC {
		if errIn == io.EOF {
			continue
		}
		span.SetTag("error", true)
		span.LogKV("event", "error", "message", errIn)
		err = errors.Wrap(err, errIn.Error())
	}

	return 0, err
}
